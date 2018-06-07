package models

import (
	"github.com/astaxie/beego/orm"
)

type ServerInfo struct {
	ServerID      int64  `orm:"pk;column(ServerID)"`
	ServerName    string `orm:"column(ServerName)"`
	Type          string `orm:"column(Type)"`
	ServerAddress string `orm:"column(ServerAddress)"`
	Port          string `orm:"column(Port)"`
	ReqInterval   string `orm:"column(ReqInterval)"`
	Sstatus       string `orm:"column(Sstatus)"`
	Info          string `orm:"column(Info)"`
	//MonitorData   []*MonitorData `orm:"reverse(many)"`
}

var (
	ServerList map[string]*ServerInfo
)

func (s *ServerInfo) TableName() string {
	return TableName("serverInfo")
}

func AddSver(serverinfo *ServerInfo) (int64, error) {

	return orm.NewOrm().Insert(serverinfo)
}
func GetSverByAddress(saddress string) *ServerInfo {
	s := new(ServerInfo)
	orm.NewOrm().Raw("SELECT * FROM serverInfo where Address = ?", saddress).QueryRow(&s)
	return s
}

//条件获取
func GetSverFilter(status string) []orm.Params {
	var op []orm.Params
	o := orm.NewOrm()
	//	orm.NewOrm().QueryTable(serverInfo).All(&serarry, "ServerAddress", "Delay")
	o.Raw("SELECT *  FROM serverInfo WHERE Sstatus=?", status).Values(&op, "ServerAddress", "ServerName", "ServerID", "ReqInterval", "Port", "Type")
	return op
}

//获取全部
func GetSverAll() ([]orm.Params, int) {
	var op []orm.Params
	orm.NewOrm().Raw("SELECT *  FROM serverInfo").Values(&op)
	count := len(op)
	return op, count

}
func UpdateSver(serverinfo *ServerInfo) (int64, error) {
	serverInfo := new(ServerInfo)
	num, err := orm.NewOrm().QueryTable(serverInfo).Filter("ServerID", serverinfo.ServerID).Update(orm.Params{
		"ServerName": serverinfo.ServerName, "ServerAddress": serverinfo.ServerAddress, "Port": serverinfo.Port, "Type": serverinfo.Type, "ReqInterval": serverinfo.ReqInterval, "Sstatus": serverinfo.Sstatus, "Info": serverinfo.Info,
	})
	return num, err
}
