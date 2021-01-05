package test

import (
	"PhoenixOracle/gophoenix/core/adapters"
	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
	"testing"
)

func TestHttpGetNotAUrlError(t *testing.T) {
	httpGet := adapters.HttpGet{Endpoint: "NotAUrl"}
	input := adapters.RunResult{}
	result := httpGet.Perform(input)
	assert.Nil(t, result.Output)
	assert.NotNil(t, result.Error)
}

func TestHttpGetResponseError(t *testing.T) {
	defer gock.Off()
	url := `https://example.com/api`

	gock.New(url).
		Get("").
		Reply(400).
		JSON(`Invalid request`)

	httpGet := adapters.HttpGet{Endpoint: url}
	input := adapters.RunResult{}
	result := httpGet.Perform(input)
	assert.Nil(t, result.Output)
	assert.NotNil(t, result.Error)
}
