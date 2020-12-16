package fixtureprepare

import (
	"PhoenixOracle/gophoenix/core/orm"
	"PhoenixOracle/gophoenix/core/web"
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

var server *httptest.Server

func SetUpDB() *storm.DB {
	orm.InitTest()
	return orm.GetDB()
}

func TearDownDB() {
	orm.Close()
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
