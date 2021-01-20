package test

import (
	"PhoenixOracle/gophoenix/core/store/models"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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
	app := NewApplication()
	server := app.NewServer()
	defer app.Stop()

	j := NewJobWithSchedule("schedule test")
	err := app.Store.Save(&j)
	assert.Nil(t, err)
	jr := j.NewRun()
	err2 := app.Store.Save(&jr)
	assert.Nil(t, err)
	assert.Nil(t, err2)

	resp, err := BasicAuthGet(server.URL + "/jobs/" + j.ID + "/runs")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	var respJSON JobRunsJSON
	json.Unmarshal(b, &respJSON)
	assert.Equal(t, 1, len(respJSON.Runs), "expected no runs to be created")
	assert.Equal(t, jr.ID, respJSON.Runs[0].ID, "expected the run IDs to match")
}

func TestJobRunsCreate(t *testing.T) {
	t.Parallel()

	app := NewApplication()
	server := app.NewServer()
	defer app.Stop()

	j := NewJobWithWebInitiator()
	assert.Nil(t, app.Store.SaveJob(j))

	url := server.URL + "/jobs/" + j.ID + "/runs"
	resp, err := BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	respJSON := JobJSONFromResponse(resp.Body)

	jr := models.JobRun{}
	assert.Nil(t, app.Store.One("ID", respJSON.ID, &jr))
	assert.Equal(t, jr.ID, respJSON.ID)
}

func TestJobRunsCreateWithoutWebInitiator(t *testing.T) {
	t.Parallel()

	app := NewApplication()
	server := app.NewServer()
	defer app.Stop()

	j := NewJobWithSchedule("* * * * *")
	assert.Nil(t, app.Store.SaveJob(j))

	url := server.URL + "/jobs/" + j.ID + "/runs"
	resp, err := BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Nil(t, err)
	assert.Equal(t, 403, resp.StatusCode, "Response should be forbidden")
}

func TestJobRunsCreateNotFound(t *testing.T) {
	t.Parallel()

	app := NewApplication()
	server := app.NewServer()
	defer app.Stop()

	url := server.URL + "/jobs/garbageID/runs"
	resp, err := BasicAuthPost(url, "application/json", bytes.NewBuffer([]byte{}))
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}