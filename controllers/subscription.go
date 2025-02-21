package controllers

import (
	"bondi-push-notification/models"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type SubscriptionController struct {
	beego.Controller
}

func (c *SubscriptionController) Post() {
	var subscription models.PushNotificationSubscription

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subscription)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"msg": err.Error()}
		c.ServeJSON()
		return
	}

	if subscription.Notification.Endpoint == "" || subscription.Notification.Keys.P256Dh == "" || subscription.Notification.Keys.Auth == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"msg": "Endpoint, P256Dh and Auth are required"}
		c.ServeJSON()
		return
	}

}
