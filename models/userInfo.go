package models

import (
	"MonitoringSystemAPI/lib"
	"fmt"
	"gree/GrihCommon/security"

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
	query.OrderBy("UserID").Limit(pageSize, offset).All(&list)

	return list, total
}

func UserInfoGetListByIds(userId int64) ([]*UserInfo, error) {
	list1 := make([]*UserInfo, 0)
	var list []orm.Params
	//list:=[]orm.Params
	var err error
	_, err = orm.NewOrm().Raw("select UserID,LoginName,Password,UserName,Telphone,Email,RoleID,AccessToken,LoginTime from userInfo where UserID=? order by UserID asc", userId).Values(&list)
	for k, v := range list {
		fmt.Println(k, v)
	}

	fmt.Println(list)
	return list1, err
}

func UserInfoAdd(userinfo *UserInfo) (int64, error) {
	return orm.NewOrm().Insert(userinfo)
}

func UserInfoGetById(UserID int) (*UserInfo, error) {
	a := new(UserInfo)

	err := orm.NewOrm().QueryTable(TableName("userInfo")).Filter("UserID", UserID).One(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func Login(loginname, password, token string) bool {
	var booler bool
	a := new(UserInfo)
	o := orm.NewOrm()
	o.Raw("SELECT * from userInfo where LoginName = ?", loginname).QueryRow(&a)

	if a.Password == security.Md5(password+security.Md5(password)) {
		_, err := o.Raw("UPDATE userInfo Set AccessToken = ? WHERE LoginName = ? AND Password = ?", token, a.LoginName, a.Password).Exec()
		if err == nil {
			booler = true
		}
	} else {
		booler = false
	}
	return booler
}
func (a *UserInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(a, fields...); err != nil {
		return err
	}
	return nil
}
func UserExist(loginname string) int64 {
	var user []UserInfo
	num, err := orm.NewOrm().Raw("SELECT * FROM userInfo WHERE LoginName = ?", loginname).QueryRows(&user)
	var RowAffect int64
	if err == nil {

		RowAffect = num
	}
	return RowAffect
}
func VerifyUser(token, loginname string) bool {
	user := new(UserInfo)
	orm.NewOrm().Raw("SELECT * FROM userInfo WHERE LoginName = ?", loginname).QueryRow(&user)
	if user.AccessToken == token {
		return true
	} else {
		return false
	}
}
func UpdateKey(loginname, oldkey, newkey string) int64 {
	user := new(UserInfo)
	rowaffect, err := orm.NewOrm().Raw("UPDATE userInfo SET Password = ? WHERE LoginName=? AND Password =?", newkey, loginname, oldkey).QueryRows(&user)
	lib.FailOnErr(err, "Update error")
	return rowaffect
}
