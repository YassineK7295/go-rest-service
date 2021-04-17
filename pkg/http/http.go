package http

import (
	"bytes"
	"fmt"
	"net/http"
)

// Sends POST request and returns status code or error
func SendPostRequest(url, endpoint, jsonStr string) (int, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", url, endpoint), bytes.NewBuffer([]byte(jsonStr)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

// Sends DEL request and returns status code or error
func SendDelRequest(url, endpoint, key string) (int, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s/%s", url, endpoint, key), nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

// Sends PUT request and returns status code or error
func SendPutRequest(url, endpoint, key, jsonStr string) (int, error) {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s%s/%s", url, endpoint, key), bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
