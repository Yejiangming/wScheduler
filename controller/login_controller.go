package controller

import (
	"wScheduler/common"
	"wScheduler/entity"
)

type LoginController struct {
	BaseController
}

func (this *RegisterController) ToLogin() {
	this.TplName = "login/login.html"
}

func (this *RegisterController) Login() {
	username := this.GetString("Username")
	password := this.GetString("Password")
	userInfo := new(entity.UserInfo)
	userInfo.Username = username
	userInfo.GetUserInfoByUsername()
	if userInfo.Password == password {
		this.SetSession("username", userInfo.Username)
		this.Redirect("/", 302)
	} else {
		res := new(common.Result)
		res.Message = " 用户名或密码错误"
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
	}
}
