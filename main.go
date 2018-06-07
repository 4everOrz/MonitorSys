package main

import (
	"MonitorSys/controllers"
	"MonitorSys/models"
	_ "MonitorSys/routers"

	"github.com/astaxie/beego"
)

func main() {

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	models.OrmInit()        //Mysql数据库初始化
	controllers.RedisInit() //向redis里初始化存储数据
	beego.Run()
}
