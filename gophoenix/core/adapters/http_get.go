package adapters

import (
	"PhoenixOracle/gophoenix/core/models"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpGet struct {
	AdapterBase
	Endpoint string `json:"endpoint"`
}

func (self HttpGet) Perform(input models.RunResult) models.RunResult{
	response, err := http.Get(self.Endpoint)
	fmt.Println("***********************")
	fmt.Println(response)
	fmt.Println(self.Endpoint)
	fmt.Println("***********************")
	if err != nil{
		return models.RunResult{
			Error: err,
		}
	}
	defer response.Body.Close()
	bytes, err:= ioutil.ReadAll(response.Body)
	body := string(bytes)
	if err != nil{
		return models.RunResult{
			Error: err,
		}
	}
	if response.StatusCode >= 300{
		return models.RunResult{
			Error: fmt.Errorf(body),
		}
	}
	fmt.Println("!!!!!!!!!!!!!!!!!!!")
	rs :=  models.RunResultWithValue(body)
	fmt.Println(rs)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!")
	return rs
}

