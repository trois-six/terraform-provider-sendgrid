package sendgrid

import (
	"errors"
	"fmt"
)

var (
	// ErrCreateRateLimit error displayed when we reach the API calls rate limit.
	ErrCreateRateLimit = errors.New("expected instance to be created but we were rate limited")

	// ErrInvalidImportFormat error displayed when the string passed to import a template version
	// doesn't have the good format.
	ErrInvalidImportFormat = errors.New("invalid import. Supported import format: {{templateID}}/{{templateVersionID}}")

	// ErrSubUserNotFound error displayed when the subUser can not be found.
	ErrSubUserNotFound = errors.New("subUser wasn't found")
)

func subUserNotFound(name string) error {
	return fmt.Errorf("%w: %s", ErrSubUserNotFound, name)
}
