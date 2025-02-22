package controllers

import (
	cruds "bondi-push-notification/curds"
	"bondi-push-notification/models"
	"encoding/json"
)

type SubscriptionController struct {
	baseController
}

type SendPushNotification struct {
	baseController
}

func (c *SubscriptionController) Post() {
	var subscription models.PushNotificationSubscription

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subscription)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"Message": err.Error()}
		c.ServeJSON()
		return
	}

	if subscription.Notification.Endpoint == "" || subscription.Notification.Keys.P256Dh == "" || subscription.Notification.Keys.Auth == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"Message": "Endpoint, P256Dh and Auth are required"}
		c.ServeJSON()
		return
	}

	if err := cruds.StoreSubscriptionData(subscription); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"Message": err.Error()}
		c.ServeJSON()
		return
	}
	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = map[string]string{"Message": "Subscription data stored successfully"}
	c.ServeJSON()
}

func (c *SendPushNotification) Get() {
	var (
		messageId    = c.Ctx.Input.Param(":messageId")
		params       = c.Ctx.Request.URL.Query()
		requestType  = params.Get("type")
		courseId     = params.Get("id")
		studentIds   = []interface{}{}
		notification models.PushNotification
	)
	var err error
	if messageId == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"Message": "Message ID is required"}
		c.ServeJSON()
		return
	} else {
		notification, err = cruds.GetNotificationData(messageId)
		if err != nil {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]string{"Message": err.Error()}
			c.ServeJSON()
			return
		}
	}
	if requestType != "" && courseId != "" {
		studentIds, err = cruds.GetCourseWiseStudentIds(courseId)
		if err != nil {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]string{"Message": err.Error()}
			c.ServeJSON()
			return
		}
	}
	err = cruds.SendNotificationToRegisteredStudent(studentIds, notification)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"Message": err.Error()}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = studentIds
	c.ServeJSON()
}
