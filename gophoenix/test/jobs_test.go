package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	SetUpDB()
	defer TearDownDB()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}
	models.Save(&j1)

	var j2 models.Job
	models.Find("ID",j1.ID,&j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}

func TestFind(t *testing.T) {
	SetUpDB()
	defer TearDownDB()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}
	models.Save(&j1)

	var j2 models.Job
	models.Find("ID",j1.ID,&j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}

func TestWhere(t *testing.T) {
	SetUpDB()
	defer TearDownDB()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}
	models.Save(&j1)

	jobs := []models.Job{j1}
	models.Where("ID",j1.ID,&jobs)

	assert.Equal(t, len(jobs), 1)
	assert.Equal(t, jobs[0].Schedule, j1.Schedule)

}

func TestAllIndexedNotFound(t *testing.T) {
	SetUpDB()
	defer TearDownDB()

	j1 := models.NewJob()
	jobs := []models.Job{j1}

	err := models.AllIndexed("Cron", &jobs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestWhereNotFound(t *testing.T) {
	SetUpDB()
	defer TearDownDB()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}
	models.Save(&j1)

	jobs := []models.Job{j1}
	models.Where("ID","1",&jobs)

	assert.Equal(t, len(jobs), 0,"id not found")
}

func TestAll(t *testing.T) {
	SetUpDB()
	defer TearDownDB()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}

	j2 := models.NewJob()
	j2.Schedule = models.Schedule{Cron: "2 * * * *"}
	models.Save(&j1)
	models.Save(&j2)

	jobs := []models.Job{}
	models.All(&jobs)

	assert.Equal(t, 2, len(jobs))
}