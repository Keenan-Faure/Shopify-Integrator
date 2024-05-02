package main

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/jarcoal/httpmock"
)

func TestOAuthGoogleCallback(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, MOCK_APP_API_URL+"/api/google/callback",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(303, "")
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	/* Test 1 - invalid authentication */
	res, _ := FetchHelper(MOCK_APP_API_URL, "api/google/callback", http.MethodGet, nil)
	assert.Equal(t, http.StatusSeeOther, res.StatusCode)
}

func TestOAuthGoogleLogin(t *testing.T) {

}

func TestOAuthGoogleOAuthn(t *testing.T) {

}

func TestGenerateStateOauthCookie(t *testing.T) {

}

func TestGetUserDataFromGoogle(t *testing.T) {

}

func FetchHelper(host, endpoint, method string, body io.Reader) (*http.Response, error) {
	httpClient := http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(method, host+"/"+endpoint, body)
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return res, nil
}
