package ngrok

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"objects"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/jarcoal/httpmock"
)

const NGROK_URL = "http://ngrok.api.localhost:8888"

func TestFetchNgrokTunnels(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	response := CreateTestNgrokPayload("test-case-valid-data.json")
	httpmock.RegisterResponder("GET", NGROK_URL+"/api/tunnels",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, response)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	res, err := fetchHelper(NGROK_URL, "api/tunnels", http.MethodGet, nil)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}
	if res.StatusCode != 200 {
		t.Errorf("expected '200' but found: " + fmt.Sprint(res.StatusCode))
	}
	tunnels := objects.NgrokTunnelResponse{}
	err = json.Unmarshal(respBody, &tunnels)
	if err != nil {
		t.Errorf("expected 'nil' but found: " + err.Error())
	}

	// get count info
	httpmock.GetTotalCallCount()
	requestInfo := httpmock.GetCallCountInfo()

	// Test 1 - Assert if the amount of API calls that is made is correct.
	assert.Equal(t, 1, requestInfo["GET "+NGROK_URL+"/api/tunnels"])

	// Test 2 - Assert if the response that we receive is correct and is the correct body
	assert.Equal(t, 1, len(response.Tunnels))
	assert.Equal(t, "website", response.Tunnels[0].Name)
	assert.Equal(t, "/api/tunnels/website", response.Tunnels[0].URI)
	assert.Equal(t, "http://host.docker.internal:8080", response.Tunnels[0].Config.Addr)
	assert.Equal(t, "/api/tunnels", response.URI)
}

/* Returns a test user RequestBodyRegister struct */
func CreateTestNgrokPayload(fileName string) objects.NgrokTunnelResponse {
	fileBytes := payload("./test_payloads/" + fileName)
	ngrokTunnelResponse := objects.NgrokTunnelResponse{}
	err := json.Unmarshal(fileBytes, &ngrokTunnelResponse)
	if err != nil {
		log.Println(err)
	}
	return ngrokTunnelResponse
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
