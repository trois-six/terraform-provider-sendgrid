package sendgrid

import (
	"encoding/json"
	"fmt"
)

// Template is a Sendgrid transactional template.
type Template struct {
	ID         string            `json:"id,omitempty"`
	Name       string            `json:"name,omitempty"`
	Generation string            `json:"generation,omitempty"`
	UpdatedAt  string            `json:"updated_at,omitempty"`
	Versions   []TemplateVersion `json:"versions,omitempty"`
	Warnings   []string          `json:"warnings,omitempty"`
}

func parseTemplate(respBody string) (*Template, error) {
	var body Template
	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

// CreateTemplate creates a transactional template and returns it.
func (c *Client) CreateTemplate(name, generation string) (*Template, error) {
	if name == "" {
		return nil, fmt.Errorf("[CreateTemplate] a name is required")
	}
	if generation == "" {
		generation = "dynamic"
	}
	respBody, _, err := c.Post("POST", "/templates", Template{
		Name:       name,
		Generation: generation,
	})
	if err != nil {
		return nil, err
	}
	return parseTemplate(respBody)
}

// ReadTemplate retreives a transactional template and returns it.
func (c *Client) ReadTemplate(id string) (*Template, error) {
	if id == "" {
		return nil, fmt.Errorf("[ReadTemplate] an ID is required")
	}
	respBody, _, err := c.Get("GET", "/templates/"+id)
	if err != nil {
		return nil, err
	}
	return parseTemplate(respBody)
}

// UpdateTemplate edits a transactional template and returns it. We can't change the "generation" of a transactional template.
func (c *Client) UpdateTemplate(id, name string) (*Template, error) {
	if id == "" {
		return nil, fmt.Errorf("[UpdateTemplate] an ID is required")
	}
	if name == "" {
		return nil, fmt.Errorf("[UpdateTemplate] a name is required")
	}
	respBody, _, err := c.Post("PATCH", "/templates/"+id, Template{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return parseTemplate(respBody)
}

// DeleteTemplate deletes a transactional template.
func (c *Client) DeleteTemplate(id string) (bool, error) {
	if id == "" {
		return false, fmt.Errorf("[DeleteTemplate] an ID is required")
	}
	_, _, err := c.Get("DELETE", "/templates/"+id)
	if err != nil {
		return false, err
	}
	return true, nil
}
