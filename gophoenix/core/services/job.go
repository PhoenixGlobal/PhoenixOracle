package services

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/models"
	"PhoenixOracle/gophoenix/core/store"
	"fmt"
)

func StartJob(run models.JobRun, store *store.Store) error {
	run.Status = "in progress"
	if err := store.Save(&run); err != nil {
		return runJobError(run, err)
	}

	logger.GetLogger().Infow("starting job", run.ForLogger()...)
	var prevRun models.TaskRun
	for i, taskRun := range run.TaskRuns {
		prevRun = startTask(taskRun, prevRun.Result, store)
		run.TaskRuns[i] = prevRun

		err:= store.Save(&run); if err != nil {
			return runJobError(run, err)
		}

		if prevRun.Result.HasError() {
			break
		}
	}

	run.Result = prevRun.Result
	if run.Result.HasError() {
		run.Status = "errored"
	} else {
		run.Status = "completed"
	}

	return runJobError(run, store.Save(&run))
}

func startTask(run models.TaskRun, input models.RunResult,store *store.Store) models.TaskRun {
	run.Status = "in progress"
	adapter, err := adapters.For(run.Task, store)

	if err != nil {
		run.Status = "errored"
		run.Result.SetError(err)
		return run
	}
	run.Result = adapter.Perform(input)

	if run.Result.HasError() {
		run.Status = "errored"
	} else {
		run.Status = "completed"
	}

	return run
}

func runJobError(run models.JobRun, err error) error {
	if err != nil {
		return fmt.Errorf("StartJob#%v: %v", run.JobID, err)
	}
	return nil
}