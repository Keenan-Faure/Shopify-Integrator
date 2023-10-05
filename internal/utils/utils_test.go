package utils

import (
	"errors"
	"fmt"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	fmt.Println("Test Case 1 - Key exists")
	key := "port"
	value := LoadEnv(key)
	if value != "8080" {
		t.Errorf("Expected '8080' but found" + value)
	}
	fmt.Println("Test Case 2 - Key does not exist")
	key = "porters"
	value = LoadEnv(key)
	if value != "" {
		t.Errorf("Expected '' but found" + value)
	}
}

func TestExtractApiKey(t *testing.T) {
	fmt.Println("Test Case 1 - Empty Header")
	expected := ""
	actual, err := ExtractAPIKey("")
	if err == nil {
		t.Errorf("Expected error but found " + actual)
	}
	fmt.Println("Test Case 2 - ApiKey exists in header")
	expected = "ApiKey erbaj7ia8nasdasd7"
	actual, err = ExtractAPIKey(expected)
	if err != nil {
		t.Errorf("Expected 'nil' but found " + actual)
	}
	fmt.Println("Test Case 3 - Malformed ApiKey")
	expected = "ApiKe"
	_, err = ExtractAPIKey(expected)
	if err == nil {
		t.Errorf("Expected 'error' but found " + err.Error())
	}
	fmt.Println("Test Case 4 - Malformed 2nd part of API")
	_ = "KeyApi 8nasdasd7"
	_, err = ExtractAPIKey(expected)
	if err == nil {
		t.Errorf("Expected 'error' but found " + err.Error())
	}
}

func TestConvertStringToLike(t *testing.T) {
	fmt.Println("Test case 1 - Valid string")
	arg := "string"
	results := ConvertStringToLike(arg)
	fmt.Println(results[len(arg)+1:])
	if results[0:1] != "%" || results[len(arg)+1:] != "%" {
		t.Errorf("Unexpected result")
	}
	fmt.Println("Test case 2 - Invalid string")
	arg = ""
	results = ConvertStringToLike(arg)
	if results[0:1] != "%" || results[len(arg)+1:] != "%" {
		t.Errorf("Unexpected result")
	}
}

func TestConfirmFilters(t *testing.T) {
	fmt.Println("Test case 1 - Valid filter")
	arg := "string"
	results := ConfirmFilters(arg)
	if results != arg {
		t.Errorf("Unexpected result")
	}
	fmt.Println("Test case 2 - Invalid filter")
	arg = ""
	results = ConfirmFilters(arg)
	if results != arg {
		t.Errorf("Unexpected result")
	}
}

func TestConvertStringToSQL(t *testing.T) {
	fmt.Println("Test case 1 - Valid string")
	arg := "string"
	results := ConvertStringToSQL(arg)
	if !results.Valid {
		t.Errorf("Expected 'true' but found 'false'")
	}
	fmt.Println("Test case 2 - Invalid (empty) string")
	arg = ""
	results = ConvertStringToSQL(arg)
	if results.Valid {
		t.Errorf("Expected 'false' but found 'true")
	}
}

func TestConvertIntToSQL(t *testing.T) {
	fmt.Println("Test case 1 - Valid Integer")
	arg := 531
	results := ConvertIntToSQL(arg)
	if !results.Valid {
		t.Errorf("Expected 'true' but found 'false'")
	}
	fmt.Println("Test case 2 - Invalid (nil value) Integer")
	arg = 0
	results = ConvertIntToSQL(arg)
	if !results.Valid {
		t.Errorf("Expected 'false' but found 'true")
	}
}

func TestConfirmError(t *testing.T) {
	fmt.Println("Test case 1 - Valid Duplicate Error")
	err := errors.New("pq: duplicate key value violates unique constraint")
	results := ConfirmError(err)
	if results != "duplicate fields not allowed" {
		t.Errorf("Unexpected results, expected 'duplicate fields not allowed' but found " + results)
	}
	fmt.Println("Test case 2 - None Duplicate Error")
	err = errors.New("Invalid database credentials")
	results = ConfirmError(err)
	if results == "duplicate fields not allowed" {
		t.Errorf("Unexpected results, expected 'Invalid database credentials' but found " + results)
	}
}

func TestExtractVID(t *testing.T) {
	fmt.Println("Test Case 1 - Valid VID")
	variable := ExtractVID("gid://shopify/ProductVariant/40466067357761")
	if variable != "40466067357761" {
		t.Errorf("Expected '40466067357761', but found " + variable)
	}

	fmt.Println("Test Case 2 - Short VID")
	variable = ExtractVID("gid://shopify/Produc")
	if variable != "" {
		t.Errorf("Expected '', but found " + variable)
	}

	fmt.Println("Test Case 3 - Invalid VID")
	variable = ExtractVID("")
	if variable != "" {
		t.Errorf("Expected '', but found " + variable)
	}
}

func TestExtractPID(t *testing.T) {
	fmt.Println("Test Case 1 - Valid PID")
	variable := ExtractPID("gid://shopify/Product/6971324465217")
	if variable != "6971324465217" {
		t.Errorf("Expected '6971324465217', but found " + variable)
	}

	fmt.Println("Test Case 2 - Short PID")
	variable = ExtractPID("gid://shopify/Produ")
	if variable != "" {
		t.Errorf("Expected '', but found " + variable)
	}

	fmt.Println("Test Case 1 - Invalid PID")
	variable = ExtractPID("")
	if variable != "" {
		t.Errorf("Expected '', but found " + variable)
	}
}
