package cluster

import (
	"context"
	"fmt"
	"net/http"

	c87camunda8v2 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda8/v2"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/editors"
	"github.com/grafvonb/camunder/internal/services/auth"
)

type Service struct {
	c       *c87camunda8v2.ClientWithResponses
	auth    *auth.Service
	cfg     *config.Config
	isQuiet bool
}

type Option func(*Service)

func WithQuietEnabled(enabled bool) Option {
	return func(s *Service) {
		s.isQuiet = enabled
	}
}

func New(cfg *config.Config, httpClient *http.Client, auth *auth.Service, opts ...Option) (*Service, error) {
	c, err := c87camunda8v2.NewClientWithResponses(
		cfg.APIs.Camunda8.BaseURL,
		c87camunda8v2.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	s := &Service{
		c:    c,
		cfg:  cfg,
		auth: auth,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) GetClusterTopology(ctx context.Context) (*c87camunda8v2.Topology, error) {
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.Camunda8ApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving camunda8 token: %w", err)
	}
	resp, err := s.c.GetClusterTopologyWithResponse(ctx,
		editors.BearerTokenEditorFn[c87camunda8v2.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}
