package web

import (
	"PhoenixOracle/core/models"
	"PhoenixOracle/core/orm"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var r *gin.Engine

func init() {
	r = Router()
}

type TaskJSON struct {
	ID string `json:"id"`
}

func TestCreateTasks(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()

	orm.Init()
	defer orm.Close()
	db := orm.GetDB()
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

//func TestCreateJobs2(t *testing.T) {
//	server := httptest.NewServer(r)
//	defer server.Close()
//	orm.Init()
//	defer orm.Close()
//	db := orm.GetDB()
//
//	jsonStr := []byte(`{"subtasks":[{"adapterType": "httpJSON", "adapterParams": {"endpoint": "https://bitstamp.net/api/ticker/", "fields": ["last"]}}], "schedule": "* 7 * * *","version":"1.0.0"}`)
//	resp, err := http.Post(server.URL+"/jobs", "application/json", bytes.NewBuffer(jsonStr))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	assert.Equal(t, 200, resp.StatusCode, "Response should be success")
//
//	defer resp.Body.Close()
//	b, err := ioutil.ReadAll(resp.Body)
//	var respJSON JobJSON
//	json.Unmarshal(b, &respJSON)
//
//	var j models.Task
//	db.One("ID", respJSON.ID, &j)
//	assert.Equal(t, j.ID, respJSON.ID, "Wrong job returned")
//	assert.Equal(t, j.Schedule, "* 7 * * *", "Wrong schedule saved")
//}

func TestCreateInvalidTasks(t *testing.T) {
	server := httptest.NewServer(r)
	defer server.Close()

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


