package main

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const MOCK_EMAIL = "test@tester.com"
const MOCK_NAME = "testertests"

func TestEmail(t *testing.T) {
	// Test Case 1 - invalid email
	err := Email(uuid.New(), false, "", MOCK_NAME)
	assert.NotEqual(t, nil, err)

	// Test Case 2 - invalid name
	err = Email(uuid.New(), false, MOCK_EMAIL, "")
	assert.NotEqual(t, nil, err)

	// Test Case 3 - invalid token
	err = Email(uuid.Nil, false, MOCK_EMAIL, MOCK_NAME)
	assert.NotEqual(t, nil, err)

	// Test Case 4 - valid email
	err = Email(uuid.New(), false, MOCK_EMAIL, MOCK_NAME)
	assert.Equal(t, err, nil)
}

func TestEmailValidation(t *testing.T) {
	// Test Case 1 - invalid email
	valid, err := emailValidation(uuid.New(), "", MOCK_NAME)
	assert.NotEqual(t, nil, err)
	assert.NotEqual(t, true, valid)

	// Test Case 2 - invalid name
	valid, err = emailValidation(uuid.New(), MOCK_EMAIL, "")
	assert.NotEqual(t, nil, err)
	assert.NotEqual(t, true, valid)

	// Test Case 3 - invalid token
	valid, err = emailValidation(uuid.Nil, MOCK_EMAIL, MOCK_NAME)
	assert.NotEqual(t, nil, err)
	assert.NotEqual(t, true, valid)

	// Test Case 4 - valid parameters
	valid, err = emailValidation(uuid.New(), MOCK_EMAIL, MOCK_NAME)
	assert.Equal(t, err, nil)
	assert.Equal(t, valid, true)
}
