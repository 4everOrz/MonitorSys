package models

type RoleInfo struct {
	RoleID   int64  `orm:"pk;column(RoleID)"`
	RoleName string `orm:"column(RoleName)"`
	Level    int    `orm:"column(Level)"`
}
