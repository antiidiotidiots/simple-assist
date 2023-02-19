package core

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func Request(url string, requestType string, payloadMap map[string]interface{}, headers map[string]string) (string, error) {
	requestType = strings.ToUpper(requestType)
	// Send a request to the URL, with the URL which was passed to the function
	var req *http.Request
	var err error
	// If payloadMap is nil, don't send a payload
	if payloadMap == nil {
		req, err = http.NewRequest(requestType, url, nil)
		if err != nil {
			return "", err
		}
	} else {
		// logrus.Debugf("Payload map: %v", payloadMap)
		payloadBytes, err := json.Marshal(payloadMap)
		if err != nil {
			return "", err
		}
		payload := strings.NewReader(string(payloadBytes))
		req, err = http.NewRequest(requestType, url, payload)
		if err != nil {
			return "", err
		}
	}

	// If headers is nil, don't send any headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Make the actual request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	// Convert the result body to a string and then return it
	resultBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = res.Body.Close()
	if err != nil {
		return "", err
	}
	// Debug the status code
	// logrus.Debugf("Sending %s request to %s with payload %v, bearer authorization %s. Got status code %d", requestType, url, payloadMap, bearerAuth, res.StatusCode)
	return string(resultBodyBytes), nil
}
