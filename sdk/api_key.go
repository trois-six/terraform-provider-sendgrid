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
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, fmt.Errorf("failed parsing API key: %w", err)
	}

	return &body, nil
}

// CreateAPIKey creates an APIKey and returns it.
func (c *Client) CreateAPIKey(name string, scopes []string) (*APIKey, error) {
	if name == "" {
		return nil, ErrNameRequired
	}

	respBody, _, err := c.Post("POST", "/api_keys", APIKey{
		Name:   name,
		Scopes: scopes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating API key: %w", err)
	}

	return parseAPIKey(respBody)
}

// ReadAPIKey retreives an APIKey and returns it.
func (c *Client) ReadAPIKey(id string) (*APIKey, error) {
	if id == "" {
		return nil, ErrAPIKeyIDRequired
	}

	respBody, _, err := c.Get("GET", "/api_keys/"+id)
	if err != nil {
		return nil, fmt.Errorf("failed reading API key: %w", err)
	}

	return parseAPIKey(respBody)
}

// UpdateAPIKey edits an APIKey and returns it.
func (c *Client) UpdateAPIKey(id, name string, scopes []string) (*APIKey, error) {
	if id == "" {
		return nil, ErrAPIKeyIDRequired
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
		return nil, fmt.Errorf("failed updating API key: %w", err)
	}

	return parseAPIKey(respBody)
}

// DeleteAPIKey deletes an APIKey.
func (c *Client) DeleteAPIKey(id string) (bool, error) {
	if id == "" {
		return false, ErrAPIKeyIDRequired
	}

	if _, statusCode, err := c.Get("DELETE", "/api_keys/"+id); statusCode > 299 || err != nil {
		return false, fmt.Errorf("failed deleting API key: %w", err)
	}

	return true, nil
}
