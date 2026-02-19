package appcontext

import (
	"context"
	"errors"
)

type contextKey string

const appContextKey contextKey = "app_context"

type AppContext struct {
	UserID   uint   `json:"user_id"`
	TenantID uint   `json:"tenant_id"`
	Role     string `json:"role"`
}

func WithAppContext(ctx context.Context, ac *AppContext) context.Context {
	return context.WithValue(ctx, appContextKey, ac)
}

func FromContext(ctx context.Context) (*AppContext, error) {
	ac, ok := ctx.Value(appContextKey).(*AppContext)
	if !ok || ac == nil {
		return nil, errors.New("app context not found")
	}
	return ac, nil
}