package web

import (
	"PhoenixOracle/gophoenix/core/web/controllers"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	t := controllers.JobsController{}
	r.POST("/jobs", t.Create)
	r.GET("/jobs/:id", t.Show)
	return r
}
