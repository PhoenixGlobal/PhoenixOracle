package main

import (
	"PhoenixOracle/core/orm"
	"PhoenixOracle/core/web"
	"log"
)

func main() {
	orm.Init()
	defer orm.Close()
	r := web.Router()
	log.Fatal(r.Run())
}