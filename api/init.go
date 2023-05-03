package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	HTTPClient    *http.Client
	Host          string
	Authorization string
}

func NewClient(host string, authorization string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{
			Timeout: time.Duration(10) * time.Second,
		},
		Host:          host,
		Authorization: authorization,
	}
	return &c, nil
}

func (c *Client) doRequest(method string, url string, body []byte, headers map[string]string) ([]byte, int, string, map[string][]string, error) {
	payload := bytes.NewBuffer(body)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, 500, "500 Internal Server Error", nil, err
	}

	if c.Authorization != "" {
		c.Authorization = c.getBearerToken()
	}

	req.Header.Add("Authorization", c.Authorization)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "PCT")
	req.Header.Add("Content-Type", "application/json")

	for header, value := range headers {
		req.Header.Add(header, value)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 500, "500 Internal Server Error", nil, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 500, "500 Internal Server Error", nil, err
	}

	defer res.Body.Close()
	return b, res.StatusCode, res.Status, res.Header, nil
}
func (c *Client) getBearerToken() string {
	if c.Authorization == "" {
		return ""
	}

	if strings.HasPrefix(c.Authorization, "Bearer") {
		return c.Authorization

	}
	return fmt.Sprintf("Bearer %s", c.Authorization)
}
