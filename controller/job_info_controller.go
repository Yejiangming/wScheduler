package controller

import (
	"fmt"
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
func (this *JobInfoController) Add() {
	name := this.GetString("Name")
	group := this.GetString("Group")
	corn := this.GetString("Cron")
	urls := this.GetString("Urls")
	params := this.GetString("Params")

	res := new(common.Result)

	if name == "" || group == "" || corn == "" || urls == "" {
		res.Message = " 有未填字段"
	} else {
		jobInfo := new(entity.JobInfo)
		jobInfo.Name = name
		jobInfo.Group = group
		jobInfo.Cron = corn
		jobInfo.Urls = urls
		jobInfo.Params = params
		jobInfo.CreateTime = time.Now()
		jobInfo.ModifyTime = jobInfo.CreateTime
		if err := jobInfo.SaveJobInfo(); err != nil {
			res.Message = " 存入数据库失败"
		} else {
			res.Message = " 存入数据库成功"
			res.Success = true
		}
	}
	this.Data["res"] = res
	this.TplName = "jobinfo/res.html"
}

// 根据name与group从数据库中查询
func (this *JobInfoController) List() {
	name := this.GetString("Name")
	group := this.GetString("Group")
	jobInfo := &entity.JobInfo{Name: name, Group: group}
	jobs, err := jobInfo.GetList()
	common.PanicIf(err)
	this.Data["name"] = name
	this.Data["group"] = group
	this.Data["jobs"] = jobs
	this.TplName = "jobinfo/list.html"
}

// 先删除内存中job 再删除数据库中job
func (this *JobInfoController) Delete() {
	res := new(common.Result)
	id, err := this.GetInt("Id")
	common.PanicIf(err)
	jm := job.NewJobManager()
	if jm.QZ.ExistsJob(id) != -1 {
		err = jm.QZ.RemoveJob(id)
		if err != nil {
			res.Message = " 从内存删除失败"
		} else {
			res.Message = " 从内存删除成功"
		}
	} else {
		res.Message = " 尚未激活"
	}
	jobInfo := new(entity.JobInfo)
	jobInfo.Id = id
	err = jobInfo.DeleteJobInfo()
	if err != nil {
		res.Message += " 从数据库删除失败"
	} else {
		res.Message += " 从数据库删除成功"
		res.Success = true
	}
	this.Data["res"] = res
	this.TplName = "jobinfo/res.html"
}

// 现将任务调入内存 再更改数据库中任务信息
func (this *JobInfoController) Activate() {
	res := new(common.Result)
	id, err := this.GetInt("Id")
	common.PanicIf(err)
	jobInfo := new(entity.JobInfo)
	jobInfo.Id = id
	jobInfo.GetJobInfo()
	if jobInfo.IsActive == 1 {
		res.Message = " 已被激活"
	} else {
		jm := job.NewJobManager()
		err := jm.AddJobInfoToQuartz(jobInfo)
		if err != nil {
			res.Message = " 激活失败"
		} else {
			res.Message = " 激活成功"
			jobInfo.IsActive = 1
			jobInfo.ModifyTime = time.Now()
			err := jobInfo.UpdateJobInfo()
			if err != nil {
				res.Message += " 更新数据库失败"
			} else {
				res.Message += " 更新数据库成功"
				res.Success = true
			}
		}
	}
	this.Data["res"] = res
	this.TplName = "jobinfo/res.html"
}

// 先将任务从内存调出 再更改数据库中任务信息
func (this *JobInfoController) InActivate() {
	res := new(common.Result)
	id, err := this.GetInt("Id")
	common.PanicIf(err)
	jobInfo := new(entity.JobInfo)
	jobInfo.Id = id
	jobInfo.GetJobInfo()
	if jobInfo.IsActive == 0 {
		res.Message = " 尚未被激活"
	} else {
		jm := job.NewJobManager()
		err := jm.QZ.RemoveJob(id)
		if err != nil {
			res.Message = " 在内存中注销失败"
		} else {
			res.Message = " 在内存中注销成功"
			jobInfo.IsActive = 0
			jobInfo.ModifyTime = time.Now()
			err := jobInfo.UpdateJobInfo()
			if err != nil {
				res.Message += " 更新数据库失败"
			} else {
				res.Message += " 更新数据库成功"
				res.Success = true
			}
		}
	}
	this.Data["res"] = res
	this.TplName = "jobinfo/res.html"
}

func (this *JobInfoController) ToEdit() {
	id, err := this.GetInt("Id")
	common.PanicIf(err)
	jobInfo := new(entity.JobInfo)
	jobInfo.Id = id
	err = jobInfo.GetJobInfo()
	common.PanicIf(err)
	this.Data["job"] = jobInfo
	this.TplName = "jobinfo/edit.html"
}

// 先更新数据库 再根据是否已经激活更改内存中的jobPool
func (this *JobInfoController) Edit() {
	id, err := this.GetInt("Id")
	common.PanicIf(err)
	name := this.GetString("Name")
	group := this.GetString("Group")
	corn := this.GetString("Cron")
	urls := this.GetString("Urls")
	params := this.GetString("Params")
	isActive, err := this.GetInt("IsActive")
	common.PanicIf(err)

	res := new(common.Result)

	if name == "" || group == "" || corn == "" || urls == "" {
		res.Message = " 有未填字段"
	} else {
		jobInfo := new(entity.JobInfo)
		jobInfo.Id = id
		jobInfo.Name = name
		jobInfo.Group = group
		jobInfo.Cron = corn
		jobInfo.Urls = urls
		jobInfo.Params = params
		jobInfo.IsActive = isActive
		jobInfo.ModifyTime = time.Now()
		if err := jobInfo.UpdateJobInfo(); err != nil {
			res.Message = " 数据库保存失败"
		} else {
			res.Message = " 数据库保存成功"
			jm := job.NewJobManager()
			// 内存中存在
			if jm.QZ.ExistsJob(id) != -1 {
				jm.QZ.RemoveJob(id)
				err := jm.AddJobInfoToQuartz(jobInfo)
				if err != nil {
					res.Message += " 重新激活失败"
				} else {
					res.Message += " 重新激活成功"
					res.Success = true
				}
			} else {
				res.Message += " 尚未激活"
				res.Success = true
			}
		}
	}
	fmt.Println(res)
	this.Data["res"] = res
	this.TplName = "jobinfo/res.html"
}
