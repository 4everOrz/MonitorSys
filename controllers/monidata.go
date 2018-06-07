package controllers

import (
	"MonitorSys/models"
	"strconv"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

type MdController struct {
	beego.Controller
}

// @router / [post]
func (md *MdController) Post() {
	//	header:= md.Ctx.Request.Header
	//	appid:=header["AppID"]
	//apptoken:=header["AppToken"]
	ms := md.Ctx.Input.RequestBody

	length := len(ms)
	if length > 120 && length < 280 && ms[0] == 123 && ms[length-1] == 125 {
		if catchdata(ms) {
			result := &Result{0, "success", 0, nil}
			md.Data["json"] = result
		} else {
			result := &Result{1, "failed", 0, nil}
			md.Data["json"] = result
		}
	} else {
		result := &Result{1, "Unlawful data", 0, nil}
		md.Data["json"] = result
	}
	md.ServeJSON() /**/

}

// @router /getone [post]
func (md *MdController) GetData() {
	dataid := md.GetString("DataID")
	id64, _ := strconv.ParseInt(dataid, 10, 64)
	maps := models.GetOne(id64)
	md.Data["json"] = maps
	md.ServeJSON()
}

// @router /getall [post]
func (md *MdController) GetAll() {
	userid, _ := md.GetInt("UserID")
	accesstoken := md.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		page, _ := md.GetInt("Page")
		pagesize, _ := md.GetInt("PageSize")
		//appid, _ := md.GetInt("AppID")
		region := md.GetString("Region")
		serverid, _ := md.GetInt("ServerID")
		time1 := md.GetString("Time1")
		time2 := md.GetString("Time2")
		//	data, total := models.GetMDbyPage2(page, pagesize, appid, serverid, time1, time2)
		data, total := models.GetMDbyPage3(page, pagesize, serverid, region, time1, time2)
		result := &Result{0, "success", total, data}
		md.Data["json"] = result
	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		md.Data["json"] = result
	}
	md.ServeJSON()
}

//@router /getfilter [post]
func (md *MdController) GetAllFilter() {
	var data []orm.Params
	var total interface{}
	userid, _ := md.GetInt("UserID")
	accesstoken := md.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		page, _ := md.GetInt("Page")
		pagesize, _ := md.GetInt("PageSize")
		//appid, _ := md.GetInt("AppID")
		region := md.GetString("Region")
		operator := md.GetString("Operator")
		networktype := md.GetString("NetworkType")
		networkprotocol := md.GetString("NetworkProtocol")
		timedelay := md.GetString("TimeDelay")
		serverid, _ := md.GetInt("ServerID")
		time1 := md.GetString("Time1")
		time2 := md.GetString("Time2")
		data, total = models.GetMDFliter8(page, pagesize, serverid, time1, time2, region, networktype, networkprotocol, operator, timedelay) //多条件筛选
		result := &Result{0, "success", total, data}
		md.Data["json"] = result
	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		md.Data["json"] = result
	}
	md.ServeJSON()
}
