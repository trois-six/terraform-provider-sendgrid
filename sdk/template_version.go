package sendgrid

import (
	"encoding/json"
	"fmt"
)

// TemplateVersion is a Sendgrid transactional template version.
type TemplateVersion struct {
	ID                   string   `json:"id,omitempty"`
	TemplateID           string   `json:"template_id,omitempty"`
	UpdatedAt            string   `json:"updated_at,omitempty"`
	ThumbnailURL         string   `json:"thumbnail_url,omitempty"`
	Warnings             []string `json:"warnings,omitempty"`
	Active               int      `json:"active,omitempty"`
	Name                 string   `json:"name,omitempty"`
	HTMLContent          string   `json:"html_content,omitempty"`
	PlainContent         string   `json:"plain_content,omitempty"`
	GeneratePlainContent bool     `json:"generate_plain_content,omitempty"`
	Subject              string   `json:"subject,omitempty"`
	Editor               string   `json:"editor,omitempty"`
	TestData             string   `json:"test_data,omitempty"`
}

func parseTemplateVersion(respBody string) (*TemplateVersion, error) {
	var body TemplateVersion
	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

// CreateTemplateVersion creates a new version of a transactional template and returns it.
func (c *Client) CreateTemplateVersion(t TemplateVersion) (*TemplateVersion, error) {
	if t.TemplateID == "" {
		return nil, fmt.Errorf("[CreateTemplateVersion] a template ID is required")
	}
	if t.Name == "" {
		return nil, fmt.Errorf("[CreateTemplateVersion] a template Name is required")
	}
	if t.Subject == "" {
		return nil, fmt.Errorf("[CreateTemplateVersion] a template Subject is required")
	}
	respBody, _, err := c.Post("POST", "/templates/"+t.TemplateID+"/versions", t)
	if err != nil {
		return nil, err
	}
	return parseTemplateVersion(respBody)
}

// ReadTemplateVersion retreives a version of a transactional template and returns it.
func (c *Client) ReadTemplateVersion(templateID, id string) (*TemplateVersion, error) {
	if templateID == "" {
		return nil, fmt.Errorf("[ReadTemplateVersion] a template ID is required")
	}
	if id == "" {
		return nil, fmt.Errorf("[ReadTemplateVersion] a version ID is required")
	}
	respBody, _, err := c.Get("GET", "/templates/"+templateID+"/versions/"+id)
	if err != nil {
		return nil, err
	}
	return parseTemplateVersion(respBody)
}

// UpdateTemplateVersion edits a version of a transactional template and returns it.
func (c *Client) UpdateTemplateVersion(t TemplateVersion) (*TemplateVersion, error) {
	if t.ID == "" {
		return nil, fmt.Errorf("[UpdateTemplateVersion] a template ID is required")
	}
	if t.TemplateID == "" {
		return nil, fmt.Errorf("[UpdateTemplateVersion] a template ID is required")
	}
	respBody, _, err := c.Post("PATCH", "/templates/"+t.TemplateID+"/versions/"+t.ID, t)
	if err != nil {
		return nil, err
	}
	return parseTemplateVersion(respBody)
}

// DeleteTemplateVersion deletes a version of a transactional template.
func (c *Client) DeleteTemplateVersion(templateID, id string) (bool, error) {
	_, _, err := c.Get("DELETE", "/templates/"+templateID+"/versions/"+id)
	if err != nil {
		return false, err
	}
	return true, nil
}
