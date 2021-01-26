package store

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts"
	"strconv"
	"strings"
)

type Eth struct {
	Caller
}

type Caller interface {
	Call(result interface{}, method string, args ...interface{}) error
}

func (self *Eth) GetNonce(account accounts.Account) (uint64, error) {
	var result string
	err := self.Call(&result, "eth_getTransactionCount", account.Address.Hex())
	if err != nil {
		return 0, err
	}
	if strings.ToLower(result[0:2]) == "0x" {
		result = result[2:]
	}
	return strconv.ParseUint(result, 16, 64)
}

func (self *Eth) SendRawTx(hex string) (string, error) {
	var result string
	err := self.Call(&result, "eth_sendRawTransaction", hex)
	return result, err
}

func (self *Eth) GetTxReceipt(txid string) (TxReceipt, error) {
	receipt := TxReceipt{}
	err := self.Call(&receipt, "eth_getTransactionReceipt", txid)
	return receipt, err
}

type TxReceipt struct {
	BlockNumber uint64 `json:"blockNumber,string"`
	TXID        string `json:"transactionHash"`
}

func (self *TxReceipt) UnmarshalJSON(b []byte) error {
	type Rcpt struct {
		BlockNumber string `json:"blockNumber"`
		TXID        string `json:"transactionHash"`
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
	self.TXID = rcpt.TXID
	return nil
}

func (self *TxReceipt) Unconfirmed() bool {
	return self.TXID == ""
}