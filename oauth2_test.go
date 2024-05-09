package main

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/jarcoal/httpmock"
)

// TODO not sure how to properly test OAuth2, need to think about it

func TestOAuthGoogleCallback(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, MOCK_APP_API_URL+"/api/google/callback",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusSeeOther, "")
			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
			}
			return resp, nil
		},
	)

	/* Test 1 - invalid authentication */
	res, _ := FetchHelper(MOCK_APP_API_URL, "api/google/callback", http.MethodGet, nil)
	assert.Equal(t, http.StatusSeeOther, res.StatusCode)
}

func TestOAuthGoogleLogin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, MOCK_APP_API_URL+"/api/google/oauth2/login",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusTemporaryRedirect, "")
			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
			}
			return resp, nil
		},
	)

	/* Test 1 - invalid authentication */
	res, _ := FetchHelper(MOCK_APP_API_URL, "api/google/oauth2/login", http.MethodGet, nil)
	assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
}

func TestOAuthGoogleOAuthn(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, MOCK_APP_API_URL+"/api/google/oauth2/login",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(http.StatusOK, "")
			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
			}
			return resp, nil
		},
	)

	/* Test 1 - invalid authentication */
	res, _ := FetchHelper(MOCK_APP_API_URL, "api/google/oauth2/login", http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGenerateStateOauthCookie(t *testing.T) {
	// Test 1 - invalid writer

	// Test 2 - valid writer
}

func TestGetUserDataFromGoogle(t *testing.T) {
	// cant test this
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
