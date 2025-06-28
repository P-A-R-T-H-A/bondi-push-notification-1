package cruds

import (
	"bondi-push-notification/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

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
		o                   = orm.NewOrm()
	)

	// 1) load all subscribers
	qs := o.QueryTable(new(models.PushSubscribers))
	if len(studentIds) > 0 {
		qs = qs.Filter("student_id__in", studentIds)
	}
	if _, err := qs.All(&results); err != nil {
		return err
	}

	// 2) build payload
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

	// 3) concurrent send with worker‐pool
	const maxWorkers = 10
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	errCh := make(chan error, len(results))
	expiredSubs := make(chan models.PushSubscribers, len(results))

	for i := range results {
		data := results[i]
		wg.Add(1)
		sem <- struct{}{} // acquire slot

		go func(data models.PushSubscribers) {
			defer wg.Done()
			defer func() { <-sem }() // release slot
			// prepare subscription
			keys := webpush.Keys{Auth: data.Auth, P256dh: data.P256dh}
			subscription := webpush.Subscription{Endpoint: data.Endpoint, Keys: keys}

			// send it
			resp, err := webpush.SendNotification(payloadBytes, &subscription, &webpush.Options{
				Subscriber:      "mailto:admin@bondipathshala.education",
				VAPIDPrivateKey: privateKey,
				VAPIDPublicKey:  publicKey,
				TTL:             60,
			})
			if err != nil {
				errCh <- fmt.Errorf("student %v: %w", data.StudentId, err)
				return
			}
			fmt.Println(resp.StatusCode, data.StudentId)

			// collect expired for later delete
			if resp.StatusCode == http.StatusGone {
				expiredSubs <- data
			}
		}(data)
	}

	wg.Wait()
	close(errCh)
	close(expiredSubs)

	// 4) if any send error, return it
	if sendErr := <-errCh; sendErr != nil {
		return sendErr
	}

	// 5) batch-delete expired subscriptions
	for data := range expiredSubs {
		if _, err := o.
			QueryTable(new(models.PushSubscribers)).
			Filter("endpoint", data.Endpoint).
			Filter("auth", data.Auth).
			Filter("p256dh", data.P256dh).
			Delete(); err != nil {
			return err
		}
	}

	return nil
}
