package errs

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrTooManyRequests = errors.New("too many requests")
	ErrPanic           = errors.New("panic")
	ErrMediaType       = errors.New("unsupported media type")

	ErrUnsupportedDBType = errors.New("unsupported database type")

	ErrSetPair      = errors.New("failed to set pair")
	ErrGetPair      = errors.New("failed to get pair")
	ErrCastValue    = errors.New("failed to cast value")
	ErrKeyNotFound  = errors.New("key not found")
	ErrUserNotFound = errors.New("user not found")

	ErrJSONDecode = errors.New("failed to decode JSON")
	ErrSetPairs   = errors.New("failed to set pairs")
	ErrGetPairs   = errors.New("failed to get pairs")
	ErrLogin      = errors.New("failed to login")

	ErrEmptyPairs = errors.New("empty pairs")

	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
	ErrUnauthorized         = errors.New("unauthorized")
)

func WrapError(wrap, err error) error {
	if err == nil {
		return wrap
	}

	return fmt.Errorf("%s: %w", wrap.Error(), err)
}

type MultiError struct {
	Errors []error
}

func (m *MultiError) Error() string {
	msgs := make([]string, len(m.Errors))
	for i, err := range m.Errors {
		msgs[i] = err.Error()
	}

	return fmt.Sprintf("multiple errors: %s", strings.Join(msgs, " & "))
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}
