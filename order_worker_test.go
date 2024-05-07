package main

import (
	"objects"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddOrder(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid order body
	result := dbconfig.AddOrder(objects.RequestBodyOrder{})
	assert.NotEqual(t, result, nil)
	assert.Equal(t, "data validation error", result.Error())

	// Test 2 - valid order request body
	result = dbconfig.AddOrder(OrderPayload("test-case-valid-order.json"))
	assert.Equal(t, result, nil)
	ClearOrderTestData(&dbconfig)
}

func TestUpdateOrder(t *testing.T) {
	dbconfig := setupDatabase("", "", "", false)

	// Test 1 - invalid order body
	result := dbconfig.UpdateOrder(objects.RequestBodyOrder{})
	assert.NotEqual(t, result, nil)

	// Test 2 - valid request body | order ID do not exist
	ClearOrderTestData(&dbconfig)
	result = dbconfig.UpdateOrder(OrderPayload("test-case-valid-order.json"))
	assert.NotEqual(t, result, nil)

	// Test 3 - valid request body | order exist
	createDatabaseOrder(&dbconfig)
	result = dbconfig.UpdateOrder(OrderPayload("test-case-valid-order.json"))
	assert.Equal(t, result, nil)
	ClearOrderTestData(&dbconfig)
}
