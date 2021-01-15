package store

import (
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/models"
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
}

func NewStore(config Config) *Store {
	err := os.MkdirAll(config.RootDir, os.FileMode(0700))
	if err != nil {
		logger.Fatal(err)
	}
	orm := models.NewORM(config.RootDir)
	return &Store{
		ORM:       orm,
		Config:    config,
		KeyStore:  NewKeyStore(config.KeysDir()),
		Exiter:    os.Exit,
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








