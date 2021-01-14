package test

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/models"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
	"testing"
)

func TestSendingEthereumTx(t *testing.T) {
	defer CloseGock(t)

	address := "0x1234567890"
	fid := "0x12345678"
	value := "0000abcdef"
	input := models.RunResultWithValue(value)
	config := NewConfig()

	response := `{"result": "0x0100"}`
	gock.New(config.EthereumURL).
		Post("/api").
		Reply(200).
		JSON(response)

	adapter := adapters.EthSendTx{
		Address:    address,
		FunctionID: fid,
		AdapterBase: adapters.AdapterBase{config},
	}
	result := adapter.Perform(input)
	assert.Equal(t, "0x0100", result.Value())
}

