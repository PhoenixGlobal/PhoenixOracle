package scheduler

import (
	"PhoenixOracle/gophoenix/core/models"
	"fmt"
	cronlib "github.com/robfig/cron"
)

type Scheduler struct {
	cron *cronlib.Cron
}

func Start() (*Scheduler, error) {
	sched := New()
	err := sched.Start()
	return sched, err
}

func New() *Scheduler {
	return &Scheduler{cronlib.New()}
}

func (self *Scheduler) Start() error {
	jobs,_ := models.JobsWithCron()
	err := models.AllIndexed("Cron",&jobs)
	if err != nil {
		return fmt.Errorf("Scheduler: ", err)
	}

	for _, j := range jobs {
		self.AddJob(j)
	}

	self.cron.Start()
	return nil
}

func (self *Scheduler) AddJob(job models.Job) {
	cronStr := string(job.Schedule.Cron)
	self.cron.AddFunc(cronStr, func() { job.Run() })
}


func (self *Scheduler) Stop() {
	self.cron.Stop()
}

