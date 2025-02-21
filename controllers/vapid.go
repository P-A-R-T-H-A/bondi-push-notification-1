package controllers

import (
	"bondi-push-notification/models"

	beego "github.com/beego/beego/v2/server/web"
)

type VapidController struct {
	baseController
}

func (c *VapidController) Get() {
	publicKey, _ := beego.AppConfig.String("PUSH::VapidPublicKey")
	c.Data["json"] = models.VapidKeyResponse{
		PublicKey: publicKey,
	}
	c.ServeJSON()
}
