package model_test

import (
	"PhoenixOracle/core/models"
	"PhoenixOracle/core/orm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	j1 := models.NewTask()
	j1.Schedule = "1 * * * *"
	orm.Init()
	defer orm.Close()

	db := orm.GetDB()
	db.Save(&j1)

	var j2 models.Task
	db.One("ID",j1.ID,&j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}