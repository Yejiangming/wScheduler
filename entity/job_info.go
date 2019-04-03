package entity

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type JobInfo struct {
	Id         int
	Name       string
	Group      string
	Cron       string
	Urls       string
	Params     string
	IsActive   int
	CreateTime time.Time
	ModifyTime time.Time
}

func (this *JobInfo) SaveJobInfo() error {
	_, err := orm.NewOrm().Insert(this)
	return err
}

func (this *JobInfo) GetJobInfo() error {
	o := orm.NewOrm()
	return o.Read(this)
}

func (this *JobInfo) GetList() ([]*JobInfo, error) {
	var jobs []*JobInfo
	o := orm.NewOrm()
	qs := o.QueryTable("job_info")
	if this.Name != "" {
		qs = qs.Filter("name", this.Name)
	}
	if this.Group != "" {
		qs = qs.Filter("group", this.Group)
	}
	_, err := qs.All(&jobs)
	return jobs, err
}

func (this *JobInfo) GetActiveList() ([]*JobInfo, error) {
	var jobs []*JobInfo
	qs := orm.NewOrm().QueryTable("job_info").Filter("is_active", 1)
	_, err := qs.All(&jobs)
	return jobs, err
}

func (this *JobInfo) DeleteJobInfo() error {
	o := orm.NewOrm()
	_, err := o.Delete(this)
	return err
}

func (this *JobInfo) UpdateJobInfo() error {
	_, err := orm.NewOrm().Update(this, "name", "group", "cron", "urls", "params", "is_active", "modify_time")
	return err
}
