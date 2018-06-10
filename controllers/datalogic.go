package controllers

import (
	"fmt"
	"math"
	"strings"

	"github.com/astaxie/beego"
)

/*可在配置文件中配置数据采集窗口的大小，以及过界阈值*/
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
func mapCalculate(maparry map[string][]string, percentmap map[string]int, key, flagbit string) int {
	var rightcount int
	var errcount int
	if _, ok := maparry[key]; ok {

	} else {
		maparry[key] = make([]string, arrysize)
		arryex := maparry[key]
		for m := 0; m < arrysize; m++ {
			arryex[m] = "1"
		}
	}
	slice := maparry[key]
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
	percentmap[key] = percent //存当前百分比
	return percent
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
func chadu(percentmap map[string]int, keys ...string) float64 {
	var total, difference float64
	//	var sqrtcount float64
	var h float64
	if len(keys) == 2 {
		fmt.Println("222")
		for key, value := range percentmap {
			splitarry := strings.Split(key, "/")
			if splitarry[0] == keys[0] {
				h++
				total = total + float64(value)
			}
		}
		//	log.Println("regionsvrcount:", regionsvrcount)
		average := total / h //同地区不同服务器异常平均百分比
		pop := float64(percentmap[keys[0]+"/"+keys[1]]) - average
		if pop <= 0 {
			difference = 0
		} else {
			difference = pop //当前值与平均值之差
		}
	} else if len(keys) == 3 {
		fmt.Println("333")
		for key, value := range percentmap {
			splitarry := strings.Split(key, "/")
			if splitarry[0]+splitarry[1] == keys[0]+keys[1] {
				h++
				total = total + float64(value)
			}
		}
		//	log.Println("regionsvrcount:", regionsvrcount)
		average := total / h //同地区不同服务器异常平均百分比
		pop := float64(percentmap[keys[0]+"/"+keys[1]+"/"+keys[2]]) - average
		if pop <= 0 {
			difference = 0
		} else {
			difference = pop //当前值与平均值之差
		}
	}

	/*	for key, value := range percentmap {
			splitarry2 := strings.Split(key, "/")
			if splitarry2[0] == key1 {
				tata = float64(value) - average
				if tata < 0 {
					sqrtcount = sqrtcount + math.Sqrt(average-float64(value))
				} else {
					sqrtcount = sqrtcount + math.Sqrt(tata)
				}
			}
		}
		log.Println("sqrtcount:", sqrtcount)
		Variance := svrsqrtcount / h //方差
		StanDev := Square(Variance) //标准差*/

	return difference
}

//地区,服务器
func regionsvrjudge(serverID, flagbit, region string) int {
	return mapCalculate(RegionSvrArry, RegionSvr, region+"/"+serverID, flagbit)
}

//服务器,端口
func svrportjudge(serverID, port, flagbit string) int {
	return mapCalculate(SvrPortArry, SvrPort, serverID+"/"+port, flagbit)
}

//地区,App
func regionappjudge(appID, flagbit, region string) int {
	return mapCalculate(RegionAppArry, RegionApp, region+"/"+appID, flagbit)
}

//服务器,APP
func svrappjudge(serverID, appID, flagbit string) int {
	return mapCalculate(SvrAppArry, SvrApp, serverID+"/"+appID, flagbit)
}

//服务器,地区
func svrregionjudge(serverID, flagbit, region string) int {
	return mapCalculate(SvrRegionArry, SvrRegion, serverID+"/"+region, flagbit)
}

//服务器，网络协议，端口

func svrnetypeportjudge(serverID, port, flagbit, netprotocol string) int {
	return mapCalculate(SvrNetPortArry, SvrNetPort, serverID+"/"+netprotocol+"/"+port, flagbit)
}

//服务器，网络协议，端口，app
func svrnetypeportappjudge(serverID, port, flagbit, netprotocol, appid string) int {
	return mapCalculate(SvrNetPortAppArry, SvrNetPortApp, serverID+"/"+netprotocol+"/"+port+"/"+appid, flagbit)
}

//数据收集器
func judgecount(serverid, port, flagbit, netprotocol, appIDstr, region string) int {
	regionsvrjudge(serverid, flagbit, region)
	regionappjudge(appIDstr, flagbit, region)
	svrappjudge(serverid, appIDstr, flagbit)
	svrregionjudge(serverid, flagbit, region)
	svrnetypeportjudge(serverid, port, flagbit, netprotocol)
	return svrnetypeportappjudge(serverid, port, flagbit, netprotocol, appIDstr)
}

//差度统计
func chaducount(serverid, appIDstr, port, flagbit, netprotocol, region string) (float64, float64, float64, float64, float64) {
	/**********同地区不同服务器******************/
	regionsvrdef := chadu(RegionSvr, region, serverid)

	/*************同服务器不同地区***************/
	svrregiondef := chadu(SvrRegion, serverid, region)

	/***************同地区不同App*************/
	regionappdef := chadu(RegionApp, region, appIDstr)

	/**************同服务器不同网络协议端口**************/
	svrnetportdef := chadu(SvrNetPort, serverid, netprotocol, port)
	/***************同服务器不同app***************/
	svrappdef := chadu(SvrApp, serverid, appIDstr)

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
	/* fmt.Println(" regionsvrdef:", regionsvrdef)
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
