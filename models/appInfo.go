package models

import "github.com/astaxie/beego/orm"

type AppInfo struct {
	AppID    int64  `orm:"pk;column(AppID)"`
	AppName  string `orm:"column(AppName)"`
	Password string `orm:"column(Password)"`
	AppToken string `orm:"column(AppToken)"`
	Region   string `orm:"column(Region)"`
	RegTime  string `orm:"column(RegTime)"`
}

var (
	AppList map[string]*AppInfo
)

func (a *AppInfo) TableName() string {
	return TableName("appInfo")
}
func AppInfoAdd(appinfo *AppInfo) (int64, error) {
	return orm.NewOrm().Insert(appinfo)
}
