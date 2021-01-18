package test

import (
	"PhoenixOracle/gophoenix/core/models"
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

func TestAllIndexedNotFound(t *testing.T) {
	t.Parallel()
	store := NewStore()
	defer store.Close()

	j1 := models.NewJob()
	jobs := []models.Job{j1}

	err := store.AllByIndex("Cron", &jobs)
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

