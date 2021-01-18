package adapters

import (
	"PhoenixOracle/gophoenix/core/store/models"
)

type NoOp struct {
	AdapterBase
}

func (self *NoOp) Perform(input models.RunResult) models.RunResult {
	return models.RunResult{}
}