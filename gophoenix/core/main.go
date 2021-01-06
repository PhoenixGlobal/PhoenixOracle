package main

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/store"
	"PhoenixOracle/gophoenix/core/web"
	"log"
)

func main() {
	store := store.New()
	sugarLogger := adapters.GetLogger()
	defer sugarLogger.Sync()

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