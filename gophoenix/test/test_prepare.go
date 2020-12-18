package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/web"
	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"time"
)

var server *httptest.Server

func SetUpDB() {
	models.InitDBTest()
}

func TearDownDB() {
	models.CloseDB()
}

func SetUpWeb() *httptest.Server {
	gin.SetMode(gin.TestMode)
	server = httptest.NewServer(web.Router())
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

