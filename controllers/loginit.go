package controllers

import (
	"github.com/astaxie/beego"
)

func init() {
	beego.SetLogger("file", `{"filename":"logs/monitorsys.log"}`) //输出到文件和控制台
	beego.SetLogFuncCall(true)                                    //是否带有文件名和行号
	beego.SetLevel(beego.LevelDebug)                              //设置日志级别
	/*LevelEmergency    级别依次降低
	  LevelAlert
	  LevelCritical
	  LevelError
	  LevelWarning
	  LevelNotice
	  LevelInformational
	  LevelDebug*/
}
