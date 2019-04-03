package entity

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type JobSnapshot struct {
	Id          int
	JobId       int
	Name        string
	Group       string
	Status      string
	Url         string
	Params      string
	CreateTime  time.Time
	ModifyTime  time.Time
	TimeConsume string
}

func (this *JobSnapshot) SaveJobSnapshot() error {
	_, err := orm.NewOrm().Insert(this)
	return err
}

func (this *JobSnapshot) UpdateSnapshot() error {
	_, err := orm.NewOrm().Update(this, "status", "modify_time", "time_consume")
	return err
}

func (this *JobSnapshot) GetList() ([]*JobSnapshot, error) {
	var jobSnapshots []*JobSnapshot
	o := orm.NewOrm()
	qs := o.QueryTable("job_snapshot")
	if this.Name != "" {
		qs = qs.Filter("name", this.Name)
	}
	if this.Group != "" {
		qs = qs.Filter("group", this.Group)
	}
	_, err := qs.All(&jobSnapshots)
	return jobSnapshots, err
}
