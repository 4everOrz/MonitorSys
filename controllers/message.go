package controllers

import (
	"MonitorSys/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/go-gomail/gomail"
)

var (
	MlAddress    string
	MlSender     string
	MlSecureCode string
	ErrorArry    = make(map[string]int64)
	MessageDelay int64
	MlReceiver   string
	MlServer     string
	MlPort       int
	MsgUrl       string
	MsgUser      string
	MsgPassword  string
)

func init() {
	MlAddress = beego.AppConfig.String("ml.address")
	MlSender = beego.AppConfig.String("ml.sender")
	MlSecureCode = beego.AppConfig.String("ml.scode")
	MlServer = beego.AppConfig.String("ml.server")
	MlPort, _ = beego.AppConfig.Int("ml.port")
	MessageDelay = beego.AppConfig.DefaultInt64("messageDelay", 1800) //娶不到值，默认半小时
	MsgUrl = beego.AppConfig.String("msg.Url")
	MsgUser = beego.AppConfig.String("msg.User")
	MsgPassword = beego.AppConfig.String("msg.Password")
}
func messageSend(serverid, port, statuscode, possibility string) {
	if parsestr(serverid, port, statuscode) {
		messageCenter(serverid, port, statuscode, possibility)
	}
}

//判断解析
func parsestr(serverID, port, statuscode string) bool {
	fixstr := serverID + "/" + statuscode + "/" + port
	if ErrorArry[fixstr] == 0 {
		ErrorArry[fixstr] = time.Now().Unix() + MessageDelay
		return true
	} else {
		t1 := time.Now().Unix()
		t2 := ErrorArry[fixstr]
		if t2-t1 > 0 {
			return false
		} else {
			ErrorArry[fixstr] = t1 + MessageDelay
			return true
		}
	}
}

//消息发送中心
func messageCenter(serverid, port, statuscode, possibility string) {
	var message string
	var telmsg string
	var textTel string
	serverInf, _ := redisHMGET("ServerID:" + serverid)
	statuscodeInf, err1 := redisGet("StatusCode:" + statuscode)
	if err1 != nil {
		message = "运维人员,测试服务器疑似故障，请尽快排查！<br /> 详细信息:<br />  服务器名称=><a>" + serverInf["ServerName"] + "</a>,<br /> 端口=><a>" + port + "</a>,<br />  返回状态码=> <a>" + statuscode + "</a>,<br />可能原因=> <a>" + possibility + "</a>"
		telmsg = "您好,服务器疑似故障，请排查!详情,服务器:" + serverInf["ServerName"] + ",端口:" + port + ",状态码:" + statuscode + ",可能原因：" + possibility
	} else {
		message = "运维人员,疑似测试服务器发生故障，请尽快排查！<br /> 详细信息:<br />  服务器名称=><a>" + serverInf["ServerName"] + "</a>,<br /> 端口=><a>" + port + "</a>,<br />  返回状态码=> <a>" + statuscode + " : " + statuscodeInf + "</a>,<br />可能原因=> <a>" + possibility + "</a>"
		telmsg = "您好,服务器疑似故障,请排查!详情,服务器:" + serverInf["ServerName"] + ",端口:" + port + ",状态码:" + statuscode + "-" + statuscodeInf + ",可能原因：" + possibility
	}
	arry, count, _ := models.GetUserByRole("OperationUser")
	strarry := make([]string, count)

	for i := count; i > 0; i-- {
		strarry[i-1] = arry[i-1]["Mail"].(string)
		if i == count {
			textTel = arry[i-1]["Telphone"].(string)
		} else {
			textTel += "," + arry[i-1]["Telphone"].(string)
		}
	}
	TelMsgHttpGet(textTel, telmsg, count) //发送短信
	SendEMail(strarry, message)           //发送邮件
	return
}

//发送邮件
func SendEMail(receiver []string, message string) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", MlAddress, MlSender)
	m.SetHeader2("To", receiver)
	m.SetHeader("Subject", "重要信息")
	m.SetBody("text/html", message)
	d := gomail.NewPlainDialer(MlServer, MlPort, MlAddress, MlSecureCode) // 发送邮件服务器、端口、发件人账号、发件人密码
	if err2 := d.DialAndSend(m); err2 != nil {
		beego.Error("Send Email error:", err2)
	}
}

//发送手机短信 post
func TelMsgHttpPost(mobiles, msg string, mobilescount int64) {
	resp, err := http.Post(MsgUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader("userId="+MsgUser+"&password="+MsgPassword+"&pszMobis="+mobiles+"&pszMsg="+msg+"&iMobiCount="+strconv.FormatInt(mobilescount, 10)+"&pszSubPort=*"))
	if err != nil {
		beego.Error("post for sending message was failed")
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(" response error")
		return
	}
	fmt.Println(string(body))

}

//发送手机短信 get
func TelMsgHttpGet(mobiles, msg string, mobilescount int64) {
	str := "userId=" + MsgUser + "&password=" + MsgPassword + "&pszMobis=" + mobiles + "&pszMsg=" + msg + "&iMobiCount=1&pszSubPort=*"
	resp, err := http.Get(MsgUrl + str)
	if err != nil {
		beego.Error("get for sending message was failed")
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(" response error")
		return
	}
	fmt.Println(string(body))

}
