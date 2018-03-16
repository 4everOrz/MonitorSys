package models

import (
	"net/url"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func Init() {
	/*****************Mysql初始化*********************/
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	dbtimezone := beego.AppConfig.String("db.timezone")

	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	// fmt.Println(dsn)

	if dbtimezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(dbtimezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModel(new(UserInfo))
	orm.RegisterModel(new(AppInfo))
	orm.RegisterModel(new(ServerInfo))
	orm.RegisterModel(new(MonitorData))

	/*	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = false
	}*/
	/**********************redis初始化**************************************/
	/*	rdhost := beego.AppConfig.String("rd.host")
		rdport := beego.AppConfig.String("rd.port")
		//rdpassword := beego.AppConfig.String("rd.password")
		rddb, _ := beego.AppConfig.Int("rd.db")
		// 建立连接池
		RedisClient := &redis.Pool{
			// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
			MaxIdle:     beego.AppConfig.DefaultInt("redis.maxidle", 1),
			MaxActive:   beego.AppConfig.DefaultInt("redis.maxactive", 10),
			IdleTimeout: 180 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", rdhost+":"+rdport)
				if err != nil {
					return nil, err
				}
				// 选择db
				c.Do("SELECT", rddb)
				return c, nil
			},
		}*/
}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}
