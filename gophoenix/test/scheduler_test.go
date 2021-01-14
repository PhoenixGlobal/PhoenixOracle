package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/services"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadingSavedSchedules(t *testing.T) {
	t.Parallel()
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

	sched := services.NewScheduler(store.ORM,store.Config)
	_ = sched.Start()


	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))

	sched.Stop()
}

func TestSchedulesWithEmptyCron(t *testing.T) {
	RegisterTestingT(t)
	store := Store()
	defer store.Close()

	j := models.NewJob()
	_ = store.Save(&j)

	sched := services.NewScheduler(store.ORM,store.Config)
	_ = sched.Start()
	defer sched.Stop()

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(0))
}

func TestAddJob(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store := Store()
	defer store.Close()

	sched := services.NewScheduler(store.ORM,store.Config)
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

func TestAddJobWhenStopped(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store := Store()
	defer store.Close()

	sched := services.NewScheduler(store.ORM,store.Config)

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	_ = store.Save(&j)
	sched.AddJob(j)

	jobRuns := []models.JobRun{}
	Consistently(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(0))

	sched.Start()
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}
