package main

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestInitConn(t *testing.T) {
	// Test 1 - invalid connection
	result, err := InitConn("postgres://testuser:testpassword@127.0.0.1:5432/database")
	assert.Equal(t, nil, err)
	assert.Equal(t, result.Valid, false)

	// Cannot test valid connection as I cannot guess what the env file will contain
}

func TestInitCustomConnection(t *testing.T) {
	// Test 1 - invalid connection
	_, err := InitCustomConnection("postgres://testuser:testpassword@127.0.0.1:5432/database")
	assert.NotEqual(t, nil, err)

	// Cannot test valid connection as I cannot guess what the env file will contain
}

func TestInitConnectionString(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)
	useLocalhost, _ := dbconfig.GetFlagValue(HOST_RUNTIME_FLAG_NAME)

	// Test 1 - mock
	result := InitConnectionString(useLocalhost, true)
	assert.Equal(t, "?sslmode=disable", result[len(result)-16:])
}

func TestStoreConfig(t *testing.T) {

}
