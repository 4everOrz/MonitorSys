package main

import (
	_ "MonitoringSystemAPI/routers"
	"MonitoringSystemAPI/lib"
	"MonitoringSystemAPI/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/http"
)

func main() {

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	models.Init()
	beego.InsertFilter("/platform/*", beego.BeforeRouter, func(ctx *context.Context) {
		cookie, err := ctx.Request.Cookie("Authorization")
		if err != nil || !lib.CheckToken(cookie.Value) {
			http.Redirect(ctx.ResponseWriter, ctx.Request, "/", http.StatusMovedPermanently)
		}
	})
	beego.Run()
}
