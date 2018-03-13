package Data

import (
	"MonitoringSystemAPI/lib"
	"MonitoringSystemAPI/models"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

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
func mqConnect(m *Mqfeild) {
	var err error
	conn, err = amqp.Dial(m.Mqurl)
	lib.FailOnErr(err, "failed to connect tp rabbitmq")

	channel, err = conn.Channel()
	lib.FailOnErr(err, "failed to open a channel")
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
func push(m *Mqfeild, msgContent []byte) {
	if channel == nil {
		mqConnect(m)
	}
	channel.Publish(m.Exchange, m.Key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msgContent,
	})
}

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
		dataid, _ := models.AddOne(&md)
		ma := strconv.FormatInt(dataid, 10)
		fmt.Printf("Data:", ma)
		fmt.Printf("receve msg is :%s -- %d\n", *s, count)
	}

}
func MqSend(ms []byte) {
	m := mqstr()
	mqConnect(m)
	push(m, ms)
	defer mqClose()
}
func MqReceive() {
	m := mqstr()
	mqConnect(m)
	receive(m)
	defer mqClose()
}
