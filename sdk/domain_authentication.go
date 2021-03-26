package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DomainAuthenticationDns struct {
	MailCNAME DomainAuthenticationDnsValue `json:"mail_cname,omitempty"`
	DKIM1     DomainAuthenticationDnsValue `json:"dkim1,omitempty"`
	DKIM2     DomainAuthenticationDnsValue `json:"dkim2,omitempty"`
}

type DomainAuthenticationDnsValue struct {
	Valid bool   `json:"valid,omitempty"`
	Type  string `json:"type,omitempty"`
	Host  string `json:"host,omitempty"`
	Data  string `json:"data,omitempty"`
}

// DomainAuthentication is a Sendgrid domain authentication
type DomainAuthentication struct {
	ID                 int32                   `json:"id,omitempty"`
	UserID             int32                   `json:"user_id,omitempty"`
	Domain             string                  `json:"domain,omitempty"`
	Subdomain          string                  `json:"subdomain,omitempty"`
	Username           string                  `json:"username,omitempty"`
	IPs                []string                `json:"ips,omitempty"`
	CustomSPF          bool                    `json:"custom_spf"`
	IsDefault          bool                    `json:"default"`
	AutomaticSecurity  bool                    `json:"automatic_security"`
	CustomDKIMSelector string                  `json:"custom_dkim_selector"`
	Legacy             bool                    `json:"legacy,omitempty"`
	Valid              bool                    `json:"valid,omitempty"`
	Dns                DomainAuthenticationDns `json:"dns,omitempty"`
}

func parseDomainAuthentication(respBody string) (*DomainAuthentication, RequestError) {
	var body DomainAuthentication
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing API key: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateDomainAuthentication creates an DomainAuthentication and returns it.
func (c *Client) CreateDomainAuthentication(domain string, subdomain string, ips []string, customSpf bool, isDefault bool, automaticSecurity bool, customDKIMSelector string) (*DomainAuthentication, RequestError) {
	if domain == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrNameRequired,
		}
	}

	respBody, statusCode, err := c.Post("POST", "/whitelabel/domains", DomainAuthentication{
		Domain:             domain,
		Subdomain:          subdomain,
		IPs:                ips,
		CustomSPF:          customSpf,
		IsDefault:          isDefault,
		AutomaticSecurity:  automaticSecurity,
		CustomDKIMSelector: customDKIMSelector,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed creating API key: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedCreatingDomainAuthentication, statusCode, respBody),
		}
	}

	return parseDomainAuthentication(respBody)
}

// ReadDomainAuthentication retrieves an DomainAuthentication and returns it.
func (c *Client) ReadDomainAuthentication(id string) (*DomainAuthentication, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrDomainAuthenticationIDRequired,
		}
	}

	respBody, _, err := c.Get("GET", "/whitelabel/domains/"+id)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseDomainAuthentication(respBody)
}

// UpdateDomainAuthentication edits an DomainAuthentication and returns it.
func (c *Client) UpdateDomainAuthentication(id string, isDefault bool, customSPF bool) (*DomainAuthentication, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrDomainAuthenticationIDRequired,
		}
	}

	t := DomainAuthentication{}
	t.IsDefault = isDefault
	t.CustomSPF = customSPF

	respBody, _, err := c.Post("PATCH", "/whitelabel/domains/"+id, t)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseDomainAuthentication(respBody)
}

func (c *Client) ValidateDomainAuthentication(id string) RequestError {
	if id == "" {
		return RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrDomainAuthenticationIDRequired,
		}
	}

	_, statusCode, err := c.Post("POST", "/whitelabel/domains/"+id+"/validate", nil)
	if err != nil || statusCode != 200 {
		return RequestError{
			StatusCode: statusCode,
			Err:        err,
		}
	}

	return RequestError{StatusCode: http.StatusOK, Err: nil}
}

// DeleteDomainAuthentication deletes an DomainAuthentication.
func (c *Client) DeleteDomainAuthentication(id string) (bool, RequestError) {
	if id == "" {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrDomainAuthenticationIDRequired,
		}
	}

	responseBody, statusCode, err := c.Get("DELETE", "/whitelabel/domains/"+id)
	if err != nil {
		return false, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	if statusCode >= http.StatusMultipleChoices && statusCode != http.StatusNotFound { // ignore not found
		return false, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedDeletingDomainAuthentication, statusCode, responseBody),
		}
	}

	return true, RequestError{StatusCode: http.StatusOK, Err: nil}
}
