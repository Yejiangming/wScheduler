package quartz

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var (
	Error_EXISTS_JOB     = errors.New("exists job error !")
	Error_NOT_EXISTS_JOB = errors.New("not exists job error !")
	Error_WRONG_CRON     = errors.New("wrong cron")
)

// defind the job call method
type JobFunc func(jobId int)

type Schedule interface {
	Next(time.Time) time.Time
}

// job  多余字段为了将来进行扩展
// 字段暴露可以直接修改
type Job struct {
	// id
	Id int
	// 任务名称
	Name string
	// 分组
	Group string
	// ip
	Ip string
	//目标服务器地址
	Urls string
	// 上次触发时间
	Pre time.Time
	// 下次执行时间
	Next time.Time
	//schedule
	schedule Schedule
	// 执行的job 方法
	JobFunc JobFunc
	// corn 表达式
	Expression string
	// 状态 0 init  1 invoking 2 COMPLETED 3
	State  uint
	Params string
	//执行策略
	InvokePolicy string
	//是否激活  1 激活 0 非激活
	IsActive int
}

type Quartz struct {
	jobPool   []*Job
	lk        sync.Mutex
	running   bool
	addJob    chan *Job
	removeJob chan int
	snapshot  chan []*Job
	modifyJob chan *Job
	Stop      chan int
	jobError  chan error
}

func New() *Quartz {
	return &Quartz{
		Stop:      make(chan int),
		addJob:    make(chan *Job),
		removeJob: make(chan int),
		modifyJob: make(chan *Job),
		snapshot:  make(chan []*Job),
		jobError:  make(chan error),
	}
}

func (qz *Quartz) AddJob(job *Job) error {
	qz.lk.Lock()
	defer qz.lk.Unlock()
	if qz.ExistsJob(job.Id) == -1 {
		//log.Println(job.Expression)
		schedule, err := Parse(job.Expression)
		if err != nil {
			return err
		}
		job.schedule = schedule
		if qz.running == false {
			qz.jobPool = append(qz.jobPool, job)
			return nil
		} else {
			qz.addJob <- job
			return <-qz.jobError
		}
	} else {
		return Error_EXISTS_JOB
	}
}

func (qz *Quartz) RemoveJob(jobId int) error {
	qz.lk.Lock()
	defer qz.lk.Unlock()
	index := qz.ExistsJob(jobId)
	if index == -1 {
		return Error_NOT_EXISTS_JOB
	}
	if !qz.running {
		qz.jobPool = qz.jobPool[:index+copy(qz.jobPool[index:], qz.jobPool[index+1:])]
		return nil
	} else {
		qz.removeJob <- jobId
		return <-qz.jobError
	}
}

func (qz *Quartz) SnapshotJob() ([]*Job, error) {
	if !qz.running {
		return qz.jobPool, nil
	} else {
		qz.snapshot <- nil
		list := <-qz.snapshot
		return list, nil
	}
}

func (qz *Quartz) ModifyJob(job *Job) error {
	qz.lk.Lock()
	defer qz.lk.Unlock()
	index := qz.ExistsJob(job.Id)
	if index == -1 {
		return Error_NOT_EXISTS_JOB
	}
	schedule, err := Parse(job.Expression)
	if err != nil {
		return err
	}
	job.schedule = schedule
	if !qz.running {
		qz.jobPool[index] = job
		return nil
	} else {
		qz.modifyJob <- job
		return <-qz.jobError
	}
}

func (qz *Quartz) BootStrap() {
	go qz.run()
}

func (qz *Quartz) StopJob() {
	qz.lk.Lock()
	defer qz.lk.Unlock()
	qz.Stop <- 1
}

func (qz *Quartz) ExistsJob(jobId int) int {
	for i, v := range qz.jobPool {
		if v.Id == jobId {
			return i
		}
	}
	return -1
}

func (qz *Quartz) run() {
	qz.running = true
	now := time.Now()
	for _, v := range qz.jobPool {
		v.Next = v.schedule.Next(now)
	}

	for {
		sort.Sort(JobPool(qz.jobPool))
		var effective time.Time
		if len(qz.jobPool) == 0 || qz.jobPool[0].Next.IsZero() {
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = qz.jobPool[0].Next
		}

		select {
		case now = <-time.After(effective.Sub(now)):
			for _, v := range qz.jobPool {
				if v.Next != effective {
					break
				}
				v.Pre = effective
				v.Next = v.schedule.Next(effective)
				go v.JobFunc(v.Id)
			}
		case newJob := <-qz.addJob:
			if qz.ExistsJob(newJob.Id) == -1 {
				now = time.Now()
				qz.jobPool = append(qz.jobPool, newJob)
				newJob.Next = newJob.schedule.Next(now)
				qz.jobError <- nil
			} else {
				qz.jobError <- Error_EXISTS_JOB
			}
		case removeJobId := <-qz.removeJob:
			index := qz.ExistsJob(removeJobId)
			if index != -1 {
				qz.jobPool = qz.jobPool[:index+copy(qz.jobPool[index:], qz.jobPool[index+1:])]
				qz.jobError <- nil
			} else {
				qz.jobError <- Error_NOT_EXISTS_JOB
			}
		case <-qz.snapshot:
			qz.snapshot <- qz.jobPool
		case modifyJob := <-qz.modifyJob:
			index := qz.ExistsJob(modifyJob.Id)
			if index == -1 {
				qz.jobError <- Error_NOT_EXISTS_JOB
			} else {
				now = time.Now()
				nextTime := modifyJob.schedule.Next(now)
				modifyJob.Next = nextTime
				qz.jobPool[index] = modifyJob
				qz.jobError <- nil
			}
		case <-qz.Stop:
			qz.running = false
			return
		}
	}
}

type JobPool []*Job

func (jp JobPool) Len() int {
	return len(jp)
}

func (jp JobPool) Less(i, j int) bool {
	if jp[i].Next.IsZero() {
		return false
	}
	if jp[j].Next.IsZero() {
		return true
	}
	return jp[i].Next.Before(jp[j].Next)
}

func (jp JobPool) Swap(i, j int) {
	temp := jp[i]
	jp[i] = jp[j]
	jp[j] = temp
}
