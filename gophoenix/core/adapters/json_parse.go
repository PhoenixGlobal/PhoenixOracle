package adapters

import (
	"errors"

	simplejson "github.com/bitly/go-simplejson"
)

type JsonParse struct {
	Path []string `json:"path"`
}

func (self *JsonParse) Perform(input RunResult) RunResult {
	js, err := simplejson.NewJson([]byte(input.Value()))
	if err != nil {
		return RunResult{Error: err}
	}

	js, err = checkEarlyPath(js, self.Path)
	if err != nil {
		return RunResult{Error: err}
	}

	rval, _ := js.Get(self.Path[len(self.Path)-1]).String()
	return RunResult{
		Output: map[string]string{"value": rval},
		Error:  nil,
	}
}

func checkEarlyPath(js *simplejson.Json, path []string) (*simplejson.Json, error) {
	var ok bool
	for _, k := range path[:len(path)-1] {
		js, ok = js.CheckGet(k)
		if !ok {
			return js, errors.New("No value could be found for the key '"+ k + "'")
		}
	}
	return js, nil
}