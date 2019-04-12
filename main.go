package main

import (
	"wScheduler/controller"
	"wScheduler/entity"
	"wScheduler/job"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbUser := beego.AppConfig.String("dbUser")
	dbPassword := beego.AppConfig.String("dbPassword")
	dbName := beego.AppConfig.String("dbName")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dbUser+":"+dbPassword+"@/"+dbName+"?charset=utf8")
	orm.RegisterModel(&entity.JobInfo{}, &entity.JobSnapshot{}, &entity.UserInfo{})

	beego.BConfig.WebConfig.Session.SessionOn = true

	beego.Router("/", &controller.IndexController{})
	beego.Router("/jobinfo/add", &controller.JobInfoController{}, "get:ToAdd")
	beego.Router("/jobinfo/add", &controller.JobInfoController{}, "post:Add")
	beego.Router("/jobinfo/list", &controller.JobInfoController{}, "*:List")
	beego.Router("/jobinfo/delete", &controller.JobInfoController{}, "*:Delete")
	beego.Router("/jobinfo/activate", &controller.JobInfoController{}, "*:Activate")
	beego.Router("/jobinfo/inactivate", &controller.JobInfoController{}, "*:InActivate")
	beego.Router("/jobinfo/edit", &controller.JobInfoController{}, "get:ToEdit")
	beego.Router("/jobinfo/edit", &controller.JobInfoController{}, "post:Edit")
	beego.Router("/register", &controller.RegisterController{}, "get:ToRegister")
	beego.Router("/register", &controller.RegisterController{}, "post:Register")
	beego.Router("/login", &controller.RegisterController{}, "get:ToLogin")
	beego.Router("/login", &controller.RegisterController{}, "post:Login")
	beego.Router("/logout", &controller.LogoutController{}, "*:Logout")
	beego.Router("/jobsnapshot/list", &controller.JobSnapshotController{}, "*:List")
	beego.Router("/mailbox/active", &controller.MailboxController{}, "*:Active")
	jm := job.NewJobManager()
	jm.PushAllJobInfo()
	beego.Run()
}
