package main

import (
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/scheduler"
	"PhoenixOracle/gophoenix/core/web"
	"log"
)

func main() {
	models.InitDB()
	defer models.CloseDB()
	sched, err := scheduler.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer sched.Stop()
	r := web.Router()

	log.Fatal(r.Run())
}