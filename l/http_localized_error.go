package l

import (
	"fmt"
	"time"
)

type HTTPLocalizedError struct {
	LE     LocalizedError
	Status int
	Time   time.Time
}

func (e HTTPLocalizedError) Error() string {
	return fmt.Sprintf("[%v] HTTP error: %v (status: %d).", e.Time, e.LE.ErrID, e.Status)
}

// ErrorL returns localized error message.
func (e HTTPLocalizedError) ErrorL(p *Printer) string {
	return e.LE.ErrorL(p)
}

func (e HTTPLocalizedError) Unwrap() error {
	return e.LE
}
