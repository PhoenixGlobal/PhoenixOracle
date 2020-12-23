package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	store := Store()
	defer store.Close()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}
	store.Save(&j1)

	var j2 models.Job
	store.One("ID",j1.ID,&j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}

func TestFind(t *testing.T) {
	store := Store()
	defer store.Close()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}
	store.Save(&j1)

	var j2 models.Job
	store.One("ID",j1.ID,&j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}

func TestWhere(t *testing.T) {
	store := Store()
	defer store.Close()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}
	store.Save(&j1)

	jobs := []models.Job{j1}
	store.Where("ID",j1.ID,&jobs)

	assert.Equal(t, len(jobs), 1)
	assert.Equal(t, jobs[0].Schedule, j1.Schedule)

}

func TestAllIndexed(t *testing.T) {
	store := Store()
	defer store.Close()

	j1 := models.NewJob()
	jobs := []models.Job{j1}

	err := store.AllByIndex("Cron", &jobs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")

	j2 := models.NewJob()
	j2.Schedule = models.Schedule{Cron: "2 * * * *"}
	store.Save(&j2)

	err2 := store.AllByIndex("Cron", &jobs)
	assert.Nil(t, err2)
	assert.Equal(t, 1, len(jobs), "Queried array should be 1")
}

func TestWhereNotFound(t *testing.T) {
	store := Store()
	defer store.Close()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}
	store.Save(&j1)

	jobs := []models.Job{j1}
	store.Where("ID","1",&jobs)

	assert.Equal(t, len(jobs), 0,"id not found")
}

func TestAll(t *testing.T) {
	store := Store()
	defer store.Close()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}

	j2 := models.NewJob()
	j2.Schedule = models.Schedule{Cron: "2 * * * *"}
	store.Save(&j1)
	store.Save(&j2)

	jobs := []models.Job{}
	store.All(&jobs)

	assert.Equal(t, 2, len(jobs))
}