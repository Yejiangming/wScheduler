package controller

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare() {
	var m map[string]string = make(map[string]string)
	m["navigation"] = "layout/navigation.html"
	this.LayoutSections = m
	this.Layout = "layout/layout.html"

	sess := this.StartSession()
	username := sess.Get("username")
	uri := this.Ctx.Request.RequestURI
	if username == nil {
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
	}
}
