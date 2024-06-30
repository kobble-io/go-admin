package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpClientConfig struct {
	BaseURL string
	Secret  string
}

const SdkSecretHeaderName = "Kobble-Sdk-Secret"

type HttpClient struct {
	config    HttpClientConfig
	userAgent string
}

func NewHttpClient(config HttpClientConfig) *HttpClient {
	return &HttpClient{
		config:    config,
		userAgent: "Kobble Go SDK/1.x",
	}
}

func (c *HttpClient) makeURL(path string, params map[string]string) (string, error) {
	base, err := url.Parse(c.config.BaseURL)
	if err != nil {
		return "", err
	}

	pathURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	finalURL := base.ResolveReference(pathURL)
	query := finalURL.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	finalURL.RawQuery = query.Encode()

	return finalURL.String(), nil
}

func (c *HttpClient) GetJson(path string, params map[string]string, result any, expectedStatus int) error {
	fullURL, err := c.makeURL(path, params)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set(SdkSecretHeaderName, c.config.Secret)
	req.Header.Set("User-Agent", c.userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error: %s", string(bodyBytes))
	}

	if result == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(&result)
}

func (c *HttpClient) PostJson(path string, payload any, result any, expectedStatus int) error {
	fullURL, err := c.makeURL(path, nil)
	if err != nil {
		return err
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set(SdkSecretHeaderName, c.config.Secret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("expectedStatus code %d, got %d\n", expectedStatus, resp.StatusCode)

	if resp.StatusCode != expectedStatus {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error: %s", string(bodyBytes))
	}

	if result == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(&result)
}
