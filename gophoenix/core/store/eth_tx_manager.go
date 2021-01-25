package store

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type EthTxManager struct {
	KeyStore *KeyStore
	Eth      *Eth
	Config   Config
}

func (self *EthTxManager) NewSignedTx(to, data string) (*types.Transaction, error) {
	account := self.KeyStore.GetAccount()
	nonce, err := self.Eth.GetNonce(account)
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(to),
		big.NewInt(0),
		500000,
		big.NewInt(20000000000),
		common.FromHex(data),
	)
	return self.KeyStore.SignTx(tx, self.Config.ChainID)
}

func (self *EthTxManager) SendRawTx(hex string) (string, error) {
	return self.Eth.SendRawTx(hex)
}

func (self *EthTxManager) TxConfirmed(txid string) (bool, error) {
	receipt, err := self.Eth.GetTxReceipt(txid)
	if err != nil {
		return false, err
	} else if receipt.TxHash.Hex() == "" {
		return false, nil
	}
	return true, nil
}
