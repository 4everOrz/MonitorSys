package models

import (
	"github.com/astaxie/beego/orm"
)

type RoleInfo struct {
	RoleID   int64  `orm:"pk;column(RoleID)"`
	RoleName string `orm:"column(RoleName)"`
	Level    int    `orm:"column(Level)"`
}

func GetRoleInfo() ([]orm.Params, int64, error) {
	var role []orm.Params
	num, err := orm.NewOrm().Raw("Select * From roleInfo").Values(&role, "RoleID", "RoleName", "Level")
	return role, num, err
}
