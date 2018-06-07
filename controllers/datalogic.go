package controllers

import (
	"fmt"
	"math"
	"strings"

	"github.com/astaxie/beego"
)

/*可在配置文件中配置数据采集窗口的大小，以及过界阈值*/

var (
	ServerArry           = make(map[int][]string)
	SvrPortArry          = make(map[string][]string)
	RegionSvrArry        = make(map[string][]string)
	RegionAppArry        = make(map[string][]string)
	SvrAppArry           = make(map[string][]string)
	SvrRegionArry        = make(map[string][]string)
	SvrNetPortArry       = make(map[string][]string) //服务器，协议，端口
	SvrNetPortAppArry    = make(map[string][]string) //服务器，协议，端口,app
	RegionSvr            = make(map[string]int)      //地区，服务器
	SvrRegion            = make(map[string]int)      //服务器，地区
	RegionApp            = make(map[string]int)
	SvrPort              = make(map[string]int)
	SvrApp               = make(map[string]int)
	SvrNetPort           = make(map[string]int)
	SvrNetPortApp        = make(map[string]int)
	SvrNetPortAppCount   = make(map[string]int)
	arrysize, filtersize int
)

//var SvrPortArryMap = make(map[string][]string)
//var RegionSvrArryMap = make(map[string][]string)
//var RegionAppArryMap = make(map[string][]string)
//var SvrAppArryMap = make(map[string][]string)
//var SvrRegionArryMap = make(map[string][]string)
//var SvrNetPortArryMAP = make(map[string][]string)
var SvrNetPortFilter = make(map[string][]string)

func init() {
	arrysize = beego.AppConfig.DefaultInt("arrysize", 100)
	filtersize = beego.AppConfig.DefaultInt("filtersize", 5)
}

func samplefilter(serverID, port, flagbit, netprotocol string) string { //简单滤波
	var rightcount int
	var errcount int
	svnetprotocolport := serverID + "/" + netprotocol + "/" + port
	if _, ok := SvrNetPortFilter[svnetprotocolport]; ok {

	} else {
		SvrNetPortFilter[svnetprotocolport] = make([]string, filtersize)
		arryex := SvrNetPortFilter[svnetprotocolport]
		for m := 0; m < filtersize; m++ {
			arryex[m] = "1"
		}
	}
	slice := SvrNetPortFilter[svnetprotocolport]
	for i := (len(slice) - 1); i > 0; i-- {
		slice[i] = slice[i-1]
	}
	slice[0] = flagbit
	for _, value := range slice {
		if value == "1" {
			rightcount++
		} else {
			errcount++
		}
	}
	percent := rightcount * 100 / (errcount + rightcount)
	if percent >= 50 {
		flagbit = "1"
	} else {
		flagbit = "0"
	}
	return flagbit
} /**/

//地区,服务器
func regionsvrjudge(serverID, appID, flagbit string) int {
	var rightcount int
	var errcount int
	appinfo, _ := redisHMGET("AppID:" + appID)
	regionsvr := appinfo["Region"] + "/" + serverID //服务器+地区
	if _, ok := RegionSvrArry[regionsvr]; ok {

	} else {

		RegionSvrArry[regionsvr] = make([]string, arrysize)
		arryex := RegionSvrArry[regionsvr]
		for m := 0; m < arrysize; m++ {
			arryex[m] = "1"
		}
	}
	slice := RegionSvrArry[regionsvr]
	for i := (len(slice) - 1); i > 0; i-- {
		slice[i] = slice[i-1]
	}
	slice[0] = flagbit
	for _, value := range slice {
		if value == "1" {
			rightcount++
		} else {
			errcount++
		}
	}
	percent := errcount * 100 / (errcount + rightcount)
	//	beego.Informational("Region/Server:", regionsvr, ",Percent:", percent, "%, RegionSvrInfo:", RegionSvrArry[regionsvr])
	RegionSvr[regionsvr] = percent //存当前百分比
	return percent
}

//服务器,端口
func svrportjudge(serverID, port, flagbit string) int {
	var rightcount int
	var errcount int
	svport := serverID + "/" + port
	if _, ok := SvrPortArry[svport]; ok {

	} else {
		SvrPortArry[svport] = make([]string, arrysize)
		arryex := SvrPortArry[svport]
		for m := 0; m < arrysize; m++ {
			arryex[m] = "1"
		}
	}
	slice := SvrPortArry[svport]
	for i := (len(slice) - 1); i > 0; i-- {
		slice[i] = slice[i-1]
	}
	slice[0] = flagbit
	for _, value := range slice {
		if value == "1" {
			rightcount++
		} else {
			errcount++
		}
	}
	percent := errcount * 100 / (errcount + rightcount)
	//	beego.Informational("Server/Port:", svport, ",Percent:", percent, "%, ServerPortInfo:", SvrPortArry[svport])
	SvrPort[svport] = percent //存当前百分比
	return percent
}

//地区,App
func regionappjudge(appID, flagbit string) int {
	var rightcount int
	var errcount int
	appinfo, _ := redisHMGET("AppID:" + appID)
	regionapp := appinfo["Region"] + "/" + appID
	if _, ok := RegionAppArry[regionapp]; ok {

	} else {
		RegionAppArry[regionapp] = make([]string, arrysize)
		arryex := RegionAppArry[regionapp]
		for m := 0; m < arrysize; m++ {
			arryex[m] = "1"
		}
	}
	slice := RegionAppArry[regionapp]
	for i := (len(slice) - 1); i > 0; i-- {
		slice[i] = slice[i-1]
	}
	slice[0] = flagbit
	for _, value := range slice {
		if value == "1" {
			rightcount++
		} else {
			errcount++
		}
	}
	percent := errcount * 100 / (errcount + rightcount)
	//	beego.Informational("Region/AppID:", regionapp, ",Percent:", percent, "%, RegionAppInfo:", RegionAppArry[regionapp])
	RegionApp[regionapp] = percent //存当前百分比
	return percent
}

//服务器,APP
func svrappjudge(serverID, appID, flagbit string) int {
	var rightcount int
	var errcount int
	svrapp := serverID + "/" + appID //服务器+地区
	if _, ok := SvrAppArry[svrapp]; ok {

	} else {

		SvrAppArry[svrapp] = make([]string, arrysize)
		arryex := SvrAppArry[svrapp]
		for m := 0; m < arrysize; m++ {
			arryex[m] = "1"
		}
	}
	slice := SvrAppArry[svrapp]
	for i := (len(slice) - 1); i > 0; i-- {
		slice[i] = slice[i-1]
	}
	slice[0] = flagbit
	for _, value := range slice {
		if value == "1" {
			rightcount++
		} else {
			errcount++
		}
	}
	percent := errcount * 100 / (errcount + rightcount)
	//	beego.Informational("Server/App:", svrapp, ",Percent:", percent, "%, SvrAppInfo:", SvrAppArry[svrapp])
	SvrApp[svrapp] = percent //存当前百分比
	return percent
}

//服务器,地区
func svrregionjudge(serverID, appID, flagbit string) int {
	var rightcount int
	var errcount int
	appinfo, _ := redisHMGET("AppID:" + appID)
	svrregion := serverID + "/" + appinfo["Region"] //服务器+地区
	if _, ok := SvrRegionArry[svrregion]; ok {

	} else {

		SvrRegionArry[svrregion] = make([]string, arrysize)
		arryex := SvrRegionArry[svrregion]
		for m := 0; m < arrysize; m++ {
			arryex[m] = "1"
		}
	}
	slice := SvrRegionArry[svrregion]
	for i := (len(slice) - 1); i > 0; i-- {
		slice[i] = slice[i-1]
	}
	slice[0] = flagbit
	for _, value := range slice {
		if value == "1" {
			rightcount++
		} else {
			errcount++
		}
	}
	percent := errcount * 100 / (errcount + rightcount)
	//	beego.Informational("Server/Region:", svrregion, ",Percent:", percent, "%, SvrRegionInfo:", SvrRegionArry[svrregion])
	SvrRegion[svrregion] = percent //存当前百分比
	return percent
}

//服务器，网络协议，端口

func svrnetypeportjudge(serverID, port, flagbit, netprotocol string) int {
	var rightcount int
	var errcount int
	svnetprotocolport := serverID + "/" + netprotocol + "/" + port
	if _, ok := SvrNetPortArry[svnetprotocolport]; ok {

	} else {
		SvrNetPortArry[svnetprotocolport] = make([]string, arrysize)
		arryex := SvrNetPortArry[svnetprotocolport]
		for m := 0; m < arrysize; m++ {
			arryex[m] = "1"
		}
	}
	slice := SvrNetPortArry[svnetprotocolport]
	for i := (len(slice) - 1); i > 0; i-- {
		slice[i] = slice[i-1]
	}
	slice[0] = flagbit
	for _, value := range slice {
		if value == "1" {
			rightcount++
		} else {
			errcount++
		}
	}
	percent := errcount * 100 / (errcount + rightcount)
	//	beego.Informational("Server/netprotocol/Port:", svnetprotocolport, ",Percent:", percent, "%, ServerNetprotocolPortInfo:", SvrNetPortArry[svnetprotocolport])
	SvrNetPort[svnetprotocolport] = percent //存当前百分比
	return percent
}

//服务器，网络协议，端口，app
func svrnetypeportappjudge(serverID, port, flagbit, netprotocol, appid string) int {
	var rightcount int
	var errcount int
	svnetprotocolportapp := serverID + "/" + netprotocol + "/" + port + "/" + appid
	if _, ok := SvrNetPortAppArry[svnetprotocolportapp]; ok {

	} else {
		SvrNetPortAppArry[svnetprotocolportapp] = make([]string, arrysize)
		arryex := SvrNetPortAppArry[svnetprotocolportapp]
		for m := 0; m < arrysize; m++ {
			arryex[m] = "1"
		}
	}
	slice := SvrNetPortAppArry[svnetprotocolportapp]
	for i := (len(slice) - 1); i > 0; i-- {
		slice[i] = slice[i-1]
	}
	slice[0] = flagbit
	for _, value := range slice {
		if value == "1" {
			rightcount++
		} else {
			errcount++
		}
	}
	percent := errcount * 100 / (errcount + rightcount)

	beego.Informational("ServerID/Netprotocol/Port/AppID:", svnetprotocolportapp, ",Percent:", percent, "%, ServerNetprotocolPortAppInfo:", SvrNetPortAppArry[svnetprotocolportapp])
	SvrNetPortApp[svnetprotocolportapp] = percent
	return percent
}

//综合统计
func mathBoss(serverid, appIDstr, port, flagbit, netprotocol string) (float64, float64, float64, float64, float64) {
	var h, e, l, o, w float64
	var regionsvrcount, svrregioncount, regionappcount, svrnetportcount, svrappcount float64
	var regionsvrdef, svrregiondef, regionappdef, svrnetportdef, svrappdef float64
	//	var regionsvrsqrtcount, svrregionsqrtcount, regionappsqrtcount, svrnetportsqrtcount, svrappsqrtcount float64

	regionsvrjudge(serverid, appIDstr, flagbit)
	regionappjudge(appIDstr, flagbit)
	svrappjudge(serverid, appIDstr, flagbit)
	svrregionjudge(serverid, appIDstr, flagbit)
	svrnetypeportjudge(serverid, port, flagbit, netprotocol)
	//	fmt.Println("RegionSvr:", RegionSvr)
	//	fmt.Println("SvrNetPort:", SvrNetPort)
	//	fmt.Println("RegionApp:", RegionApp)
	//	fmt.Println("SvrRegion:", SvrRegion)
	//	fmt.Println("SvrApp:", SvrApp)

	appinfo, _ := redisHMGET("AppID:" + appIDstr)

	/**********同地区不同服务器******************/
	for key, value := range RegionSvr {
		Region1 := strings.Split(key, "/")
		if Region1[0] == appinfo["Region"] {
			h++
			regionsvrcount = regionsvrcount + float64(value)
		}
	}
	//	log.Println("regionsvrcount:", regionsvrcount)
	regionsvravr := regionsvrcount / h //同地区不同服务器异常平均百分比
	pop := float64(RegionSvr[appinfo["Region"+"/"+serverid]]) - regionsvravr
	if pop <= 0 {
		regionsvrdef = 0
	} else {
		regionsvrdef = pop //当前值与平均值之差
	}

	/*	for key, value := range RegionSvr {
			Region1 := strings.Split(key, "/")
			if Region1[0] == appinfo["Region"] {
				tata = float64(value) - regionsvravr
				if tata < 0 {
					regionsvrsqrtcount = regionsvrsqrtcount + math.Sqrt(regionsvravr-float64(value))
				} else {
					regionsvrsqrtcount = regionsvrsqrtcount + math.Sqrt(tata)
				}
			}
		}
		log.Println("regionsvrsqrtcount:", regionsvrsqrtcount)
		regionsvrVariance := regionsvrsqrtcount / h //方差
		regionsvrStanDev := Square(regionsvrVariance) //标准差*/
	/*************同服务器不同地区***************/
	for key, value := range SvrRegion {
		Server1 := strings.Split(key, "/")
		if Server1[0] == serverid {
			e++
			svrregioncount = svrregioncount + float64(value)
		}
	}
	//	log.Println("svrregioncount:", svrregioncount)
	svrregionavr := svrregioncount / e //同服务器不同地区平均百分比
	pop = float64(SvrRegion[serverid+"/"+appinfo["Region"]]) - svrregionavr
	if pop <= 0 {
		svrregiondef = 0
	} else {
		svrregiondef = pop
	}

	/*	for key, value := range SvrRegion {
			Server1 := strings.Split(key, "/")
			if Server1[0] == serverid {
				tata = float64(value) - svrregionavr
				if tata < 0 {
					svrregionsqrtcount = svrregionsqrtcount + math.Sqrt(svrregionavr-float64(value))
				} else {
					svrregionsqrtcount = svrregionsqrtcount + math.Sqrt(tata)
				}

			}
		}
		log.Println("svrregionsqrtcount:", svrregionsqrtcount)
		svrregionVariance := svrregionsqrtcount / e
		svrregionStanDev := Square(svrregionVariance) */
	/***************同地区不同App*************/
	for key, value := range RegionApp {
		Region1 := strings.Split(key, "/")
		if Region1[0] == appinfo["Region"] {
			l++
			regionappcount = regionappcount + float64(value)
		}
	}
	//	log.Println(" regionappcount:", regionappcount)
	regionappavr := regionappcount / l //同地区不同App平均百分比
	pop = float64(RegionApp[appinfo["Region"]+"/"+appIDstr]) - regionappavr
	if pop <= 0 {
		regionappdef = 0
	} else {
		regionappdef = pop
	}
	/*	for key, value := range RegionApp {
			Region1 := strings.Split(key, "/")
			if Region1[0] == appinfo["Region"] {
				tata = float64(value) - regionappavr
				if tata < 0 {
					regionappsqrtcount = regionappsqrtcount + math.Sqrt(regionappavr-float64(value))
				} else {
					regionappsqrtcount = regionappsqrtcount + math.Sqrt(tata)
				}

			}
		}
		log.Println(" regionappsqrtcount:", regionappsqrtcount)
		regionappVariance := regionappsqrtcount / l
		regionappStanDev := Square(regionappVariance)*/
	/**************同服务器不同网络协议端口**************/
	for key, value := range SvrNetPort {
		Server1 := strings.Split(key, "/")
		if Server1[0]+"/"+Server1[1] == serverid+"/"+netprotocol {
			o++
			svrnetportcount = svrnetportcount + float64(value)
		}
	}
	//log.Println("svrnetportcount:", svrnetportcount)
	svrnetportavr := svrnetportcount / o //同服务器不同地区平均百分比
	pop = float64(SvrNetPort[serverid+"/"+netprotocol+"/"+port]) - svrnetportavr
	if pop <= 0 {
		svrnetportdef = 0
	} else {
		svrnetportdef = pop
	}
	/*	for key, value := range SvrNetPort {
			Server1 := strings.Split(key, "/")
			if Server1[0] == serverid {
				tata = float64(value) - svrnetportavr
				if tata < 0 {
					svrnetportsqrtcount = svrnetportsqrtcount + math.Sqrt(svrnetportavr-float64(value))
				} else {
					svrnetportsqrtcount = svrnetportsqrtcount + math.Sqrt(tata)
				}

			}
		}
		log.Println("svrnetportsqrtcount:", svrnetportsqrtcount)
		svrnetportVariance := svrnetportsqrtcount / o //同服务器不同地区平均百分比
		svrnetportStanDev := Square(svrnetportVariance)*/
	/*************同服务器不同app***************/
	for key, value := range SvrApp {
		Server1 := strings.Split(key, "/")
		if Server1[0] == serverid {
			w++
			svrappcount = svrappcount + float64(value)
		}
	}
	//log.Println("svrappcount:", svrappcount)
	svrappavr := svrappcount / w //同服务器不同App平均百分比
	pop = float64(SvrApp[serverid+"/"+appIDstr]) - svrappavr
	if pop <= 0 {
		svrappdef = 0
	} else {
		svrappdef = pop
	}
	/*	for key, value := range SvrApp {
			Server1 := strings.Split(key, "/")
			if Server1[0] == serverid {
				tata = float64(value) - svrappavr
				if tata < 0 {
					svrappsqrtcount = svrappsqrtcount + math.Sqrt(svrappavr-float64(value))
				} else {
					svrappsqrtcount = svrappsqrtcount + math.Sqrt(tata)
				}

			}
		}
		log.Println("svrappsqrtcount:", svrappsqrtcount)
		svrappVariance := svrappsqrtcount / w
		svrappStanDev := Square(svrappVariance)*/
	/********************************************/
	/*fmt.Println(" regionsvrVariance:", regionsvrVariance)
	fmt.Println(" svrregionVariance:", svrregionVariance)
	fmt.Println(" regionappVariance:", regionappVariance)
	fmt.Println(" svrnetportVariance:", svrnetportVariance)
	fmt.Println("svrappVariance:", svrappVariance)
	return regionsvrVariance, svrregionVariance, regionappVariance, svrnetportVariance, svrappVariance */
	/*	fmt.Println(" regionsvravr:", regionsvravr)
		fmt.Println(" svrregionVariance:", svrregionavr)
		fmt.Println(" regionappVariance:", regionappavr)
		fmt.Println(" svrnetportVariance:", svrnetportavr)
		fmt.Println("svrappVariance:", svrappavr)
		return regionsvravr, svrregionavr, regionappavr, svrnetportavr, svrappavr */
	/*	fmt.Println(" regionsvrStanDev:", regionsvrStanDev)
		fmt.Println(" svrregionStanDev:", svrregionStanDev)
		fmt.Println(" regionappStanDev:", regionappStanDev)
		fmt.Println(" svrnetportStanDeve:", svrnetportStanDev)
		fmt.Println("svrappStanDev:", svrappStanDev)
		return regionsvrStanDev, svrregionStanDev, regionappStanDev, svrnetportStanDev, svrappStanDev */
	/*	fmt.Println(" regionsvrdef:", regionsvrdef)
		fmt.Println(" svrregiondef:", svrregiondef)
		fmt.Println(" regionappdef:", regionappdef)
		fmt.Println(" svrnetportdef:", svrnetportdef)
		fmt.Println("svrappdef:", svrappdef)*/
	return regionsvrdef, svrregiondef, regionappdef, svrnetportdef, svrappdef
}

//开平方
func Square(x float64) float64 {
	z := 1.0
	for {
		tmp := z - (z*z-x)/(2*z)
		if tmp == z || math.Abs(tmp-z) < 0.000000000001 {
			break
		}
		z = tmp
	}
	fmt.Println("z:", z)
	return z
}

//冒泡排序法
func BubbleSort(values []float64) {
	flag := true
	for i := 0; i < len(values)-1; i++ {
		flag = true
		for j := 0; j < len(values)-i-1; j++ {
			if values[j] > values[j+1] {
				values[j], values[j+1] = values[j+1], values[j]
				flag = false
			}
		}
		if flag == true {
			break
		}
	}
}
