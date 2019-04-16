package controller

import (
	"wScheduler/common"
	"wScheduler/entity"
)

type JobSnapshotController struct {
	BaseController
}

func (this *JobSnapshotController) List() {
	name := this.GetString("Name")
	group := this.GetString("Group")
	username, ok := this.GetSession("username").(string)
	if !ok {
		res := &common.Result{
			Message: "获取用户信息失败",
		}
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
		return
	}
	jobSnapshot := entity.JobSnapshot{Name: name, Group: group, UserName: username}
	jobSnapshots, err := jobSnapshot.GetList()
	common.PanicIf(err)
	this.Data["jobsnapshots"] = jobSnapshots
	this.Data["name"] = name
	this.Data["group"] = group
	this.TplName = "jobsnapshot/list.html"
	this.Render()
}
