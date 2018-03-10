package controllers

import (
	"MonitoringSystemAPI/models"
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
// @router / [post]
func (a *AppController) Post() {
	var app models.AppInfo
	app.AppName = a.GetString("AppName")
	app.Password = a.GetString("Password")
	app.AppToken = a.GetString("AppToken")
	app.Region = a.GetString("Region")
	app.RegTime = time.Now().Format("2006-01-02 15:04:05")
	aid, _ := models.AppInfoAdd(&app)
	a.Data["json"] = map[string]string{"aid": strconv.FormatInt(aid, 10)}
	a.ServeJSON()
}
