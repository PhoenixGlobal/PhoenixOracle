package store

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type KeyStore struct {
	*keystore.KeyStore
}

func NewKeyStore(keyDir string) *KeyStore {
	ks := keystore.NewKeyStore(
		keyDir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	return &KeyStore{ks}
}

func (self *KeyStore) HasAccounts() bool {
	return len(self.Accounts()) > 0
}

func (self *KeyStore) Unlock(phrase string) error {
	for _, account := range self.Accounts() {
		err := self.KeyStore.Unlock(account, phrase)
		if err != nil {
			return fmt.Errorf("Invalid password for account: %s\n\nPlease try again...\n", account.Address.Hex())
		}
	}
	return nil
}

func (self *KeyStore) SignTx(tx *types.Transaction, chainID int64) (*types.Transaction, error) {
	acc := self.Accounts()[0]
	return self.KeyStore.SignTx(acc, tx, big.NewInt(chainID))
}