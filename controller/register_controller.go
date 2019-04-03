package controller

import (
	"wScheduler/common"
	"wScheduler/entity"
)

type RegisterController struct {
	BaseController
}

func (this *RegisterController) ToRegister() {
	this.TplName = "register/register.html"
}

func (this *RegisterController) Register() {
	res := new(common.Result)
	username := this.GetString("Username")
	password := this.GetString("Password")
	if username == "" || password == "" {
		res.Message = " 有未填字段"
	} else {
		userInfo := new(entity.UserInfo)
		userInfo.Username = username
		userInfo.GetUserInfo()
		if userInfo.Id != 0 {
			res.Message = " 用户名已被注册"
		} else {
			userInfo.Password = password
			err := userInfo.SaveUserInfo()
			if err != nil {
				res.Message = " 注册失败"
			} else {
				res.Success = true
				res.Message = " 注册成功"
			}
		}
	}
	this.Data["res"] = res
	this.TplName = "jobinfo/res.html"
}
