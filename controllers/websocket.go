package controllers

//ws://172.28.162.234:8080/v1/ws
import (
	"MonitorSys/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	beego.Controller
}

// @Title wscoket
// @Description link
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {string} return
// @Failure 403 body is empty
// @router / [get].
func (this *WebSocketController) Join() {
	userid := this.GetString("UserID")
	useridint, _ := strconv.Atoi(userid)
	token := this.GetString("AccessToken")
	if userid != "" && token != "" && VerifyFromRedis(token, useridint) {
		ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			beego.Error("Cannot setup WebSocket connection:", err)
			return
		}

		Join(userid, ws) //新增客户端连接
		defer Leave(userid)
		for {
			_, serveridbyte, err := ws.ReadMessage() //读取客户端消息
			if err == nil && serveridbyte != nil {
				serverid := string(serveridbyte)
				if _, err := strconv.Atoi(serverid); err == nil {
					servertype <- UserServer{UserID: userid, ServerID: serverid}
				}
			} else {
				return
			}
		}
	} else {
		beego.Informational("Verify User Error,check LoginName and token,userid:", userid)
	}
}
func broadcast(event models.Event) {
	Data := event.Content
	serverdd := Data[0]["ServerID"]
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		ws := sub.Value.(Subscriber).Conn
		user := sub.Value.(Subscriber).Name
		if ws != nil {
			value := UserServerMap[user]
			if judge(value, serverdd) {
				Data[0]["MsgType"] = 1
			} else {
				Data[0]["MsgType"] = 0
			}
			datajson, _ := json.Marshal(&Data)
			err := ws.WriteMessage(websocket.TextMessage, datajson)
			if err != nil {
				unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}
func judge(value string, serverid interface{}) bool {
	var Apple bool
	if _, ok := serverid.(string); ok {
		if value == serverid {
			Apple = true
		} else {
			Apple = false
		}
	} else {
		valuestr, _ := strconv.Atoi(value)
		if valuestr == serverid {
			Apple = true
		} else {
			Apple = false
		}
	}
	return Apple
}
