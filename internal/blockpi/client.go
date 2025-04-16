package blockpi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	Token   string
	Proxy   string // (future use)
	client  *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://hub.nimbly.blockpi.org/",
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type request struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     int           `json:"id"`
}

type Response struct {
	ID      int                    `json:"id"`
	Jsonrpc string                 `json:"jsonrpc"`
	Result  map[string]interface{} `json:"result"`
}

func (c *Client) makeRequest(method string, params map[string]interface{}, client *http.Client, authRequired bool) (Response, error) {
	body, _ := json.Marshal(request{
		Method: method,
		Params: []interface{}{params},
		ID:     1,
	})

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(body))
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	if (authRequired || client != c.client) && c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	if client == nil {
		client = c.client
	}

	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		raw, _ := io.ReadAll(resp.Body)
		return Response{}, fmt.Errorf("%s failed (%d): %s", method, resp.StatusCode, string(raw))
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Response{}, err
	}

	return result, nil
}

func (c *Client) call(method string, params map[string]interface{}, authRequired bool) (Response, error) {
	return c.makeRequest(method, params, nil, authRequired)
}

func (c *Client) callWithClient(method string, params map[string]interface{}, useClient *http.Client, authRequired bool) (Response, error) {
	return c.makeRequest(method, params, useClient, authRequired)
}
