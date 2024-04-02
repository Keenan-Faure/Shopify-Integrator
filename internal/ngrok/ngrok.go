package ngrok

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"objects"
	"time"
)

const NGROK_HOST = "http://localhost:8888"
const NGROK_TUNNEL_NAME = "website"

// GET /api/tunnels
func FetchNgrokTunnels() (objects.NgrokTunnelResponse, error) {
	res, err := fetchHelper("api/tunnels", http.MethodGet, nil)
	if err != nil {
		return objects.NgrokTunnelResponse{}, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return objects.NgrokTunnelResponse{}, err
	}
	if res.StatusCode != 200 {
		return objects.NgrokTunnelResponse{}, errors.New(string(respBody))
	}
	tunnels := objects.NgrokTunnelResponse{}
	err = json.Unmarshal(respBody, &tunnels)
	if err != nil {
		log.Println(err)
		return objects.NgrokTunnelResponse{}, err
	}
	return tunnels, nil
}

// Returns the specific ngrok tunnel
func FetchWebsiteTunnel(tunnels objects.NgrokTunnelResponse) string {
	for _, tunnel := range tunnels.Tunnels {
		if tunnel.Name == NGROK_TUNNEL_NAME {
			return tunnel.PublicURL
		}
	}
	return ""
}

// Util function
// Creates the webhook url
func SetUpWebhookURL(domain, api_key, token string) string {
	return domain + "/api/orders?token=" + token + "&api_key=" + api_key
}

// Util fetch helper
func fetchHelper(endpoint, method string, body io.Reader) (*http.Response, error) {
	httpClient := http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(method, NGROK_HOST+"/"+endpoint, body)
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
