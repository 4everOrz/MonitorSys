package controllers

import (
	"MonitorSys/lib"
	"MonitorSys/models"
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
	TotalCount interface{}
	ResData    interface{}
}
type UserController struct {
	beego.Controller
}

func Userinit() {
	mapArry, _, _ := models.GetAllUser2()
	for _, value := range mapArry {
		if err := redisHMSET("UserID:"+strconv.FormatInt(value["UserID"].(int64), 10), value); err != nil {
			log.Println("error:", err)
		}
	}
}

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
	user.AccessToken = ""
	user.LoginTime = time.Now().Format("2006-01-02 15:04:05")
	user.RoleID = 2
	user.Password = security.Md5(user.Password + security.Md5(user.Password))
	user.Ustatus = "正常"
	if user.UserName != "" && user.LoginName != "" && user.Password != "" {
		if models.UserExist(user.LoginName) != 0 {
			result := &Result{1, "LoginName has exist", 0, nil}
			u.Data["json"] = result
		} else {
			UserID, _ := models.UserInfoAdd(&user)                    //存Mysql
			redisHMSET("UserID:"+strconv.FormatInt(UserID, 10), user) //存Redis
			usermap["UserID"] = UserID
			result := &Result{0, "success", 0, usermap}
			u.Data["json"] = result
		}
	} else {
		result := &Result{1, "chack field format", 0, nil}
		u.Data["json"] = result
	}
	u.ServeJSON()
}

// @router /ulogin [post]
func (u *UserController) Login() {
	loginname := u.GetString("LoginName")
	password := u.GetString("Password")
	accesstoken := lib.GenToken()
	logintime := time.Now().Format("2006-01-02 15:04:05")
	if models.UserExist(loginname) == 0 {
		result := &Result{1, "LoginName not exist", 0, nil}
		u.Data["json"] = result
	} else {
		userid, booler, arry := models.Login(loginname, password, accesstoken, logintime)
		mapuser, _ := redisHMGET("UserID:" + strconv.FormatInt(userid, 10))
		mapuser["AccessToken"] = accesstoken
		mapuser["LoginTime"] = logintime
		redisHMSET("UserID:"+strconv.FormatInt(userid, 10), mapuser) //Redis更新AccessToken
		if booler {
			arry := arry[0]
			result := &Result{0, "success", 0, arry}
			u.Data["json"] = result
		} else {
			result := &Result{1, "Wrong password or Invalid user ", 0, nil}
			u.Data["json"] = result
		}
	}
	u.ServeJSON()
}

// @router /ulogout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}

// @router /changekey [post]
func (u *UserController) Changekey() {
	newpassword := u.GetString("NewPassword")
	oldpassword := u.GetString("Password")
	userid, _ := u.GetInt("UserID")
	accesstoken := u.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		newpsw := security.Md5(newpassword + security.Md5(newpassword))
		oldpsw := security.Md5(oldpassword + security.Md5(oldpassword))
		if models.UpdateKey(userid, oldpsw, newpsw) == 1 {
			mapuser, _ := redisHMGET("UserID:" + strconv.Itoa(userid))
			mapuser["Password"] = newpassword
			redisHMSET("UserID:"+strconv.Itoa(userid), mapuser) //Redis更新密码
			result := &Result{0, "success", 0, nil}
			u.Data["json"] = result
		} else {
			result := &Result{1, "check your password", 0, nil}
			u.Data["json"] = result
		}
	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		u.Data["json"] = result
	}
	u.ServeJSON()
}

// @router /getall [post]
/*func (u *UserController) Getall() {
	userid, _ := u.GetInt("UserID")
	accesstoken := u.GetString("AccessToken")
	if models.VerifyUser(accesstoken, userid) {
		page, _ := u.GetInt("Page")
		pagesize, _ := u.GetInt("PageSize")
		var tt []interface{} = []interface{}{"", 0}
		users, tcount := models.UserInfoGetList(page, pagesize, tt)
		Res := &Result{0, "success", tcount, users}
		u.Data["json"] = Res
	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		u.Data["json"] = result
	}

	u.ServeJSON()
}*/
// @router /getall [post]
func (u *UserController) Getall() {
	userid, _ := u.GetInt("UserID")
	accesstoken := u.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		users, tcount, _ := models.GetAllUser()
		Res := &Result{0, "success", tcount, users}
		u.Data["json"] = Res
	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		u.Data["json"] = result
	}

	u.ServeJSON()
}

// @router /uupdate [post]
func (u *UserController) UpdateOne() {
	var err error
	userid, _ := u.GetInt("UserID")
	accesstoken := u.GetString("AccessToken")
	if VerifyFromRedis(accesstoken, userid) {
		var user models.UserInfo
		user.UserID, _ = u.GetInt64("UserID_C")
		user.UserName = u.GetString("UserName")
		user.Telphone = u.GetString("Telphone")
		user.Mail = u.GetString("Mail")
		user.RoleID, _ = u.GetInt("RoleID")
		user.Ustatus = u.GetString("Ustatus")
		newpassword := u.GetString("Password")
		/***********Redis 匹配字段********************/
		usermap, _ := redisHMGET("UserID:" + strconv.FormatInt(user.UserID, 10))
		usermap["UserName"] = user.UserName
		usermap["Telphone"] = user.Telphone
		usermap["Mail"] = user.Mail
		usermap["RoleID"] = strconv.Itoa(user.RoleID)
		usermap["Ustatus"] = user.Ustatus
		/*********************************************/
		if newpassword == "" {
			_, err = models.UserUpdate2(&user)
		} else {
			user.Password = security.Md5(newpassword + security.Md5(newpassword))
			_, err = models.UserUpdate1(&user)
			usermap["Password"] = user.Password //redis更新密码
		}
		redisHMSET("UserID:"+strconv.FormatInt(user.UserID, 10), usermap) //redis存入更改后的信息
		lib.FailOnErr(err, "UserUpdate error")
		if err != nil {
			result := &Result{1, "update failed", 0, nil}
			u.Data["json"] = result
		} else {
			result := &Result{0, "success", 0, nil}
			u.Data["json"] = result
		}
	} else {
		result := &Result{1, "user verify error check userid and token", 0, nil}
		u.Data["json"] = result
	}
	u.ServeJSON()
}

//CheckToken(redis)
func VerifyFromRedis(AccessToken string, UserID int) bool {
	usermap, err := redisHMGET("UserID:" + strconv.Itoa(UserID))
	if err != nil {
		return models.VerifyUser(AccessToken, UserID)
	} else {
		if usermap["AccessToken"] == AccessToken {
			return true
		} else {
			return false
		}
	}
}
