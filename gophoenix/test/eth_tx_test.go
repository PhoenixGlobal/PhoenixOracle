package test

import (
	"PhoenixOracle/gophoenix/core/adapters"
	strpkg "PhoenixOracle/gophoenix/core/store"
	"PhoenixOracle/gophoenix/core/store/models"
	"PhoenixOracle/gophoenix/core/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEthTxAdapterConfirmed(t *testing.T) {
	t.Parallel()
	app := NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := store.Config
	app.Store.KeyStore.Unlock(Password)
	eth := app.MockEthClient()
	eth.Register("eth_getTransactionCount", `0x0100`)
	txid := NewTxID()
	confed := uint64(23456)
	eth.Register("eth_sendRawTransaction", txid)
	eth.Register("eth_getTransactionReceipt", strpkg.TxReceipt{TxID: txid, BlockNumber: confed})
	eth.Register("eth_blockNumber", utils.Uint64ToHex(confed+config.EthMinConfirmations))

	adapter := adapters.EthTx{
		Address:    NewEthAddress(),
		FunctionID: "12345678",
	}
	input := models.RunResultWithValue("")
	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())

	from := store.KeyStore.GetAccount().Address.String()
	txs := []models.EthTx{}
	assert.Nil(t, store.Where("From", from, &txs))
	assert.Equal(t, 1, len(txs))
	assert.Equal(t, 1, len(txs[0].Attempts))

	assert.True(t, eth.AllCalled())
}

func TestEthTxAdapterFromPending(t *testing.T) {
	t.Parallel()
	app := NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold-1))

	from := store.KeyStore.GetAccount().Address.String()
	txr := NewEthTx(from, sentAt)
	assert.Nil(t, store.SaveTx(txr))
	adapter := adapters.EthTx{Address: NewEthAddress(), FunctionID: "12345678"}
	input := models.RunResultPending(models.RunResultWithValue(txr.TxID()))

	output := adapter.Perform(input, store)

	assert.True(t, output.Pending)
	assert.Nil(t, store.One("ID", txr.ID, txr))
	assert.Equal(t, 1, len(txr.Attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterFromPendingBumpGas(t *testing.T) {
	t.Parallel()
	app := NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := store.Config

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
	sentAt := uint64(23456)
	ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold))
	ethMock.Register("eth_sendRawTransaction", NewTxID())

	from := store.KeyStore.GetAccount().Address.String()
	txr := NewEthTx(from, sentAt)
	assert.Nil(t, store.SaveTx(txr))
	adapter := adapters.EthTx{Address: NewEthAddress(), FunctionID: "12345678"}
	input := models.RunResultPending(models.RunResultWithValue(txr.TxID()))

	output := adapter.Perform(input, store)

	assert.True(t, output.Pending)
	assert.Nil(t, store.One("ID", txr.ID, txr))
	assert.Equal(t, 2, len(txr.Attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthTxAdapterWithError(t *testing.T) {
	t.Parallel()
	app := NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	eth := app.MockEthClient()
	eth.RegisterError("eth_getTransactionCount", "Cannot connect to nodes")

	adapter := adapters.EthTx{
		Address:    "recipient",
		FunctionID: "fid",
	}
	input := models.RunResultWithValue("Hello World!")
	output := adapter.Perform(input, store)

	assert.True(t, output.HasError())
	assert.Equal(t, output.Error(), "Cannot connect to nodes")
}



