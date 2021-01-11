package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type JobRunsJSON struct {
	Runs []JobRun `json:"runs"`
}

type JobRun struct {
	ID string `json:"id"`
}

func TestJobRunsIndex(t *testing.T) {
	t.Parallel()
	store := Store()
	defer store.Close()
	server := store.SetUpWeb()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "schedule test"}
	err := store.Save(&j)
	assert.Nil(t, err)
	jr := j.NewRun()
	err2 := store.Save(&jr)
	assert.Nil(t, err)
	assert.Nil(t, err2)

	resp, err := http.Get(server.URL + "/jobs/" + j.ID + "/runs")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	var respJSON JobRunsJSON
	json.Unmarshal(b, &respJSON)
	assert.Equal(t, 1, len(respJSON.Runs), "expected no runs to be created")
	assert.Equal(t, jr.ID, respJSON.Runs[0].ID, "expected the run IDs to match")
}
