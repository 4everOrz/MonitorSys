package controllers

import (
	"MonitoringSystemAPI/models"

	"github.com/astaxie/beego"
)

type ServerController struct {
	beego.Controller
}

// @Title get server
// @Description get server all by filter
// @Param	body		body 	models.ServerInfo	true "DSD"
// @Success 200 body is empty
// @Failure 403 body is empty
// @router /filter [get]
func (s *ServerController) GetServerFilter() {
	serverarry := models.GetSverFilter()
	s.Data["json"] = serverarry
	s.ServeJSON()
}

// @Title get server
// @Description get server all
// @Param	body		body 	models.ServerInfo	true "DSD"
// @Success 200 body is empty
// @Failure 403 body is empty
// @router /all [get]
func (s *ServerController) GetServerAll() {
	serverarry := models.GetSverAll()
	s.Data["json"] = serverarry
	s.ServeJSON()
}

// @Title 添加服务器信息
// @Description 添加服务器信息
// @Param	body		body 	models.ServerInfo	true "DSD"
// @Success 200 body is empty
// @Failure 403 body is empty
// @router /add [post]
func (s *ServerController) Addone() {
	var serv models.ServerInfo
	serv.ServerAddress = s.GetString("ServerAddress")
	serv.ServerName = s.GetString("ServerName")
	serv.Port = s.GetString("Port")
	serv.Delay = s.GetString("Delay")
	servID, err := models.AddSver(&serv)
	if err == nil {
		var sv map[string]interface{}
		sv = make(map[string]interface{})
		sv["ServerID"] = servID
		result := &Result{0, "success", 0, sv}
		s.Data["json"] = result
	} else {
		result := &Result{1, "failed", 0, nil}
		s.Data["json"] = result
	}
	s.ServeJSON()
}
