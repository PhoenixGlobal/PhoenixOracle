package services

import (
	configlib "PhoenixOracle/gophoenix/core/config"
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/models"
	"errors"
	"fmt"
	cronlib "github.com/mrwonko/cron"
)

type Scheduler struct {
	cron *cronlib.Cron
	orm  *models.ORM
	config  configlib.Config
	started bool
}


func NewScheduler(orm *models.ORM, cf configlib.Config) *Scheduler {
	return &Scheduler{orm: orm,config: cf}
}

func (self *Scheduler) Start() error {
	if self.started {
		return errors.New("Scheduler already started")
	}
	self.started = true
	self.cron = cronlib.New()
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
	if !self.started {
		return
	}
	cronStr := string(job.Schedule.Cron)
	self.cron.AddFunc(cronStr, func() {
		err := StartJob(job.NewRun(), self.orm, self.config)
		if err != nil{
			logger.GetLogger().Panic(err.Error())
		}
	})
}


func (self *Scheduler) Stop() {
	if self.started{
		self.cron.Stop()
		self.cron.Wait()
		self.started = false
	}
}

