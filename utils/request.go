package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Header struct {
	Key   string
	Value string
}

//Request data json
func Request(method string, url string, body string, header []Header) (*http.Response, error) {

	var jsonStr = []byte(body)
	var result interface{}

	client := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	if len(header) > 0 {
		for _, h := range header {
			req.Header.Set(h.Key, h.Value)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Error method or header from : %v", url)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Error to request data from : %v", url)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, fmt.Errorf("Error to request data from : %v", url)
	}

	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.(map[string]interface{})) == 0 {
		return nil, fmt.Errorf("the result is empty : %v", url)
	}

	if result == nil {
		return nil, fmt.Errorf("the result is empty : %v", url)
	}

	return resp, nil
}
