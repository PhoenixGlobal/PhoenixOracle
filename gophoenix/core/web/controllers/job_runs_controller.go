package controllers

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/services"
	"github.com/gin-gonic/gin"
)

type JobRunsController struct{
	App *services.Application
}

func (self *JobRunsController) Index(c *gin.Context) {
	id := c.Param("id")
	jobRuns := []models.JobRun{}

	err := self.App.Store.Where("JobID", id, &jobRuns)

	if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"runs": jobRuns})
	}
}

