package main

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/web"
	"log"
)

func main() {
	models.InitDB()
	defer models.CloseDB()
	r := web.Router()
	log.Fatal(r.Run())
}