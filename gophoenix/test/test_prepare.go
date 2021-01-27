package test

import (
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/services"
	"PhoenixOracle/gophoenix/core/store"
	"PhoenixOracle/gophoenix/core/store/models"
	"PhoenixOracle/gophoenix/core/web"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"reflect"
	"testing"
	"time"
)

type TestStore struct {
	*store.Store
	Server *httptest.Server
}
const testRootDir = "./tmp/test"
const testUsername = "testusername"
const testPassword = "testpassword"
const Password = "password"

func init() {
	dir, err := homedir.Expand(testRootDir)
	if err != nil {
		logger.Fatal(err)
	}

	if err = os.RemoveAll(dir); err != nil {
		log.Println(err)
	}

	gomega.SetDefaultEventuallyTimeout(3 * time.Second)
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

func StoreWithConfig(config store.Config)  *TestStore{
	if err := os.MkdirAll(config.RootDir, os.FileMode(0700)); err != nil {
		log.Fatal(err)
	}
	logger.SetLoggerDir(config.RootDir)
	store := store.NewStore(config)
	return &TestStore{
		Store: store,
	}
}

func NewConfig() store.Config {
	return store.Config{
		RootDir:           path.Join(testRootDir, fmt.Sprintf("%d", time.Now().UnixNano())),
		BasicAuthUsername: testUsername,
		BasicAuthPassword: testPassword,
		EthereumURL:       "https://example.com/api",
		ChainID:           3,
		EthMinConfirmations: 6,
		EthGasBumpWei:      big.NewInt(5000000000),
		EthGasBumpThreshold: 3,
		EthGasPriceDefault:  big.NewInt(20000000000),
		PollingSchedule:   "* * * * * *",
	}
}

type TestApplication struct {
	*services.Application
	Server *httptest.Server
}

func NewApplication() *TestApplication {
	return NewApplicationWithConfig(NewConfig())
}

func NewApplicationWithConfig(config store.Config) *TestApplication {
	return  &TestApplication{Application: services.NewApplication(config)}
}

func NewApplicationWithKeyStore() *TestApplication {
	app := NewApplication()
	if _, err := app.Store.KeyStore.NewAccount(Password); err != nil {
		logger.Fatal(err)
	}
	if err := app.Store.KeyStore.Unlock(Password); err != nil {
		logger.Fatal(err)
	}
	return app
}

func (self *TestApplication)NewServer() *httptest.Server {
	gin.SetMode(gin.TestMode)
	server := httptest.NewServer(web.Router(self.Application))
	self.Server = server
	return server
}

func (self *TestApplication)Stop()() {
	self.Application.Stop()
	CleanUpStore(self.Store)
	if self.Server != nil {
		gin.SetMode(gin.DebugMode)
		self.Server.Close()
	}
}

func (self *TestApplication) MockEthClient() *EthMock {
	mock := NewMockGethRpc()
	eth := &store.EthClient{mock}
	self.Store.Eth.EthClient = eth
	return mock
}

func NewMockGethRpc() *EthMock {
	return &EthMock{}
}

type EthMock struct {
	Responses []MockResponse
}

type MockResponse struct {
	methodName string
	response   interface{}
	errMsg     string
	hasError   bool
}

func (self *EthMock) Register(method string, response interface{}) {
	res := MockResponse{
		methodName: method,
		response:   response,
	}
	self.Responses = append(self.Responses, res)
}

func (self *EthMock) RegisterError(method, errMsg string) {
	res := MockResponse{
		methodName: method,
		errMsg:     errMsg,
		hasError:   true,
	}
	self.Responses = append(self.Responses, res)
}

func (self *EthMock) AllCalled() bool {
	return len(self.Responses) == 0
}

func copyWithoutIndex(s []MockResponse, index int) []MockResponse {
	return append(s[:index], s[index+1:]...)
}

func (self *EthMock) Call(result interface{}, method string, args ...interface{}) error {
	for i, resp := range self.Responses {
		if resp.methodName == method {
			self.Responses = copyWithoutIndex(self.Responses, i)
			if resp.hasError {
				return fmt.Errorf(resp.errMsg)
			} else {
				ref := reflect.ValueOf(result)
				reflect.Indirect(ref).Set(reflect.ValueOf(resp.response))
				return nil
			}
		}
	}
	return fmt.Errorf("EthMock: Method %v not registered", method)
}


func NewStore() *store.Store {
	return store.NewStore(NewConfig())
}

func CleanUpStore(store *store.Store) {
	store.Close()
	if err := os.RemoveAll(store.Config.RootDir); err != nil {
		log.Println(err)
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

func copyFile(src, dst string) {
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

func AddPrivateKey(config store.Config, src string) {
	err := os.MkdirAll(config.KeysDir(), os.FileMode(0700))
	if err != nil {
		log.Fatal(err)
	}

	dst := config.KeysDir() + "/testwallet.json"
	copyFile(src, dst)
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

func NewJob() models.Job {
	j := models.NewJob()
	j.Tasks = []models.Task{{Type: "NoOp"}}
	return j
}

func NewJobWithSchedule(sched string) models.Job {
	j := NewJob()
	j.Initiators = []models.Initiator{{Type: "cron", Schedule: models.Cron(sched)}}
	return j
}

func NewJobWithWebInitiator() models.Job {
	j := NewJob()
	j.Initiators = []models.Initiator{{Type: "web"}}
	return j
}

func NewEthTx(from string, sentAt uint64) *models.EthTx {
	return &models.EthTx{
		From:     from,
		Nonce:    0,
		Data:     "deadbeef",
		Value:    big.NewInt(0),
		GasLimit: uint64(250000),
		Attempts: []*models.EthTxAttempt{&models.EthTxAttempt{
			TxID:     NewTxID(),
			GasPrice: big.NewInt(20000000000),
			Hex:      "0x0000",
			SentAt:   sentAt,
		}},
	}
}

func NewTxID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hexutil.Encode(b)
}

func NewEthAddress() string {
	b := make([]byte, 20)
	rand.Read(b)
	return hexutil.Encode(b)
}
