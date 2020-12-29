package scheduler

import (
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/services"
	"fmt"
	cronlib "github.com/robfig/cron"
)

type Scheduler struct {
	cron *cronlib.Cron
	orm  models.ORM
}


func NewScheduler(orm models.ORM) *Scheduler {
	return &Scheduler{cronlib.New(),orm}
}

func (self *Scheduler) Start() error {
	jobs,err := self.orm.JobsWithCron()
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
	self.cron.AddFunc(cronStr, func() {
		err := services.StartJob(job.NewRun(), self.orm)
		if err != nil{
			logger.GetLogger().Panic(err.Error())
		}
	})
}


func (self *Scheduler) Stop() {
	self.cron.Stop()
}

