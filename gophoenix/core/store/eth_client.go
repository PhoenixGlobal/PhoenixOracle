package store

import (
	"PhoenixOracle/gophoenix/core/utils"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts"
	"strconv"
)

type EthClient struct {
	Caller
}

type Caller interface {
	Call(result interface{}, method string, args ...interface{}) error
}

func (self *EthClient) GetNonce(account accounts.Account) (uint64, error) {
	var result string
	err := self.Call(&result, "eth_getTransactionCount", account.Address.Hex())
	if err != nil {
		return 0, err
	}
	return utils.HexToUint64(result)
}

func (self *EthClient) SendRawTx(hex string) (string, error) {
	var result string
	err := self.Call(&result, "eth_sendRawTransaction", hex)
	return result, err
}

func (self *EthClient) GetTxReceipt(hash string) (*TxReceipt, error) {
	receipt := TxReceipt{}
	err := self.Call(&receipt, "eth_getTransactionReceipt", hash)
	return &receipt, err
}

func (self *EthClient) BlockNumber() (uint64, error) {
	result := ""
	if err := self.Call(&result, "eth_blockNumber"); err != nil {
		return 0, err
	}
	return utils.HexToUint64(result)
}

type TxReceipt struct {
	BlockNumber uint64 `json:"blockNumber,string"`
	Hash        string `json:"transactionHash"`
}

func (self *TxReceipt) UnmarshalJSON(b []byte) error {
	type Rcpt struct {
		BlockNumber string `json:"blockNumber"`
		Hash        string `json:"transactionHash"`
	}
	var rcpt Rcpt
	if err := json.Unmarshal(b, &rcpt); err != nil {
		return err
	}
	block, err := strconv.ParseUint(rcpt.BlockNumber[2:], 16, 64)
	if err != nil {
		return err
	}
	self.BlockNumber = block
	self.Hash = rcpt.Hash
	return nil
}

func (self *TxReceipt) Unconfirmed() bool {
	return self.Hash == ""
}