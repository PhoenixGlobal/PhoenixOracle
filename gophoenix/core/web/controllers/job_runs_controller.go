package controllers

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/store"
	"github.com/gin-gonic/gin"
)

type JobRunsController struct{
	Store store.Store
}

func (self *JobRunsController) Index(c *gin.Context) {
	id := c.Param("id")
	jobRuns := []models.JobRun{}

	err := self.Store.Where("JobID", id, &jobRuns)

	if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"runs": jobRuns})
	}
}

