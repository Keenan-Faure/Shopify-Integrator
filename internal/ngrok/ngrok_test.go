package ngrok

import (
	"fmt"
	"net/http"
	"objects"
	"testing"

	"github.com/jarcoal/httpmock"
)

const NGROK_URL = "http://localhost:8888"

func TestFetchNgrokTunnels(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	response := objects.NgrokTunnelResponse{
		Tunnels: []struct {
			Name      string "json:\"name\""
			ID        string "json:\"ID\""
			URI       string "json:\"uri\""
			PublicURL string "json:\"public_url\""
			Proto     string "json:\"proto\""
			Config    struct {
				Addr    string "json:\"addr\""
				Inspect bool   "json:\"inspect\""
			} "json:\"config\""
			Metrics struct {
				Conns struct {
					Count  int "json:\"count\""
					Gauge  int "json:\"gauge\""
					Rate1  int "json:\"rate1\""
					Rate5  int "json:\"rate5\""
					Rate15 int "json:\"rate15\""
					P50    int "json:\"p50\""
					P90    int "json:\"p90\""
					P95    int "json:\"p95\""
					P99    int "json:\"p99\""
				} "json:\"conns\""
				HTTP struct {
					Count  int "json:\"count\""
					Rate1  int "json:\"rate1\""
					Rate5  int "json:\"rate5\""
					Rate15 int "json:\"rate15\""
					P50    int "json:\"p50\""
					P90    int "json:\"p90\""
					P95    int "json:\"p95\""
					P99    int "json:\"p99\""
				} "json:\"http\""
			} "json:\"metrics\""
		}{},
		URI: "",
	}

	// Exact URL match
	httpmock.RegisterResponder("GET", NGROK_URL+"api/tunnels",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, response)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})

	// get count info
	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()
	fmt.Println(info)
}
