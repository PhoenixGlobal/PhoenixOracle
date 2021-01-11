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

func Router(store *services.Store) *gin.Engine {
	r := gin.New()
	r.Use(handlerFunc(logger.LoggerWriter()), gin.Recovery())
	t := controllers.JobsController{store}
	r.POST("/jobs", t.Create)
	r.GET("/jobs/:id", t.Show)

	jr := controllers.JobRunsController{store}
	r.GET("/jobs/:id/runs", jr.Index)
	return r
}

func handlerFunc(logger *logger.Logger) gin.HandlerFunc {
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
