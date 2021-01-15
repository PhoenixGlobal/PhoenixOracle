package test

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/models"
	storelib "PhoenixOracle/gophoenix/core/store"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
	"testing"
)

func TestSendingEthereumTx(t *testing.T) {
	store := NewStore()
	defer store.Close()
	defer CloseGock(t)

	value := "0000abcdef"
	input := models.RunResultWithValue(value)
	config := store.Config
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
	config := NewConfig()
	AddPrivateKey(config, "./fixtures/3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea.json")
	password := "password"

	store := storelib.NewStore(config)
	defer CleanUpStore(store)

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
}


