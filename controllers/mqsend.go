package controllers

import (
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
)

type MqProducter struct {
	text        chan []byte
	channel     *amqp.Channel
	conn        *amqp.Connection
	connflag    chan int
	RabbitError chan *amqp.Error
}

var ProducterConnTag = true //生产者断网重连标志
var mq1 *MqProducter
var MQueueName string
var MExchange string
var MMqurl string
var MKey string

func init() {
	MQueueName = beego.AppConfig.String("rm.queueName")
	MExchange = beego.AppConfig.String("rm.exchange")
	MMqurl = beego.AppConfig.String("rm.mqurl")
	MKey = beego.AppConfig.String("rm.key")

}

func (mq *MqProducter) pingmq() {
	var err error

	for {
		mq.conn, err = amqp.Dial(MMqurl)

		if err == nil {
			mq.connflag <- 1
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (mq *MqProducter) mqConnect() error {
	var err error
	mq.conn, err = amqp.Dial(MMqurl)
	if err == nil {
		mq.channel, err = mq.conn.Channel()
	}
	return err
}
func (mq *MqProducter) mqClose() {
	mq.channel.Close()
	mq.conn.Close()
}

func (mq *MqProducter) push(msgContent []byte) {
	if mq.channel == nil {
		mq.channel, _ = mq.conn.Channel()
	}

	err := mq.channel.Publish(MExchange, MKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msgContent,
	})
	if err != nil {
		log.Println("发送失败")
	}

	/*if err := mq.channel.Confirm(false); err != nil {
		mq1.conn = nil
		mq1.channel = nil
		mq1.mqConnect()
	} */ //丢掉一条数据
}

/*func MqfeildSend() *MqProducter {

	if mq1 == nil {
		mq1 = &MqProducter{
			text: make(chan []byte),
		}
		err1 := mq1.mqConnect()
		if err1 != nil {
			log.Println("mqconnect error:", err1)

		}
		go mq1.translate()
		go producterConnListener()
	}
	return mq1
}
*/
func MqfeildSend() *MqProducter {
	connSucessFlag := make(chan int, 1)
	mq1 := GetProducterConn(connSucessFlag)
	go func() {
		for {
			<-mq1.connflag
			go mq1.translate()
		}
	}()

	return mq1
}

//转发数据
func (mq *MqProducter) translate() {
	var saveData []byte
	for {
		select {
		case saveData = <-mq.text:
			mq.push(saveData)

		}
	}
	defer mq.mqClose()
}

/*****************************************/

func (c *MqProducter) connectToRabbitMQ() {
	var err error
	for {
		c.conn, err = amqp.Dial(MMqurl)

		if err == nil {
			c.connflag <- 1
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}
func GetProducterConn(connFlag chan int) *MqProducter {
	if mq1 == nil {
		mq1 = &MqProducter{
			RabbitError: make(chan *amqp.Error),
			connflag:    connFlag,
			text:        make(chan []byte),
		}
		mq1.connectToRabbitMQ()
		mq1.conn.NotifyClose(mq1.RabbitError) //set a connection error listener
		go mq1.rabbitConnector()              //goroutine about reconnect
	}
	return mq1
}
func (c *MqProducter) rabbitConnector() {
	var rabbitErr *amqp.Error
	for {
		rabbitErr = <-c.RabbitError //connection closed error reader
		if rabbitErr != nil {
			ProducterConnTag = false
			log.Println("mqproducter connect error,try reconnect!")
			c.conn = nil
			c.channel = nil
			c.connectToRabbitMQ()
			c.channel, _ = c.conn.Channel()
			ProducterConnTag = true
			log.Println("mqproducter connect ok!")
			c.RabbitError = make(chan *amqp.Error) //reinitialize the connection close error
			c.conn.NotifyClose(c.RabbitError)      //reset the connection listener
		}
	}
}
