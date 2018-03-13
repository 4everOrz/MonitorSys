package models

import (
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
