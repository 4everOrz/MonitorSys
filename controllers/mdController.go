package controllers

import (
	"MonitoringSystemAPI/Data"
	"MonitoringSystemAPI/models"
	"log"
	"strconv"

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

// @Title getdata
// @Description get data
// @Param	body		body 	models.monitorData	true "DSD"
// @Success 200 body is empty
// @Failure 403 body is empty
// @router /getone [post]
func (md *MdController) GetData() {
	dataid := md.GetString("DataID")
	id64, _ := strconv.ParseInt(dataid, 10, 64)
	log.Println("dataid", id64)
	maps := models.GetOne(id64)
	md.Data["json"] = maps
	md.ServeJSON()
}

// @Title getdata
// @Description get data
// @Param	body		body 	models.monitorData	true "DSD"
// @Success 200 body is empty
// @Failure 403 body is empty
// @router /getall [post]
func (md *MdController) GetAll() {
	pa1ge := md.GetString("Page")
	pa1gesize := md.GetString("PageSize")
	page, _ := strconv.Atoi(pa1ge)
	pagesize, _ := strconv.Atoi(pa1gesize)
	var tt []interface{} = []interface{}{"", 0}
	data, total := models.GetMDbyPage(page, pagesize, tt)
	Res := &Result{0, "success", total, data}

	md.Data["json"] = Res
	md.ServeJSON()
}
