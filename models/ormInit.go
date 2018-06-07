package models

import (
	"log"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func OrmInit() {
	/*****************Mysql初始化*********************/
	db1host := beego.AppConfig.String("db1.host")
	db1port := beego.AppConfig.String("db1.port")
	db1user := beego.AppConfig.String("db1.user")
	db1password := beego.AppConfig.String("db1.password")
	db1name := beego.AppConfig.String("db1.name")
	db1timezone := beego.AppConfig.String("db1.timezone") //区域

	if db1port == "" {
		db1port = "3306"
	}
	dsn1 := db1user + ":" + db1password + "@tcp(" + db1host + ":" + db1port + ")/" + db1name + "?charset=utf8"

	if db1timezone != "" {
		dsn1 = dsn1 //+ "&loc=" + url.QueryEscape(db1timezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn1)
	orm.RegisterModel(new(UserInfo))
	orm.RegisterModel(new(AppInfo))
	orm.RegisterModel(new(ServerInfo))
	orm.RegisterModel(new(MonitorData))
	orm.RegisterModel(new(RoleInfo))

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = false
	}
	//设置连接池
	mdb1, _ := orm.GetDB("default")
	mdb1.SetConnMaxLifetime(time.Second * 20)
	mdb1.SetMaxIdleConns(10)
	mdb1.SetMaxOpenConns(30)

	/*******************数据库切换**************************
	db2host := beego.AppConfig.String("db2.host")
	db2port := beego.AppConfig.String("db2.port")
	db2user := beego.AppConfig.String("db2.user")
	db2password := beego.AppConfig.String("db2.password")
	db2name := beego.AppConfig.String("db2.name")
	db2timezone := beego.AppConfig.String("db2.timezone")

	if db2port == "" {
		db2port = "3306"
	}
	dsn2 := db2user + ":" + db2password + "@tcp(" + db2host + ":" + db2port + ")/" + db2name + "?charset=utf8"

	if db2timezone != "" {
		dsn2 = dsn2 + "&loc=" + url.QueryEscape(db2timezone)
	}
	orm.RegisterDataBase("db2", "mysql", dsn2)
	orm.RegisterModel(new(Api_mo_db01))
	orm.RegisterModel(new(Api_mt_db01))
	orm.RegisterModel(new(Api_rpt_db01))
	//设置连接池
	mdb2, _ := orm.GetDB("db2")
	mdb2.SetConnMaxLifetime(time.Second * 20)
	mdb2.SetMaxIdleConns(10)
	mdb2.SetMaxOpenConns(30)*/
	/*************************************************/
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover", r)
		}
	}()
}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}
