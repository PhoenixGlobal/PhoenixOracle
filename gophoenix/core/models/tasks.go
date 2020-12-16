package models

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Task struct {
	ID        string    `storm:"id,index,unique"`
	Schedule string    `json:"schedule"`
	Subtasks []Subtask `json:"subtasks"`
	CreatedAt time.Time `storm:"index"`
}

type Subtask struct {
	Type   string                 `json:"adapterType"`
	Params map[string]interface{} `json:"adapterParams"`
}

func (a *Task) Valid() (bool, error) {
	for _, s := range a.Subtasks {
		if s.Type != "httpJSON" {
			return false, errors.New(`"` + s.Type + `" is not a supported adapter type.`)
		}
	}
	return true, nil
}

func NewTask() Task {
	return Task{ID: uuid.NewV4().String(), CreatedAt: time.Now()}
}
