package test

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/store/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatingAdapterWithConfig(t *testing.T) {
	store := NewStore()
	defer store.Close()
	task := models.Task{Type: "NoOp"}
	adapter, err := adapters.For(task, store)
	adapter.Perform(models.RunResult{})
	assert.Nil(t, err)
	rval := adapter.(*adapters.NoOp).Store.Config
	assert.Equal(t, "", rval.EthereumURL)
}