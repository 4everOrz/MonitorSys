package controllers

import (
	"MonitorSys/models"
	"time"

	"github.com/astaxie/beego"
)

//定期清理过期数据
var Durint64 int64

func init() {
	Durint64 = beego.AppConfig.DefaultInt64("duringtime", 2592000)
	go timer()
}
func timer() {
	ticker := time.NewTicker(86400 * time.Second) //15天 1296000  一周 604800  1天 86400
	for {
		select {
		case <-ticker.C:
			cleardata()
		}
	}
}
func cleardata() {
	timenow := time.Now().Unix()
	timebefore := time.Unix(timenow-Durint64, 0).Format("2006-01-02 15:04:05")
	for key, _ := range ServerArry {
		if err := models.DeleteData(key, timebefore); err != nil {
			beego.Error("delete data failed:", err)
		}
	}
}
