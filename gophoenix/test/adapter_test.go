package test

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatingAdapterWithConfig(t *testing.T) {
	config := NewConfig()
	task := models.Task{Type: "NoOp"}
	adapter, err := adapters.For(task, config)
	adapter.Perform(models.RunResult{})
	assert.Nil(t, err)
	rval := adapter.(*adapters.NoOp).Config
	assert.Equal(t, "", rval.EthereumURL)
}