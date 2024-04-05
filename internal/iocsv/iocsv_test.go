package iocsv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"objects"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

const MOCK_REQUEST_URL = "http://mock.localhost:9876"
const MOCK_REQUEST_METHOD = http.MethodPost
const MOCK_REQUEST_FORM_KEY = "file"

func TestUploadFile(t *testing.T) {
	// Test Case 1 - invalid content-type
	mpw, multiPartFormData := CreateMultiPartFormData("test-case-invalid-content-type.csv", MOCK_REQUEST_FORM_KEY)
	requestHeaders := make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType()}
	httpRequest := Init(MOCK_REQUEST_URL, MOCK_REQUEST_METHOD, requestHeaders, multiPartFormData)
	fileName, err := UploadFile(httpRequest)
	RemoveFile(fileName)
	if err == nil {
		t.Errorf("Expected 'only CSV extensions are supported' but found: nil")
	}
	if err.Error() != "only CSV extensions are supported" {
		t.Errorf("Expected 'only CSV extensions are supported' but found: " + err.Error())
	}
	if fileName != "" {
		t.Errorf("Expected '' but found: " + fileName)
	}

	// Test Case 2 - invalid file-type
	mpw, multiPartFormData = CreateMultiPartFormData("test-case-invalid-file-type.json", MOCK_REQUEST_FORM_KEY)
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType(), "text/csv"}
	httpRequest = Init(MOCK_REQUEST_URL, MOCK_REQUEST_METHOD, requestHeaders, multiPartFormData)
	fileName, err = UploadFile(httpRequest)
	RemoveFile(fileName)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	if fileName == "" {
		t.Errorf("Expected 'upload-***' but found: " + fileName)
	}

	// Test Case 3 - invalid headers
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{"text/csv"}
	httpRequest = Init(MOCK_REQUEST_URL, MOCK_REQUEST_METHOD, requestHeaders, multiPartFormData)
	fileName, err = UploadFile(httpRequest)
	RemoveFile(fileName)
	if err == nil {
		t.Errorf("Expected 'request Content-Type isn't multipart/form-data' but found: 'nil'")
	}
	if err.Error() != "request Content-Type isn't multipart/form-data" {
		t.Errorf("Expected 'request Content-Type isn't multipart/form-data' but found: " + err.Error())
	}
	if fileName != "" {
		t.Errorf("Expected '' but found: " + fileName)
	}

	// Test Case 4 - invalid formKey used in multipart/form
	mpw, multiPartFormData = CreateMultiPartFormData("test-case-valid-data.csv", "ABC")
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType(), "text/csv"}
	httpRequest = Init(MOCK_REQUEST_URL, MOCK_REQUEST_METHOD, requestHeaders, multiPartFormData)
	fileName, err = UploadFile(httpRequest)
	RemoveFile(fileName)
	if err == nil {
		t.Errorf("Expected 'http: no such file' but found: 'nil'")
	}
	if err.Error() != "http: no such file" {
		t.Errorf("Expected 'http: no such file' but found: " + err.Error())
	}
	if fileName != "" {
		t.Errorf("Expected '' but found: " + fileName)
	}

	// Test Case 5 - valid http request
	mpw, multiPartFormData = CreateMultiPartFormData("test-case-valid-data.csv", MOCK_REQUEST_FORM_KEY)
	requestHeaders = make(map[string][]string)
	requestHeaders["Content-Type"] = []string{mpw.FormDataContentType(), "text/csv"}
	httpRequest = Init(MOCK_REQUEST_URL, MOCK_REQUEST_METHOD, requestHeaders, multiPartFormData)
	fileName, err = UploadFile(httpRequest)
	RemoveFile(fileName)
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	if fileName == "" {
		t.Errorf("Expected 'upload-***' but found: ''")
	}
}

func TestCSVProductHeaders(t *testing.T) {
	// Test Case 1 - invalid (empty) struct
	product := objects.Product{}
	result := CSVProductHeaders(product)
	if len(result) != 16 {
		t.Errorf("Expected '16' but found: " + fmt.Sprint(len(result)))
	}
	assert.Equal(t, "id", result[0])
	assert.Equal(t, "product_code", result[1])
	assert.Equal(t, "active", result[2])

	// Test Case 2- valid struct
	product = ProductPayload("test-case-valid-product.json")
	result = CSVProductHeaders(product)
	if len(result) != 16 {
		t.Errorf("Expected '16' but found: " + fmt.Sprint(len(result)))
	}
	assert.Equal(t, "vendor", result[6])
	assert.Equal(t, "product_type", result[7])
	assert.Equal(t, "sku", result[8])
}

func TestCSVProductValuesByVariant(t *testing.T) {
	// Test Case 1 - valid request parameters | no pricing | no qty | no images | empty product
	result := CSVProductValuesByVariant(objects.Product{}, objects.ProductVariant{}, 0, 0, 0)

	assert.Equal(t, result[1], "")
	assert.Equal(t, result[3], "")
	assert.Equal(t, len(result), 16)

	// Test Case 2 - valid request parameters | varying maximums
	product := ProductPayload("test-case-valid-product.json")

	result = CSVProductValuesByVariant(product, product.Variants[0], 1, 0, 1)
	assert.Equal(t, result[1], "product_code")
	assert.Equal(t, result[3], "product_title")
	assert.Equal(t, len(result), 19)

	// Test Case 3 - valid request
	result = CSVProductValuesByVariant(product, product.Variants[0], 1, 1, 1)
	assert.Equal(t, result[2], "1")
	assert.Equal(t, result[4], "<p>I am a body_html</p>")
	assert.Equal(t, result[6], "product_vendor")
	assert.Equal(t, len(result), 19)
}

func TestCSVProductVariant(t *testing.T) {
	// Test Case 1 - invalid variant data
	productVariant := objects.ProductVariant{}
	result := CSVProductVariant(productVariant)

	assert.Equal(t, result[0], "")
	assert.Equal(t, result[1], "")
	assert.Equal(t, len(result), 2)

	// Test Case 2 - valid variant data
	result = CSVProductVariant(ProductPayload("test-case-valid-product.json").Variants[0])

	assert.Equal(t, result[0], "product_sku")
	assert.Equal(t, result[1], "2347234-9824")
	assert.Equal(t, len(result), 2)
}

func TestCSVVariantOptions(t *testing.T) {
	// Test Case 1 - invalid data | empty struct
	result := CSVVariantOptions(objects.Product{}, objects.ProductVariant{})

	assert.Equal(t, result[0], "")
	assert.Equal(t, result[1], "")
	assert.Equal(t, len(result), 6)

	// Test Case 2 - valid data
	product := ProductPayload("test-case-valid-product.json")
	result = CSVVariantOptions(product, product.Variants[0])

	assert.Equal(t, result[0], "Size")
	assert.Equal(t, result[1], "option1")
	assert.Equal(t, result[2], "Colour")
	assert.Equal(t, result[3], "option2")
	assert.Equal(t, result[4], "Fabric")
	assert.Equal(t, result[5], "option3")
	assert.Equal(t, len(result), 6)
}

func TestGenerateProductOptions(t *testing.T) {
	// static function, nothing to test
	// Test Case 1
	result := generateProductOptions()
	assert.Equal(t, result[0], "option1_name")
	assert.Equal(t, result[1], "option1_value")
	assert.Equal(t, result[2], "option2_name")
	assert.Equal(t, result[3], "option2_value")
	assert.Equal(t, result[4], "option3_name")
	assert.Equal(t, result[5], "option3_value")
	assert.Equal(t, len(result), 6)
}

func TestGetVariantQtyCSV(t *testing.T) {
	// Test Case 1 - invalid data
	result := getVariantQtyCSV(objects.ProductVariant{}, 0)
	assert.Equal(t, len(result), 0)

	// Test Case 2 - valid data | lower qty max
	productVariant := ProductPayload("test-case-valid-product.json").Variants[0]
	result = getVariantQtyCSV(productVariant, 1)
	assert.Equal(t, result[0], "123")
	assert.Equal(t, len(result), 1)

	// Test Case 3  - valid data | higher qty max
	result = getVariantQtyCSV(productVariant, 3)
	assert.Equal(t, result[0], "123")
	assert.Equal(t, result[1], "0")
	assert.Equal(t, result[2], "0")
	assert.Equal(t, len(result), 3)
}

func TestGetVariantPricingCSV(t *testing.T) {
	// Test Case 1 - invalid data
	result := getVariantPricingCSV(objects.ProductVariant{}, 0)
	assert.Equal(t, len(result), 0)

	// Test Case 2 - valid data | lower pricing max
	productVariant := ProductPayload("test-case-valid-product.json").Variants[0]
	result = getVariantPricingCSV(productVariant, 1)
	assert.Equal(t, result[0], "1500.99")
	assert.Equal(t, len(result), 1)

	// Test Case 3  - valid data | higher pricing max
	result = getVariantPricingCSV(productVariant, 3)
	assert.Equal(t, result[0], "1500.99")
	assert.Equal(t, result[1], "0.00")
	assert.Equal(t, result[2], "0.00")
	assert.Equal(t, len(result), 3)
}

func TestGetProductImagesCSV(t *testing.T) {
	// Test Case 1 - invalid product images | no key
	result := GetProductImagesCSV([]objects.ProductImages{}, 0, false)
	assert.Equal(t, len(result), 0)

	// Test Case 2 - invalid product images | key
	result = GetProductImagesCSV([]objects.ProductImages{}, 0, true)
	assert.Equal(t, len(result), 0)

	// Test Case 3 - invalid product images | key | lower max
	result = GetProductImagesCSV([]objects.ProductImages{}, 1, false)
	assert.Equal(t, result[0], "")
	assert.Equal(t, len(result), 1)

	// Test Case 4 - invalid product images | key | higher max
	result = GetProductImagesCSV([]objects.ProductImages{}, 3, true)
	assert.Equal(t, result[0], "image_1")
	assert.Equal(t, result[1], "image_2")
	assert.Equal(t, result[2], "image_3")
	assert.Equal(t, len(result), 3)

	// Test Case 5 - valid product images | no key
	productImages := ProductPayload("test-case-valid-product.json").ProductImages
	result = GetProductImagesCSV(productImages, 0, false)
	assert.Equal(t, result[0], "")
	assert.Equal(t, len(result), 1)

	// Test Case 6 - valid product images | key | zero max
	result = GetProductImagesCSV(productImages, 0, true)
	assert.Equal(t, len(result), 0)

	// Test Case 7 - invalid product images | low max
	result = GetProductImagesCSV(productImages, 1, true)
	assert.Equal(t, result[0], "image_1")
	assert.Equal(t, len(result), 1)

	// Test Case 8 - invalid product images | high max
	result = GetProductImagesCSV(productImages, 3, true)
	assert.Equal(t, result[0], "image_1")
	assert.Equal(t, result[1], "image_2")
	assert.Equal(t, result[2], "image_3")
	assert.Equal(t, len(result), 3)
}

func TestWriteFile(t *testing.T) {
	// Test Case 1 - with empty data

	// Test Case 2 - with valid data
}

func TestReadFile(t *testing.T) {

}

func TestRemoveFile(t *testing.T) {

}

func TestGetKeysByMatcher(t *testing.T) {

}

func TestLoopRemoveCSV(t *testing.T) {

}

// Creates a mock HTTP request
func Init(
	requestURL,
	requestMethod string,
	additionalRequestHeaders map[string][]string,
	buffer bytes.Buffer,
) *http.Request {
	req, _ := http.NewRequest(requestMethod, requestURL, &buffer)
	for key, value := range additionalRequestHeaders {
		for _, sub_value := range value {
			req.Header.Add(key, sub_value)
		}
	}
	req.Header.Add("Content-Type", "application/json")
	return req
}

/* Function that creates a multipart/form request to be used in the import handle */
func CreateMultiPartFormData(fileName, formKey string) (*multipart.Writer, bytes.Buffer) {
	// Create a buffer to store the request body
	var buf bytes.Buffer

	// Create a new multipart writer with the buffer
	w := multipart.NewWriter(&buf)

	// Add a file to the request
	file, err := os.Open("./test_payloads/upload/" + fileName)
	if err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}
	defer file.Close()

	// Create a new form field
	fw, err := w.CreateFormFile(formKey, "./test_payloads/upload/"+fileName)
	if err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}

	// Copy the contents of the file to the form field
	if _, err := io.Copy(fw, file); err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}
	w.Close()
	return w, buf
}

/*
Returns a byte array representing the file data that was read

Data is retrived from the project directory `test_payloads`
*/
func payload(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	respBody, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	return respBody
}

// Creates a product payload
func ProductPayload(fileName string) objects.Product {
	fileBytes := payload("./test_payloads/product/" + fileName)
	productData := objects.Product{}
	err := json.Unmarshal(fileBytes, &productData)
	if err != nil {
		log.Println(err)
	}
	return productData
}
