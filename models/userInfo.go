package models

import (
	"gree/GrihCommon/security"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserInfo struct {
	UserID      int64  `orm:"pk;column(UserID)"`
	LoginName   string `orm:"column(LoginName)"`
	Password    string `orm:"column(Password)"`
	UserName    string `orm:"column(UserName)"`
	Telphone    string `orm:"column(Telphone)"`
	Mail        string `orm:"column(Mail)"`
	AccessToken string `orm:"column(AccessToken)"`
	LoginTime   string `orm:"column(LoginTime)"`
	RoleID      int    `orm:"column(RoleID)"`
	Ustatus     string `orm:"column(Ustatus)"`
}

var (
	UserList map[string]*UserInfo
)

func (a *UserInfo) TableName() string {
	return TableName("userInfo")
}

func UserInfoGetList(page, pageSize int, filters ...interface{}) ([]*UserInfo, int64) {
	user := new(UserInfo)
	offset := (page - 1) * pageSize
	list := make([]*UserInfo, 0)
	query := orm.NewOrm().QueryTable(user)
	/*	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}*/
	total, _ := query.Count()
	query.OrderBy("UserID").Limit(pageSize, offset).All(&list, "UserID", "UserName", "LoginName", "Telphone", "Mail", "LoginTime", "RoleID", "Ustatus")

	return list, total
}

func UserInfoGetListByIds(userId int64) ([]orm.Params, error) {
	var list []orm.Params
	_, err := orm.NewOrm().Raw("select UserID,LoginName,Password,UserName,Telphone,Email,RoleID,AccessToken,LoginTime from userInfo where UserID=? order by UserID asc", userId).Values(&list)
	return list, err
}

//insert one
func UserInfoAdd(userinfo *UserInfo) (int64, error) {
	return orm.NewOrm().Insert(userinfo)
}

//get one by ID
func UserInfoGetById(UserID int) (*UserInfo, error) {
	a := new(UserInfo)
	err := orm.NewOrm().QueryTable(TableName("userInfo")).Filter("UserID", UserID).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

//user login
func Login(loginname, password, token, logintime string) (int64, bool, []orm.Params) {
	var booler bool
	var arry []orm.Params
	a := new(UserInfo)
	o := orm.NewOrm()
	o.Raw("SELECT * from userInfo where LoginName = ?", loginname).QueryRow(&a)
	if a.Ustatus == "正常" {
		if a.Password == security.Md5(password+security.Md5(password)) {
			_, err := o.Raw("UPDATE userInfo Set AccessToken = ?,LoginTime = ? WHERE UserID = ?", token, logintime, a.UserID).Exec()
			if err == nil {
				booler = true
				arry, _ = UserJoinRole(a.UserID)
			}
		} else {
			booler = false
		}
	} else {
		booler = false
	}
	return a.UserID, booler, arry
}

//check loginname
func UserExist(loginname string) int64 {
	var user []UserInfo
	num, err := orm.NewOrm().Raw("SELECT * FROM userInfo WHERE LoginName = ?", loginname).QueryRows(&user)
	var RowAffect int64
	RowAffect = 0
	if err == nil {

		RowAffect = num
	}
	return RowAffect
}

//check token(Mysql)
func VerifyUser(token string, userid int) bool {
	user := new(UserInfo)
	orm.NewOrm().Raw("SELECT * FROM userInfo WHERE UserID = ?", userid).QueryRow(&user)
	if user.AccessToken == token {
		return true
	} else {
		return false
	}
}

//update password
func UpdateKey(userid int, oldkey, newkey string) int64 {
	userInfo := new(UserInfo)
	rowaffect, err := orm.NewOrm().QueryTable(userInfo).Filter("UserID", userid).Filter("Password", oldkey).Update(orm.Params{
		"Password": newkey})
	beego.Error("error:", err)
	return rowaffect
}

//userInfo inner join roleInfo
func UserJoinRole(userid int64) ([]orm.Params, error) {
	var maps []orm.Params
	_, err := orm.NewOrm().Raw("SELECT * FROM userInfo AS u Inner join roleInfo AS r on u.RoleID = r.RoleID WHERE UserID=?", userid).Values(&maps, "UserID", "AccessToken", "RoleName", "RoleID")
	return maps, err
}

//get users by RoleName
func GetUserByRole(rolename string) ([]orm.Params, int64, error) {
	var maps []orm.Params
	affectrows, err := orm.NewOrm().Raw("SELECT * FROM userInfo AS u INNER JOIN roleInfo AS r on u.RoleID=r.RoleID WHERE RoleName=?", rolename).Values(&maps, "UserName", "Mail", "Telphone")
	return maps, affectrows, err
}

//update one with password
func UserUpdate1(user *UserInfo) (int64, error) {
	userInfo := new(UserInfo)
	rowaffect, err := orm.NewOrm().QueryTable(userInfo).Filter("UserID", user.UserID).Update(orm.Params{
		"UserName": user.UserName, "Password": user.Password, "Telphone": user.Telphone, "Mail": user.Mail, "RoleID": user.RoleID, "Ustatus": user.Ustatus})
	return rowaffect, err
}

//only update userInfo
func UserUpdate2(user *UserInfo) (int64, error) {
	userInfo := new(UserInfo)
	rowaffect, err := orm.NewOrm().QueryTable(userInfo).Filter("UserID", user.UserID).Update(orm.Params{
		"UserName": user.UserName, "Telphone": user.Telphone, "Mail": user.Mail, "RoleID": user.RoleID, "Ustatus": user.Ustatus})
	return rowaffect, err
}

//get all userInfo inner join roleInfo
func GetAllUser() ([]orm.Params, int64, error) {
	var arry []orm.Params
	affectrows, err := orm.NewOrm().Raw("SELECT u.UserID,u.UserName,u.LoginName,u.Telphone,u.Mail,u.LoginTime,u.Ustatus,r.RoleName,r.RoleID FROM userInfo AS u INNER JOIN roleInfo AS r on u.RoleID=r.RoleID ").Values(&arry, "UserID", "LoginName", "UserName", "Mail", "Telphone", "LoginTime", "RoleID", "RoleName", "Ustatus")
	return arry, affectrows, err
}

//get all userInfo
func GetAllUser2() ([]orm.Params, int64, error) {
	var paramsArry []orm.Params
	userInfo := new(UserInfo)
	affectrows, err := orm.NewOrm().QueryTable(userInfo).Values(&paramsArry)
	return paramsArry, affectrows, err
}
