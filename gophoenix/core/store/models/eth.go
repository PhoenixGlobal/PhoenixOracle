package models

import (
	"PhoenixOracle/gophoenix/core/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type EthTx struct {
	ID       uint64 `storm:"id,increment,index"`
	From     string
	To       string
	Data     string
	Nonce    uint64
	Value    *big.Int
	GasLimit uint64
	Attempts []*EthTxAttempt `storm:"inline"`
}

func (self *EthTx) NewAttempt(tx *types.Transaction) (*EthTxAttempt, error) {
	hex, err := utils.EncodeTxToHex(tx)
	if err != nil {
		return nil, err
	}
	attempt := &EthTxAttempt{
		TxID:      tx.Hash().String(),
		GasPrice:  tx.GasPrice(),
		Confirmed: false,
		Hex:       hex,
	}

	self.Attempts = append(self.Attempts, attempt)
	return attempt, nil
}

func (self *EthTx) TxID() string {
	return self.Attempts[len(self.Attempts)-1].TxID
}

func (self *EthTx) GasPrice() *big.Int {
	return self.Attempts[len(self.Attempts)-1].GasPrice
}

func (self *EthTx) Signable(gasPrice *big.Int) *types.Transaction {
	return types.NewTransaction(
		self.Nonce,
		common.HexToAddress(self.To),
		self.Value,
		self.GasLimit,
		gasPrice,
		common.FromHex(self.Data),
	)
}

type EthTxAttempt struct {
	TxID      string `storm:"id,index,unique"`
	EthTxID   uint64 `storm:"index"`
	GasPrice  *big.Int
	Confirmed bool
	Hex       string
	SentAt    uint64
}

