package main

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
)

const MOCK_RUNTIME_FLAG_NAME = "test_flag"

func TestSaveFlags(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1
	err := dbconfig.AddRuntimeFlags(MOCK_RUNTIME_FLAG_NAME, false)
	defer ClearRuntimeFlagData(&dbconfig)

	assert.Equal(t, nil, err)
}

func TestGetFlagValue(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	ClearRuntimeFlagData(&dbconfig)

	// Test 1 - get invalid flag
	result, err := dbconfig.GetFlagValue(MOCK_RUNTIME_FLAG_NAME)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, false, result)

	// Test 2 - get valid flag
	dbconfig.AddRuntimeFlags(MOCK_RUNTIME_FLAG_NAME, true)
	defer ClearRuntimeFlagData(&dbconfig)
	result, err = dbconfig.GetFlagValue(MOCK_RUNTIME_FLAG_NAME)

	assert.Equal(t, nil, err)
	assert.Equal(t, true, result)
}

func ClearRuntimeFlagData(dbconfig *DbConfig) {
	dbconfig.DB.RemoveRuntimeFlag(context.Background(), MOCK_RUNTIME_FLAG_NAME)
}
