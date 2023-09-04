package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"objects"
)

// Product: validation
func ProductIDValidation(id string) error {
	if(id == "" || len(id) <= 0 || len(id) > 16) {
		return errors.New("Invalid product id")
	}
	return nil
}

// User: validation
func UserValidation(user objects.RequestBodyUser) error {
	if user.Name == "" {
		return errors.New("empty name not allowed")
	}
	return nil
}

// User: decodes the request body
func DecodeUserRequestBody(r *http.Request) (objects.RequestBodyUser, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyUser{}
	err := decoder.Decode(&params)
	if err != nil {
		return params, err
	}
	return params, nil
}
