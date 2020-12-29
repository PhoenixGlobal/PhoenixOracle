package test

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/scheduler"
	"bytes"
	"encoding/json"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	//gock "github.com/h2non/gock.git"
	gock "gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"testing"
)


func TestCreateTasks(t *testing.T) {
	store := Store()
	defer store.Close()

	server := SetUpWeb(store)
	defer TearDownWeb()

	jsonStr := LoadJSON("./fixture/create_jobs.json")
	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	//if err != nil {
	//	t.Fatal(err)
	//}
	assert.Equal(t, 200, resp.StatusCode, "Response should be success")

	respJSON := JobJSONFromResponse(resp)
	var j models.Job
	store.One("ID", respJSON.ID, &j)
	assert.Equal(t, j.ID, respJSON.ID, "Wrong job returned")
	assert.Equal(t, j.Tasks[0].Type, "HttpGet")

	adapter1,_ := j.Tasks[0].Adapter()
	httpGet := adapter1.(*adapters.HttpGet)
	assert.Nil(t, err)
	assert.Equal(t, httpGet.Endpoint, "https://bitstamp.net/api/ticker/")


	adapter2,_ := j.Tasks[1].Adapter()
	jsonParse := adapter2.(*adapters.JsonParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})

	adapter3,_ := j.Tasks[2].Adapter()
	bytes32 := adapter3.(*adapters.EthBytes32)
	assert.Equal(t, bytes32.Address, "0x356a04bce728ba4c62a30294a55e6a8600a320b3")
	assert.Equal(t, bytes32.FunctionID, "12345679")

	schedule := j.Schedule
	assert.Equal(t, schedule.Cron, models.Cron("* 7 * * *"))
	assert.Equal(t, (*models.Time)(nil), schedule.StartAt, "Wrong start at saved")
	endAt := models.Time{TimeParse("2020-12-17T14:05:29Z")}
	assert.Equal(t, endAt, *schedule.EndAt, "Wrong end at saved")
	runAt0 := models.Time{TimeParse("2020-12-17T14:05:19Z")}
	assert.Equal(t, runAt0, schedule.RunAt[0], "Wrong run at saved")


}

func TestCreateJobsIntegration(t *testing.T) {
	RegisterTestingT(t)
	defer gock.Off()

	store := Store()
	store.Start()
	defer store.Close()
	server := SetUpWeb(store)
	defer TearDownWeb()

	jsonStr := LoadJSON("./fixtures/create_http_get_job.json")
	resp, _ := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	respJSON := JobJSONFromResponse(resp)

	expectedResponse := `{"high": "10744.00", "last": "10583.75", "timestamp": "1512156162", "bid": "10555.13", "vwap": "10097.98", "volume": "17861.33960013", "low": "9370.11", "ask": "10583.00", "open": "9927.29"}`
	gock.New("https://www.bitstamp.net").
		Get("/api/ticker").
		Reply(200).
		JSON(expectedResponse)
	sched:= scheduler.NewScheduler(store.ORM)
	sched.Start()
	defer sched.Stop()

	jobRuns := []models.JobRun{}
	Eventually(func() []models.JobRun {
		_ = store.Where("JobID", respJSON.ID, &jobRuns)
		return jobRuns
	}).Should(HaveLen(1))

	var job models.Job
	err := store.One("ID",respJSON.ID, &job)
	assert.Nil(t,err)

	jobRuns, err = store.JobRunsFor(job)
	assert.Nil(t, err)
	jobRun := jobRuns[0]
	assert.Equal(t, jobRun.Result.Output["value"], expectedResponse)
}


func TestCreateInvalidTasks(t *testing.T) {
	//fixtureprepare.SetUpDB()
	//defer fixtureprepare.TearDownDB()
	store := Store()
	defer store.Close()
	server := SetUpWeb(store)
	defer TearDownWeb()

	jsonStr := LoadJSON("./fixture/create_invalid_jobs.json")
	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["jobs_not_exist is not a supported adapters type"]}`, string(body), "Repsonse should return JSON")
}

func TestCreateInvalidCron(t *testing.T) {
	store := Store()
	defer store.Close()
	server := SetUpWeb(store)
	defer TearDownWeb()

	jsonStr := LoadJSON("./fixture/create_invalid_cron.json")
	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 500, resp.StatusCode, "Response should be internal error")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, `{"errors":["Cron: Failed to parse int from !: strconv.Atoi: parsing \"!\": invalid syntax"]}`, string(body), "Response should return JSON")
}

func TestShowJobs(t *testing.T) {
	store := Store()
	defer store.Close()
	server := SetUpWeb(store)
	defer TearDownWeb()

	j := models.NewJob()
	j.Schedule = models.Schedule{Cron: "1 * * * *"}

	store.Save(&j)

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
	store := Store()
	defer store.Close()
	server := SetUpWeb(store)
	defer TearDownWeb()
	resp, err := http.Get(server.URL + "/jobs/" + "garbage")
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode, "Response should be not found")
}


