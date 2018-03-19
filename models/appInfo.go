package models

import (
	"gree/GrihCommon/security"

	"github.com/astaxie/beego/orm"
)

type AppInfo struct {
	AppID    int64  `orm:"pk;column(AppID)"`
	AppName  string `orm:"column(AppName)"`
	Password string `orm:"column(Password)"`
	AppToken string `orm:"column(AppToken)"`
	Region   string `orm:"column(Region)"`
	RegTime  string `orm:"column(RegTime)"`
	//MonitorData []*MonitorData `orm:"reverse(many)"`
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
func AppLogin(appname, password string) int64 {
	var appid int64
	a := new(AppInfo)
	o := orm.NewOrm()
	password = security.Md5(password + security.Md5(password))
	o.Raw("SELECT * from appInfo where AppName = ?", appname).QueryRow(&a)
	if a.Password == password {
		appid = a.AppID
	} else {
		appid = 0
	}
	return appid
}
func AppExist(AppName string) int64 {
	var app []AppInfo
	var affectrows int64
	num, err := orm.NewOrm().Raw("SELECT * FROM appInfo Where AppName = ?", AppName).QueryRows(&app)
	if err == nil {
		affectrows = num
	}
	return affectrows
}
func Appbyid(appid int64) *AppInfo {
	a := new(AppInfo)
	orm.NewOrm().Raw("SELECT * FROM appInfo where AppID = ?", appid).QueryRow(&a)
	return a
}
