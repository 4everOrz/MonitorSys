package controllers

import (
	"MonitorSys/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

var AppOnlineMap = make(map[int64]int64)
var AppOuttime int64

func init() {
	AppOuttime = beego.AppConfig.DefaultInt64("app.outtime", 300) //取不到默认5分钟
	go timer2()
}

func timer2() {
	ticker := time.NewTicker(1 * time.Second) //15天 1296000  一周 604800  1天 86400
	for {
		select {
		case <-ticker.C:
			for appid, _ := range AppOnlineMap {
				AppOnlineMap[appid]++
				if AppOnlineMap[appid] >= AppOuttime {
					AppOnlineMap[appid] = 0
					_, err := models.AppOnlineHeart(appid, "下线")
					if err == nil {
						V, err := redisHMGET("AppID:" + strconv.FormatInt(appid, 10)) //redis中获取数据
						if err != nil {
							beego.Error("controllers/appcheck error:", err)
						}
						V["Online"] = "下线"
						if err := redisHMSET("AppID:"+strconv.FormatInt(appid, 10), V); err != nil {
							beego.Error("error:", err)
						}
					}
				}
			}
		}
	}
}
