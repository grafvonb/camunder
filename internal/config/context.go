package config

import "context"

type ctxKey struct{}

func IntoContext(ctx context.Context, c Config) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

func FromContext(ctx context.Context) (Config, bool) {
	c, ok := ctx.Value(ctxKey{}).(Config)
	return c, ok
}

// Convenience: panic if missing (useful if you always load in root)
func MustFrom(ctx context.Context) Config {
	c, ok := FromContext(ctx)
	if !ok {
		panic("config: not found in context")
	}
	return c
}
