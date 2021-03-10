package sendgrid

import "errors"

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
