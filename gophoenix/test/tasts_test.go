package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/test/fixtureprepare"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	db := fixtureprepare.SetUpDB()
	defer fixtureprepare.TearDownDB()
	j1 := models.NewTask()
	j1.Schedule = "1 * * * *"
	db.Save(&j1)

	var j2 models.Task
	db.One("ID",j1.ID,&j2)

	assert.Equal(t, j1.Schedule, j2.Schedule)
}