package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/scheduler"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadingSavedSchedules(t *testing.T) {
	RegisterTestingT(t)
	store := Store()
	defer store.Close()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	_ = store.Save(&j)

	jobs := []models.Job{}
	e := store.AllByIndex("Cron", &jobs)
	assert.Equal(t,nil , e)
	assert.Equal(t, 1, len(jobs))

	sched := scheduler.New(store.ORM)
	_ = sched.Start()
	defer sched.Stop()

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}

func TestSchedulesWithEmptyCron(t *testing.T) {
	RegisterTestingT(t)
	store := Store()
	defer store.Close()

	j := models.NewJob()
	_ = store.Save(&j)

	sched := scheduler.New(store.ORM)
	_ = sched.Start()
	defer sched.Stop()

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(0))
}

func TestAddJob(t *testing.T) {
	RegisterTestingT(t)
	store := Store()
	defer store.Close()

	sched := scheduler.New(store.ORM)
	sched.Start()
	defer sched.Stop()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	//_ = store.Save(&j)
	sched.AddJob(j)

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}
