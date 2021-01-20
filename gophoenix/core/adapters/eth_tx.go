package adapters

import (
	"PhoenixOracle/gophoenix/core/store/models"
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type EthSendRawTx struct {
	AdapterBase
}

func (self *EthSendRawTx) Perform(input models.RunResult) models.RunResult {
	result, err := self.Store.Eth.SendRawTx(input.Value())
	if err != nil {
		return models.RunResultWithError(err)
	}
	return models.RunResultWithValue(result)
}

type EthSignTx struct {
	AdapterBase
	Address    string `json:"address"`
	FunctionID string `json:"functionID"`
}


func (self *EthSignTx) Perform(input models.RunResult) models.RunResult {
	str := self.FunctionID + input.Value()
	data := common.FromHex(str)
	keyStore := self.Store.KeyStore
	nonce, err := keyStore.GetAccount().GetNonce(self.Store.Config)
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

	signedTx, err := keyStore.SignTx(tx, self.Store.Config.ChainID)
	if err != nil {
		return models.RunResultWithError(err)
	}

	buffer := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(buffer); err != nil {
		return models.RunResultWithError(err)
	}
	return models.RunResultWithValue(common.Bytes2Hex(buffer.Bytes()))
}
