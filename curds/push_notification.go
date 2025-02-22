package cruds

import (
	"bondi-push-notification/models"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
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
		results []models.PushSubscribers
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

	fmt.Print(results)
	return nil
}
