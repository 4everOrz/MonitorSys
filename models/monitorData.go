package models

import (
	"MonitorSys/lib"
	"fmt"
	"log"
	"strconv"

	"github.com/astaxie/beego"

	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type MonitorData struct {
	DataID          int64  `orm:"pk;column(DataID)"`
	AppID           int64  `orm:"column(AppID)"`
	ServerID        int    `orm:"column(ServerID)"`
	Port            string `orm:"column(Port)"`
	NetworkType     string `orm:"column(NetworkType)"`
	NetworkProtocol string `orm:"column(NetworkProtocol)"`
	StatusCode      string `orm:"column(StatusCode)"`
	FlagBit         string `orm:"column(FlagBit)"`
	SubTime         string `orm:"column(SubTime)"`
	TimeDelay       string `orm:"column(TimeDelay)"`

	//	ServerInfo      *ServerInfo `orm:"rel(fk)"`
	//	AppInfo         *AppInfo    `orm:"rel(fk)"`
}

var (
	MoniDataList map[string]*MonitorData
)

func (m *MonitorData) TableName() string {
	return TableName("monitorData")
}
func AddOne(monitordata *MonitorData) (int64, error) {
	return orm.NewOrm().Insert(monitordata)
}
func GetOne(dataid int64) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	o.Raw("SELECT D.DataID,D.AppID,A.AppName,D.ServerAddress,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.SubTime,D.TimeDelay From monitorData as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID WHERE DataID = ?", dataid).Values(&maps)
	/**/ for k, v := range maps {
		fmt.Println(k, v)
	}
	fmt.Println("data1", maps[0])
	return maps
}

/***********************************************************/
func Add1(monitordata *MonitorData) (int64, error) {
	number := strconv.Itoa(monitordata.ServerID)
	sqlresult, err1 := orm.NewOrm().Raw("INSERT INTO Data_"+number+" (AppID,ServerID,Port,NetworkType,NetworkProtocol,StatusCode,FlagBit,SubTime,TimeDelay)VALUES(?,?,?,?,?,?,?,?,?);", monitordata.AppID, monitordata.ServerID, monitordata.Port, monitordata.NetworkType, monitordata.NetworkProtocol, monitordata.StatusCode, monitordata.FlagBit, monitordata.SubTime, monitordata.TimeDelay).Exec()
	lastinsertid, err2 := sqlresult.LastInsertId()
	lib.FailOnErr(err1, "insert into error")
	lib.FailOnErr(err2, "get lastinsertID error")
	return lastinsertid, err1
}

/********************************************************/
func GetMDbyID(dataid int64, serverid int) []orm.Params {
	tablenum := strconv.Itoa(serverid)
	o := orm.NewOrm()
	var maps []orm.Params
	affect, err := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID WHERE D.DataID=? ORDER BY D.SubTime desc", dataid).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	lib.FailOnErr(err, "models/MonitorData/GetMDbyPage ")
	fmt.Println("affect", affect)
	return maps
}

//分页获取监控数据2
func GetMDbyPage2(page, pageSize, appID, serverID int, time1, time2 string) ([]orm.Params, interface{}) {
	var maps []orm.Params
	var maps2 []orm.Params
	var total interface{}
	o := orm.NewOrm()
	tablenum := strconv.Itoa(serverID)
	offset := (page - 1) * pageSize
	if appID == 0 {
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" Where SubTime BETWEEN ? AND ?;", time1, time2).Values(&maps2)
		total = maps2[0]["COUNT(*)"]
		_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  Where D.SubTime BETWEEN ? AND ? ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
		lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")

	} else {
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+"   WHERE AppID= ? AND (SubTime BETWEEN ? AND ?);", appID, time1, time2).Values(&maps2)
		_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID WHERE D.AppID= ? AND (D.SubTime BETWEEN ? AND ? ) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", appID, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
		lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
		total = maps2[0]["COUNT(*)"]
	}
	return maps, total
}

//分页获取数据3
func GetMDbyPage3(page, pageSize, serverID int, region, time1, time2 string) ([]orm.Params, interface{}) {
	var maps []orm.Params
	var maps2 []orm.Params
	var total interface{}
	o := orm.NewOrm()
	tablenum := strconv.Itoa(serverID)
	offset := (page - 1) * pageSize
	if region == "" {
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" Where SubTime BETWEEN ? AND ?;", time1, time2).Values(&maps2)
		total = maps2[0]["COUNT(*)"]
		_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  Where D.SubTime BETWEEN ? AND ? ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
		lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
	} else {
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region= ? AND (SubTime BETWEEN ? AND ?);", region, time1, time2).Values(&maps2)
		_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID WHERE A.Region = ? AND (D.SubTime BETWEEN ? AND ? ) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
		lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
		total = maps2[0]["COUNT(*)"]
	}
	return maps, total
}

//分页获取数据筛选(地区)
func GetMDFliter1(page, pageSize, serverID int, time1, time2, region string) ([]orm.Params, interface{}) {
	var maps []orm.Params
	var maps2 []orm.Params
	var total interface{}
	o := orm.NewOrm()
	tablenum := strconv.Itoa(serverID)
	offset := (page - 1) * pageSize
	o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region= ? AND (SubTime BETWEEN ? AND ?);", region, time1, time2).Values(&maps2)
	_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID WHERE A.Region = ? AND (D.SubTime BETWEEN ? AND ? ) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
	total = maps2[0]["COUNT(*)"]
	return maps, total
}

//分页获取数据筛选(运营商)
func GetMDFliter2(page, pageSize, serverID int, time1, time2, operator string) ([]orm.Params, interface{}) {
	var maps []orm.Params
	var maps2 []orm.Params
	var total interface{}
	o := orm.NewOrm()
	tablenum := strconv.Itoa(serverID)
	offset := (page - 1) * pageSize
	o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Operator= ? AND (SubTime BETWEEN ? AND ?);", operator, time1, time2).Values(&maps2)
	_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID WHERE A.Operator= ? AND (D.SubTime BETWEEN ? AND ? ) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
	total = maps2[0]["COUNT(*)"]
	return maps, total
}

//分页获取数据筛选(网络类型)
func GetMDFliter3(page, pageSize, serverID int, time1, time2, networkType string) ([]orm.Params, interface{}) {
	var maps []orm.Params
	var maps2 []orm.Params
	var total interface{}
	o := orm.NewOrm()
	tablenum := strconv.Itoa(serverID)
	offset := (page - 1) * pageSize
	o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE D.NetworkType= ? AND (SubTime BETWEEN ? AND ?);", networkType, time1, time2).Values(&maps2)
	_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID WHERE D.NetworkType = ? AND (D.SubTime BETWEEN ? AND ? ) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkType, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
	total = maps2[0]["COUNT(*)"]
	return maps, total
}

//分页获取数据筛选(网络协议)
func GetMDFliter4(page, pageSize, serverID int, time1, time2, networkProtocol string) ([]orm.Params, interface{}) {
	var maps []orm.Params
	var maps2 []orm.Params
	var total interface{}
	o := orm.NewOrm()
	tablenum := strconv.Itoa(serverID)
	offset := (page - 1) * pageSize
	o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE D.NetworkProtocol= ? AND (SubTime BETWEEN ? AND ?);", networkProtocol, time1, time2).Values(&maps2)
	_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID WHERE D.NetworkProtocol = ? AND (D.SubTime BETWEEN ? AND ? ) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkProtocol, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
	total = maps2[0]["COUNT(*)"]

	return maps, total
}

//分页获取数据筛选(延时)
func GetMDFliter5(page, pageSize, serverID int, time1, time2, timeDelay string) ([]orm.Params, interface{}) {
	var maps []orm.Params
	var maps2 []orm.Params
	var total interface{}
	timeDelayInt, _ := strconv.Atoi(timeDelay)
	o := orm.NewOrm()
	tablenum := strconv.Itoa(serverID)
	offset := (page - 1) * pageSize
	o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", timeDelayInt, time1, time2).Values(&maps2)
	_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID WHERE D.TimeDelay > ? AND (D.SubTime BETWEEN ? AND ? ) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
	total = maps2[0]["COUNT(*)"]
	return maps, total
}

//分页获取数据筛选(全条件1)
func GetMDFliter6(page, pageSize, serverID int, time1, time2, region, networkType, networkProtocol, operator, timeDelay string) ([]orm.Params, interface{}) {
	var maps []orm.Params
	var maps2 []orm.Params
	var total interface{}
	o := orm.NewOrm()
	timeDelayInt, _ := strconv.Atoi(timeDelay)
	tablenum := strconv.Itoa(serverID)
	offset := (page - 1) * pageSize
	o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND D.NetworkType=? AND D.NetworkProtocol=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", region, networkType, networkProtocol, operator, timeDelayInt, time1, time2).Values(&maps2)
	_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkType=? AND D.NetworkProtocol=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkType, networkProtocol, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
	total = maps2[0]["COUNT(*)"]
	return maps, total
}

//分页获取数据筛选(全条件2)~~好恶心的筛选,笨办法
func GetMDFliter7(page, pageSize, serverID int, time1, time2, region, networkType, networkProtocol, operator, timeDelay string) ([]orm.Params, interface{}) {
	var maps []orm.Params
	var maps2 []orm.Params
	var total interface{}
	var err2 error
	o := orm.NewOrm()
	timeDelayInt, _ := strconv.Atoi(timeDelay)
	tablenum := strconv.Itoa(serverID)
	offset := (page - 1) * pageSize
	if region == "" && networkType == "" && networkProtocol == "" && operator == "" && timeDelayInt == 0 { //00000
		//fmt.Println("00000")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE SubTime BETWEEN ? AND ?;", time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE  SubTime BETWEEN ? AND ? ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol == "" && operator == "" && timeDelayInt != 0 { //00001
		//fmt.Println("00001")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE   D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol == "" && operator != "" && timeDelayInt == 0 { //00010
		//	fmt.Println("00010")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE  A.Operator=? AND (SubTime BETWEEN ? AND ?);", operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Operator=?  AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol == "" && operator != "" && timeDelayInt != 0 { //00011
		//	fmt.Println("00011")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol != "" && operator == "" && timeDelayInt == 0 { //00100
		//	fmt.Println("00100")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE  D.NetworkProtocol=?  AND (SubTime BETWEEN ? AND ?);", networkProtocol, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE  D.NetworkProtocol=?  AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkProtocol, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol != "" && operator == "" && timeDelayInt != 0 { //00101
		//	fmt.Println("00101")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE D.NetworkProtocol=?  AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", networkProtocol, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE  D.NetworkProtocol=?  AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkProtocol, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol != "" && operator != "" && timeDelayInt == 0 { //00110
		//	fmt.Println("00110")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE D.NetworkProtocol=? AND A.Operator=? AND (SubTime BETWEEN ? AND ?);", networkProtocol, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE  D.NetworkProtocol=? AND A.Operator=?  AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkProtocol, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol != "" && operator != "" && timeDelayInt != 0 { //00111
		//	fmt.Println("00111")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE  D.NetworkProtocol=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", networkProtocol, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE D.NetworkProtocol=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkProtocol, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol == "" && operator == "" && timeDelayInt == 0 { //01000
		//	fmt.Println("01000")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE D.NetworkType=?  AND (SubTime BETWEEN ? AND ?);", networkType, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE  D.NetworkType=?  AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkType, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol == "" && operator == "" && timeDelayInt != 0 { //01001
		//	fmt.Println("01001")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE  D.NetworkType=? AND D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", networkType, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE  D.NetworkType=?  AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkType, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol == "" && operator != "" && timeDelayInt == 0 { //01010
		//	fmt.Println("01010")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE  D.NetworkType=?  AND A.Operator=? AND (SubTime BETWEEN ? AND ?);", networkType, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE  D.NetworkType=?  AND A.Operator=?  AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkType, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol == "" && operator != "" && timeDelayInt != 0 { //01011
		//	fmt.Println("01011")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE D.NetworkType=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", networkType, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE  D.NetworkType=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkType, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol != "" && operator == "" && timeDelayInt == 0 { //01100
		//	fmt.Println("01100")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE D.NetworkType=? AND D.NetworkProtocol=?  AND (SubTime BETWEEN ? AND ?);", networkType, networkProtocol, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE D.NetworkType=? AND D.NetworkProtocol=? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkType, networkProtocol, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol != "" && operator == "" && timeDelayInt != 0 { //01101
		//	fmt.Println("01101")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE D.NetworkType=? AND D.NetworkProtocol=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", networkType, networkProtocol, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE D.NetworkType=? AND D.NetworkProtocol=?  AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkType, networkProtocol, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol != "" && operator != "" && timeDelayInt == 0 { //01110
		//	fmt.Println("01110")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE  D.NetworkType=? AND D.NetworkProtocol=? AND A.Operator=?  AND (SubTime BETWEEN ? AND ?);", networkType, networkProtocol, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE  D.NetworkType=? AND D.NetworkProtocol=? AND A.Operator=? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkType, networkProtocol, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol != "" && operator != "" && timeDelayInt != 0 { //01111
		//	fmt.Println("01111")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE D.NetworkType=? AND D.NetworkProtocol=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", networkType, networkProtocol, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE  D.NetworkType=? AND D.NetworkProtocol=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", networkType, networkProtocol, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol == "" && operator == "" && timeDelayInt == 0 { //10000
		//	fmt.Println("10000")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND (SubTime BETWEEN ? AND ?);", region, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=?  AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol == "" && operator == "" && timeDelayInt != 0 { //10001
		//	fmt.Println("10001")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", region, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=?  AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol == "" && operator != "" && timeDelayInt == 0 { //10010
		//	fmt.Println("10010")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=?  AND A.Operator=?  AND (SubTime BETWEEN ? AND ?);", region, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=?  AND A.Operator=?  AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol == "" && operator != "" && timeDelayInt != 0 { //10011
		//	fmt.Println("10011")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=?  AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", region, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=?  AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol != "" && operator == "" && timeDelayInt == 0 { //10100
		//	fmt.Println("10100")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=?  AND D.NetworkProtocol=?  AND (SubTime BETWEEN ? AND ?);", region, networkProtocol, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkProtocol=? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkProtocol, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol != "" && operator == "" && timeDelayInt != 0 { //10101
		//	fmt.Println("10101")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=?  AND D.NetworkProtocol=?  AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", region, networkProtocol, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=?  AND D.NetworkProtocol=?  AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkProtocol, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol != "" && operator != "" && timeDelayInt == 0 { //10110
		//	fmt.Println("10110")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=?  AND D.NetworkProtocol=? AND A.Operator=?  AND (SubTime BETWEEN ? AND ?);", region, networkProtocol, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkProtocol=? AND A.Operator=? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkProtocol, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol != "" && operator != "" && timeDelayInt != 0 { //10111
		//	fmt.Println("10111")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=?  AND D.NetworkProtocol=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", region, networkProtocol, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=?  AND D.NetworkProtocol=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkProtocol, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol == "" && operator == "" && timeDelayInt == 0 { //11000
		//	fmt.Println("11000")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND D.NetworkType=? AND (SubTime BETWEEN ? AND ?);", region, networkType, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkType=?  AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkType, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol == "" && operator == "" && timeDelayInt != 0 { //11001
		//	fmt.Println("11001")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND D.NetworkType=?  AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", region, networkType, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkType=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkType, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol == "" && operator != "" && timeDelayInt == 0 { //11010
		//	fmt.Println("11010")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND D.NetworkType=? AND A.Operator=? AND (SubTime BETWEEN ? AND ?);", region, networkType, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkType=? AND A.Operator=? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkType, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol == "" && operator != "" && timeDelayInt != 0 { //11011
		//	fmt.Println("11011")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND D.NetworkType=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", region, networkType, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkType=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkType, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol != "" && operator == "" && timeDelayInt == 0 { //11100
		//	fmt.Println("11100")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND D.NetworkType=? AND D.NetworkProtocol=?  AND (SubTime BETWEEN ? AND ?);", region, networkType, networkProtocol, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkType=? AND D.NetworkProtocol=?  AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkType, networkProtocol, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol != "" && operator == "" && timeDelayInt != 0 { //11101
		//	fmt.Println("11101")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND D.NetworkType=? AND D.NetworkProtocol=?  AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", region, networkType, networkProtocol, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkType=? AND D.NetworkProtocol=?  AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkType, networkProtocol, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol != "" && operator != "" && timeDelayInt == 0 { //11110
		//	fmt.Println("11110")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND D.NetworkType=? AND D.NetworkProtocol=? AND A.Operator=?  AND (SubTime BETWEEN ? AND ?);", region, networkType, networkProtocol, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkType=? AND D.NetworkProtocol=? AND A.Operator=?  AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkType, networkProtocol, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol != "" && operator != "" && timeDelayInt != 0 { //11111
		//	fmt.Println("11111")
		o.Raw("SELECT COUNT(*) FROM Data_"+tablenum+" AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE A.Region=? AND D.NetworkType=? AND D.NetworkProtocol=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?);", region, networkType, networkProtocol, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw("SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_"+tablenum+" as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE A.Region=? AND D.NetworkType=? AND D.NetworkProtocol=? AND A.Operator=? AND  D.TimeDelay > ? AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", region, networkType, networkProtocol, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	}
	if err2 != nil {
		beego.Error(err2, "models/MonitorData/GetMDbyPage ")
	}
	total = maps2[0]["COUNT(*)"]
	return maps, total
}

//分页获取数据筛选(全条件2) sql语句z自动拼接
func GetMDFliter8(page, pageSize, serverID int, time1, time2, region, networkType, networkProtocol, operator, timeDelay string) ([]orm.Params, interface{}) {
	var maps []orm.Params
	var maps2 []orm.Params
	var total interface{}
	var sqlstr string
	var sqlstr1 string
	var sqlstr2 string
	var err2 error
	o := orm.NewOrm()
	timeDelayInt, _ := strconv.Atoi(timeDelay)
	tablenum := strconv.Itoa(serverID)
	offset := (page - 1) * pageSize
	sqlstr1 = "SELECT COUNT(*) FROM Data_" + tablenum + " AS D INNER JOIN appInfo AS A ON D.AppID=A.AppID  WHERE"
	sqlstr2 = "SELECT D.DataID,D.AppID,A.AppName,S.ServerAddress,S.ServerID,S.ServerName,A.Region,D.Port,A.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime,D.TimeDelay  From Data_" + tablenum + " as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerID=S.ServerID  WHERE"
	if region != "" {
		if sqlstr != "" {
			sqlstr = sqlstr + " AND A.Region=? "
		} else {
			sqlstr = " A.Region=? "
		}
	}
	if networkType != "" {
		if sqlstr != "" {
			sqlstr = sqlstr + " AND D.NetworkType=? "
		} else {
			sqlstr = " D.NetworkType=? "
		}
	}
	if networkProtocol != "" {
		if sqlstr != "" {
			sqlstr = sqlstr + " AND D.NetworkProtocol=? "
		} else {
			sqlstr = " D.NetworkProtocol=? "
		}
	}
	if operator != "" {
		if sqlstr != "" {
			sqlstr = sqlstr + " AND A.Operator=? "
		} else {
			sqlstr = " A.Operator=? "
		}
	}
	if timeDelayInt != 0 {
		if sqlstr != "" {
			sqlstr = sqlstr + " AND D.TimeDelay > ? "
		} else {
			sqlstr = " D.TimeDelay > ? "
		}
	}
	if region == "" && networkType == "" && networkProtocol == "" && operator == "" && timeDelayInt == 0 {
		sqlstr1 = sqlstr1 + " SubTime BETWEEN ? AND ?;"
		sqlstr2 = sqlstr2 + " SubTime BETWEEN ? AND ? ORDER BY D.SubTime desc LIMIT ? OFFSET ?"
	} else {
		sqlstr1 = sqlstr1 + sqlstr + " AND (SubTime BETWEEN ? AND ?);"
		sqlstr2 = sqlstr2 + sqlstr + " AND (SubTime BETWEEN ? AND ?) ORDER BY D.SubTime desc LIMIT ? OFFSET ?"
	}

	//fmt.Println("sql1:", sqlstr1)
	//fmt.Println("sql2:", sqlstr2)
	/****************************************************/
	if region == "" && networkType == "" && networkProtocol == "" && operator == "" && timeDelayInt == 0 { //00000
		//fmt.Println("00000")
		o.Raw(sqlstr1, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol == "" && operator == "" && timeDelayInt != 0 { //00001
		//fmt.Println("00001")
		o.Raw(sqlstr1, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol == "" && operator != "" && timeDelayInt == 0 { //00010
		//	fmt.Println("00010")
		o.Raw(sqlstr1, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol == "" && operator != "" && timeDelayInt != 0 { //00011
		//	fmt.Println("00011")
		o.Raw(sqlstr1, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol != "" && operator == "" && timeDelayInt == 0 { //00100
		//	fmt.Println("00100")
		o.Raw(sqlstr1, networkProtocol, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkProtocol, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol != "" && operator == "" && timeDelayInt != 0 { //00101
		//	fmt.Println("00101")
		o.Raw(sqlstr1, networkProtocol, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkProtocol, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol != "" && operator != "" && timeDelayInt == 0 { //00110
		//	fmt.Println("00110")
		o.Raw(sqlstr1, networkProtocol, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkProtocol, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType == "" && networkProtocol != "" && operator != "" && timeDelayInt != 0 { //00111
		//	fmt.Println("00111")
		o.Raw(sqlstr1, networkProtocol, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkProtocol, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol == "" && operator == "" && timeDelayInt == 0 { //01000
		//	fmt.Println("01000")
		o.Raw(sqlstr1, networkType, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkType, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol == "" && operator == "" && timeDelayInt != 0 { //01001
		//	fmt.Println("01001")
		o.Raw(sqlstr1, networkType, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkType, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol == "" && operator != "" && timeDelayInt == 0 { //01010
		//	fmt.Println("01010")
		o.Raw(sqlstr1, networkType, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkType, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol == "" && operator != "" && timeDelayInt != 0 { //01011
		//	fmt.Println("01011")
		o.Raw(sqlstr1, networkType, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkType, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol != "" && operator == "" && timeDelayInt == 0 { //01100
		//	fmt.Println("01100")
		o.Raw(sqlstr1, networkType, networkProtocol, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkType, networkProtocol, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol != "" && operator == "" && timeDelayInt != 0 { //01101
		//	fmt.Println("01101")
		o.Raw(sqlstr1, networkType, networkProtocol, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkType, networkProtocol, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol != "" && operator != "" && timeDelayInt == 0 { //01110
		//	fmt.Println("01110")
		o.Raw(sqlstr1, networkType, networkProtocol, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkType, networkProtocol, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region == "" && networkType != "" && networkProtocol != "" && operator != "" && timeDelayInt != 0 { //01111
		//	fmt.Println("01111")
		o.Raw(sqlstr1, networkType, networkProtocol, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, networkType, networkProtocol, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol == "" && operator == "" && timeDelayInt == 0 { //10000
		//	fmt.Println("10000")
		o.Raw(sqlstr1, region, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol == "" && operator == "" && timeDelayInt != 0 { //10001
		//	fmt.Println("10001")
		o.Raw(sqlstr1, region, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol == "" && operator != "" && timeDelayInt == 0 { //10010
		//	fmt.Println("10010")
		o.Raw(sqlstr1, region, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol == "" && operator != "" && timeDelayInt != 0 { //10011
		//	fmt.Println("10011")
		o.Raw(sqlstr1, region, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol != "" && operator == "" && timeDelayInt == 0 { //10100
		//	fmt.Println("10100")
		o.Raw(sqlstr1, region, networkProtocol, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkProtocol, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol != "" && operator == "" && timeDelayInt != 0 { //10101
		//	fmt.Println("10101")
		o.Raw(sqlstr1, region, networkProtocol, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkProtocol, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol != "" && operator != "" && timeDelayInt == 0 { //10110
		//	fmt.Println("10110")
		o.Raw(sqlstr1, region, networkProtocol, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkProtocol, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType == "" && networkProtocol != "" && operator != "" && timeDelayInt != 0 { //10111
		//	fmt.Println("10111")
		o.Raw(sqlstr1, region, networkProtocol, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkProtocol, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol == "" && operator == "" && timeDelayInt == 0 { //11000
		//	fmt.Println("11000")
		o.Raw(sqlstr1, region, networkType, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkType, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol == "" && operator == "" && timeDelayInt != 0 { //11001
		//	fmt.Println("11001")
		o.Raw(sqlstr1, region, networkType, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkType, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol == "" && operator != "" && timeDelayInt == 0 { //11010
		//	fmt.Println("11010")
		o.Raw(sqlstr1, region, networkType, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkType, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol == "" && operator != "" && timeDelayInt != 0 { //11011
		//	fmt.Println("11011")
		o.Raw(sqlstr1, region, networkType, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkType, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol != "" && operator == "" && timeDelayInt == 0 { //11100
		//	fmt.Println("11100")
		o.Raw(sqlstr1, region, networkType, networkProtocol, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkType, networkProtocol, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol != "" && operator == "" && timeDelayInt != 0 { //11101
		//	fmt.Println("11101")
		o.Raw(sqlstr1, region, networkType, networkProtocol, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkType, networkProtocol, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol != "" && operator != "" && timeDelayInt == 0 { //11110
		//	fmt.Println("11110")
		o.Raw(sqlstr1, region, networkType, networkProtocol, operator, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkType, networkProtocol, operator, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	} else if region != "" && networkType != "" && networkProtocol != "" && operator != "" && timeDelayInt != 0 { //11111
		//	fmt.Println("11111")
		o.Raw(sqlstr1, region, networkType, networkProtocol, operator, timeDelayInt, time1, time2).Values(&maps2)
		_, err2 = o.Raw(sqlstr2, region, networkType, networkProtocol, operator, timeDelayInt, time1, time2, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerID", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime", "TimeDelay")
	}
	/******************************************************/
	if err2 != nil {
		beego.Error(err2, "models/MonitorData/GetMDbyPage ")
	}
	total = maps2[0]["COUNT(*)"]
	return maps, total
}

//删除过期数据
func DeleteData(serverid int, time string) error {
	o := orm.NewOrm()
	tablenum := strconv.Itoa(serverid)
	_, err := o.Raw("DELETE FROM Data_"+tablenum+" WHERE SubTime <= ?", time).Exec()
	//	afecrows, _ := result.RowsAffected()
	return err
}
func CreateTable(serverid int64) {
	servid := strconv.FormatInt(serverid, 10)
	sqlstr := `CREATE TABLE IF NOT EXISTS Data_` + servid + ` (
	DataID  INT NOT NULL AUTO_INCREMENT,
	AppID   INT  NULL,
	ServerID  INT  NULL,
	Port  VARCHAR(10)  NULL,
	NetworkType VARCHAR(10)  NULL,
	NetworkProtocol VARCHAR(10)  NULL,
	StatusCode   VARCHAR(10)  NULL,
	FlagBit  VARCHAR(10)   NULL,
	SubTime   datetime    NULL,
	TimeDelay  VARCHAR(20) NULL,
	PRIMARY KEY(DataID)
	)`
	fmt.Println("sqlstr:", sqlstr)
	_, err := orm.NewOrm().Raw(sqlstr).Exec()
	log.Println("create table error:", err)
}
