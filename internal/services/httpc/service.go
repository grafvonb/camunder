package httpc

import (
	"context"
	"errors"
	nethttp "net/http"
	"time"

	"github.com/grafvonb/camunder/internal/config"
)

var (
	ErrNoHttpServiceInContext  = errors.New("no http service in context")
	ErrInvalidServiceInContext = errors.New("invalid service in context")
)

type Service struct {
	c   *nethttp.Client
	cfg *config.Config

	isQuiet bool
}

type Option func(*Service)

func WithQuietEnabled(enabled bool) Option {
	return func(s *Service) {
		s.isQuiet = enabled
	}
}

// WithTimeout sets the timeout directly.
func WithTimeout(d time.Duration) Option {
	return func(s *Service) {
		s.c.Timeout = d
	}
}

// WithTimeoutString parses a string like "5s" or "2m" and sets the timeout.
func WithTimeoutString(v string) Option {
	return func(s *Service) {
		if v == "" {
			return
		}
		// swallow error, if parsing fails, just don't set the timeout
		if d, err := time.ParseDuration(v); err == nil {
			s.c.Timeout = d
		}
	}
}

func New(cfg *config.Config, opts ...Option) (*Service, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}
	timeout, err := time.ParseDuration(cfg.HTTP.Timeout)
	if err != nil {
		return nil, err
	}
	httpClient := &nethttp.Client{
		Timeout: timeout,
	}

	s := &Service{
		c:   httpClient,
		cfg: cfg,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

// Client returns the underlying http client
func (s *Service) Client() *nethttp.Client {
	return s.c
}

type ctxHttpServiceKey struct{}

func (s *Service) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxHttpServiceKey{}, s)
}

func FromContext(ctx context.Context) (*Service, error) {
	v := ctx.Value(ctxHttpServiceKey{})
	if v == nil {
		return nil, ErrNoHttpServiceInContext
	}
	s, ok := v.(*Service)
	if !ok || s == nil {
		return nil, ErrInvalidServiceInContext
	}
	return s, nil
}

// MustClient retrieves the http client from the context or returns the default http client
func MustClient(ctx context.Context) *nethttp.Client {
	if s, err := FromContext(ctx); err == nil && s != nil {
		return s.c
	}
	return nethttp.DefaultClient
}
