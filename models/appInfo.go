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
	Operator string `orm:"column(Operator)"`
	Astatus  string `orm:"column(Astatus)"`
	Online   string `orm:"column(Online)"`
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
func AppLogin(appname, password, apptoken, online string) int64 {
	var appid int64
	a := new(AppInfo)
	o := orm.NewOrm()
	password = security.Md5(password + security.Md5(password))
	o.Raw("SELECT * from appInfo where AppName = ?", appname).QueryRow(&a)
	if a.Password == password {
		appid = a.AppID
		o.Raw("Update appInfo Set AppToken=?,Online=? where AppID=?", apptoken, online, a.AppID).Exec()
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

//获取app信息不包含token
func GetAppAll() ([]orm.Params, int) {
	var app []orm.Params
	appInfo := new(AppInfo)
	orm.NewOrm().QueryTable(appInfo).Values(&app, "AppID", "AppName", "Region", "Operator", "RegTime", "Astatus", "Online")
	len := len(app)

	return app, len
}

func UpdateApp(app *AppInfo) (int64, error) {
	var num int64
	var err error
	appInfo := new(AppInfo)
	if app.Password != "" {
		num, err = orm.NewOrm().QueryTable(appInfo).Filter("AppID", app.AppID).Update(orm.Params{
			"AppName": app.AppName, "Password": app.Password, "Operator": app.Operator, "Region": app.Region, "RegTime": app.RegTime, "Astatus": app.Astatus,
		})
	} else {
		num, err = orm.NewOrm().QueryTable(appInfo).Filter("AppID", app.AppID).Update(orm.Params{
			"AppName": app.AppName, "Operator": app.Operator, "Region": app.Region, "RegTime": app.RegTime, "Astatus": app.Astatus,
		})
	}
	return num, err
}
func GetAppAllinside() ([]orm.Params, int) {
	var app []orm.Params
	appInfo := new(AppInfo)
	orm.NewOrm().QueryTable(appInfo).Values(&app, "AppID", "AppName", "Password", "AppToken", "Region", "Operator", "RegTime", "Astatus", "Online")
	len := len(app)
	return app, len
}
func VerrifyAppFromMysql(appid int64, accesstoken string) bool {
	appinfo := Appbyid(appid)
	if appinfo.AppToken == accesstoken {
		return true
	}
	return false
}
func AppOnlineHeart(appId int64, online string) (int64, error) {
	var num int64
	var err error
	appInfo := new(AppInfo)
	num, err = orm.NewOrm().QueryTable(appInfo).Filter("AppID", appId).Update(orm.Params{
		"Online": online,
	})
	return num, err
}
