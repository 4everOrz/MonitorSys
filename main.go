package main

import (
	"MonitoringSystemAPI/lib"
	"MonitoringSystemAPI/models"
	_ "MonitoringSystemAPI/routers"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
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
