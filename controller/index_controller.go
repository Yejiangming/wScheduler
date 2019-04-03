package controller

type IndexController struct {
	BaseController
}

func (this *IndexController) Get() {
	this.TplName = "index/index.html"
}
