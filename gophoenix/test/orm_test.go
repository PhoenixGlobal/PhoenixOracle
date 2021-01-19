package test

import (
	"PhoenixOracle/gophoenix/core/store/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWhereNotFound(t *testing.T) {
	t.Parallel()
	store := NewStore()
	defer store.Close()

	j1 := models.NewJob()
	jobs := []models.Job{j1}

	err := store.Where("ID", "bogus", &jobs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestAllNotFound(t *testing.T) {
	t.Parallel()
	store := NewStore()
	defer store.Close()

	var jobs []models.Job
	err := store.All(&jobs)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(jobs), "Queried array should be empty")
}

func TestORMSaveJob(t *testing.T) {
	t.Parallel()
	store := NewStore()
	defer store.Close()

	j1 := NewJobWithSchedule("* * * * *")
	store.SaveJob(j1)

	var j2 models.Job
	store.One("ID", j1.ID, &j2)
	assert.Equal(t, j1.ID, j2.ID)

	var initr models.Initiator
	store.One("JobID", j1.ID, &initr)
	assert.Equal(t, models.Cron("* * * * *"), initr.Schedule)
}

