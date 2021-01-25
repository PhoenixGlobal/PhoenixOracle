package test

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestEthGetTxReceipt(t *testing.T) {
	store := NewStore()
	defer CleanUpStore(store)
	eth := store.Eth

	response := LoadJSON("./fixture/eth_getTransactionReceipt.json")
	gock.New(store.Config.EthereumURL).
		Post("").
		Reply(200).
		JSON(response)

	txid := "0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"
	receipt, err := eth.GetTxReceipt(txid)
	assert.Nil(t, err)
	assert.Equal(t, txid, receipt.TXID)
	assert.Equal(t, uint64(11), receipt.BlockNumber)
}
