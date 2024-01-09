package opensubtitles

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type Client struct {
	key   string
	token string
}

func NewClient(key string) *Client {
	return &Client{
		key: key,
	}
}

func (c Client) Request(method, url string, body []byte, token string) ([]byte, error) {

	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, URL+url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, URL+url, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", c.key)
	req.Header.Set("User-Agent", "go-find-subititle v0.0.1")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: timeout * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
