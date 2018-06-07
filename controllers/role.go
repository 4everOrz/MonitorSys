package controllers

import (
	"MonitorSys/models"

	"github.com/astaxie/beego"
)

type RoleController struct {
	beego.Controller
}

func Roleinit() {
	roleinfo, _, _ := models.GetRoleInfo()
	for _, value := range roleinfo {
		if err := redisHMSET("RoleID:"+value["RoleID"].(string), value); err != nil {
			beego.Error("error:", err)
		}
	}
}

// @router /rgetall [get]
func (r *RoleController) Getall() {
	roleinfo, _, _ := models.GetRoleInfo()
	r.Data["json"] = roleinfo
	r.ServeJSON()
}
