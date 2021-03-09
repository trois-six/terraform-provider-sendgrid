package sendgrid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

type creditAllocation struct {
	Type string `json:"credit_allocation,omitempty"`
}

type RequestError struct {
	StatusCode int

	Err error
}

// Subuser is a Sendgrid Subuser.
type Subuser struct {
	ID                 int              `json:"id,omitempty"`
	UserID             int              `json:"user_id,omitempty"`
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
	log.Printf("[DEBUG] Parsing sub users, response body: %s", respBody)
	err := json.Unmarshal([]byte(respBody), &body)
	if err != nil {
		log.Printf("[ERROR] Couldn't parse sub users")
		return nil, err
	}
	return &body[0], nil
}

// CreateSubuser creates a subuser and returns it.
func (c *Client) CreateSubuser(username, email, password string, ips []string) (*Subuser, *RequestError) {
	if username == "" {
		return nil, GenericError(fmt.Errorf("[CreateSubuser] a username is required"))
	}
	if email == "" {
		return nil, GenericError(fmt.Errorf("[CreateSubuser] an email is required"))
	}
	if password == "" {
		return nil, GenericError(fmt.Errorf("[CreateSubuser] a password is required"))
	}
	if len(ips) < 1 {
		return nil, GenericError(fmt.Errorf("[CreateSubuser] at least one ip address is required"))
	}
	respBody, statusCode, err := c.Post("POST", "/subusers", Subuser{
		UserName: username,
		Email:    email,
		Password: password,
		IPs:      ips,
	})

	log.Printf("[DEBUG] [CreateSubuser] status: %d, body: %s", statusCode, respBody)

	if err != nil {
		return nil, &RequestError{
			StatusCode: statusCode,
			Err:        err,
		}
	}

	if statusCode >= 300 {
		return nil, &RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("[CreateSubuser] error returned from API, status: %d, body: %s", statusCode, respBody),
		}
	}

	subuser, err := parseSubuser(respBody)
	return subuser, GenericError(err)
}

// ReadSubuser retreives a subuser and returns it.
func (c *Client) ReadSubuser(username string) (*Subuser, error) {
	if username == "" {
		return nil, fmt.Errorf("[ReadSubuser] a username is required")
	}

	endpoint := "/subusers?username=" + url.QueryEscape(username)

	log.Printf("[DEBUG] Searching for user at: %s", endpoint)

	respBody, _, err := c.Get("GET", endpoint)
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

func (r *RequestError) Error() string {
	return fmt.Sprintf("status: %d, err: %s", r.StatusCode, r.Err)
}

func GenericError(error error) *RequestError {
	return &RequestError{
		StatusCode: 500,
		Err:        error,
	}
}

// DeleteSubuser deletes a subuser.
func (c *Client) DeleteSubuser(username string) (bool, *RequestError) {
	if username == "" {
		return false, GenericError(fmt.Errorf("[DeleteSubuser] a username is required"))
	}
	responseBody, statusCode, err := c.Get("DELETE", "/subusers/"+username)
	if err != nil {
		return false, GenericError(err)
	}

	if statusCode > 299 && statusCode != 404 {
		return false, &RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("[DeleteSubuser] error deleting user: %s", username),
		}
	}

	log.Printf("[DEBUG] [DeleteSubuser] status code: %d, responseBody: %s", statusCode, responseBody)

	return true, nil
}
