package v87

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	camundav87 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda/v87"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/editors"
	"github.com/grafvonb/camunder/internal/services/auth"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/cluster"
)

type Service struct {
	c    GenClusterClient
	auth *auth.Service
	cfg  *config.Config
	log  *slog.Logger
}

type Option func(*Service)

func WithClient(c GenClusterClient) Option {
	return func(s *Service) {
		s.c = c
	}
}

func New(cfg *config.Config, httpClient *http.Client, auth *auth.Service, log *slog.Logger, opts ...Option) (*Service, error) {
	c, err := camundav87.NewClientWithResponses(
		cfg.APIs.Camunda.BaseURL,
		camundav87.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	s := &Service{
		c:    c,
		cfg:  cfg,
		auth: auth,
		log:  log,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) GetClusterTopology(ctx context.Context) (cluster.Topology, error) {
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.CamundaApiKeyConst)
	if err != nil {
		return cluster.Topology{}, fmt.Errorf("error retrieving camunda token: %w", err)
	}
	resp, err := s.c.GetClusterTopologyWithResponse(ctx,
		editors.BearerTokenEditorFn[camundav87.RequestEditorFn](token))
	if err != nil {
		return cluster.Topology{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return cluster.Topology{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	ret, err := resp.JSON200.ToStable()
	if err != nil {
		return cluster.Topology{}, fmt.Errorf("convert to stable topology: %w", err)
	}
	return ret, nil
}

func (s *Service) Capabilities(ctx context.Context) camunda.Capabilities {
	panic("not implemented in v87")
}
