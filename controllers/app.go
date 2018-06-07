package controllers

import (
	"MonitorSys/lib"
	"MonitorSys/models"
	"encoding/json"
	"security"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

type AppController struct {
	beego.Controller
}

func Appinit() {
	apparry, _ := models.GetAppAllinside()
	for _, value := range apparry {
		if err := redisHMSET("AppID:"+strconv.FormatInt(value["AppID"].(int64), 10), value); err != nil {
			beego.Error("error:", err)
		}
	}
}

// @router /aregist [post]
func (a *AppController) AppRegist() {
	userid, _ := a.GetInt("UserID")
	accesstoken := a.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		var app models.AppInfo
		var appmap map[string]interface{}
		appmap = make(map[string]interface{})
		app.AppName = a.GetString("AppName")
		app.Password = a.GetString("Password")
		app.Region = a.GetString("Region")
		app.Operator = a.GetString("Operator")
		app.RegTime = time.Now().Format("2006-01-02 15:04:05")
		app.Astatus = "正常"
		app.Online = "下线"
		app.Password = security.Md5(app.Password + security.Md5(app.Password))
		if models.AppExist(app.AppName) != 0 {
			result := &Result{1, "AppName has exist", 0, nil}
			a.Data["json"] = result
		} else {
			AppID, _ := models.AppInfoAdd(&app) //数据库存入
			/********数据存入redis*********/
			if err := redisHMSET("AppID:"+strconv.FormatInt(AppID, 10), app); err != nil {
				beego.Error("error:", err)
			}
			/******************************/
			appmap["AppID"] = AppID
			result := &Result{0, "success", 0, appmap}
			a.Data["json"] = result
		}
	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		a.Data["json"] = result
	}

	a.ServeJSON()
}

// @router /alogin [post]
func (a *AppController) AppLogin() {
	var app models.AppInfo
	var appmap map[string]interface{}
	appmap = make(map[string]interface{})
	json.Unmarshal(a.Ctx.Input.RequestBody, &app)
	app.RegTime = time.Now().Format("2006-01-02 15:04:05")
	app.AppToken = lib.GenToken() ///////
	app.Online = "在线"
	if models.AppExist(app.AppName) == 1 {
		appid := models.AppLogin(app.AppName, app.Password, app.AppToken, app.Online) //////
		if appid == 0 {
			result := &Result{1, "wrong password", 0, nil}
			a.Data["json"] = result
		} else {
			appmap["AppID"] = appid
			appinfo, _ := redisHMGET("AppID:" + strconv.FormatInt(appid, 10))
			appinfo["AppToken"] = app.AppToken
			appinfo["Online"] = app.Online
			redisHMSET("AppID:"+strconv.FormatInt(appid, 10), appinfo) //更新token in redis
			//appmap["AppToken"] = app.AppToken //////
			result := &Result{0, "success", 0, appmap}
			a.Data["json"] = result
		}
	} else {
		result := &Result{1, "app not exist ", 0, nil}
		a.Data["json"] = result
	}
	a.ServeJSON()
}

// @router /getaall [post]
func (a *AppController) GetAll() {
	userid, _ := a.GetInt("UserID")
	accesstoken := a.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		apparry, countq := models.GetAppAll()
		result := &Result{0, "succeed ", countq, apparry}
		a.Data["json"] = result
	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		a.Data["json"] = result
	}
	a.ServeJSON()
}

// @router /update [post]
func (a *AppController) UpdateApp() {
	userid, _ := a.GetInt("UserID")
	accesstoken := a.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		var app models.AppInfo
		app.AppID, _ = a.GetInt64("AppID")
		app.AppName = a.GetString("AppName")
		app.Region = a.GetString("Region")
		app.Operator = a.GetString("Operator")
		app.RegTime = time.Now().Format("2006-01-02 15:04:05")
		app.Astatus = a.GetString("Astatus")
		app.Password = a.GetString("Password")
		V, err := redisHMGET("AppID:" + strconv.FormatInt(app.AppID, 10)) //redis中获取数据
		if err != nil {
			beego.Error("error:", err)
		}
		if app.Password != "999999" {
			app.Password = security.Md5(app.Password + security.Md5(app.Password))
			V["Password"] = app.Password
		}
		/***********修改后重新存入redis(不可直接覆盖)***************/
		V["Region"] = app.Region
		V["Operator"] = app.Operator
		V["Astatus"] = app.Astatus

		if err := redisHMSET("AppID:"+strconv.FormatInt(app.AppID, 10), V); err != nil {
			beego.Error("error:", err)
		}
		/*************************/
		affectrows, err := models.UpdateApp(&app)
		if affectrows == 1 && err == nil {
			result := &Result{0, "success", affectrows, nil}
			a.Data["json"] = result
		} else {
			result := &Result{1, "failed", 0, err}
			a.Data["json"] = result
		}
	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		a.Data["json"] = result
	}
	a.ServeJSON()
}
func VerifyAppFromRedis(appid int64, accesstoken string) bool {
	V, err := redisHMGET("AppID:" + strconv.FormatInt(appid, 10))
	if err != nil {
		beego.Error("error:", err)
	}
	if V["AppToken"] == accesstoken {
		return true
	}
	return false
}

// @router /heart [get]
func (a *AppController) AppHeart() {
	AppID, err := a.GetInt64("AppID")
	if err != nil {
		result := &Result{1, "Parameter format error", 0, nil}
		a.Data["json"] = result
	} else {
		AppOnlineMap[AppID] = 0
		if _, err := models.AppOnlineHeart(AppID, "在线"); err == nil {
			V, err := redisHMGET("AppID:" + strconv.FormatInt(AppID, 10)) //redis中获取数据
			if err != nil {
				beego.Error("controllers/appcheck error:", err)
			}
			V["Online"] = "在线"
			if err := redisHMSET("AppID:"+strconv.FormatInt(AppID, 10), V); err != nil {
				beego.Error("error:", err)
			}
		}
		result := &Result{0, "success", 0, nil}
		a.Data["json"] = result
	}
	a.ServeJSON()
}
