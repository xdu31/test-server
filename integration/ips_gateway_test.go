// +build integration

package integration

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/bitly/go-simplejson"
	types1 "github.com/infobloxopen/protoc-gen-gorm/types"
	pb "github.com/xdu31/test-server/pkg/pb"
)

const (
	appId = "test-server/ips"
)

// TestGetVersion_REST gets the version of the service
func TestGetVersion_REST(t *testing.T) {
	dbTest.Reset(t)
	resGet, err := MakeRequestWithDefaults(
		http.MethodGet,
		"http://localhost:8080/test-server/v1/version",
		nil,
	)
	if err != nil {
		t.Fatalf("unable to get version: %v", err)
	}
	ValidateResponseCode(t, resGet, http.StatusOK)
	getJSON, err := simplejson.NewFromReader(resGet.Body)
	if err != nil {
		t.Fatalf("unable to marshal json response: %v", err)
	}
	var tests = []struct {
		name   string
		json   *simplejson.Json
		expect string
	}{
		{
			name:   "service version",
			json:   getJSON.GetPath("version"),
			expect: `"0.0.1"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ValidateJSONSchema(t, test.json, test.expect)
		})
	}
}

// TestCreateIp_REST uses the REST gateway to create a new ip and
// ensure the JSON response matches what is expected
// 1. Create an ip entry with a POST request
// 2. Unmarshal the JSON into a simplejson struct
// 3. Ensure the JSON fields match what is expected
func TestCreateIps_REST(t *testing.T) {
	dbTest.Reset(t)
	ip := pb.Ip{
		IpAddress: &types1.InetValue{Value: "1.2.3.4/24"},
	}
	resCreate, err := MakeRequestWithDefaultsIP(
		http.MethodPost,
		"http://localhost:8080/test-server/v1/ips",
		ip,
	)
	if err != nil {
		t.Fatalf("unable to list ip: %v", err)
	}
	ValidateResponseCode(t, resCreate, http.StatusOK)
	createJSON, err := simplejson.NewFromReader(resCreate.Body)
	if err != nil {
		t.Fatalf("unable to marshal json response: %v", err)
	}
	var tests = []struct {
		name   string
		json   *simplejson.Json
		expect string
	}{
		{
			name:   "ip address",
			json:   createJSON.GetPath("result", "ip_address"),
			expect: `"1.2.3.4/24"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ValidateJSONSchema(t, test.json, test.expect)
		})
	}
}

// TestReadIps_REST uses the REST gateway to create a new ip and
// then read that ip from the application
// 1. Create an ip entry with a POST request
// 2. Get the ip from the applicaiton
// 2. Unmarshal the JSON into a simplejson struct
// 3. Ensure the JSON fields match what is expected
func TestReadIps_REST(t *testing.T) {
	dbTest.Reset(t)
	ip := pb.Ip{
		IpAddress: &types1.InetValue{Value: "1.2.3.5/24"},
	}
	resCreate, err := MakeRequestWithDefaultsIP(
		http.MethodPost,
		"http://localhost:8080/test-server/v1/ips",
		ip,
	)
	createJSON, err := simplejson.NewFromReader(resCreate.Body)
	if err != nil {
		t.Fatalf("unable to unmarshal create ip response body: %v", err)
	}
	id, err := createJSON.GetPath("result", "id").String()
	if err != nil {
		t.Fatalf("unable to get contact id from response json: %v", err)
	}
	id = strings.TrimPrefix(id, appId)
	resRead, err := MakeRequestWithDefaults(
		http.MethodGet, fmt.Sprintf("http://localhost:8080/test-server/v1/ips/%s", id),
		nil,
	)
	if err != nil {
		t.Fatalf("unable to get ip: %v", err)
	}
	ValidateResponseCode(t, resRead, http.StatusOK)
	readJSON, err := simplejson.NewFromReader(resRead.Body)
	var tests = []struct {
		name   string
		json   *simplejson.Json
		expect string
	}{
		{
			name:   "ip address",
			json:   readJSON.GetPath("result", "ip_address"),
			expect: `"1.2.3.5/24"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ValidateJSONSchema(t, test.json, test.expect)
		})
	}
}

// ValidateResponseCode checks the http status of a given request and will
// fail the current test if it doesn't match the expected status code
func ValidateResponseCode(t *testing.T, res *http.Response, expected int) {
	if expected != res.StatusCode {
		t.Errorf("validation error: unexpected http response status: have %d; want %d",
			res.StatusCode, expected,
		)
	}
}

func ValidateServiceVersion(t *testing.T, version string, expected string) {
	if version != expected {
		t.Errorf("actual version does not match expected version: have %s; want %v",
			version, expected,
		)
	}
}

// ValidateJSONSchema ensures a given json field matches an expcted json
// string
func ValidateJSONSchema(t *testing.T, json *simplejson.Json, expected string) {
	if json == nil {
		t.Fatalf("validation error: json schema for is nil")
	}
	encoded, err := json.Encode()
	if err != nil {
		t.Fatalf("validation error: unable to encode expected json: %v", err)
	}
	if actual := string(encoded); actual != expected {
		t.Errorf("actual json schema does not match expected schema: have %s; want %v",
			actual, expected,
		)
	}
}
