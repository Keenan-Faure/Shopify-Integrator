package main

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestIgnoreDefaultOption(t *testing.T) {
	// Test 1 - zero valued function param
	result := IgnoreDefaultOption("")
	assert.Equal(t, "", result)

	// Test 2 - valid function param values
	result = IgnoreDefaultOption("MOCK_OPTION_VALUE")
	assert.Equal(t, "MOCK_OPTION_VALUE", result)

	// Test 3 - valid function param values | with special characters
	result = IgnoreDefaultOption("MOCK_OPTION_VALUE\"/pokemon/$easy^")
	assert.Equal(t, "MOCK_OPTION_VALUE\"/pokemon/$easy^", result)
}
