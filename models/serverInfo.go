package models

import "github.com/astaxie/beego/orm"

type ServerInfo struct {
	ServerID      int64  `orm:"pk;column(ServerID)"`
	ServerName    string `orm:"column(ServerName)"`
	ServerAddress string `orm:"column(ServerAddress)"`
	Port          string `orm:"column(Port)"`
	Delay         string `orm:"column(Delay)"`
}

var (
	ServerList map[string]*ServerInfo
)

func (s *ServerInfo) TableName() string {
	return TableName("serverInfo")
}
func ServerInfoAdd(serverinfo *ServerInfo) (int64, error) {
	return orm.NewOrm().Insert(serverinfo)
}
func Sbyaddress(saddress string) *ServerInfo {
	s := new(ServerInfo)
	orm.NewOrm().Raw("SELECT * FROM serverInfo where Address = ?", saddress).QueryRow(&s)
	return s
}
