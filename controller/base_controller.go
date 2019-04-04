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

	sess := this.StartSession()
	username := sess.Get("username")
	userinfo := new(entity.UserInfo)
	uri := this.Ctx.Request.RequestURI
	if username == nil {
		userinfo.Username = "未登陆"
		this.Data["userinfo"] = userinfo
		if uri == "/register" {
			return
		}
		if uri == "/login" {
			return
		}
		if uri == "/logout" {
			return
		}
		this.Redirect("/login", 302)
	} else {
		userinfo.Username = username.(string)
		this.Data["userinfo"] = userinfo
	}
}
