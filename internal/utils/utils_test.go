package utils

import (
	"errors"
	"fmt"
	"strings"
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

func TestGetAppSettings(t *testing.T) {
	fmt.Println("Test Case 1 - Returning All Keys for app settings")
	settings_map := GetAppSettings("app")
	if len(settings_map) == 0 {
		t.Errorf("Expected non-zero value but found " + fmt.Sprint(len(settings_map)))
	}
	if settings_map["APP_ENABLE_SHOPIFY_FETCH"] == "" {
		t.Errorf("Expected 'description' value but found " + settings_map["APP_ENABLE_SHOPIFY_FETCH"])
	}
	fmt.Println("Test Case 2 - Returning All Keys for shopify settings")
	shopify_settings_map := GetAppSettings("shopify")
	if len(shopify_settings_map) == 0 {
		t.Errorf("Expected non-zero value but found " + fmt.Sprint(len(shopify_settings_map)))
	}
}

func TestRandomPassword(t *testing.T) {
	fmt.Println("Test Case 1 - Generating a random password of 10 length")
	first_rand_psw := RandStringBytes(10)
	if first_rand_psw == "" || len(first_rand_psw) == 0 {
		t.Errorf("expected 10 but found " + fmt.Sprint(len(first_rand_psw)))
	}
	fmt.Println("Test Case 2 - Generating a random password of 20 length")
	second_rand_psw := RandStringBytes(20)
	if second_rand_psw == "" || len(second_rand_psw) == 0 {
		t.Errorf("expected 20 but found " + fmt.Sprint(len(second_rand_psw)))
	}
	fmt.Println("Test Case 3 - Generating a random password of 10 length and compare")
	third_rand_psw := RandStringBytes(10)
	if third_rand_psw == "" || len(third_rand_psw) == 0 {
		t.Errorf("expected 20 but found " + fmt.Sprint(len(third_rand_psw)))
	}
	if first_rand_psw == third_rand_psw {
		t.Errorf("expected non-equality but found " + fmt.Sprint(first_rand_psw == third_rand_psw))
	}
}

func TestGetNextURL(t *testing.T) {
	fmt.Println("Test Case 1 - String with semicolon and html tags <>")
	strng_case_1 := "<https://keenan-faure.myshopify.com/admin/api/2023-10/products.json?limit=50&page_info=eyJsYXN0X2lkIjo3MDczNTAwNzI1MzA5LCJsYXN0X3ZhbHVlIjoiU2F2aW9yIEZyb20gQW5vdGhlciBXb3JsZCAtIEFsb3kiLCJkaXJlY3Rpb24iOiJuZXh0In0>; rel='next'"
	strng_case_2 := "https://keenan-faure.myshopify.com/admin/api/2023-10/products.json?limit=50&page_info=eyJsYXN0X2lkIjo3MDczNTAwNzI1MzA5LCJsYXN0X3ZhbHVlIjoiU2F2aW9yIEZyb20gQW5vdGhlciBXb3JsZCAtIEFsb3kiLCJkaXJlY3Rpb24iOiJuZXh0In0; rel='next'"
	strng_case_3 := "<https://keenan-faure.myshopify.com/admin/api/2023-10/products.json?limit=50&page_info=eyJsYXN0X2lkIjo3MDczNTAwNzI1MzA5LCJsYXN0X3ZhbHVlIjoiU2F2aW9yIEZyb20gQW5vdGhlciBXb3JsZCAtIEFsb3kiLCJkaXJlY3Rpb24iOiJuZXh0In0>"
	result := GetNextURL(strng_case_1)
	if len(result) == len(strng_case_1) {
		t.Errorf("Expected string length to be reduced, but found same length")
	}
	if strings.Contains(result, "rel='next'") {
		t.Errorf("Expected string rel='next' to be removed from result")
	}
	if strings.Contains(result, "<") {
		t.Errorf("Expected string '<' to be removed from result")
	}
	if strings.Contains(result, ">") {
		t.Errorf("Expected string '>' to be removed from result")
	}
	fmt.Println("Test Case 2 - String with semicolon and without html tags")
	result = GetNextURL(strng_case_2)
	if len(result) == len(strng_case_2) {
		t.Errorf("Expected string length to be reduced, but found same length")
	}
	if strings.Contains(result, "rel='next'") {
		t.Errorf("Expected string rel='next' to be removed from result")
	}
	fmt.Println("Test Case 3 - String without semicolon and with html tags <>")
	result = GetNextURL(strng_case_3)
	if len(result) == len(strng_case_3) {
		t.Errorf("Expected string length to be reduced, but found same length")
	}
	if strings.Contains(result, "<") {
		t.Errorf("Expected string '<' to be removed from result")
	}
	if strings.Contains(result, ">") {
		t.Errorf("Expected string '>' to be removed from result")
	}
}
