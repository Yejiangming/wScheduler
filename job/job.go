package job

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"wScheduler/common"
	"wScheduler/entity"
	"wScheduler/quartz"
)

const (
	EXECUTING     = "EXECUTING"
	CONNECT_ERROR = "CONNECT_ERROR"
	SUCCESS       = "SUCCESS"
	FAIL          = "FAIL"
)

type JobManager struct {
	QZ *quartz.Quartz
}

var JM *JobManager

func NewJobManager() *JobManager {
	if JM == nil {
		qz := quartz.New()
		qz.BootStrap()
		JM = &JobManager{QZ: qz}
	}
	return JM
}

func (this *JobManager) PushAllJobInfo() {
	jobInfo := new(entity.JobInfo)
	jobs, err := jobInfo.GetActiveList()
	common.PanicIf(err)
	if len(jobs) > 0 {
		for _, job := range jobs {
			j := new(quartz.Job)
			j.Id = job.Id
			j.Name = job.Name
			j.Group = job.Group
			j.Expression = job.Cron
			j.Urls = job.Urls
			j.Params = job.Params
			j.IsActive = job.IsActive
			j.JobFunc = this.Execute
			err := this.QZ.AddJob(j)
			common.PanicIf(err)
		}
	}
}

// 远程被执行函数返回common.res类型的结果
func (this *JobManager) Execute(JobId int) {
	jobInfo := &entity.JobInfo{
		Id: JobId,
	}
	jobInfo.GetJobInfo()

	// 参数以json格式传送
	var parameter map[string]interface{}
	json.Unmarshal([]byte(jobInfo.Params), &parameter)
	content, err := json.Marshal(parameter)
	if err != nil {
		log.Println(err)
		return
	}

	res := new(common.Result)

	urls := strings.Split(jobInfo.Urls, ";")
	for _, url := range urls {
		now := time.Now()
		jobsnap := &entity.JobSnapshot{
			JobId:      jobInfo.Id,
			Name:       jobInfo.Name,
			Group:      jobInfo.Group,
			Status:     EXECUTING,
			Url:        url,
			Params:     jobInfo.Params,
			CreateTime: now,
			ModifyTime: now,
			UserName:   jobInfo.UserName,
		}
		jobsnap.SaveJobSnapshot()
		resp, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(content))
		if err != nil {
			jobsnap.Status = err.Error()
			jobsnap.ModifyTime = time.Now()
			jobsnap.TimeConsume = jobsnap.ModifyTime.Sub(jobsnap.CreateTime).String()
			jobsnap.UpdateSnapshot()
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, res)

		jobsnap.Status = res.Message
		jobsnap.ModifyTime = time.Now()
		jobsnap.TimeConsume = jobsnap.ModifyTime.Sub(jobsnap.CreateTime).String()
		jobsnap.UpdateSnapshot()
	}
}

func (this *JobManager) AddJobInfoToQuartz(job *entity.JobInfo) error {
	j := new(quartz.Job)
	j.Id = job.Id
	j.Name = job.Name
	j.Group = job.Group
	j.Expression = job.Cron
	j.Urls = job.Urls
	j.Params = job.Params
	j.JobFunc = this.Execute
	err := this.QZ.AddJob(j)
	return err
}
