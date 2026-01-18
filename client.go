package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchAndDecode[T any](client *http.Client, req *http.Request) (*T, error) {

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("Api Error: status=%d, body=%s", response.StatusCode, body)
	}

	var result T
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
