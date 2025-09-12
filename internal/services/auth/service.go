package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/grafvonb/camunder/internal/api/gen/clients/auth"
	"github.com/grafvonb/camunder/internal/config"
)

const (
	formCT = "application/x-www-form-urlencoded"
)

var (
	ErrNoAuthServiceInContext  = errors.New("no auth service in context")
	ErrInvalidServiceInContext = errors.New("invalid service in context")
)

type Service struct {
	c   *auth.ClientWithResponses
	cfg *config.Config

	mu      sync.Mutex
	cache   map[string]string
	isQuiet bool
}

type Option func(*Service)

func WithQuietEnabled(enabled bool) Option {
	return func(s *Service) {
		s.isQuiet = enabled
	}
}

func New(cfg *config.Config, httpClient *http.Client, opts ...Option) (*Service, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c, err := auth.NewClientWithResponses(cfg.Auth.TokenURL, auth.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("init auth client: %w", err)
	}
	s := &Service{
		c:     c,
		cfg:   cfg,
		cache: make(map[string]string),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) Warmup(ctx context.Context) error {
	_, err := s.retrieveTokensForAPIs(ctx, config.ValidAPIKeys...)
	return err
}

func (s *Service) ClearCache() {
	s.mu.Lock()
	s.cache = make(map[string]string)
	s.mu.Unlock()
}

func (s *Service) RetrieveTokenForAPI(ctx context.Context, target string) (string, error) {
	s.mu.Lock()
	if tok, ok := s.cache[target]; ok && tok != "" {
		s.mu.Unlock()
		return tok, nil
	}
	s.mu.Unlock() // no defer, we want to release the lock before the network call

	scope := s.cfg.Auth.Scope(target)
	tok, err := s.requestToken(ctx, s.cfg.Auth.ClientID, s.cfg.Auth.ClientSecret, scope)
	if err != nil {
		return "", fmt.Errorf("retrieve token for %s: %w", target, err)
	}

	s.mu.Lock()
	s.cache[target] = tok
	s.mu.Unlock()
	return tok, nil
}

func (s *Service) retrieveTokensForAPIs(ctx context.Context, targets ...string) (map[string]string, error) {
	out := make(map[string]string, len(targets))
	for _, t := range targets {
		tok, err := s.RetrieveTokenForAPI(ctx, t)
		if err != nil {
			return nil, err
		}
		out[t] = tok
	}
	return out, nil
}

type ctxAuthServiceKey struct{}

func (s *Service) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxAuthServiceKey{}, s)
}

func FromContext(ctx context.Context) (*Service, error) {
	v := ctx.Value(ctxAuthServiceKey{})
	if v == nil {
		return nil, ErrNoAuthServiceInContext
	}
	s, ok := v.(*Service)
	if !ok || s == nil {
		return nil, ErrInvalidServiceInContext
	}
	return s, nil
}

func (s *Service) requestToken(ctx context.Context, clientID, clientSecret, scope string) (string, error) {
	body := formBody(clientID, clientSecret, scope)
	resp, err := s.c.RequestTokenWithBodyWithResponse(ctx, formCT, body)

	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("nil token response")
	}
	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusMultipleChoices {
		return "", fmt.Errorf("token request failed: status=%d body=%s", resp.StatusCode(), string(resp.Body))
	}
	if resp.JSON200 == nil || resp.JSON200.AccessToken == "" {
		return "", fmt.Errorf("missing access token in successful response (status=%d)", resp.StatusCode())
	}
	return resp.JSON200.AccessToken, nil
}

func formBody(clientID, clientSecret, scope string) io.Reader {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	if strings.TrimSpace(scope) != "" {
		form.Set("scope", scope)
	}
	return strings.NewReader(form.Encode())
}
