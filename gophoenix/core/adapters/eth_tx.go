package adapters

import (
	"PhoenixOracle/gophoenix/core/store"
	"PhoenixOracle/gophoenix/core/store/models"
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type EthSignTx struct {
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSignTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	str := self.FunctionID + input.Value()
	data := common.FromHex(str)
	fmt.Println("&&&&&&&&&&&&&&&&&&&&&")
	keyStore := store.KeyStore
	nonce, err := store.Eth.GetNonce(keyStore.GetAccount())
	if err != nil {
		return models.RunResultWithError(err)
	}
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(self.Address),
		big.NewInt(0),
		500000,
		big.NewInt(20000000000),
		data,
	)

	signedTx, err := keyStore.SignTx(tx, store.Config.ChainID)
	if err != nil {
		return models.RunResultWithError(err)
	}

	buffer := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(buffer); err != nil {
		return models.RunResultWithError(err)
	}
	return models.RunResultWithValue(common.Bytes2Hex(buffer.Bytes()))
}

type EthSendRawTx struct {}

func (self *EthSendRawTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	result, err := store.Eth.SendRawTx(input.Value())
	if err != nil {
		return models.RunResultWithError(err)
	}
	return models.RunResultWithValue(result)
}

type EthConfirmTx struct {}

func (self *EthConfirmTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	txid := input.Value()
	for {
		receipt, err := store.Eth.GetTxReceipt(txid)
		if err != nil {
			return models.RunResultWithError(err)
		} else if receipt.TxHash.Hex() == "" {
			return models.RunResultPending(input)
		} else {
			return models.RunResultWithValue(txid)
		}
	}
}

type EthSignAndSendTx struct {
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}

func (self *EthSignAndSendTx) Perform(input models.RunResult, store *store.Store) models.RunResult {
	signer := &EthSignTx{
		Address:     self.Address,
		FunctionID:  self.FunctionID,
	}
	sender := &EthSendRawTx{
	}
	confirmer := &EthConfirmTx{
	}

	signed := signer.Perform(input,store)
	if signed.HasError() {
		return signed
	}
	sent := sender.Perform(signed,store)
	if sent.HasError() {
		return sent
	}
	return confirmer.Perform(sent,store)
}

