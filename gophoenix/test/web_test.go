package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/test/fixtureprepare"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)


type TaskJSON struct {
	ID string `json:"id"`
}

func TestCreateTasks(t *testing.T) {
	db := fixtureprepare.SetUpDB()
	defer fixtureprepare.TearDownDB()
	server := fixtureprepare.SetUpWeb()
	defer fixtureprepare.TearDownWeb()

	//jsonStr := []byte(`{"version": "1.0.0"}`)
	jsonStr := []byte(`{"subtasks":[{"adapterType": "httpJSON", "adapterParams": {"endpoint": "https://bitstamp.net/api/ticker/", "fields": ["last"]}}], "schedule": "* * * * *","version":"1.0.0"}`)
	resp, err := http.Post(server.URL+"/tasks", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode, "Response should be success")

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)


	var respJSON TaskJSON
	json.Unmarshal(b, &respJSON)


	var j models.Task
	db.One("ID", respJSON.ID, &j)
	assert.Equal(t, j.ID, respJSON.ID, "Wrong job returned")
	assert.Equal(t, j.Schedule, "* * * * *", "Wrong schedule saved")
}

type JobJSON struct {
	ID string `json:"id"`
}

func TestCreateInvalidTasks(t *testing.T) {
	server := fixtureprepare.SetUpWeb()
	defer fixtureprepare.TearDownWeb()

	jsonStr := []byte(`{"subtasks":[{"adapterType": "ethereumBytes32", "adapterParams": {}}], "schedule": "* * * * *","version":"1.0.0"}`)
	resp, err := http.Post(server.URL+"/tasks", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["\"ethereumBytes32\" is not a supported adapter type."]}`, string(body), "Repsonse should return JSON")
}


