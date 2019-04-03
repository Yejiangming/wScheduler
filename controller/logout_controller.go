package controller

type LogoutController struct {
	BaseController
}

func (this *LogoutController) Logout() {
	this.DelSession("username")
	this.Redirect("/login", 302)
}
