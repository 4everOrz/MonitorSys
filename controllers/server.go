package controllers

import (
	"MonitorSys/models"
	"strconv"

	"github.com/astaxie/beego"
)

type ServerController struct {
	beego.Controller
}

func Serverinit() {
	serverarry, _ := models.GetSverAll()
	for _, value := range serverarry {
		if err := redisHMSET("ServerID:"+value["ServerID"].(string), value); err != nil {
			beego.Error("error:", err)
		}
	}

}

// @router /filter [get]
func (s *ServerController) GetServerFilter() {
	sstatus := "正常"
	serverarry := models.GetSverFilter(sstatus)
	s.Data["json"] = serverarry
	s.ServeJSON()
}

// @router /all [post]
func (s *ServerController) GetServerAll() {
	userid, _ := s.GetInt("UserID")
	accesstoken := s.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		serverarry, num := models.GetSverAll()
		result := &Result{0, "success", num, serverarry}
		s.Data["json"] = result

	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		s.Data["json"] = result
	}
	s.ServeJSON()
}

// @router /add [post]
func (s *ServerController) Addone() {
	userid, _ := s.GetInt("UserID")
	accesstoken := s.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		var serv models.ServerInfo
		serv.ServerAddress = s.GetString("ServerAddress")
		serv.ServerName = s.GetString("ServerName")
		serv.Port = s.GetString("Port")
		serv.Type = s.GetString("Type")
		serv.ReqInterval = s.GetString("ReqInterval")
		serv.Info = s.GetString("Info")
		serv.Sstatus = "正常"
		servID, err := models.AddSver(&serv)
		if err == nil {
			var sv map[string]interface{}
			sv = make(map[string]interface{})
			sv["ServerID"] = servID
			models.CreateTable(servID) //自动创建对应数据表
			if err := redisHMSET("ServerID:"+strconv.FormatInt(servID, 10), serv); err != nil {
				beego.Error("error:", err)
			} //存入redis
			result := &Result{0, "success", 0, sv}
			s.Data["json"] = result
		} else {
			result := &Result{1, "failed", 0, nil}
			s.Data["json"] = result
		}
	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		s.Data["json"] = result
	}
	s.ServeJSON()
}

// @router /update [post]
func (s *ServerController) Updateone() {
	var serv models.ServerInfo
	userid, _ := s.GetInt("UserID")
	accesstoken := s.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		serv.ServerID, _ = s.GetInt64("ServerID")
		serv.ServerAddress = s.GetString("ServerAddress")
		serv.ServerName = s.GetString("ServerName")
		serv.Port = s.GetString("Port")
		serv.ReqInterval = s.GetString("ReqInterval")
		serv.Type = s.GetString("Type")
		serv.Sstatus = s.GetString("Sstatus")
		serv.Info = s.GetString("Info")
		_, err := models.UpdateSver(&serv) //更新Mysql
		if err := redisHMSET("ServerID:"+strconv.FormatInt(serv.ServerID, 10), serv); err != nil {
			beego.Error("error:", err)
		} //更新redis
		if err == nil {
			result := &Result{0, "success", 0, nil}
			s.Data["json"] = result
		} else {
			result := &Result{1, "failed", 0, err}
			s.Data["json"] = result
		}
	} else {
		result := &Result{1, "user verify error check userid and token ", 0, nil}
		s.Data["json"] = result
	}
	s.ServeJSON()
}
