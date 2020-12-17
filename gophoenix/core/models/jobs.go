package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Job struct {
	ID        string    `storm:"id,index,unique"`
	Schedule string     `json:"schedule"`
	Tasks []Task        `json:"tasks" storm:"inline"`
	CreatedAt time.Time `storm:"index"`
}

func NewJob() Job {
	return Job{ID: uuid.NewV4().String(), CreatedAt: time.Now()}
}
