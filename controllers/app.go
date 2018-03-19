package controllers

import (
	"MonitoringSystemAPI/models"
	"fmt"
	"security"
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
	var appmap map[string]interface{}
	appmap = make(map[string]interface{})
	app.AppName = a.GetString("AppName")
	app.Password = a.GetString("Password")
	app.Region = a.GetString("Region")
	app.RegTime = time.Now().Format("2006-01-02 15:04:05")
	app.Password = security.Md5(app.Password + security.Md5(app.Password))
	if models.AppExist(app.AppName) != 0 {
		fmt.Printf("AppName has exist")
		result := &Result{1, "AppName has exist", 0, nil}
		a.Data["json"] = result
	} else {
		AppID, _ := models.AppInfoAdd(&app)
		appmap["AppID"] = AppID
		result := &Result{0, "success", 0, appmap}
		a.Data["json"] = result
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
	var appmap map[string]interface{}
	appmap = make(map[string]interface{})
	app.AppName = a.GetString("AppName")
	app.Password = a.GetString("Password")
	app.RegTime = time.Now().Format("2006-01-02 15:04:05")
	if models.AppExist(app.AppName) == 1 {
		appid := models.AppLogin(app.AppName, app.Password)
		if appid == 0 {
			result := &Result{1, "wrong password", 0, nil}
			a.Data["json"] = result
		} else {
			appmap["AppID"] = appid
			result := &Result{0, "success", 0, appmap}
			a.Data["json"] = result
		}
	} else {
		result := &Result{1, "app not exist ", 0, nil}
		a.Data["json"] = result
	}
	a.ServeJSON()
}
