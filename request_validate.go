package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"objects"
)

// Product: validation
func IDValidation(id string) error {
	if id == "" || len(id) <= 0 || len(id) > 16 {
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

// Product: validation
func ProductValidation(product objects.RequestBodyProduct) error {
	if product.Title == "" {
		return errors.New("empty title not allowed")
	}
	if product.Variants[0].Sku == "" {
		return errors.New("empty SKU codes not allowed")
	}
	if product.Variants[0].VariantPricing[0].Name == "" {
		return errors.New("empty price tier name not allowed")
	}
	if product.Variants[0].VariantQuantity[0].Name == "" {
		return errors.New("empty warehouse name not allowed")
	}
	return nil
}

// Product: Duplicate Option validation

// Product: Duplicate SKU validation

// Product: Duplicate Option value validation (variations)

// Product: decodes the request body
func DecodeProductRequestBody(r *http.Request) (objects.RequestBodyProduct, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyProduct{}
	err := decoder.Decode(&params)
	if err != nil {
		return params, err
	}
	return params, nil
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
