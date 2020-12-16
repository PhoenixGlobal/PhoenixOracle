package models

import (
	"PhoenixOracle/gophoenix/core/models/tasks"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Job struct {
	ID        string    `storm:"id,index,unique"`
	Schedule string    `json:"schedule"`
	Tasks []  tasks.Task `json:"tasks" storm:"inline"`
	CreatedAt time.Time `storm:"index"`
}

type Subtask struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}


func NewTask() Job {
	return Job{ID: uuid.NewV4().String(), CreatedAt: time.Now()}
}

func isValidTask(t Subtask) bool {
	switch t.Type {
	case "HttpGet":
		return true
	}
	return false
}
