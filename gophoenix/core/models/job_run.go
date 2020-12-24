package models

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"time"
)

type JobRun struct {
	ID        string    `storm:"id,index,unique"`
	JobID     string    `storm:"index"`
	Status    string
	CreatedAt time.Time `storm:"index"`
	Result    adapters.RunResult `storm:"inline"`
	TaskRuns  []TaskRun          `storm:"inline"`
}
