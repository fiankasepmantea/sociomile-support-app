package utils

import "context"

type BlacklistStore interface {
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}
