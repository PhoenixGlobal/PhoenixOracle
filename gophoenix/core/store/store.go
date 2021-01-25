package store

import (
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/store/models"
	"github.com/ethereum/go-ethereum/rpc"
	"os"
	"os/signal"
	"syscall"
)

type Store struct {
	*models.ORM
	Config    Config
	KeyStore  *KeyStore
	sigs      chan os.Signal
	Exiter    func(int)
	Eth      *Eth
	Tx       *EthTxManager
}

func NewStore(config Config) *Store {
	err := os.MkdirAll(config.RootDir, os.FileMode(0700))
	if err != nil {
		logger.Fatal(err)
	}
	orm := models.NewORM(config.RootDir)
	ethrpc, err := rpc.Dial(config.EthereumURL)
	if err != nil {
		logger.Fatal(err)
	}
	keyStore := NewKeyStore(config.KeysDir())
	eth := &Eth{ethrpc}
	return &Store{
		ORM:       orm,
		Config:    config,
		KeyStore:  keyStore,
		Exiter:    os.Exit,
		Eth: eth,
		Tx:       &EthTxManager{keyStore, eth, config},
	}
}

func (self *Store) Start(){
	self.sigs = make(chan os.Signal, 1)
	signal.Notify(self.sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-self.sigs
		self.Close()
		self.Exiter(1)
	}()
}








