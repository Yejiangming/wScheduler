package entity

import (
	"github.com/astaxie/beego/orm"
)

type UserInfo struct {
	Id       int
	Username string
	Password string
	Mailbox  string
}

func (this *UserInfo) GetUserInfoByUsername() error {
	o := orm.NewOrm()
	err := o.Read(this, "username")
	return err
}

func (this *UserInfo) GetUserInfoByMailbox() error {
	o := orm.NewOrm()
	err := o.Read(this, "mailbox")
	return err
}

func (this *UserInfo) SaveUserInfo() error {
	_, err := orm.NewOrm().Insert(this)
	return err
}
