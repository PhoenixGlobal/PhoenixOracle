package controllers

import (
	"PhoenixOracle/gophoenix/core/models"
	"github.com/gin-gonic/gin"
)

type JobRunsController struct{}

func (self *JobRunsController) Index(c *gin.Context) {
	id := c.Param("id")
	jobRuns := []models.JobRun{}

	err := models.Where("JobID", id, &jobRuns)

	if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"runs": jobRuns})
	}
}

