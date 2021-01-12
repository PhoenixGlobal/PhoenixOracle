package main

import (
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/services"
	"PhoenixOracle/gophoenix/core/web"
	"log"
)

func main() {
	config := services.NewConfig()
	logger.SetLoggerDir(config.RootDir)
	store := services.NewStore(config)
	sugarLogger := logger.GetLogger()
	defer sugarLogger.Sync()

	services.Authenticate(store)
	r := web.Router(store)
	err := store.Start()
	if err != nil{
		sugarLogger.Error(err)
		log.Fatal(err)
	}

	sugarLogger.Info("this is main entry")
	defer store.Close()
	log.Fatal(r.Run())
}