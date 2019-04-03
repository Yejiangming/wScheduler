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
	jobSnapshot := entity.JobSnapshot{Name: name, Group: group}
	jobSnapshots, err := jobSnapshot.GetList()
	common.PanicIf(err)
	this.Data["jobsnapshots"] = jobSnapshots
	this.Data["name"] = name
	this.Data["group"] = group
	this.TplName = "jobsnapshot/list.html"
	this.Render()
}
