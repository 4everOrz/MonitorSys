package controllers

import (
	"MonitoringSystemAPI/lib"
	"MonitoringSystemAPI/models"
	"encoding/json"
	"fmt"
	"security"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

type AppController struct {
	beego.Controller
}

// @Title CreateApp
// @Description create app
// @Param	body		body 	models.appInfo	true		"body for app content"
// @Success 200 {int} models.App.AppID
// @Failure 403 body is empty
// @router /aregist [post]
func (a *AppController) AppRegist() {
	var app models.AppInfo
	json.Unmarshal(a.Ctx.Input.RequestBody, &app)
	app.RegTime = time.Now().Format("2006-01-02 15:04:05")
	app.Password = security.Md5(app.Password + security.Md5(app.Password))
	if models.AppExist(app.AppName) != 0 {
		fmt.Printf("AppName has exist")
		a.Data["json"] = map[string]interface{}{"resCode": "1", "resMsg": "Change LoginName and try again", "resData": "null"}
	} else {
		AppID, _ := models.AppInfoAdd(&app)
		//Result{"1", "failed", 0, UserID}
		a.Data["json"] = map[string]interface{}{"resCode": "0", "resMsg": "success", "AppID": strconv.FormatInt(AppID, 10)}
	}
	a.ServeJSON()
}

// @Title AppLogin
// @Description for app login
// @Param	body		body 	models.appInfo	true		"body for app content"
// @Success 200 {string} AppToken
// @Failure 403 body is empty
// @router /alogin [post]
func (a *AppController) AppLogin() {
	var app models.AppInfo
	json.Unmarshal(a.Ctx.Input.RequestBody, &app)
	app.RegTime = time.Now().Format("2006-01-02 15:04:05")
	app.AppToken = lib.GenToken()
	if models.AppLogin(app.AppName, app.Password, app.AppToken) {
		a.Data["json"] = map[string]interface{}{"resCode": "0", "resMsg": "success", "AccessToken": app.AppToken}

	} else {
		a.Data["json"] = map[string]interface{}{"resCode": "1", "resMsg": "app not exist"}
	}
	a.ServeJSON()
}
