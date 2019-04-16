package controller

import (
	"time"
	"wScheduler/common"
	"wScheduler/entity"
	"wScheduler/job"
)

type JobInfoController struct {
	BaseController
}

func (this *JobInfoController) ToAdd() {
	this.TplName = "jobinfo/add.html"
}

// 最初加入的job都是未激活的
// 直接存入数据库中
// 存入用户名
func (this *JobInfoController) Add() {
	res := new(common.Result)
	defer func() {
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
	}()

	name := this.GetString("Name")
	group := this.GetString("Group")
	cron := this.GetString("Cron")
	urls := this.GetString("Urls")
	params := this.GetString("Params")

	if name == "" || group == "" || cron == "" || urls == "" {
		res.Message = "有未填字段"
		return
	}

	username, ok := this.GetSession("username").(string)
	if !ok {
		res.Message = "获取用户信息失败"
		return
	}

	jobInfo := &entity.JobInfo{
		Name:       name,
		Group:      group,
		Cron:       cron,
		Urls:       urls,
		Params:     params,
		CreateTime: time.Now(),
		UserName:   username,
	}
	jobInfo.ModifyTime = jobInfo.CreateTime

	if err := jobInfo.SaveJobInfo(); err != nil {
		res.Message = "存入数据库失败"
	} else {
		res.Message = "存入数据库成功"
		res.Success = true
	}

}

// 根据name与group从数据库中查询
func (this *JobInfoController) List() {
	name := this.GetString("Name")
	group := this.GetString("Group")
	username, ok := this.GetSession("username").(string)
	if !ok {
		res := new(common.Result)
		res.Message = "获取用户信息失败"
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
		return
	}
	jobInfo := &entity.JobInfo{Name: name, Group: group, UserName: username}
	jobs, err := jobInfo.GetList()
	if err != nil {
		res := new(common.Result)
		res.Message = err.Error()
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
	} else {
		this.Data["name"] = name
		this.Data["group"] = group
		this.Data["jobs"] = jobs
		this.TplName = "jobinfo/list.html"
	}
}

// 先删除内存中job 再删除数据库中job
func (this *JobInfoController) Delete() {
	res := new(common.Result)
	defer func() {
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
	}()

	username, ok := this.GetSession("username").(string)
	if !ok {
		res.Message = "获取用户信息失败"
		return
	}

	id, err := this.GetInt("Id")
	if err != nil {
		res.Message = err.Error()
		return
	}

	jobInfo := &entity.JobInfo{
		Id: id,
	}
	jobInfo.GetJobInfo()
	if !common.CheckUser(jobInfo.UserName, username) {
		res.Message = "权限不足"
		return
	}

	jm := job.NewJobManager()
	if jm.QZ.ExistsJob(id) != -1 {
		err = jm.QZ.RemoveJob(id)
		if err != nil {
			res.Message = "从内存删除失败"
		} else {
			res.Message = "从内存删除成功"
		}
	} else {
		res.Message = "尚未激活"
	}

	err = jobInfo.DeleteJobInfo()
	if err != nil {
		res.Message += " 从数据库删除失败"
	} else {
		res.Message += " 从数据库删除成功"
		res.Success = true
	}
}

// 现将任务调入内存 再更改数据库中任务信息
func (this *JobInfoController) Activate() {
	res := new(common.Result)
	defer func() {
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
	}()

	username, ok := this.GetSession("username").(string)
	if !ok {
		res.Message = "获取用户信息失败"
		return
	}

	id, err := this.GetInt("Id")
	if err != nil {
		res.Message = err.Error()
		return
	}
	jobInfo := &entity.JobInfo{
		Id: id,
	}
	jobInfo.GetJobInfo()

	if !common.CheckUser(jobInfo.UserName, username) {
		res.Message = "权限不足"
		return
	}

	if jobInfo.IsActive == 1 {
		res.Message = "已被激活"
		return
	}

	jm := job.NewJobManager()
	err = jm.AddJobInfoToQuartz(jobInfo)
	if err != nil {
		res.Message = "激活失败"
		return
	}
	res.Message = "激活成功"
	jobInfo.IsActive = 1
	jobInfo.ModifyTime = time.Now()
	err = jobInfo.UpdateJobInfo()
	if err != nil {
		res.Message += " 更新数据库失败"
	} else {
		res.Message += " 更新数据库成功"
		res.Success = true
	}
}

// 先将任务从内存调出 再更改数据库中任务信息
func (this *JobInfoController) InActivate() {
	res := new(common.Result)
	defer func() {
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
	}()

	username, ok := this.GetSession("username").(string)
	if !ok {
		res.Message = "获取用户信息失败"
		return
	}

	id, err := this.GetInt("Id")
	if err != nil {
		res.Message = err.Error()
		return
	}
	jobInfo := &entity.JobInfo{
		Id: id,
	}
	jobInfo.GetJobInfo()

	if !common.CheckUser(jobInfo.UserName, username) {
		res.Message = "权限不足"
		return
	}

	if jobInfo.IsActive == 0 {
		res.Message = "尚未被激活"
		return
	}
	jm := job.NewJobManager()
	err = jm.QZ.RemoveJob(id)
	if err != nil {
		res.Message = "在内存中注销失败"
		return
	}
	res.Message = "在内存中注销成功"
	jobInfo.IsActive = 0
	jobInfo.ModifyTime = time.Now()
	err = jobInfo.UpdateJobInfo()
	if err != nil {
		res.Message += " 更新数据库失败"
	} else {
		res.Message += " 更新数据库成功"
		res.Success = true
	}
}

func (this *JobInfoController) ToEdit() {
	res := new(common.Result)

	username, ok := this.GetSession("username").(string)
	if !ok {
		res.Message = "获取用户信息失败"
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
		return
	}

	id, err := this.GetInt("Id")
	if err != nil {
		res.Message = err.Error()
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
		return
	}

	jobInfo := &entity.JobInfo{
		Id: id,
	}
	jobInfo.GetJobInfo()

	if !common.CheckUser(jobInfo.UserName, username) {
		res.Message = "权限不足"
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
		return
	}

	this.Data["job"] = jobInfo
	this.TplName = "jobinfo/edit.html"
}

// 先更新数据库 再根据是否已经激活更改内存中的jobPool
func (this *JobInfoController) Edit() {
	res := new(common.Result)
	defer func() {
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
	}()

	username, ok := this.GetSession("username").(string)
	if !ok {
		res.Message = "获取用户信息失败"
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
		return
	}

	id, err := this.GetInt("Id")
	if err != nil {
		res.Message = err.Error()
		return
	}

	tmpJobInfo := &entity.JobInfo{
		Id: id,
	}
	tmpJobInfo.GetJobInfo()
	if !common.CheckUser(tmpJobInfo.UserName, username) {
		res.Message = "权限不足"
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
		return
	}

	name := this.GetString("Name")
	group := this.GetString("Group")
	corn := this.GetString("Cron")
	urls := this.GetString("Urls")
	params := this.GetString("Params")
	isActive, err := this.GetInt("IsActive")
	if err != nil {
		res.Message = err.Error()
		return
	}
	if name == "" || group == "" || corn == "" || urls == "" {
		res.Message = "有未填字段"
		return
	}
	jobInfo := &entity.JobInfo{
		Id:       id,
		Name:     name,
		Group:    group,
		Cron:     corn,
		Urls:     urls,
		Params:   params,
		IsActive: isActive,
	}
	jobInfo.ModifyTime = time.Now()
	if err := jobInfo.UpdateJobInfo(); err != nil {
		res.Message = "数据库保存失败"
		return
	}
	res.Message = "数据库保存成功"
	jm := job.NewJobManager()
	// 内存中不存在
	if jm.QZ.ExistsJob(id) == -1 {
		res.Message += " 尚未激活"
		res.Success = true
		return
	}
	// 内存中存在
	jm.QZ.RemoveJob(id)
	err = jm.AddJobInfoToQuartz(jobInfo)
	if err != nil {
		res.Message += " 重新激活失败"
	} else {
		res.Message += " 重新激活成功"
		res.Success = true
	}
}
