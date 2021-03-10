package sendgrid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type creditAllocation struct {
	Type string `json:"credit_allocation,omitempty"`
}

// SubUser is a Sendgrid SubUser.
type SubUser struct {
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

type subUserError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type subUserErrors struct {
	Errors []subUserError `json:"errors,omitempty"`
}

func parseSubUser(respBody string) (*SubUser, RequestError) {
	var body SubUser
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		log.Printf("[DEBUG] [parseSubUser] failed parsing subUser, response body: %s", respBody)
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing subUsers: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func parseSubUsers(respBody string) ([]SubUser, RequestError) {
	var body []SubUser
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		log.Printf("[DEBUG] [parseSubUsers] failed parsing subUsers, response body: %s", respBody)
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing subUsers: %w", err),
		}
	}

	return body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateSubuser creates a subuser and returns it.
func (c *Client) CreateSubuser(username, email, password string, ips []string) (*SubUser, RequestError) {
	if username == "" {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	if email == "" {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrEmailRequired}
	}

	if password == "" {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrPasswordRequired}
	}

	if len(ips) < 1 {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrIPRequired}
	}

	respBody, statusCode, err := c.Post("POST", "/subusers", SubUser{
		UserName: username,
		Email:    email,
		Password: password,
		IPs:      ips,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed creating subUser: %w", err),
		}
	}

	if statusCode >= 300 {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed creating subUser, status: %d, response: %s", statusCode, respBody),
		}
	}

	return parseSubUser(respBody)
}

// ReadSubuser retreives a subuser and returns it.
func (c *Client) ReadSubUser(username string) ([]SubUser, RequestError) {
	if username == "" {
		return nil, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	endpoint := "/subusers?username=" + url.QueryEscape(username)

	respBody, statusCode, err := c.Get("GET", endpoint)
	if err != nil {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed reading subUser: %w", err),
		}
	}

	return parseSubUsers(respBody)
}

// UpdateSubuser enables/disables a subuser.
func (c *Client) UpdateSubuser(username string, disabled bool) (bool, RequestError) {
	if username == "" {
		return false, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	respBody, statusCode, err := c.Post("PATCH", "/subusers/"+username, SubUser{
		Disabled: disabled,
	})
	if err != nil {
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed updating subUser: %w", err),
		}
	}

	var body subUserErrors
	if err = json.Unmarshal([]byte(respBody), &body); err != nil {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed updating subUser: %w", err),
		}
	}

	return len(body.Errors) == 0, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// DeleteSubuser deletes a subuser.
func (c *Client) DeleteSubuser(username string) (bool, RequestError) {
	if username == "" {
		return false, RequestError{StatusCode: http.StatusNotAcceptable, Err: ErrUsernameRequired}
	}

	if _, statusCode, err := c.Get("DELETE", "/subusers/"+username); statusCode > 299 ||
		err != nil {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed deleting subUser: %w", err),
		}
	}

	return true, RequestError{StatusCode: http.StatusOK, Err: nil}
}
