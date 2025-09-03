package config

import (
	"context"
	"errors"
)

var ErrNoConfigInContext = errors.New("no config in context")

type ctxKey struct{}

func IntoContext(ctx context.Context, c Config) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

func FromContext(ctx context.Context) (Config, error) {
	c, ok := ctx.Value(ctxKey{}).(Config)
	if !ok {
		return c, ErrNoConfigInContext
	}
	return c, nil
}
