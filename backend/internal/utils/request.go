package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Get(url string, result any) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}

	return nil
}

// Generic type
func GetJSON[T any](url string, result *T) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func HttpPostJSON(url string, payload any) error {

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to post request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil
}
