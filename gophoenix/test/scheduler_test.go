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
	SetUpDB()
	defer TearDownDB()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	_ = models.Save(&j)

	jobs := []models.Job{}
	e := models.AllIndexed("Cron", &jobs)
	assert.Equal(t,nil , e)
	assert.Equal(t, 1, len(jobs))

	sched := scheduler.New()
	_ = sched.Start()
	defer sched.Stop()

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = models.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}

func TestSchedulesWithEmptyCron(t *testing.T) {
	RegisterTestingT(t)
	SetUpDB()
	defer TearDownDB()

	j := models.NewJob()
	_ = models.Save(&j)

	sched := scheduler.New()
	_ = sched.Start()
	defer sched.Stop()

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = models.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(0))
}

func TestAddJob(t *testing.T) {
	RegisterTestingT(t)
	SetUpDB()
	defer TearDownDB()
	sched, _ := scheduler.Start()
	defer sched.Stop()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	_ = models.Save(&j)
	sched.AddJob(j)

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = models.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}
