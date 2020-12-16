package controllers

import (
	"PhoenixOracle/core/models"
	"PhoenixOracle/core/orm"
	"github.com/gin-gonic/gin"
)

type TasksController struct{}

func (ac *TasksController) Create(c *gin.Context) {
	db := orm.GetDB()
	j := models.NewTask()
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if _, err = j.Valid(); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = db.Save(&j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, gin.H{"id": j.ID})
	}
}


