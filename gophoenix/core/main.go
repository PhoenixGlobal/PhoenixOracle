package main

import (
	"PhoenixOracle/gophoenix/core/orm"
	"PhoenixOracle/gophoenix/core/web"
	"log"
)

func main() {
	orm.Init()
	defer orm.Close()
	r := web.Router()
	log.Fatal(r.Run())
}