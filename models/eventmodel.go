package models

import (
	"github.com/astaxie/beego/orm"
)

type Event struct {
	Timestamp int
	Content   []orm.Params
}
