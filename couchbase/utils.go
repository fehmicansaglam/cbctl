package couchbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fehmicansaglam/cbctl/shared"
)

func debugLog(format string, args ...interface{}) {
	if shared.Debug {
		fmt.Fprintf(os.Stderr, "DEBUG: "+format+"\n", args...)
	}
}

func httpRequest(method, endpoint string, body, target interface{}, expectedStatusCode int) error {
	baseURL := fmt.Sprintf("%s://%s:%d/%s", shared.CouchbaseProtocol, shared.CouchbaseHost, shared.CouchbasePort, endpoint)

	if shared.Debug {
		debugLog("Request URL: %s", baseURL)
	}

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		if shared.Debug {
			debugLog("Request Body: %s", bodyBytes)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, baseURL, bodyReader)
	if err != nil {
		return err
	}

	if shared.CouchbaseUsername != "" && shared.CouchbasePassword != "" {
		req.SetBasicAuth(shared.CouchbaseUsername, shared.CouchbasePassword)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		return fmt.Errorf("failed to get buckets, status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func getJSONResponse(endpoint string, target interface{}) error {
	return httpRequest(http.MethodGet, endpoint, nil, target, http.StatusOK)
}

func getJSONResponseWithBody(endpoint string, target interface{}, body interface{}) error {
	return httpRequest(http.MethodPost, endpoint, body, target, http.StatusOK)
}

func postWithoutBody(endpoint string, target interface{}) error {
	return httpRequest(http.MethodPost, endpoint, nil, target, http.StatusOK)
}

func getNestedPath(field string, nestedPaths []string) (string, bool) {
	for _, path := range nestedPaths {
		if strings.HasPrefix(field, path) {
			return path, true
		}
	}
	return "", false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
