package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/scheduler"
	. "github.com/onsi/gomega"
	"testing"
)

func TestLoadingSavedSchedules(t *testing.T) {
	RegisterTestingT(t)
	SetUpDB()
	defer TearDownDB()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "* * * * *"}
	_ = models.Save(&j)

	sched := scheduler.New()
	_ = sched.Start()
	defer sched.Stop()

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = models.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}

