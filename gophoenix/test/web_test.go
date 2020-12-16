package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/models/tasks"
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
	db := SetUpDB()
	defer TearDownDB()
	server := SetUpWeb()
	defer TearDownWeb()

	jsonStr := LoadJSON("./fixture/create_jobs.json")
	resp, err := http.Post(server.URL+"/tasks", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, resp.StatusCode, "Response should be success")

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)


	var respJSON TaskJSON
	json.Unmarshal(b, &respJSON)


	var j models.Job
	db.One("ID", respJSON.ID, &j)
	assert.Equal(t, j.ID, respJSON.ID, "Wrong job returned")
	assert.Equal(t, j.Schedule, "* * * * *", "Wrong schedule saved")
	assert.Equal(t, j.Tasks[0].Type, "HttpGet")

	httpGet := j.Tasks[0].Adapter.(*tasks.HttpGet)
	assert.Nil(t, err)
	assert.Equal(t, httpGet.Endpoint, "https://bitstamp.net/api/ticker/")


	jsonParse := j.Tasks[1].Adapter.(*tasks.JsonParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})


	bytes32 := j.Tasks[2].Adapter.(*tasks.EthBytes32)
	assert.Equal(t, bytes32.Address, "0x356a04bce728ba4c62a30294a55e6a8600a320b3")
	assert.Equal(t, bytes32.FunctionID, "12345679")
}

type JobJSON struct {
	ID string `json:"id"`
}

func TestCreateInvalidTasks(t *testing.T) {
	//fixtureprepare.SetUpDB()
	//defer fixtureprepare.TearDownDB()
	server := SetUpWeb()
	defer TearDownWeb()

	jsonStr := LoadJSON("./fixture/create_invalid_jobs.json")
	resp, err := http.Post(server.URL+"/tasks", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["jobs_not_exist is not a supported adapter type"]}`, string(body), "Repsonse should return JSON")
}

func TestShowJobs(t *testing.T) {
	db := SetUpDB()
	defer TearDownDB()
	server := SetUpWeb()
	defer TearDownWeb()

	j := models.NewTask()
	j.Schedule = "*****"

	db.Save(&j)

	resp, err := http.Get(server.URL + "/jobs/" + j.ID)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var respJob models.Job
	json.Unmarshal(b, &respJob)
	assert.Equal(t, respJob.Schedule, j.Schedule, "should have the same schedule")
}

func TestShowNotFoundJobs(t *testing.T) {
	SetUpDB()
	defer TearDownDB()
	server := SetUpWeb()
	defer TearDownWeb()
	resp, err := http.Get(server.URL + "/jobs/" + "garbage")
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}


