package cruds

import (
	"bondi-push-notification/models"
	"encoding/json"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

func GetCourseWiseStudentIds(courseId string) ([]interface{}, error) {

	o := orm.NewOrm()
	q := o.QueryTable("student_course").Filter("course_id", courseId)
	var QueryResultList []orm.ParamsList
	_, err := q.ValuesList(&QueryResultList, "student_id")
	if err != nil {
		return nil, err
	}
	var studentIds []interface{}
	for _, data := range QueryResultList {
		studentIds = append(studentIds, data[0])
	}
	return studentIds, nil
}

func GetNotificationData(messageId string) (models.PushNotification, error) {
	o := orm.NewOrm()
	var notification models.PushNotification
	err := o.QueryTable("push_notification").Filter("id", messageId).One(&notification)
	if err != nil {
		return notification, err
	}
	return notification, nil
}

func SendNotificationToRegisteredStudent(studentIds []interface{}, notification models.PushNotification) error {

	var (
		results             []models.PushSubscribers
		notificationPayload = map[string]interface{}{}
		privateKey, _       = beego.AppConfig.String("PUSH::VapidPrivateKey")
		publicKey, _        = beego.AppConfig.String("PUSH::VapidPublicKey")
	)

	o := orm.NewOrm()
	qs := o.QueryTable(new(models.PushSubscribers))
	if len(studentIds) > 0 {
		qs = qs.Filter("student_id__in", studentIds)
	}

	_, err := qs.All(&results)
	if err != nil {
		return err
	}

	notificationPayload = map[string]interface{}{
		"title": "Bondi Pathshala",
		"body":  notification.NotificationContent,
		"icon":  "https://www.bondipathshala.education/static/img/logo-mini.png",
		"image": notification.NotificationImage,
	}
	if notification.Url != "" {
		notificationPayload["link"] = notification.Url
	}
	payloadBytes, err := json.Marshal(notificationPayload)
	if err != nil {
		return err
	}

	for i, data := range results {
		var Keys = webpush.Keys{
			Auth:   data.Auth,
			P256dh: data.P256dh,
		}
		var subscription = webpush.Subscription{
			Endpoint: data.Endpoint,
			Keys:     Keys,
		}

		resp, err := webpush.SendNotification(payloadBytes, &subscription, &webpush.Options{
			Subscriber:      "mailto:admin@bondipathshala.education",
			VAPIDPrivateKey: privateKey,
			VAPIDPublicKey:  publicKey,
			TTL:             60,
		})
		if resp.StatusCode == http.StatusGone {
			o.Delete(&results[i])
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}
