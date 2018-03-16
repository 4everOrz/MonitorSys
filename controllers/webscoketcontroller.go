package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	beego.Controller
}

// Get method handles GET requests for WebSocketController.
func (this *WebSocketController) Get() {
	// Safe check.
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	this.TplName = "websocket.html"
	this.Data["IsWebSocket"] = true
	this.Data["UserName"] = uname
}

// @Title wscoket
// @Description link
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {string} return
// @Failure 403 body is empty
// @router / [get].
func (this *WebSocketController) Join() {

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		M, _ := json.Marshal("helloworld")
		fmt.Printf("msg", p)
		ws.WriteMessage(websocket.TextMessage, M)
		//defer ws.Close()
	}
}
