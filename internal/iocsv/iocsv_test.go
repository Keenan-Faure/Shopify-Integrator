package iocsv

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
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
	if err != nil {
		t.Errorf("Expected 'nil' but found: " + err.Error())
	}
	if fileName == "" {
		t.Errorf("Expected 'upload-***' but found: " + fileName)
	}

	// Test Case 3 - invalid headers

	// Test Case 4 - invalid formKey used in multipart/form

	// Test Case 5 - valid http request
}

func TestCSVProductHeaders(t *testing.T) {

}

func TestCSVProductValuesByVariant(t *testing.T) {

}

func TestCSVProductVariant(t *testing.T) {

}

func TestCSVVariantOptions(t *testing.T) {

}

func TestGenerateProductOptions(t *testing.T) {

}

func TestGetVariantQtyCSV(t *testing.T) {

}

func TestGetVariantPricingCSV(t *testing.T) {

}

func TestGetProductImagesCSV(t *testing.T) {

}

func TestWriteFile(t *testing.T) {

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
	file, err := os.Open("./test_payloads//" + fileName)
	if err != nil {
		log.Println(err)
		return &multipart.Writer{}, buf
	}
	defer file.Close()

	// Create a new form field
	fw, err := w.CreateFormFile(formKey, "./test_payloads/"+fileName)
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
