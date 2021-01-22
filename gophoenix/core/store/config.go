package store

import (
	"github.com/caarlos0/env"
	homedir "github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path"
)

type Config struct {
	RootDir string               `env:"ROOT" envDefault:"~/.phoenix"`
	BasicAuthUsername string	 `env:"USERNAME" envDefault:"phoenix"`
	BasicAuthPassword string	 `env:"PASSWORD" envDefault:"p@ssword"`
	EthereumURL       string `env:"ETHEREUM_URL" envDefault:"http://localhost:8545"`
	ChainID           int64  `env:"ETHEREUM_CHAIN_ID" envDefault:0`
	PollingSchedule   string `env:"POLLING_SCHEDULE" envDefault:"* * * * * *"`
}

func NewConfig() Config {
	config := Config{}
	env.Parse(&config)
	dir, err := homedir.Expand(config.RootDir)
	if err != nil {
		log.Fatal(err)
	}
	if err = os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		log.Fatal(err)
	}
	config.RootDir = dir
	return config
}

func (self Config) KeysDir() string {
	return path.Join(self.RootDir, "keys")
}