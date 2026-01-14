// internal/cache/cache.go
package cache

import "context"

type Cache interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string, dest any) (bool, error)
	Delete(ctx context.Context, key string) error
}
