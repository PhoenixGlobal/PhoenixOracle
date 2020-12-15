package main

import (
	"PhoenixOracle/core/web"
	"log"
)

func main() {
	r := web.Router()
	log.Fatal(r.Run())
}