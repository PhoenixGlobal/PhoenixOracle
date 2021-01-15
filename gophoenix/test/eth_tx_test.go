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

	value := "0000abcdef"
	input := models.RunResultWithValue(value)
	config := NewConfig()

	response := `{"result": "0x0100"}`
	gock.New(config.EthereumURL).
		Post("").
		Reply(200).
		JSON(response)

	adapter := adapters.EthSendRawTx{
		AdapterBase: adapters.AdapterBase{config},
	}

	result := adapter.Perform(input)
	assert.Equal(t, "0x0100", result.Value())
}

func TestSigningEthereumTx(t *testing.T) {
	defer CloseGock(t)

	data := "0000abcdef"
	address := "0xb70a511bac46ec6442ac6d598eac327334e634db"
	fid := "0x12345678"
	input := models.RunResultWithValue(data)
	config := NewConfig()

	adapter := adapters.EthSignTx{
		Address:     address,
		FunctionID:  fid,
		AdapterBase: adapters.AdapterBase{config},
	}
	result := adapter.Perform(input)
	assert.Contains(t, result.Value(), data)
	assert.Contains(t, result.Value(), address[2:len(address)])
}


