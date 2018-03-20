package controllers

import (
	"MonitoringSystemAPI/lib"
	"MonitoringSystemAPI/models"
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
	TotalCount int
	ResData    interface{}
}

//var Resmap map[string]interface{}

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title 添加用户
// @Description 实现新用户的注册
// @Param  UserName	     formData     string 	  true		"用户姓名"
// @Param  LoginName    formData      string      true       "登录名"
// @Param  Password     formData      string      true      "密码"
// @Param  Telphone     formData      string     false       "电话"
// @Param  Mail         formData      string     false       "邮件"
// @Success 200 {object} models.User.Result
// @Failure 403 body is empty
// @router /uregist [post]
func (u *UserController) Regist() {
	var user models.UserInfo
	var usermap map[string]interface{}
	usermap = make(map[string]interface{})
	//json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	user.UserName = u.GetString("UserName")
	user.LoginName = u.GetString("LoginName")
	user.Password = u.GetString("Password")
	user.Telphone = u.GetString("Telphone")
	user.Mail = u.GetString("Mail")
	user.LoginTime = time.Now().Format("2006-01-02 15:04:05")
	user.RoleID = 1
	user.Password = security.Md5(user.Password + security.Md5(user.Password))
	if models.UserExist(user.LoginName) != 0 {
		result := &Result{1, "LoginName has exist", 0, nil}
		u.Data["json"] = result

	} else {
		UserID, _ := models.UserInfoAdd(&user)
		//Result{"1", "failed", 0, UserID}//strconv.FormatInt(UserID, 10, 64)
		usermap["UserID"] = UserID
		result := &Result{0, "success", 0, usermap}
		u.Data["json"] = result
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
	//var usermap map[string]interface{}
	//usermap = make(map[string]interface{})
	loginname := u.GetString("LoginName")
	password := u.GetString("Password")
	accesstoken := lib.GenToken()
	if models.UserExist(loginname) == 0 {
		result := &Result{1, "LoginName not exist", 0, nil}
		u.Data["json"] = result
	} else {
		booler, arry := models.Login(loginname, password, accesstoken)
		if booler {
			arry := arry[0]
			result := &Result{0, "success", 0, arry}
			u.Data["json"] = result
		} else {
			result := &Result{1, "wrong password", 0, nil}
			u.Data["json"] = result
		}
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

	//	var usermap map[string]interface{}
	//	usermap = make(map[string]interface{})
	newpassword := u.GetString("NewPassword")
	oldpassword := u.GetString("Password")
	loginname := u.GetString("LoginName")
	accesstoken := u.GetString("AccessToken")

	//	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if models.VerifyUser(accesstoken, loginname) {
		newpsw := security.Md5(newpassword + security.Md5(newpassword))
		oldpsw := security.Md5(oldpassword + security.Md5(oldpassword))
		if models.UpdateKey(loginname, oldpsw, newpsw) == 1 {

			result := &Result{0, "success", 0, nil}
			u.Data["json"] = result
		} else {
			result := &Result{1, "check your password", 0, nil}
			u.Data["json"] = result
		}
	} else {
		result := &Result{1, "User verify error(Check AccessToken)", 0, nil}
		u.Data["json"] = result
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
	users, _ := models.UserInfoGetList(page, pagesize, tt)
	Res := &Result{0, "success", 0, users}
	u.Data["json"] = Res
	u.ServeJSON()
}
