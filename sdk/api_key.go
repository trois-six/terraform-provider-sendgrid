package sendgrid

import (
	"encoding/json"
	"fmt"
)

// APIKey is a Sendgrid API key.
type APIKey struct {
	ID     string   `json:"api_key_id,omitempty"`
	APIKey string   `json:"api_key,omitempty"`
	Name   string   `json:"name,omitempty"`
	Scopes []string `json:"scopes,omitempty"`
}

func parseAPIKey(respBody string) (*APIKey, error) {
	var body APIKey
	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

// CreateAPIKey creates an APIKey and returns it.
func (c *Client) CreateAPIKey(name string, scopes []string) (*APIKey, error) {
	if name == "" {
		return nil, fmt.Errorf("[CreateAPIKey] a name is required")
	}
	respBody, _, err := c.Post("POST", "/api_keys", APIKey{
		Name:   name,
		Scopes: scopes,
	})
	if err != nil {
		return nil, err
	}
	return parseAPIKey(respBody)
}

// ReadAPIKey retreives an APIKey and returns it.
func (c *Client) ReadAPIKey(id string) (*APIKey, error) {
	if id == "" {
		return nil, fmt.Errorf("[ReadAPIKey] an ID is required")
	}
	respBody, _, err := c.Get("GET", "/api_keys/"+id)
	if err != nil {
		return nil, err
	}
	return parseAPIKey(respBody)
}

// UpdateAPIKey edits an APIKey and returns it.
func (c *Client) UpdateAPIKey(id, name string, scopes []string) (*APIKey, error) {
	if id == "" {
		return nil, fmt.Errorf("[UpdateAPIKey] an ID is required")
	}
	t := APIKey{}
	if name != "" {
		t.Name = name
	}
	if len(scopes) > 0 {
		t.Scopes = scopes
	}
	respBody, _, err := c.Post("PUT", "/api_keys/"+id, t)
	if err != nil {
		return nil, err
	}
	return parseAPIKey(respBody)
}

// DeleteAPIKey deletes an APIKey.
func (c *Client) DeleteAPIKey(id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("[DeleteAPIKey] an ID is required")
	}
	_, _, err := c.Get("DELETE", "/api_keys/"+id)
	if err != nil {
		return false, err
	}
	return true, nil
}
