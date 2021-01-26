package store

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Eth struct {
	*EthClient
	KeyStore *KeyStore
	Config   Config
}

func (self *Eth) CreateTx(to, data string) (*types.Transaction, error) {
	tx, err := self.NewSignedTx(to, data)
	if err != nil {
		return tx, err
	}
	hex, err := encodeTxToHex(tx)
	if err != nil {
		return tx, err
	}
	if _, err = self.SendRawTx(hex); err != nil {
		return tx, err
	}
	return tx, nil
}

func (self *Eth) NewSignedTx(to, data string) (*types.Transaction, error) {
	account := self.KeyStore.GetAccount()
	nonce, err := self.GetNonce(account)
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

func (self *Eth) TxConfirmed(txid string) (bool, error) {
	receipt, err := self.GetTxReceipt(txid)
	if err != nil {
		return false, err
	} else if receipt.Unconfirmed() {
		return false, nil
	}

	min := receipt.BlockNumber + self.Config.EthConfMin
	current, err := self.BlockNumber()
	if err != nil {
		return false, err
	}
	return (min <= current), nil
}

func encodeTxToHex(tx *types.Transaction) (string, error) {
	rlp := new(bytes.Buffer)
	if err := tx.EncodeRLP(rlp); err != nil {
		return "", err
	}
	return common.Bytes2Hex(rlp.Bytes()), nil
}