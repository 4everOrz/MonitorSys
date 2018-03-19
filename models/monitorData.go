package models

import (
	"MonitoringSystemAPI/lib"
	"fmt"

	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type MonitorData struct {
	DataID          int64  `orm:"pk;column(DataID)"`
	AppID           int64  `orm:"column(AppID)"`
	ServerAddress   string `orm:"column(ServerAddress)"`
	Port            string `orm:"column(Port)"`
	Operator        string `orm:"column(Operator)"`
	NetworkType     string `orm:"column(NetworkType)"`
	NetworkProtocol string `orm:"column(NetworkProtocol)"`
	StatusCode      string `orm:"column(StatusCode)"`
	FlagBit         string `orm:"column(FlagBit)"`
	SubTime         string `orm:"column(SubTime)"`

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
	o.Raw("SELECT D.DataID,D.AppID,A.AppName,D.ServerAddress,S.ServerName,A.Region,D.Port,D.NetworkType,D.NetworkProtocol,D.SubTime From monitorData as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerAddress=S.ServerAddress WHERE DataID = ?", dataid).Values(&maps)
	/**/ for k, v := range maps {
		fmt.Println(k, v)
	}
	fmt.Println("data1", maps[0])
	return maps
}

//分页获取监控数据，方式有待优化
func GetMDbyPage1(page, pageSize int, filters ...interface{}) ([]orm.Params, int64) {
	/* 	o := orm.NewOrm()
	   	var maps []orm.Params //
	   	rawseter := o.Raw("SELECT D.DataID,D.AppID,A.AppName,D.ServerAddress,S.ServerName,A.Region,D.Port,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime From monitorData as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerAddress=S.ServerAddress ORDER BY D.SubTime desc;")
	   	total, _ := rawseter.Values(&maps, "DataID", "AppID", "AppName", "ServerAddress", "ServerName", "Region", "Port", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime")
	   	for k, v := range maps {
	   		fmt.Println(k, v)
	   	}
	   	//fmt.Println("data1", maps[0])
	   	return maps, total */
	/******************************************/
	/*var maps []orm.Params
	monitor := new(MonitorData)
	offset := (page - 1) * pageSize
	//list := make([]*ResultData, 0)
	query := orm.NewOrm().QueryTable(monitor)
	query = query.RelatedSel()

	total, _ := query.Count()
	query.OrderBy("SubTime").Limit(pageSize, offset).Values(&maps)

	return maps, total  */
	/********************************************/
	o := orm.NewOrm()
	var maps []orm.Params
	monitor := new(MonitorData)
	offset := (page - 1) * pageSize
	total, err1 := o.QueryTable(monitor).Count()
	if err1 == nil {
		_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,D.ServerAddress,S.ServerName,A.Region,D.Port,D.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime From monitorData as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerAddress=S.ServerAddress ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime")
		lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
	}
	return maps, total
}

//分页获取监控数据
func GetMDbyPage2(page, pageSize int, appname, servername string) ([]orm.Params, int) {
	o := orm.NewOrm()
	var maps []orm.Params
	var total int
	monitor := new(MonitorData)
	offset := (page - 1) * pageSize
	if appname == "" && servername == "" {
		_, err1 := o.QueryTable(monitor).Count()
		if err1 == nil {
			_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,D.ServerAddress,S.ServerName,A.Region,D.Port,D.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime From monitorData as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerAddress=S.ServerAddress ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime")
			lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
		}
		total = len(maps)

	} else if appname != "" && servername == "" {
		_, err1 := o.QueryTable(monitor).Count()
		if err1 == nil {
			_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,D.ServerAddress,S.ServerName,A.Region,D.Port,D.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime From monitorData as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerAddress=S.ServerAddress  WHERE AppName = ? ORDER BY D.SubTime desc LIMIT ? OFFSET ? ", appname, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime")
			lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
		}
		total = len(maps)
	} else if appname == "" && servername != "" {
		_, err1 := o.QueryTable(monitor).Count()
		if err1 == nil {
			_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,D.ServerAddress,S.ServerName,A.Region,D.Port,D.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime From monitorData as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerAddress=S.ServerAddress WHERE ServerName = ? ORDER BY D.SubTime  desc LIMIT ? OFFSET ? ", servername, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime")
			lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
		}

		total = len(maps)
	} else if appname != "" && servername != "" {
		_, err1 := o.QueryTable(monitor).Count()
		if err1 == nil {
			_, err2 := o.Raw("SELECT D.DataID,D.AppID,A.AppName,D.ServerAddress,S.ServerName,A.Region,D.Port,D.Operator,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime From monitorData as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerAddress=S.ServerAddress WHERE AppName = ? AND ServerName = ? ORDER BY D.SubTime desc LIMIT ? OFFSET ?  ", appname, servername, pageSize, offset).Values(&maps, "DataID", "AppID", "AppName", "ServerAddress", "ServerName", "Region", "Port", "Operator", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime")
			lib.FailOnErr(err2, "models/MonitorData/GetMDbyPage ")
		}
		total = len(maps)
	}
	return maps, total
}
