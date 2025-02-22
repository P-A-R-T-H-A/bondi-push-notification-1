package cruds

import (
	"bondi-push-notification/models"
	"errors"

	"github.com/beego/beego/v2/client/orm"
)

func StoreSubscriptionData(subscription models.PushNotificationSubscription) error {
	o := orm.NewOrm()
	subscribedStudent := models.PushSubscribers{
		StudentId: subscription.StudentId,
		Endpoint:  subscription.Notification.Endpoint,
		P256dh:    subscription.Notification.Keys.P256Dh,
		Auth:      subscription.Notification.Keys.Auth,
	}

	_, err := o.Insert(&subscribedStudent)
	if err != nil && errors.Is(err, orm.ErrLastInsertIdUnavailable) {
		return nil
	}

	return err
}
