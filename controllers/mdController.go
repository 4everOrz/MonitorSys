package controllers

import (
	"MonitoringSystemAPI/Data"

	"github.com/astaxie/beego"
)

type MdController struct {
	beego.Controller
}

// @Title Senddata
// @Description send data
// @Param	body		body 	models.monitorData	true "DSD"
// @Success 200 {string} context
// @Failure 403 body is empty
// @router / [post]
func (md *MdController) Post() {
	//	var m lib.Mqfeild
	Data.MqSend(md.Ctx.Input.RequestBody)
	//go lib.MqReceive()
	ms := "send to mq OK!"
	md.Data["json"] = ms
	md.ServeJSON()

}

// @Title recdata
// @Description receive data
// @Param	body		body 	models.monitorData	true "DSD"
// @Success 200 body is empty
// @Failure 403 body is empty
// @router / [get]
func (md *MdController) Get() {

	Data.MqReceive()
	ms := "got"
	md.Data["json"] = ms
	md.ServeJSON()

}
