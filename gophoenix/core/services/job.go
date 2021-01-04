package services

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"PhoenixOracle/gophoenix/core/models"
	"fmt"
	"time"
)

func StartJob(run models.JobRun, orm models.ORM) error {
	run.Status = "in progress"
	err := orm.Save(&run)
	if err != nil {
		return runJobError(run, err)
	}

	GetLogger().Infow("starting job", run.ForLogger()...)
	var prevRun models.TaskRun
	for i, taskRun := range run.TaskRuns {
		prevRun = startTask(taskRun, prevRun.Result)
		fmt.Println("333333333333333333333")
		fmt.Println(prevRun.Result.Output)
		run.TaskRuns[i] = prevRun
		fmt.Println("before ......................")
		fmt.Println(orm)
		fmt.Println(time.Now())

		err1 := orm.Save(&run)
		fmt.Println("after ......................")
		if err1 != nil {
			fmt.Println("err hafdahdfhsahfhdsdh")
			return runJobError(run, err)
		}

		if prevRun.Result.Error != nil {
			fmt.Println("err hafdahdfhsahfhdsdh")
			break
		}
	}

	run.Result = prevRun.Result
	if run.Result.Error != nil {
		run.Status = "errored"
	} else {
		run.Status = "completed"
	}
	fmt.Println("444444444444444444")
	fmt.Println(run.TaskRuns[0].Result)

	return runJobError(run, orm.Save(&run))
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
	fmt.Println("2222222222222222222222")
	fmt.Println(run.Result.Output)

	return run
}

func runJobError(run models.JobRun, err error) error {
	if err != nil {
		return fmt.Errorf("StartJob#%v: %v", run.JobID, err)
	}
	return nil
}