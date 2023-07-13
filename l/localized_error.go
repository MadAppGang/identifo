package l

import (
	"fmt"
)

type LocalizedError struct {
	Locale  string
	ErrID   LocalizedString
	Details []any
}

// NewError creates localized error with details with no locale.
func NewError(errID LocalizedString, details ...any) LocalizedError {
	return LocalizedError{
		ErrID:   errID,
		Details: details,
	}
}

// Error returns raw error message. We are missing locale to print the localized version.
func (e LocalizedError) Error() string {
	return fmt.Sprintf("localized error: %v. Details: %v.", e.ErrID, e.Details)
}

// ErrorL returns localized error message.
func (e LocalizedError) ErrorL(p *Printer) string {
	return p.SL(e.Locale, e.ErrID, e.Details...)
}

// Unwrap returns real error to be identified by the new error.As: https://pkg.go.dev/errors#As .
func (e LocalizedError) Unwrap() error {
	if len(e.ErrID) > 0 {
		return e.ErrID
	}
	return nil
}

// ErrorL returns localized error message.
func (e LocalizedError) ErrorWithLocale(locale string) LocalizedError {
	e.Locale = locale
	return e
}

func ErrorWithLocale(err error, locale string) error {
	e, ok := err.(LocalizedError)
	if ok {
		return e.ErrorWithLocale(locale)
	}

	ee, ok := err.(HTTPLocalizedError)
	if ok {
		ee.LE = ee.LE.ErrorWithLocale(locale)
		return ee
	}

	return err
}
