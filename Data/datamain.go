package Data

import (
	"MonitoringSystemAPI/lib"
	"MonitoringSystemAPI/models"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel
var count = 0
var Msg chan string

type Mqfeild struct {
	QueueName string
	Exchange  string
	Mqurl     string
	Key       string
}

func mqstr() *Mqfeild {

	m := new(Mqfeild)
	m.QueueName = beego.AppConfig.String("rm.queueName")
	m.Exchange = beego.AppConfig.String("rm.exchange")
	m.Mqurl = beego.AppConfig.String("rm.mqurl")
	m.Key = beego.AppConfig.String("rm.key")
	return m
} /**/
func mqConnect(m *Mqfeild) error {
	var err error
	conn, err = amqp.Dial(m.Mqurl)
	lib.FailOnErr(err, "failed to connect tp rabbitmq")

	channel, err = conn.Channel()
	lib.FailOnErr(err, "failed to open a channel")
	return err
}
func mqClose() {
	channel.Close()
	conn.Close()
}
func bytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}
func push(m *Mqfeild, msgContent []byte) error {
	if channel == nil {
		err := mqConnect(m)
		lib.FailOnErr(err, "connect mq wrong")
	}
	err := channel.Publish(m.Exchange, m.Key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msgContent,
	})
	return err
}

func MqSend(ms []byte) error {
	var err error
	m := mqstr()
	err = mqConnect(m)
	lib.FailOnErr(err, "connect mq wrong")
	if err == nil {
		err := push(m, ms)
		lib.FailOnErr(err, "something wrong on pushing to mq")
	}
	defer mqClose()
	return err
}
func MqReceive() {
	m := mqstr()
	mqConnect(m)
	receive(m)
	defer mqClose()
}

/****************************************************/
func receive(m *Mqfeild) {
	var s *string
	var md models.MonitorData
	if channel == nil {
		mqConnect(m)
	}
	msgs, err := channel.Consume(m.QueueName, "", true, false, false, false, nil)
	lib.FailOnErr(err, "")

	for d := range msgs {
		s = bytesToString(&(d.Body))
		count++
		json.Unmarshal(d.Body, &md)
		md.SubTime = time.Now().Format("2006-01-02 15:04:05")
		dataid, err1 := models.AddOne(&md) //监控数据直接存入Mysql
		if err1 != nil {
			lib.FailOnErr(err1, "监控数据存入Mysql出现错误")
			return
		}
		dataarry := models.GetMDbyID(dataid) //关系查询
		json.Marshal(dataarry)               //转为json格式
		//WebSocket传递给前端
		//统计预警
		fmt.Println("dataarry", dataarry)
		ma := strconv.FormatInt(dataid, 10)
		fmt.Printf("Data:", ma)
		fmt.Printf("receve msg is :%s -- %d\n", *s, count)
	}

}
