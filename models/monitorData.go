package models

import (
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
	NetworkType     string `orm:"column(NetworkType)"`
	NetworkProtocol string `orm:"column(NetworkProtocol)"`
	StatusCode      string `orm:"column(StatusCode)"`
	FlagBit         string `orm:"column(FlagBit)"`
	SubTime         string `orm:"column(SubTime)"`
	//	ServerInfo      *ServerInfo `orm:"rel(one)"`
	//	AppInfo         *AppInfo    `orm:"rel(one)"`
}
type ResultData struct {
	DataID          int64  `orm:"column(DataID)"`
	AppID           int64  `orm:"column(AppID)"`
	ServerName      string `orm:"column(ServerName)"`
	AppName         string `orm:"column(AppName)"`
	Region          string `orm:"column(Region)"`
	ServerAddress   string `orm:"column(ServerAddress)"`
	Port            string `orm:"column(Port)"`
	NetworkType     string `orm:"column(NetworkType)"`
	NetworkProtocol string `orm:"column(NetworkProtocol)"`
	StatusCode      string `orm:"column(StatusCode)"`
	FlagBit         string `orm:"column(FlagBit)"`
	SubTime         string `orm:"column(SubTime)"`
}
type BigMonidata struct {
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
func GetMDbyPage(page, pageSize int, filters ...interface{}) ([]orm.Params, int64) {
	o := orm.NewOrm()
	var maps []orm.Params //
	rawseter := o.Raw("SELECT D.DataID,D.AppID,A.AppName,D.ServerAddress,S.ServerName,A.Region,D.Port,D.NetworkType,D.NetworkProtocol,D.StatusCode,D.FlagBit,D.SubTime From monitorData as D INNER JOIN appInfo as A on D.AppID = A.AppID INNER JOIN serverInfo as S ON D.ServerAddress=S.ServerAddress ORDER BY D.SubTime desc;")
	total, _ := rawseter.Values(&maps, "DataID", "AppID", "AppName", "ServerAddress", "ServerName", "Region", "Port", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime")
	for k, v := range maps {
		fmt.Println(k, v)
	}
	//fmt.Println("data1", maps[0])
	return maps, total /**/
	/******************************************/
	/*	var maps []orm.Params
		monitor := new(MonitorData)
		offset := (page - 1) * pageSize
		//list := make([]*ResultData, 0)
		query := orm.NewOrm().QueryTable(monitor)
		query = query.RelatedSel()

		total, _ := query.Count()
		query.OrderBy("SubTime").Limit(pageSize, offset).Values(&maps, "DataID", "AppID", "ServerName", "AppName", "Region", "ServerAddress", "Port", "NetworkType", "NetworkProtocol", "StatusCode", "FlagBit", "SubTime")

		return maps, total */
	/********************************************/

}
