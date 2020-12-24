package services

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/models"
)

func StartJob(run models.JobRun, orm models.ORM) error {
	run.Status = "in progress"
	err := orm.Save(&run)
	if err != nil {
		return err
	}

	var prevRun models.TaskRun
	for i, taskRun := range run.TaskRuns {
		prevRun = startTask(taskRun, prevRun.Result)
		run.TaskRuns[i] = prevRun
		err = orm.Save(&run)
		if err != nil {
			return err
		}

		if prevRun.Result.Error != nil {
			break
		}
	}

	run.Result = prevRun.Result
	if run.Result.Error != nil {
		run.Status = "errored"
	} else {
		run.Status = "completed"
	}

	return orm.Save(&run)
}

func startTask(run models.TaskRun, input adapters.RunResult) models.TaskRun {
	run.Status = "in progress"
	adapter, err := run.Adapter()

	if err != nil {
		run.Status = "errored"
		run.Result.Error = err
		return run
	}
	run.Result = adapter.Perform(input)

	if run.Result.Error != nil {
		run.Status = "errored"
	} else {
		run.Status = "completed"
	}

	return run
}