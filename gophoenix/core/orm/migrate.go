package orm

import (
	"PhoenixOracle/gophoenix/core/models"
	"log"
)

func migrate() {
	err := GetDB().Init(&models.Job{})
	if err != nil {
		log.Fatal(err)
	}
}

