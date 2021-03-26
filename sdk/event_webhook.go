package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// EventWebhook is a Sendgrid event webhook settings.
type EventWebhook struct {
	Enabled           bool   `json:"enabled"`
	Url               string `json:"url,omitempty"`
	GroupResubscribe  bool   `json:"group_resubscribe"`
	Delivered         bool   `json:"delivered"`
	GroupUnsubscribe  bool   `json:"group_unsubscribe"`
	SpamReport        bool   `json:"spam_report"`
	Bounce            bool   `json:"bounce"`
	Deferred          bool   `json:"deferred"`
	Unsubscribe       bool   `json:"unsubscribe"`
	Processed         bool   `json:"processed"`
	Open              bool   `json:"open"`
	Click             bool   `json:"click"`
	Dropped           bool   `json:"dropped"`
	OAuthClientId     string `json:"oauth_client_id,omitempty"`
	OAuthClientSecret string `json:"oauth_client_secret,omitempty"`
	OAuthTokenUrl     string `json:"oauth_token_url,omitempty"`
}

func parseEventWebhook(respBody string) (*EventWebhook, RequestError) {
	var body EventWebhook
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing event webhook: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateEventWebhook creates an EventWebhook and returns it.
func (c *Client) PatchEventWebhook(enabled bool, url string, groupResubscribe bool, delivered bool, groupUnsubscribe bool, spamReport bool, bounce bool, deferred bool, unsubscribe bool, processed bool, open bool, click bool, dropped bool, oauthClientId string, oauthClientSecret string, oauthTokenUrl string) (*EventWebhook, RequestError) {
	if url == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrUrlRequired,
		}
	}

	respBody, statusCode, err := c.Post("PATCH", "/user/webhooks/event/settings", EventWebhook{
		Enabled:           enabled,
		Url:               url,
		GroupResubscribe:  groupResubscribe,
		Delivered:         delivered,
		GroupUnsubscribe:  groupUnsubscribe,
		SpamReport:        spamReport,
		Bounce:            bounce,
		Deferred:          deferred,
		Unsubscribe:       unsubscribe,
		Processed:         processed,
		Open:              open,
		Click:             click,
		Dropped:           dropped,
		OAuthClientId:     oauthClientId,
		OAuthClientSecret: oauthClientSecret,
		OAuthTokenUrl:     oauthTokenUrl,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed patching event webhook: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedPatchingEventWebhook, statusCode, respBody),
		}
	}

	return parseEventWebhook(respBody)
}

// ReadEventWebhook retreives an EventWebhook and returns it.
func (c *Client) ReadEventWebhook() (*EventWebhook, RequestError) {
	respBody, _, err := c.Get("GET", "/user/webhooks/event/settings")
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseEventWebhook(respBody)
}
