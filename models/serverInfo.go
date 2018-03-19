package models

import "github.com/astaxie/beego/orm"

type ServerInfo struct {
	ServerID      int64  `orm:"pk;column(ServerID)"`
	ServerName    string `orm:"column(ServerName)"`
	ServerAddress string `orm:"column(ServerAddress)"`
	Port          string `orm:"column(Port)"`
	Delay         string `orm:"column(Delay)"`
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
func GetSverFilter() []orm.Params {
	var op []orm.Params
	//	orm.NewOrm().QueryTable(serverInfo).All(&serarry, "ServerAddress", "Delay")
	orm.NewOrm().Raw("SELECT *  FROM serverInfo").Values(&op, "ServerAddress", "Delay")
	return op
}

//获取全部
func GetSverAll() []orm.Params {
	var op []orm.Params
	orm.NewOrm().Raw("SELECT *  FROM serverInfo").Values(&op)
	return op

}
