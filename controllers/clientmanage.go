package controllers

import (
	"MonitorSys/models"
	"container/list"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/gorilla/websocket"
)

type Subscription struct {
	Archive []models.Event //缓存数组
	New     <-chan models.Event
}
type UserServer struct {
	UserID   string
	ServerID string
}

func newEvent(msg []orm.Params) models.Event {
	return models.Event{int(time.Now().Unix()), msg}
}

func Join(user string, ws *websocket.Conn) {
	subscribe <- Subscriber{Name: user, Conn: ws}
}

func Leave(user string) {
	unsubscribe <- user
}

type Subscriber struct {
	Name string
	Conn *websocket.Conn
}

var (
	subscribe     = make(chan Subscriber, 20)
	unsubscribe   = make(chan string, 20)
	publish       = make(chan models.Event, 20)
	servertype    = make(chan UserServer, 20)
	subscribers   = list.New()
	UserServerMap = make(map[string]string)
)

func init() {
	go manage()
}
func manage() {
	for {
		select {
		case sub := <-subscribe: //读取到 new client 加入
			if !isUserExist(subscribers, sub.Name) {
				subscribers.PushBack(sub) //new client 放入队列尾
			}
		case event := <-publish: //有数据传入
			broadcast(event) //向各个客户端广播数据
		case userserver := <-servertype:
			UserServerMap[userserver.UserID] = userserver.ServerID
		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == unsub {
					subscribers.Remove(sub)
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
					}
					break
				}
			}
		}
	}
}

//判断用户是否已经在队列中
func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}
