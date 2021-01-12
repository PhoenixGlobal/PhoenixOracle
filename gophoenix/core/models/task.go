package models

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"encoding/json"
	"fmt"
)

type Task struct {
	Type   string          `json:"type" storm:"index"`
	Params json.RawMessage `json:"params,omitempty"`
}

type TaskRun struct {
	ID string `storm: "id"`
	Task
	Status string
	Result adapters.RunResult
}

func (self Task) Validate() error {
	_, err := self.Adapter()
	return err
}

func (self Task) Adapter() (adapters.Adapter, error) {
	switch self.Type {
	case "HttpGet":
		temp := &adapters.HttpGet{}
		err := json.Unmarshal(self.Params, temp)
		return temp, err
	case "JsonParse":
		temp := &adapters.JsonParse{}
		err := json.Unmarshal(self.Params, temp)
		return temp, err
	case "EthBytes32":
		temp := &adapters.EthBytes32{}
		err := unmarshalOrEmpty(self.Params, temp)
		return temp, err
	case "EthSendTx":
		temp := &adapters.EthSendTx{}
		err := json.Unmarshal(self.Params, temp)
		return temp, err
	case "NoOp":
		return &adapters.NoOp{}, nil
	}


	return nil, fmt.Errorf("%s is not a supported adapters type", self.Type)
}

func unmarshalOrEmpty(params json.RawMessage, dst interface{}) error {
	if len(params) > 0 {
		return json.Unmarshal(params, dst)
	}
	return nil
}
