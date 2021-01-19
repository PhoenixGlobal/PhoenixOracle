package test

import (
	"PhoenixOracle/gophoenix/core/services"
	"PhoenixOracle/gophoenix/core/store/models"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadingSavedSchedules(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store := NewStore()
	defer store.Close()

	j := models.NewJob()
	j.Initiators = []models.Initiator{{Type: "cron", Schedule: "* * * * *"}}
	jobWoCron := models.NewJob()
	assert.Nil(t, store.SaveJob(j))
	assert.Nil(t, store.SaveJob(jobWoCron))

	sched := services.NewScheduler(store)
	_ = sched.Start()


	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))

	sched.Stop()
}

func TestSchedulesWithEmptyCron(t *testing.T) {
	RegisterTestingT(t)
	store := NewStore()
	defer store.Close()

	j := models.NewJob()
	_ = store.Save(&j)

	sched := services.NewScheduler(store)
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
	store := NewStore()
	defer store.Close()

	sched := services.NewScheduler(store)
	sched.Start()
	defer sched.Stop()

	j := NewJobWithSchedule("* * * * *")
	err := store.SaveJob(j)
	assert.Nil(t, err)
	sched.AddJob(j)

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		err = store.Where("JobID", j.ID, &jobRuns)
		assert.Nil(t, err)
		return jobRuns
	}).Should(HaveLen(1))
}

func TestAddJobWhenStopped(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store := NewStore()
	defer store.Close()

	sched := services.NewScheduler(store)

	j := NewJobWithSchedule("* * * * *")
	assert.Nil(t, store.SaveJob(j))
	sched.AddJob(j)

	jobRuns := []models.JobRun{}
	Consistently(func() []models.JobRun {
		store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(0))

	assert.Nil(t, sched.Start())
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", j.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))
}
