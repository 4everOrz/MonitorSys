package controllers

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

//redis-server.exe redis.windows.conf
var (
	RedisClient *redis.Pool
	RDhost      string
	RDport      string
	RDdb        int
)

func init() {
	/**********************redis初始化**************************************/
	RDhost := beego.AppConfig.String("rd.host")
	RDport := beego.AppConfig.String("rd.port")
	RDpassword := beego.AppConfig.String("rd.password")
	RDdb, _ := beego.AppConfig.Int("rd.db")
	RedisClient = &redis.Pool{ // 建立连接池
		MaxIdle:     beego.AppConfig.DefaultInt("rd.maxidle", 1), //从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxActive:   beego.AppConfig.DefaultInt("rd.maxactive", 10),
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", RDhost+":"+RDport, redis.DialPassword(RDpassword))
			if err != nil {
				return nil, err
			}
			c.Do("SELECT", RDdb)  // 选择db
			if RDpassword != "" { //密码校验
				if _, err := c.Do("AUTH", RDpassword); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
	}
}
func RedisInit() { //向redis里初始化存储数据
	Appinit()
	Serverinit()
	Roleinit()
	Userinit()
	StatusCodeInit()
}

func redisSet(key, value interface{}) error {
	rconn := RedisClient.Get()
	_, err := rconn.Do("SET", key, value)
	defer rconn.Close()
	return err

}
func redisGet(key interface{}) (string, error) {
	rconn := RedisClient.Get()
	v, err := redis.String(rconn.Do("GET", key))
	defer rconn.Close()
	return v, err
}
func redisExit(key interface{}) (bool, error) {
	rconn := RedisClient.Get()
	booler, err := redis.Bool(rconn.Do("EXISTS", key))
	defer rconn.Close()
	return booler, err
}
func redisDelete(key interface{}) error {
	rconn := RedisClient.Get()
	_, err := rconn.Do("DEL", key)
	defer rconn.Close()
	return err
}
func redisHMSET(key string, obj interface{}) error {
	rconn := RedisClient.Get()
	_, err := rconn.Do("HMSET", redis.Args{}.Add(key).AddFlat(obj)...)
	defer rconn.Close()
	return err
}
func redisHMGET(key string) (map[string]string, error) {
	rconn := RedisClient.Get()
	v, err := redis.StringMap(rconn.Do("HGETALL", key))
	defer rconn.Close()
	return v, err
}

func StatusCodeInit() {
	redisSet("StatusCode:200", "执行成功")
	redisSet("StatusCode:400", "api验证失败")
	redisSet("StatusCode:401", "json包解析失败")
	redisSet("StatusCode:402", "token无效")
	redisSet("StatusCode:403", "数据校验错误")
	redisSet("StatusCode:404", "上传文件时出错")
	redisSet("StatusCode:405", "未上传文件")
	redisSet("StatusCode:406", "数据非法")
	redisSet("StatusCode:407", "备份删除失败")
	redisSet("StatusCode:408", "无效操作")
	redisSet("StatusCode:409", "保存失败")
	redisSet("StatusCode:500", "服务器内部错误")
	redisSet("StatusCode:11006", "获取服务器时间失败")
	redisSet("StatusCode:11007", "Error in server type")
	redisSet("StatusCode:1", "connect fail")
	redisSet("StatusCode:2", "disconnect")
	redisSet("StatusCode:3", "time out")
	redisSet("StatusCode:4", "normal")
	redisSet("StatusCode:x", "测试数据，请忽略此条信息")
}
