package services

import (
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/store"
	"PhoenixOracle/gophoenix/core/store/models"
	"errors"
	"fmt"
	cronlib "github.com/mrwonko/cron"
	"time"
)

type Scheduler struct {
	Recurring *Recurring
	OneTime   *OneTime
	store   *store.Store
	started bool
}


func NewScheduler(store *store.Store) *Scheduler {
	return &Scheduler{
		Recurring: &Recurring{store: store},
		OneTime:   &OneTime{store: store},
		store:     store,
	}
}

func (self *Scheduler) Start() error {
	if self.started {
		return errors.New("Scheduler already started")
	}
	if err := self.Recurring.Start(); err != nil {
		return err
	}
	self.started = true
	jobs,err := self.store.Jobs()
	if err != nil {
		return fmt.Errorf("Scheduler: ", err)
	}
	for _, j := range jobs {
		self.AddJob(j)
	}
	return nil
}

func (self *Recurring) AddJob(job models.Job) {
	for _, initr := range job.InitiatorsFor("cron"){
		cronStr := string(initr.Schedule)
		self.cron.AddFunc(cronStr, func() {
			_, err := StartJob(job.NewRun(), self.store)
			if err != nil{
				logger.GetLogger().Panic(err.Error())
			}
		})
	}
}

func (self *Scheduler) AddJob(job models.Job) {
	if !self.started {
		return
	}
	self.Recurring.AddJob(job)
	self.OneTime.AddJob(job)
}

type Recurring struct {
	cron  *cronlib.Cron
	store *store.Store
}

func (self *Recurring) Start() error {
	self.cron = cronlib.New()
	self.addResumer()
	self.cron.Start()
	return nil
}

func (self *Recurring) Stop() {
	self.cron.Stop()
	self.cron.Wait()
}


func (self *Scheduler) Stop() {
	if self.started{
		self.Recurring.Stop()
		self.started = false
	}
}

func (self *Recurring) addResumer() {
	self.cron.AddFunc(self.store.Config.PollingSchedule, func() {
		pendingRuns, err := self.store.PendingJobRuns()
		if err != nil {
			logger.Panic(err.Error())
		}
		for _, jobRun := range pendingRuns {
			_, err := StartJob(jobRun, self.store)
			if err != nil {
				logger.Panic(err.Error())
			}
		}
	})
}

type Sleeper interface {
	Sleep(d time.Duration)
}

type Clock struct{}

func (self *Clock) Sleep(d time.Duration) {
	time.Sleep(d)
}

type OneTime struct {
	store *store.Store
	Clock Sleeper
}

func (self *OneTime) AddJob(job models.Job) {
	for _, initr := range job.InitiatorsFor("runAt") {
		go func() {
			self.Clock.Sleep(initr.Time.DurationFromNow())
			_, err := StartJob(job.NewRun(), self.store)
			if err != nil {
				logger.Panic(err.Error())
			}
		}()
	}
}

