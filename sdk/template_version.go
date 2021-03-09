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
		return nil, fmt.Errorf("failed parsing template version: %w", err)
	}

	return &body, nil
}

// CreateTemplateVersion creates a new version of a transactional template and returns it.
func (c *Client) CreateTemplateVersion(t TemplateVersion) (*TemplateVersion, error) {
	if t.TemplateID == "" {
		return nil, ErrTemplateVersionIDRequired
	}

	if t.Name == "" {
		return nil, ErrTemplateVersionNameRequired
	}

	if t.Subject == "" {
		return nil, ErrTemplateVersionSubjectRequired
	}

	respBody, _, err := c.Post("POST", "/templates/"+t.TemplateID+"/versions", t)
	if err != nil {
		return nil, fmt.Errorf("failed creating template version: %w", err)
	}

	return parseTemplateVersion(respBody)
}

// ReadTemplateVersion retreives a version of a transactional template and returns it.
func (c *Client) ReadTemplateVersion(templateID, id string) (*TemplateVersion, error) {
	if templateID == "" {
		return nil, ErrTemplateVersionIDRequired
	}

	if id == "" {
		return nil, ErrTemplateIDRequired
	}

	respBody, _, err := c.Get("GET", "/templates/"+templateID+"/versions/"+id)
	if err != nil {
		return nil, fmt.Errorf("failed reading template version: %w", err)
	}

	return parseTemplateVersion(respBody)
}

// UpdateTemplateVersion edits a version of a transactional template and returns it.
func (c *Client) UpdateTemplateVersion(t TemplateVersion) (*TemplateVersion, error) {
	if t.ID == "" {
		return nil, ErrTemplateVersionIDRequired
	}

	if t.TemplateID == "" {
		return nil, ErrTemplateIDRequired
	}

	respBody, _, err := c.Post("PATCH", "/templates/"+t.TemplateID+"/versions/"+t.ID, t)
	if err != nil {
		return nil, fmt.Errorf("failed updating template version: %w", err)
	}

	return parseTemplateVersion(respBody)
}

// DeleteTemplateVersion deletes a version of a transactional template.
func (c *Client) DeleteTemplateVersion(templateID, id string) (bool, error) {
	if templateID == "" {
		return false, ErrTemplateVersionIDRequired
	}

	if _, statusCode, err := c.Get("DELETE", "/templates/"+templateID+"/versions/"+id); statusCode > 299 ||
		err != nil {
		return false, fmt.Errorf("failed deleting template version: %w", err)
	}

	return true, nil
}
