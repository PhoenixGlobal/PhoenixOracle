package adapters

import (
	"PhoenixOracle/gophoenix/core/models"
)

type NoOp struct {
	AdapterBase
}

func (self *NoOp) Perform(input models.RunResult) models.RunResult {
	return models.RunResult{}
}