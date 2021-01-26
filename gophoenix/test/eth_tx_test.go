package test

import (
	"PhoenixOracle/gophoenix/core/adapters"
	strpkg "PhoenixOracle/gophoenix/core/store"
	"PhoenixOracle/gophoenix/core/store/models"
	"PhoenixOracle/gophoenix/core/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEhxTxAdapter(t *testing.T) {
	app := NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := store.Config
	app.Store.KeyStore.Unlock(Password)

	eth := app.MockEthClient()
	eth.Register("eth_getTransactionCount", `0x0100`)
	txid := `0x83c52c31cd40a023728fbc21a570316acd4f90525f81f1d7c477fd958ffa467f`
	confed := uint64(23456)
	eth.Register("eth_sendRawTransaction", txid)
	eth.Register("eth_getTransactionReceipt", strpkg.TxReceipt{TXID: txid, BlockNumber: confed})
	eth.Register("eth_blockNumber", utils.Uint64ToHex(confed+config.EthConfMin))

	adapter := adapters.EthTx{
		Address:    "0x2C83ACd90367e7E0D3762eA31aC77F18faecE874",
		FunctionID: "12345678",
	}
	input := models.RunResultWithValue("")
	output := adapter.Perform(input, store)

	assert.False(t, output.HasError())

	from := store.KeyStore.GetAccount().Address.String()
	txs := []models.EthTx{}
	assert.Nil(t, store.Where("From", from, &txs))
	store.All(&txs)
	assert.Equal(t, 1, len(txs))
}

func TestEthTxAdapterWithError(t *testing.T) {
	app := NewApplicationWithKeyStore()
	app.Store.KeyStore.Unlock(Password)
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



