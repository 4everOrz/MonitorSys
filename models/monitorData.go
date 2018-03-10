package models

import "github.com/astaxie/beego/orm"

type MonitorData struct {
	DataID          int64  `orm:"pk;column(DataID)"`
	AppID           string `orm:"column(AppID)"`
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
func MonitorDataAdd(monitordata *MonitorData) (int64, error) {
	return orm.NewOrm().Insert(monitordata)
}
