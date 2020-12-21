package models

import (
	"PhoenixOracle/gophoenix/core/models/tasks"
	"encoding/json"
	"fmt"
)

type Adapter interface {
	Perform()
}

type TaskData struct {
	Type   string          `json:"type" storm:"index"`
	Params json.RawMessage `json:"params,omitempty"`
}

type Task struct {
	TaskData
	Adapter
}

func (self *Task) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &self.TaskData)
	if err != nil {
		return err
	}
	self.Adapter, err = self.adapterFromRaw()
	return err
}

func (self *Task) adapterFromRaw() (Adapter, error) {
	switch self.Type {
	case "HttpGet":
		temp := &tasks.HttpGet{}
		err := json.Unmarshal(self.Params, temp)
		return temp, err
	case "JsonParse":
		temp := &tasks.JsonParse{}
		err := json.Unmarshal(self.Params, temp)
		return temp, err
	case "EthBytes32":
		temp := &tasks.EthBytes32{}
		err := json.Unmarshal(self.Params, temp)
		return temp, err
	case "NoOp":
		return &tasks.NoOp{}, nil
	}


	return nil, fmt.Errorf("%s is not a supported adapter type", self.Type)
}