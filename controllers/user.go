package controllers

import (
	"MonitoringSystemAPI/models"
	"encoding/json"
	"MonitoringSystemAPI/lib"
	"net/http"
	"github.com/astaxie/beego"
	"strconv"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	var user models.UserInfo
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	uid,_ := models.UserInfoAdd(&user)
	u.Data["json"] = map[string]string{"uid": strconv.FormatInt(uid,10)  }
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *UserController) Login() {
	username := u.GetString("username")
	password := u.GetString("password")
	if models.Login(username, password) {
		token:=lib.GenToken()
		cookie := http.Cookie{Name: "Authorization", Value: token, Path: "/", MaxAge: 3600}

		http.SetCookie(u.Ctx.ResponseWriter, &cookie)
		u.Data["json"] = token

	} else {
		u.Data["json"] = "user not exist"
	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}

