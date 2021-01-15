package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetrievingJobRunsWithErrorsFromDB(t *testing.T) {
	store := NewStore()
	defer store.Close()

	job := models.NewJob()
	jr := job.NewRun()
	jr.Result = models.RunResultWithError(fmt.Errorf("bad idea"))
	err := store.Save(&jr)
	assert.Nil(t, err)

	run := models.JobRun{}
	err = store.One("ID", jr.ID, &run)
	assert.Nil(t, err)
	assert.True(t, run.Result.HasError())
	assert.Equal(t, "bad idea", run.Result.Error())
}
