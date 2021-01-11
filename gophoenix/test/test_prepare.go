package test

import (
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/services"
	"PhoenixOracle/gophoenix/core/web"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"path"
	"time"
)

type TestStore struct {
	*services.Store
	Server *httptest.Server
}
const testRootDir = "./tmp/test"

func init() {
	dir, err := homedir.Expand(testRootDir)
	if err != nil {
		logger.Fatal(err)
	}

	if err = os.RemoveAll(dir); err != nil {
		log.Println(err)
	}

	gomega.SetDefaultEventuallyTimeout(2 * time.Second)
}


type JobJSON struct {
	ID string `json:"id"`
}

func JobJSONFromResponse(body io.Reader) JobJSON {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}
	var respJSON JobJSON
	json.Unmarshal(b, &respJSON)
	return respJSON
}
func Store() *TestStore {
	config := services.NewConfig(path.Join(testRootDir, fmt.Sprintf("%d", time.Now().UnixNano())))
	logger.SetLoggerDir(config.RootDir)
	store := services.NewStore(config)
	return &TestStore{
		Store: store,
	}
}

func (self *TestStore)SetUpWeb() *httptest.Server {
	gin.SetMode(gin.TestMode)
	server := httptest.NewServer(web.Router(self.Store))
	self.Server = server
	return server
}

func (self *TestStore)Close()() {
	self.Store.Close()
	if self.Server != nil {
		gin.SetMode(gin.DebugMode)
		self.Server.Close()
	}
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

