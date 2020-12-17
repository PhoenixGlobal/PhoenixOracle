package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	db := SetUpDB()
	defer TearDownDB()
	j1 := models.NewJob()
	j1.Schedule = models.Schedule{Cron: "1 * * * *"}
	db.Save(&j1)

	var j2 models.Job
	db.One("ID",j1.ID,&j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}