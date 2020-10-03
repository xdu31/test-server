package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	pb "github.com/xdu31/test-server/pkg/pb"
)

const (
	// TestSecret dummy secret used for signing test JWTs
	bearer   = "Bearer"
	jwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2NvdW50SUQiOjF9.GsXyFDDARjXe1t9DPo2LIBKHEal3O7t3vLI3edA7dGU"
)

func addDefaultTokenToRequest(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", bearer, jwtToken))
}

// MakeRequestWithDefaults issues request that contains the necessary parameters
// request to reach the contacts application (e.g. an authorization token)
func MakeRequestWithDefaults(method, url string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	addDefaultTokenToRequest(req)

	client := &http.Client{}
	return client.Do(req)
}

// MakeRequestWithDefaultsIP issues request that contains the necessary parameters
// request to reach the contacts application (e.g. an authorization token)
func MakeRequestWithDefaultsIP(method, url string, payload pb.Ip) (*http.Response, error) {
	body, err := payload.MarshalJSON()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	addDefaultTokenToRequest(req)

	client := &http.Client{}
	return client.Do(req)
}
