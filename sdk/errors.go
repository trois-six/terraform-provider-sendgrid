package sendgrid

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrNameRequired                   = errors.New("a name is required")
	ErrAPIKeyIDRequired               = errors.New("an API Key ID is required")
	ErrBodyNotNil                     = errors.New("body must not be nil")
	ErrUsernameRequired               = errors.New("a username is required")
	ErrEmailRequired                  = errors.New("an email is required")
	ErrPasswordRequired               = errors.New("a password is required")
	ErrIPRequired                     = errors.New("at least one ip address is required")
	ErrTemplateIDRequired             = errors.New("a template ID is required")
	ErrTemplateNameRequired           = errors.New("a template name is required")
	ErrTemplateVersionIDRequired      = errors.New("a template version ID is required")
	ErrTemplateVersionNameRequired    = errors.New("a template version name is required")
	ErrTemplateVersionSubjectRequired = errors.New("a template version subject is required")
)

type RequestError struct {
	StatusCode int
	Err        error
}

func RetryOnRateLimit(
	ctx context.Context, d *schema.ResourceData, f func() (interface{}, RequestError)) (interface{}, error) {
	var resp interface{}

	err := resource.RetryContext(
		ctx,
		d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			var requestErr RequestError
			resp, requestErr = f()
			if requestErr.Err != nil {
				if requestErr.StatusCode == http.StatusTooManyRequests {
					return resource.RetryableError(requestErr.Err)
				}

				return resource.NonRetryableError(requestErr.Err)
			}

			return nil
		})
	if err != nil {
		return resp, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}
