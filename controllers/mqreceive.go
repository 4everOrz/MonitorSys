package controllers

import (
	"time"

	"github.com/streadway/amqp"
)

var Timercount = 0

type MqConsumer struct {
	channel     *amqp.Channel
	conn        *amqp.Connection
	connflag    chan int
	RabbitError chan *amqp.Error
}

var receiver = make(chan MqConsumer, 10)
var mq2 *MqConsumer

func (mq *MqConsumer) mqConnect() error {
	var err error

	mq.conn, err = amqp.Dial(MMqurl)
	if err == nil {
		mq.channel, err = mq.conn.Channel()
	}
	return err
}
func (mq *MqConsumer) mqClose() {
	mq.channel.Close()
	mq.conn.Close()
}

func (c *MqConsumer) connectToRabbitMQ() {
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
func (c *MqConsumer) rabbitConnector() {
	var rabbitErr *amqp.Error
	for {
		rabbitErr = <-c.RabbitError
		if rabbitErr != nil {
			c.connectToRabbitMQ()
			c.RabbitError = make(chan *amqp.Error)
			c.conn.NotifyClose(c.RabbitError)
		}
	}
}
func GetConsumerConn(connFlag chan int) *MqConsumer {
	if mq2 == nil {
		mq2 = &MqConsumer{
			RabbitError: make(chan *amqp.Error),
			connflag:    connFlag,
		}
		mq2.connectToRabbitMQ()
		mq2.conn.NotifyClose(mq2.RabbitError)
		go mq2.rabbitConnector()
	}
	return mq2
}
