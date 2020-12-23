package main

import (
	"PhoenixOracle/gophoenix/core/store"
	"PhoenixOracle/gophoenix/core/web"
	"log"
)

func main() {
	store := store.New()
	r := web.Router(store)
	err := store.Start()
	if err != nil{
		log.Fatal(err)
	}
	defer store.Close()

	log.Fatal(r.Run())
}