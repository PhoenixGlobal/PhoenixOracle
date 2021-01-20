package store

import (
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