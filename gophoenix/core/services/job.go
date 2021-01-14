package services

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/config"
	"PhoenixOracle/gophoenix/core/logger"
	"PhoenixOracle/gophoenix/core/models"
	"fmt"
)

func StartJob(run models.JobRun, orm *models.ORM, cf config.Config) error {
	run.Status = "in progress"
	err := orm.Save(&run)
	if err != nil {
		return runJobError(run, err)
	}

	logger.GetLogger().Infow("starting job", run.ForLogger()...)
	var prevRun models.TaskRun
	for i, taskRun := range run.TaskRuns {
		prevRun = startTask(taskRun, prevRun.Result, cf)
		run.TaskRuns[i] = prevRun

		err1 := orm.Save(&run)
		if err1 != nil {
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

	return runJobError(run, orm.Save(&run))
}

func startTask(run models.TaskRun, input models.RunResult,cf config.Config) models.TaskRun {
	run.Status = "in progress"
	adapter, err := adapters.For(run.Task,cf)

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