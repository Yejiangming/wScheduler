package entity

import (
	"github.com/astaxie/beego/orm"
)

type UserInfo struct {
	Id       int
	Username string
	Password string
}

func (this *UserInfo) GetUserInfo() error {
	o := orm.NewOrm()
	err := o.Read(this, "username")
	return err
}

func (this *UserInfo) SaveUserInfo() error {
	_, err := orm.NewOrm().Insert(this)
	return err
}
