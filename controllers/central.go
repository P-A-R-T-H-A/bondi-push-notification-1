package controllers

import (
	"net/http"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type baseController struct {
	beego.Controller
}

// func (u *baseController) Prepare() {
// 	if unAuthorized := u.apiAuth(); unAuthorized {
// 		return
// 	}
// }
func (u *baseController) getHeaderByName(headerName string) string {
	headerValue := u.Ctx.Request.Header.Get(headerName)
	if headerValue == "" {
		headerValue = u.Ctx.Request.Header.Get(strings.ToLower(headerName))
	}
	return headerValue
}

func (u *baseController) apiAuth() bool {
	unAuthorized := false
	xApiKey, _ := beego.AppConfig.String("PUSH::XApiKey")
	if u.getHeaderByName("X-Api-Key") != xApiKey {
		u.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		u.Data["json"] = map[string]string{"msg": "Unauthorized request."}
		u.ServeJSON()
		unAuthorized = true
	}
	return unAuthorized
}
