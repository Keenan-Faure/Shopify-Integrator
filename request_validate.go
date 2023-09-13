package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"objects"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

// ValidateToken: Data validtion
func ValidateTokenValidation(token_request objects.RequestBodyUser) error {
	if token_request.Name == "" || len(token_request.Name) == 0 {
		return errors.New("data validation error")
	} else if token_request.Email == "" || len(token_request.Email) == 0 {
		return errors.New("data validation error")
	}
	_, err := uuid.Parse(token_request.Token)
	if err != nil {
		return err
	}
	return nil
}

// ValidateToken: decode the request body
func DecodeValidateTokenRequestBody(r *http.Request) (objects.RequestBodyUser, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyUser{}
	err := decoder.Decode(&params)
	if err != nil {
		return params, err
	}
	return params, nil
}

// PreRegister: Data validation
func PreRegisterValidation(preorder objects.RequestBodyPreRegister) error {
	if preorder.Name == "" || len(preorder.Name) == 0 || preorder.Email == "" || len(preorder.Email) == 0 {
		return errors.New("data validation error")
	}
	return nil
}

// PreRegister: decode the request body
func DecodePreRegisterRequestBody(r *http.Request) (objects.RequestBodyPreRegister, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyPreRegister{}
	err := decoder.Decode(&params)
	if err != nil {
		return params, err
	}
	return params, nil
}

// Customer: Data validation
func CustomerValidation(order objects.RequestBodyCustomer) error {
	if order.FirstName == "" {
		return errors.New("data validation error")
	}
	return nil
}

// Customer: decode the request body
func DecodeCustomerRequestBody(r *http.Request) (objects.RequestBodyCustomer, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyCustomer{}
	err := decoder.Decode(&params)
	if err != nil {
		return params, err
	}
	return params, nil
}

// Order: decodes the request body
func DecodeOrderRequestBody(r *http.Request) (objects.RequestBodyOrder, error) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyOrder{}
	err := decoder.Decode(&params)
	if err != nil {
		return params, err
	}
	return params, nil
}

// Order: data validation
func OrderValidation(order objects.RequestBodyOrder) error {
	if order.Name == "" || order.LineItems[0].Sku == "" || order.Customer.FirstName == "" {
		return errors.New("data validation error")
	}
	return nil
}

// User: data validation
func TokenValidation(key string) error {
	if key == "" || len(key) <= 0 || len(key) > 32 {
		return errors.New("invalid product id")
	}
	return nil
}

// Product: data validation
func IDValidation(id string) error {
	if id == "" || len(id) <= 0 || len(id) > 16 {
		return errors.New("invalid product id")
	}
	return nil
}

// User: data validation
func UserValidation(user objects.RequestBodyUser) error {
	if user.Name == "" {
		return errors.New("empty name not allowed")
	}
	return nil
}

// Product: data validation
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
func ValidateDuplicateOption(product objects.RequestBodyProduct) error {
	options_names := []string{}
	if len(product.ProductOptions) > 1 {
		for _, value := range product.ProductOptions {
			if value.Value != "" && len(value.Value) > 0 {
				if slices.Contains(options_names, value.Value) {
					return errors.New("duplicate options not allowed: " + value.Value)
				}
				options_names = append(options_names, value.Value)
			}
		}
	}
	return nil
}

// Product: Duplicate SKU validation
func ValidateDuplicateSKU(
	product objects.RequestBodyProduct,
	dbconfig *DbConfig,
	r *http.Request) error {
	sku_array := []string{}
	for _, value := range product.Variants {
		if slices.Contains(sku_array, value.Sku) {
			return errors.New("duplicate SKUs not allowed: " + value.Sku)
		}
		sku_array = append(sku_array, value.Sku)
	}
	for _, value := range sku_array {
		db_sku, err := dbconfig.DB.GetVariantBySKU(r.Context(), value)
		if err != nil {
			return err
		}
		if err.Error() == "record not found" {
			return nil
		}
		if db_sku.Sku == value {
			return errors.New("SKU with code " + value + " already exists")
		}
	}
	return nil
}

// Product: Duplicate Option value validation (variations)
func DuplicateOptionValues(product objects.RequestBodyProduct) error {
	if len(product.ProductOptions) == 1 {
		option_values := []string{}
		for _, value := range product.Variants {
			if slices.Contains(option_values, value.Option1) {
				return errors.New("duplicate option value")
			}
			option_values = append(option_values, value.Option1)
		}
	} else if len(product.ProductOptions) == 2 {
		option_1_values := []string{}
		option_2_values := []string{}
		for _, value := range product.Variants {
			option_1_values = append(option_1_values, value.Option1)
			option_2_values = append(option_2_values, value.Option2)
		}
		counter := 0
		for key := range option_1_values {
			for sub_key := range option_2_values {
				if option_2_values[key] == option_2_values[sub_key] && option_1_values[key] == option_1_values[sub_key] {
					counter += 1
				}
				if counter > 1 {
					return errors.New("duplicate option values not allowed")
				}
			}
		}
	} else if len(product.ProductOptions) != 3 {
		return errors.New("too many option values")
	}

	option_1_values := []string{}
	option_2_values := []string{}
	option_3_values := []string{}
	for _, value := range product.Variants {
		option_1_values = append(option_1_values, value.Option1)
		option_2_values = append(option_2_values, value.Option2)
		option_3_values = append(option_3_values, value.Option3)
	}
	counter := 0
	for key := range option_1_values {
		for sub_key := range option_2_values {
			for primal_key := range option_3_values {
				if (option_3_values[key] == option_3_values[primal_key] &&
					option_2_values[key] == option_2_values[sub_key]) &&
					option_1_values[key] == option_1_values[sub_key] {
					counter += 1
				}
				if counter > 1 {
					return errors.New("duplicate option values not allowed")
				}
			}
		}
	}
	return nil
}

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
