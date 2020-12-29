package web

import (
	"PhoenixOracle/gophoenix/core/services"
	storelib "PhoenixOracle/gophoenix/core/store"
	"PhoenixOracle/gophoenix/core/web/controllers"
	"github.com/gin-gonic/gin"
)

func Router(store storelib.Store) *gin.Engine {
	r := gin.New()
	r.Use(gin.LoggerWithWriter(services.LoggerWriter()), gin.Recovery())
	t := controllers.JobsController{store}
	r.POST("/jobs", t.Create)
	r.GET("/jobs/:id", t.Show)

	jr := controllers.JobRunsController{store}
	r.GET("/jobs/:id/runs", jr.Index)
	return r
}
