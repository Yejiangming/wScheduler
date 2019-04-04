package controller

import (
	"wScheduler/entity"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare() {
	var m map[string]string = make(map[string]string)
	m["userstate"] = "layout/userstate.html"
	m["navigation"] = "layout/navigation.html"
	this.LayoutSections = m
	this.Layout = "layout/layout.html"

	username := this.GetSession("username")
	userinfo := new(entity.UserInfo)
	if username == nil {
		userinfo.Username = "未登陆"
	} else {
		userinfo.Username = username.(string)
	}
	this.Data["userinfo"] = userinfo

	uri := this.Ctx.Request.RequestURI
	if username == nil && uri != "/register" && uri != "/login" && uri != "/logout" {
		this.Redirect("/login", 302)
	}
}
