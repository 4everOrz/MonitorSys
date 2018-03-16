package controllers

import (
	"MonitoringSystemAPI/lib"
	"MonitoringSystemAPI/models"
	"encoding/json"
	"fmt"
	"log"
	"security"
	"strconv"
	"time"
	//"time"

	"github.com/astaxie/beego"
)

//result type
type Result struct {
	ResCode    int64
	ResMsg     string
	TotalCount int64
	ResData    interface{}
}

//var Resmap map[string]interface{}

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router /uregist [post]
func (u *UserController) Regist() {
	var user models.UserInfo
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	user.LoginTime = time.Now().Format("2006-01-02 15:04:05")
	user.RoleID = 1
	user.Password = security.Md5(user.Password + security.Md5(user.Password))
	if models.UserExist(user.LoginName) != 0 {
		fmt.Printf("LoginName has exist")

		u.Data["json"] = map[string]interface{}{"resCode": "1", "resMsg": "Change LoginName and try again", "resData": "null"}

	} else {
		UserID, _ := models.UserInfoAdd(&user)
		//Result{"1", "failed", 0, UserID}

		u.Data["json"] = map[string]interface{}{"resCode": "0", "resMsg": "success", "UserID": strconv.FormatInt(UserID, 10)}
	}
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	loginname		query 	string	true		"The Loginname for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /ulogin [post]
func (u *UserController) Login() {
	var user models.UserInfo
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	user.AccessToken = lib.GenToken()
	if models.Login(user.LoginName, user.Password, user.AccessToken) {
		Result := &Result{0, "success", 0, user.AccessToken}

		//b, _ := json.Marshal(Result)
		//u.Data["json"] = map[string]interface{}{"resCode": 0, "resMsg": "success", "AccessToken": user.AccessToken}
		u.Data["json"] = Result
	} else {
		u.Data["json"] = map[string]interface{}{"resCode": 1, "resMsg": "user not exist"}
	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /ulogout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}

// @Title changekey
// @Description change user's password
// @Param	loginname		query 	string	true		"The Loginname for login"
// @Param	password		query 	string	true		"The password for login"

// @Success 200 {string} change success
// @Failure 403 user not exist
// @router /changekey [post]
func (u *UserController) Changekey() {
	//	var user models.UserInfo
	newpassword := u.GetString("NewPassword")
	oldpassword := u.GetString("Password")
	loginname := u.GetString("loginname")
	accesstoken := u.GetString("AccessToken")
	security.Md5(newpassword + security.Md5(newpassword))
	security.Md5(oldpassword + security.Md5(oldpassword))
	//	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if models.VerifyUser(accesstoken, loginname) {
		if models.UpdateKey(loginname, oldpassword, newpassword) == 1 {
			u.Data["json"] = map[string]interface{}{"resCode": 0, "resMsg": "Change Ok"}
		}
	} else {
		u.Data["json"] = map[string]interface{}{"resCode": 1, "resMsg": "User verify error(Check AccessToken)!"}
		log.Println("Error:User verify error(Check AccessToken)!")
	}
	u.ServeJSON()
}

// @Title getbypage
// @Description get user by page
// @Param	loginname		query 	string	true		"The Loginname for login"
// @Param	password		query 	string	true		"The password for login"

// @Success 200 {string} change success
// @Failure 403 user not exist
// @router /getall [post]
func (u *UserController) Getall() {

	pa1ge := u.GetString("Page")
	pa1gesize := u.GetString("PageSize")
	page, _ := strconv.Atoi(pa1ge)
	pagesize, _ := strconv.Atoi(pa1gesize)
	var tt []interface{} = []interface{}{"", 0}
	users, total := models.UserInfoGetList(page, pagesize, tt)
	Res := &Result{0, "success", total, users}
	u.Data["json"] = Res
	u.ServeJSON()
}
