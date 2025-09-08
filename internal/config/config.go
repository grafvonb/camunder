package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
)

var (
	ErrNoBaseURL      = errors.New("no base_url provided in api configuration")
	ErrNoTokenURL     = errors.New("no token_url provided in auth configuration")
	ErrNoClientID     = errors.New("no client_id provided in auth configuration")
	ErrNoClientSecret = errors.New("no client_secret provided in auth configuration")

	ErrNoConfigInContext       = errors.New("no config in context")
	ErrInvalidServiceInContext = errors.New("invalid config in context")
)

type Config struct {
	Config string `mapstructure:"config"`

	App  App            `mapstructure:"app"`
	Auth Authentication `mapstructure:"auth"`
	APIs APIs           `mapstructure:"apis"`
	HTTP HTTP           `mapstructure:"http"`
}

func (c *Config) String() string {
	var alias Config
	alias.Config = c.Config
	alias.App = c.App
	alias.HTTP = c.HTTP
	alias.APIs.Camunda8.Key = c.APIs.Camunda8.Key
	alias.APIs.Camunda8.BaseURL = c.APIs.Camunda8.BaseURL
	alias.APIs.Operate.Key = c.APIs.Operate.Key
	alias.APIs.Operate.BaseURL = c.APIs.Operate.BaseURL
	alias.APIs.Tasklist.Key = c.APIs.Tasklist.Key
	alias.APIs.Tasklist.BaseURL = c.APIs.Tasklist.BaseURL
	alias.Auth.TokenURL = c.Auth.TokenURL
	alias.Auth.ClientID = "******"
	alias.Auth.ClientSecret = "******"
	alias.Auth.Scopes = maps.Clone(c.Auth.Scopes)

	b, err := json.MarshalIndent(alias, "", "  ")
	if err != nil {
		return fmt.Sprintf("error marshaling config: %v", err)
	}
	return string(b)
}

// Validate checks all nested sections and aggregates errors.
func (c *Config) Validate() error {
	var errs []error

	if err := c.Auth.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("auth:\n%w", err))
	}
	if err := c.APIs.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("apis:\n%w", err))
	}
	if err := c.HTTP.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("http:\n%w", err))
	}

	return errors.Join(errs...)
}

type ctxConfigKey struct{}

func (c *Config) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxConfigKey{}, c)
}

func FromContext(ctx context.Context) (*Config, error) {
	v := ctx.Value(ctxConfigKey{})
	if v == nil {
		return nil, ErrNoConfigInContext
	}
	c, ok := v.(*Config)
	if !ok || c == nil {
		return nil, ErrInvalidServiceInContext
	}
	return c, nil
}
