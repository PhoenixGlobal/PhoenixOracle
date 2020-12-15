package web

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	t := TasksController{}
	r.POST("/tasks", t.Create)
	return r
}
