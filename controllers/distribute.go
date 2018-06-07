package controllers

import (
	"SvrMonitoringSys/monistorSys/lib"
	"SvrMonitoringSys/monistorSys/models"
	"encoding/json"

	"strconv"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

var (
	Percent           int   //异常百分比阈值
	ErrCount          int   //异常连续数目阈值
	ErrTimeDelay      int64 //异常重检间隔
	SvrNetPortAppTime = make(map[string]int64)
)

func init() {
	Percent = beego.AppConfig.DefaultInt("percent", 50) //取不到值默认50%
	ErrTimeDelay = beego.AppConfig.DefaultInt64("err.timedelay", 300)
	ErrCount = beego.AppConfig.DefaultInt("err.errcount", 25)
}

/*************************接收MQ数据并处理分发***************************/
func catchdata(body []byte) bool {
	var md models.MonitorData
	json.Unmarshal(body, &md)
	if md.AppID == 0 || md.FlagBit == "" || md.StatusCode == "" || md.NetworkProtocol == "" || md.NetworkType == "" || md.ServerID == 0 {
		beego.Error("data format error or missing some field ")
		return false
	} else {
		md.SubTime = time.Now().Format("2006-01-02 15:04:05")
		dataID, err1 := models.Add1(&md)
		if err1 != nil {
			beego.Error("something wrong on pushing in Mysql", err1)
		}
		if subscribers.Len() != 0 {
			appinfo, err1 := redisHMGET("AppID:" + strconv.FormatInt(md.AppID, 10))
			appinfo["Password"] = ""
			serverinfo, err2 := redisHMGET("ServerID:" + strconv.Itoa(md.ServerID))
			if err1 != nil || err2 != nil {
				databyte := models.GetMDbyID(dataID, md.ServerID)
				publish <- newEvent(databyte)
			} else {
				var databyte = make(map[string]interface{})
				for key, value := range appinfo {
					databyte[key] = value
					if key == "AppToken" {
						databyte["AppToken"] = ""
					}
				}
				for key2, value2 := range serverinfo {
					databyte[key2] = value2
				}
				mdmap := lib.Struct2Map(md)
				for key3, value3 := range mdmap {
					databyte[key3] = value3
				}
				datamsg := []orm.Params{databyte}
				publish <- newEvent(datamsg)
			}
		}
		go judgement(md.ServerID, md.AppID, md.StatusCode, md.FlagBit, md.Port, md.NetworkProtocol)

		return true
	}
}
func judgement(serverID int, appID int64, statuscode, flagbit, port, networkprotocol string) {
	serverid := strconv.Itoa(serverID)
	appid := strconv.FormatInt(appID, 10)
	serverinfo, _ := redisHMGET("ServerID:" + serverid)
	appinfo, _ := redisHMGET("AppID:" + appid)
	statuscodeinfo, _ := redisGet("StatusCode:" + statuscode)
	//	flagbit = samplefilter(serverID, port, flagbit, networkprotocol)
	portpercent := svrnetypeportappjudge(serverid, port, flagbit, networkprotocol, appid)
	nacy := lemonfilter(serverid, appid, flagbit, port, networkprotocol)
	if portpercent >= Percent && nacy && statuscode != "normal" && statuscode != "200" && statuscode != "4" {
		if timerparse(serverid, networkprotocol, port, appid) {
			if portpercent >= Percent && nacy && statuscode != "normal" && statuscode != "200" && statuscode != "4" {
				reson := resonjudge(serverid, appid, port, flagbit, networkprotocol)
				//	go messageSend(serverid, port, statuscode, reson) //消息通知
				beego.Informational("服务器：" + serverinfo["ServerName"] + ",网络协议：" + networkprotocol + "端口：" + port + ", 状态信息：" + statuscode + "-" + statuscodeinfo + ", 数据源地区：" + appinfo["Region"] + ", App昵称：" + appinfo["AppName"] + ",可能原因：" + reson)
			}
		}
	}
}
func lemonfilter(serverID, appid, flagbit, port, netprotocol string) bool {
	svrnetportappstr := serverID + "/" + netprotocol + "/" + port + "/" + appid
	switch flagbit {
	case "0":
		SvrNetPortAppCount[svrnetportappstr]++
	case "1":
		SvrNetPortAppCount[svrnetportappstr] = 0
		SvrNetPortAppTime[svrnetportappstr] = 0
	default:
		break
	}
	return SvrNetPortAppCount[svrnetportappstr] > ErrCount
}
func resonjudge(serverID, appID, port, flagbit, networkprotocol string) string {
	possibility := ""
	regionsvrdef, svrregiondef, regionappdef, svrnetportdef, svrappdef := mathBoss(serverID, appID, port, flagbit, networkprotocol)
	chadu := []float64{regionsvrdef, svrregiondef, regionappdef, svrnetportdef, svrappdef}
	BubbleSort(chadu)
	switch chadu[4] {
	case regionappdef:
		possibility = possibility + " 服务器故障"
	case svrregiondef:
		possibility = possibility + " 地区网络异常"
	case regionappdef, svrappdef:
		possibility = possibility + " App设备异常"
	case svrnetportdef:
		possibility = possibility + " 服务器端口异常"
	default:
		break
	}
	return possibility
}

func timerparse(serverid, netprotocol, port, appid string) bool {
	svrnetportappstr := serverid + "/" + netprotocol + "/" + port + "/" + appid
	if SvrNetPortAppTime[svrnetportappstr] == 0 {
		timeafter := time.Now().Unix() + ErrTimeDelay
		SvrNetPortAppTime[svrnetportappstr] = timeafter
		return false
	} else {
		t1 := time.Now().Unix()
		t2 := SvrNetPortAppTime[svrnetportappstr]
		if t2-t1 > 0 {
			return false
		} else {
			return true
		}
	}
}
