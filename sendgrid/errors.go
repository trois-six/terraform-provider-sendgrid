package sendgrid

import "errors"

var (
	ErrCreateRateLimit     = errors.New("expected instance to be created but we were rate limited")
	ErrInvalidImportFormat = errors.New("invalid import. Supported import format: {{templateID}}/{{templateVersionID}}")
)
