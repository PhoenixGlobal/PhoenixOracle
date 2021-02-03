package web

import (
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/services"
	"PhoenixOracle/gophoenix/core/web/controllers"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"time"
)

func Router(app *services.Application) *gin.Engine {
	engine := gin.New()
	config := app.Store.Config
	basicAuth := gin.BasicAuth(gin.Accounts{config.BasicAuthUsername: config.BasicAuthPassword})
	engine.Use(loggerFunc(logger.LoggerWriter()), gin.Recovery(), basicAuth)
	//v2 := engine.Group("/v2")
	fmt.Println("11111111111111111111111")
	//{
		t := controllers.JobsController{app}
	engine.GET("/jobs", t.Index)
	fmt.Println("3333333333333333333333")
	engine.POST("/jobs", t.Create)
	engine.GET("/jobs/:id", t.Show)

		jr := controllers.JobRunsController{app}
	engine.GET("/jobs/:id/runs", jr.Index)
	engine.POST("/jobs/:jobID/runs", jr.Create)
	//}
	return engine
}

func loggerFunc(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr := bytes.NewBuffer(buf)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		start := time.Now()
		c.Next()
		end := time.Now()

		logger.Infow("Web request",
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"body", readBody(rdr),
			"clientIP", c.ClientIP(),
			"comment", c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"servedAt", end.Format("2006/01/02 - 15:04:05"),
			"latency", fmt.Sprintf("%v", end.Sub(start)),
		)
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
