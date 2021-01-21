package test

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/store/models"
	"PhoenixOracle/gophoenix/core/utils"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
	"testing"
)

func TestSendingEthereumTx(t *testing.T) {
	store := NewStore()
	defer store.Close()
	defer CloseGock(t)
	config := store.Config

	value := "0000abcdef"
	input := models.RunResultWithValue(value)
	response := `{"result": "0x0100"}`
	gock.New(config.EthereumURL).
		Post("").
		Reply(200).
		JSON(response)

	adapter := adapters.EthSendRawTx{
		AdapterBase: adapters.AdapterBase{store},
	}

	result := adapter.Perform(input)
	assert.Equal(t, "0x0100", result.Value())
}

func TestSigningEthereumTx(t *testing.T) {
	defer CloseGock(t)
	app := NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := app.Store.Config
	sender := store.KeyStore.GetAccount().Address.String()
	password := "password"

	response := `{"result": "0x11"}`
	gock.New(config.EthereumURL).
		Post("").
		Reply(200).
		JSON(response)

	err := store.KeyStore.Unlock(password)
	assert.Nil(t, err)

	data := "0000abcdef"
	recipient := "0xb70a511bac46ec6442ac6d598eac327334e634db"
	fid := "0x12345678"
	input := models.RunResultWithValue(data)

	adapter := adapters.EthSignTx{
		Address:     recipient,
		FunctionID:  fid,
		AdapterBase: adapters.AdapterBase{store},
	}
	result := adapter.Perform(input)
	assert.Contains(t, result.Value(), data)
	assert.Contains(t, result.Value(), recipient[2:len(recipient)])

	tx, err := utils.DecodeTxFromHex(result.Value(), config.ChainID)
	assert.Nil(t, err)
	assert.Equal(t, uint64(17), tx.Nonce())

	actual, err := utils.SenderFromTxHex(result.Value(), config.ChainID)
	assert.Equal(t, sender, actual.Hex())
}

func TestSigningAndSendingTx(t *testing.T) {
	defer CloseGock(t)

	app := NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	eth := app.MockEthClient()
	eth.RegisterError("eth_getTransactionCount", "Cannot connect to nodes")

	adapter := adapters.EthSignTx{
		Address:     "recipient",
		FunctionID:  "fid",
		AdapterBase: adapters.AdapterBase{store},
	}
	input := models.RunResultWithValue("Hello World!")
	output := adapter.Perform(input)

	assert.True(t, output.HasError())
	assert.Equal(t, output.Error(), "Cannot connect to nodes")
}



