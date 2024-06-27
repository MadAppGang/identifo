package logging

import (
	"log/slog"
	"strconv"
)

const (
	FieldComponent = "component"
	FieldError     = "err"
	FieldErrors    = "errs"
	FieldErrorID   = "errId"
	FieldAppID     = "appId"
	FieldUserID    = "userId"
	FieldEmail     = "email"
	FieldURL       = "url"
)

const (
	ComponentAPI        = "API"
	ComponentAdmin      = "ADMIN"
	ComponentCommon     = "COMMON"
	ComponentManagement = "MANAGEMENT"
)

type LogErrors []error

func (e LogErrors) LogValue() slog.Value {
	var errs []slog.Attr
	for i, err := range e {
		errs = append(errs, slog.String(FieldError+"_"+strconv.Itoa(i), err.Error()))
	}

	return slog.GroupValue(errs...)
}
