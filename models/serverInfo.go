package models

import "github.com/astaxie/beego/orm"

type ServerInfo struct {
	ServerID   int64  `orm:"pk;column(ServerID)"`
	ServerName string `orm:"column(ServerName)"`
	Address    string `orm:"column(Address)"`
	Port       string `orm:"column(Port)"`
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
