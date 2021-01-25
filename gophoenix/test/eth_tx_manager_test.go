package test

import (
	"PhoenixOracle/gophoenix/core/utils"
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestTxManagerNewSignedTx(t *testing.T) {
	app := NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	defer CloseGock(t)
	assert.Nil(t, store.KeyStore.Unlock("password"))
	config := store.Config
	manager := store.Tx

	data := "0000abcdef"
	to := "0xb70a511baC46ec6442aC6D598eaC327334e634dB"

	response := `{"result": "0x0100"}` //
	gock.New(config.EthereumURL).
		Post("").
		Reply(200).
		JSON(response)

	tx, err := manager.NewSignedTx(to, data)
	assert.Nil(t, err)
	assert.Equal(t, uint64(256), tx.Nonce())
	assert.Equal(t, common.FromHex(data), tx.Data())
	assert.Equal(t, to, tx.To().Hex())

	signer := store.KeyStore.GetAccount().Address.String()
	rlp := bytes.NewBuffer([]byte{})
	assert.Nil(t, tx.EncodeRLP(rlp))
	rlpHex := common.Bytes2Hex(rlp.Bytes())
	sender, err := utils.SenderFromTxHex(rlpHex, config.ChainID)
	assert.Equal(t, signer, sender.Hex())
}

func TestTxManagerSending(t *testing.T) {
	store := NewStore()
	defer CleanUpStore(store)
	defer CloseGock(t)
	config := store.Config
	manager := store.Tx

	txHex := "0000abcdef"

	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	response := `{"result": "` + txid + `"}`
	gock.New(config.EthereumURL).
		Post("").
		Reply(200).
		JSON(response)

	response, err := manager.SendRawTx(txHex)
	assert.Nil(t, err)
	assert.Equal(t, txid, response)
}

func TestTxManagerConfirmTxTrue(t *testing.T) {
	t.Parallel()
	app := NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	manager := store.Tx

	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	eth := app.MockEthClient()
	eth.Register("eth_getTransactionReceipt", types.Receipt{TxHash: common.HexToHash(txid)})

	confirmed, err := manager.TxConfirmed(txid)
	assert.Nil(t, err)
	assert.True(t, confirmed)
}

func TestTxManagerConfirmTxFalse(t *testing.T) {
	t.Parallel()
	app := NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	manager := store.Tx

	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	eth := app.MockEthClient()
	eth.Register("eth_getTransactionReceipt", types.Receipt{})

	confirmed, err := manager.TxConfirmed(txid)
	assert.Nil(t, err)
	assert.False(t, confirmed)
}
