package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type creditAllocation struct {
	Type string `json:"credit_allocation,omitempty"`
}

// Subuser is a Sendgrid Subuser.
type Subuser struct {
	ID                 string           `json:"id,omitempty"`
	UserID             string           `json:"user_id,omitempty"`
	UserName           string           `json:"username,omitempty"`
	Password           string           `json:"password,omitempty"`
	Email              string           `json:"email,omitempty"`
	IPs                []string         `json:"ips,omitempty"`
	Disabled           bool             `json:"disabled,omitempty"`
	SignupSessionToken string           `json:"signup_session_token,omitempty"`
	AuthorizationToken string           `json:"authorization_token,omitempty"`
	CreditAllocation   creditAllocation `json:"credit_allocation,omitempty"`
}

type subuserError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}
type subuserErrors struct {
	Errors []subuserError `json:"errors,omitempty"`
}

func parseSubuser(respBody string) (*Subuser, error) {
	var body Subuser
	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func parseSubusers(respBody string) (*Subuser, error) {
	var body []Subuser
	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return nil, err
	}
	return &body[0], nil
}

// CreateSubuser creates a subuser and returns it.
func (c *Client) CreateSubuser(username, email, password string, ips []string) (*Subuser, error) {
	if username == "" {
		return nil, fmt.Errorf("[CreateSubuser] a username is required")
	}
	if email == "" {
		return nil, fmt.Errorf("[CreateSubuser] an email is required")
	}
	if password == "" {
		return nil, fmt.Errorf("[CreateSubuser] a password is required")
	}
	if len(ips) < 1 {
		return nil, fmt.Errorf("[CreateSubuser] at least one ip address is required")
	}
	respBody, _, err := c.Post("POST", "/subusers", Subuser{
		UserName: username,
		Email:    email,
		Password: password,
		IPs:      ips,
	})
	if err != nil {
		return nil, err
	}
	return parseSubuser(respBody)
}

// ReadSubuser retreives a subuser and returns it.
func (c *Client) ReadSubuser(username string) (*Subuser, error) {
	if username == "" {
		return nil, fmt.Errorf("[ReadSubuser] a username is required")
	}
	respBody, _, err := c.Get("GET", "/subusers?username="+url.QueryEscape(username))
	if err != nil {
		return nil, err
	}
	return parseSubusers(respBody)
}

// UpdateSubuser enables/disables a subuser.
func (c *Client) UpdateSubuser(username string, disabled bool) (bool, error) {
	if username == "" {
		return false, fmt.Errorf("[UpdateSubuser] a name is required")
	}
	respBody, _, err := c.Post("PATCH", "/subusers/"+username, Subuser{
		Disabled: disabled,
	})
	if err != nil {
		return false, err
	}
	var body subuserErrors
	err = json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		return false, err
	}
	return len(body.Errors) == 0, nil
}

// DeleteSubuser deletes a subuser.
func (c *Client) DeleteSubuser(username string) (bool, error) {
	if username == "" {
		return false, fmt.Errorf("[DeleteSubuser] a username is required")
	}
	_, _, err := c.Get("DELETE", "/Subusers/"+username)
	if err != nil {
		return false, err
	}
	return true, nil
}
