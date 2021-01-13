package test

import (
	"PhoenixOracle/gophoenix/core/models"
	"testing"

	"PhoenixOracle/gophoenix/core/adapters"
	"github.com/stretchr/testify/assert"
)

func TestParseExistingPath(t *testing.T) {
	t.Parallel()
	input := models.RunResultWithValue(`{"high": "11850.00", "last": "11779.99", "timestamp": "1512487535", "bid": "11779.89", "vwap": "11525.17", "volume": "12916.67066094", "low": "11100.00", "ask": "11779.99", "open": 11613.07}`)

	adapter := adapters.JsonParse{Path:[]string{"last"}}
	result := adapter.Perform(input)
	assert.Equal(t, "11779.99", result.Value())
	assert.Nil(t, result.Error)
}

func TestParseNonExistingPath(t *testing.T) {
	t.Parallel()
	input := models.RunResultWithValue(
		`{"high": "11850.00", "last": "11779.99", "timestamp": "1512487535", "bid": "11779.89", "vwap": "11525.17", "volume": "12916.67066094", "low": "11100.00", "ask": "11779.99", "open": 11613.07}`)

	adapter := adapters.JsonParse{Path:[]string{"doesnotexist"}}
	result := adapter.Perform(input)
	assert.Equal(t, true, result.NullValue())
	assert.Nil(t, result.Error)

	adapter = adapters.JsonParse{Path:[]string{"doesnotexist", "noreally"}}
	result = adapter.Perform(input)
	assert.Equal(t, "", result.Value())
	assert.NotNil(t, result.Error)
}

func TestParseNullValue(t *testing.T) {
	t.Parallel()
	input := models.RunResult{}

	adapter := adapters.JsonParse{Path:[]string{"last"}}
	result := adapter.Perform(input)
	assert.Equal(t, true, result.NullValue())
	assert.NotNil(t, result.Error)
}
