package test

import (
	configlib "PhoenixOracle/gophoenix/core/config"
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/services"
	"PhoenixOracle/gophoenix/core/web"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"
)

type TestStore struct {
	*services.Store
	Server *httptest.Server
}
const testRootDir = "./tmp/test"
const testUsername = "testusername"
const testPassword = "testpassword"

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
	return StoreWithConfig(NewConfig())

}

func StoreWithConfig(config configlib.Config)  *TestStore{
	if err := os.MkdirAll(config.RootDir, os.FileMode(0700)); err != nil {
		log.Fatal(err)
	}
	logger.SetLoggerDir(config.RootDir)
	store := services.NewStore(config)
	return &TestStore{
		Store: store,
	}
}

func NewConfig() configlib.Config {
	return configlib.Config{
		RootDir:           path.Join(testRootDir, fmt.Sprintf("%d", time.Now().UnixNano())),
		BasicAuthUsername: testUsername,
		BasicAuthPassword: testPassword,
		EthereumURL:       "https://example.com/api",
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

func CloseGock(t *testing.T) {
	assert.True(t, gock.IsDone(), "Not all gock requests were fulfilled")
	gock.DisableNetworking()
	gock.Off()
}

func LoadJSON(file string) []byte {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func CopyFile(src, dst string) {
	from, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

func AddPrivateKey(config configlib.Config, src string) {
	err := os.MkdirAll(config.KeysDir(), os.FileMode(0700))
	if err != nil {
		log.Fatal(err)
	}

	dst := config.KeysDir() + "/testwallet.json"
	CopyFile(src, dst)
}

func TimeParse(s string) time.Time {
	t, err := dateparse.ParseAny(s)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func BasicAuthPost(url string, contentType string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", contentType)
	request.SetBasicAuth(testUsername, testPassword)
	resp, err := client.Do(request)
	return resp, err
}

func BasicAuthGet(url string) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(testUsername, testPassword)
	resp, err := client.Do(request)
	return resp, err
}
