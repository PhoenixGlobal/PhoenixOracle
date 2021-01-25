package test

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/store/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSigningAndSendingTx(t *testing.T) {
	defer CloseGock(t)

	app := NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	eth := app.MockEthClient()
	eth.RegisterError("eth_getTransactionCount", "Cannot connect to nodes")

	adapter := adapters.EthTx{
		Address:     "recipient",
		FunctionID:  "fid",
	}
	input := models.RunResultWithValue("Hello World!")
	output := adapter.Perform(input,store)

	assert.True(t, output.HasError())
	assert.Equal(t, output.Error(), "Cannot connect to nodes")
}



