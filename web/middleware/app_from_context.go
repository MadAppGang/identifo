package middleware

import (
	"context"

	"github.com/madappgang/identifo/model"
)

// AppFromContext returns app data from request conntext.
func AppFromContext(ctx context.Context) model.AppData {
	value := ctx.Value(model.AppDataContextKey)

	if value == nil {
		return nil
	}

	return value.(model.AppData)
}
