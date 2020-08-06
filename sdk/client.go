package sendgrid

import (
	"encoding/json"
	"fmt"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

// Client is a Sendgrid client.
type Client struct {
	apiKey     string
	host       string
	onBehalfOf string
}

// NewClient creates a Sendgrid Client.
func NewClient(apiKey, host, onBehalfOf string) *Client {
	if host == "" {
		host = "https://api.sendgrid.com/v3"
	}
	return &Client{
		apiKey:     apiKey,
		host:       host,
		onBehalfOf: onBehalfOf,
	}
}

func bodyToJSON(body interface{}) ([]byte, error) {
	if body == nil {
		return nil, fmt.Errorf("[bodyToJSON] body must no be nil")
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("[bodyToJSON] body could not be jsonified")
	}
	return jsonBody, nil
}

// Get gets a resource from Sendgrid.
func (c *Client) Get(method rest.Method, endpoint string) (string, int, error) {
	var req rest.Request
	if c.onBehalfOf != "" {
		req = sendgrid.GetRequestSubuser(c.apiKey, endpoint, c.host, c.onBehalfOf)
	} else {
		req = sendgrid.GetRequest(c.apiKey, endpoint, c.host)
	}
	req.Method = method
	resp, err := sendgrid.API(req)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("[Get] %s", err)
	}
	return resp.Body, resp.StatusCode, nil
}

// Post posts a resource to Sendgrid.
func (c *Client) Post(method rest.Method, endpoint string, body interface{}) (string, int, error) {
	var err error
	var req rest.Request
	if c.onBehalfOf != "" {
		req = sendgrid.GetRequestSubuser(c.apiKey, endpoint, c.host, c.onBehalfOf)
	} else {
		req = sendgrid.GetRequest(c.apiKey, endpoint, c.host)
	}
	req.Method = method
	req.Body, err = bodyToJSON(body)
	if err != nil {
		return "", 0, err
	}
	resp, err := sendgrid.API(req)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("[Post] %s", err)
	}
	return resp.Body, resp.StatusCode, nil
}
