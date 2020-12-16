package web

import (
	"PhoenixOracle/core/web/controllers"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	t := controllers.TasksController{}
	r.POST("/tasks", t.Create)
	return r
}
