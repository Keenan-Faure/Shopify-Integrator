package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestrictionValidation(t *testing.T) {
	// Test Case 1 - invalid push restrictions
	pushRestrictions := PushRestrictionPayload("test-case-invalid-request.json")
	valid := RestrictionValidation(pushRestrictions)
	assert.NotEqual(t, nil, valid)

	// Test Case 2 - valid push restrictions
	pushRestrictions = PushRestrictionPayload("test-case-valid-request.json")
	valid = RestrictionValidation(pushRestrictions)
	assert.Equal(t, nil, valid)

	// Test Case 3 - invalid push restrictions
	fetchRestrictions := FetchRestrictionPayload("test-case-invalid-request.json")
	valid = RestrictionValidation(fetchRestrictions)
	assert.NotEqual(t, nil, valid)

	// Test Case 4 - valid push restrictions
	fetchRestrictions = FetchRestrictionPayload("test-case-valid-request.json")
	valid = RestrictionValidation(fetchRestrictions)
	assert.Equal(t, nil, valid)
}

func TestDecodeRestriction(t *testing.T) {

}
