package orm

import (
	"PhoenixOracle/gophoenix/core/models"
	"log"
)

func migrate() {
	initializeModel(&models.Job{})
	initializeModel(&models.JobRun{})
}

func initializeModel(klass interface{}) {
	err := GetDB().Init(klass)
	if err != nil {
		log.Fatal(err)
	}
}

