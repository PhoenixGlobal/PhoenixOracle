package controllers

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/services"
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
)

type JobsController struct{
	Store *services.Store
}

func (tc *JobsController) Create(c *gin.Context) {
	j := models.NewJob()
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = tc.Store.AddJob(j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": j.ID})
	}
}

func (tc *JobsController) Show(c *gin.Context) {
	id := c.Param("id")
	var j models.Job
	err := tc.Store.One("ID", id, &j)

	if err == storm.ErrNotFound {
		c.JSON(404, gin.H{
			"errors": []string{"Job not found."},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, j)
	}
}


