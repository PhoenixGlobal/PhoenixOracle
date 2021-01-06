package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/scheduler"
	"PhoenixOracle/gophoenix/core/store"
	"PhoenixOracle/gophoenix/core/web"
	"encoding/json"
	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"
)

var server *httptest.Server
func init() {
	if err := os.RemoveAll(filepath.Dir(models.DBpath("test"))); err != nil {
		log.Println(err)
	}

	gomega.SetDefaultEventuallyTimeout(3 * time.Second)
}


type JobJSON struct {
	ID string `json:"id"`
}

func JobJSONFromResponse(resp *http.Response) JobJSON {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var respJSON JobJSON
	json.Unmarshal(b, &respJSON)
	return respJSON
}
func Store() store.Store {
	orm := models.InitORM("test")
	return store.Store{
		ORM:       orm,
		Scheduler: scheduler.NewScheduler(orm),
	}
}

func SetUpWeb(s store.Store) *httptest.Server {
	gin.SetMode(gin.TestMode)
	server = httptest.NewServer(web.Router(s))
	return server
}

func TearDownWeb() {
	gin.SetMode(gin.DebugMode)
	server.Close()
}

func LoadJSON(file string) []byte {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func TimeParse(s string) time.Time {
	t, err := dateparse.ParseAny(s)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

