package hr

import (
	"container/list"
	"fmt"
	"github.com/gorhill/cronexpr"
	"log"
	"sync"
	"time"
)

type Job struct {
	Name           string
	CronExpression string
	JobFunc        func()

	ScheduleTime *time.Time

	LastRunTime *time.Time
}

func NewManager() JobManager {
	manager := JobManager{}
	return manager
}

func NewManagerWithOutput(onStart func(job Job), onEnd func(job Job, err error)) JobManager {
	manager := JobManager{}

	if onStart != nil {
		manager.onStart = onStart
	} else {
		manager.onStart = func(job Job) {
			log.Println(job.Name, " starts at ", time.Now())
		}
	}
	if onEnd != nil {
		manager.onEnd = onEnd
	} else {
		manager.onEnd = func(job Job, err error) {
			msg := fmt.Sprintf("%s ends at %v", job.Name, time.Now())
			if err == nil {
				msg += " successfully."
			} else {
				msg += fmt.Sprintf(", with error %s", err)
			}
			log.Println(msg)
		}
	}
	return manager
}

type JobManager struct {
	running  bool
	lock     sync.Mutex
	stopSign sync.WaitGroup
	jobs     list.List

	onStart func(jobToStart Job)
	onEnd   func(jobToStart Job, err error)
}

func (manager *JobManager) Start() {
	if manager.running {
		log.Println("already started")
		return
	}
	manager.stopSign.Add(1)

	go func() {
		for {
			if manager.jobs.Len() == 0 {
				time.Sleep(5 * time.Second)
				continue
			}

			for current := manager.jobs.Front(); current != nil; current = current.Next() {
				j := current.Value.(Job)
				if j.ScheduleTime.Before(time.Now()) {
					if manager.onStart != nil {
						manager.onStart(j)
					}

					go func() { manager.runJob(j) }()

					now := time.Now()
					j.LastRunTime = &now
					nextTime := cronexpr.MustParse(j.CronExpression).Next(now)
					j.ScheduleTime = &nextTime

					current.Value = j
				}
			}
		}
	}()

	manager.stopSign.Wait()
	log.Println("clock out")
}

// register a new job
func (manager *JobManager) Register(jobName, cronExpression string, jobFunc func()) error {
	exp, validateError := cronexpr.Parse(cronExpression)
	if validateError != nil {
		return validateError
	}

	manager.lock.Lock()
	defer manager.lock.Unlock()
	newJob := Job{}
	newJob.CronExpression = cronExpression
	newJob.Name = jobName
	newJob.JobFunc = jobFunc
	newJob.LastRunTime = nil
	nextTime := exp.Next(time.Now())
	newJob.ScheduleTime = &nextTime

	manager.jobs.PushBack(newJob)
	return nil
}

func (manager *JobManager) runJob(job Job) {
	defer func() {
		manager.stopSign.Done()
		var runErr error
		if r := recover(); r != nil {
			runErr = fmt.Errorf("%v", r)
		}

		if manager.onEnd != nil {
			manager.onEnd(job, runErr)
		}
	}()

	manager.stopSign.Add(1)
	job.JobFunc()
}

func (manager *JobManager) Quit() {
	manager.stopSign.Done()
}
