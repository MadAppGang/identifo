package middleware

import (
	"context"

	"github.com/madappgang/identifo/v2/model"
)

// AppFromContext returns app data from request conntext.
func AppFromContext(ctx context.Context) model.AppData {
	value := ctx.Value(model.AppDataContextKey)

	if value == nil {
		return model.AppData{}
	}

	return value.(model.AppData)
}
