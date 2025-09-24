package v87

import (
	"context"
	"log/slog"
	"net/http"

	camundav88 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda/v88"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/cluster"
)

type Service struct {
	c   GenClusterClient
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func WithClient(c GenClusterClient) Option { return func(s *Service) { s.c = c } }

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	c, err := camundav88.NewClientWithResponses(
		cfg.APIs.Camunda.BaseURL,
		camundav88.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	s := &Service{c: c, cfg: cfg, log: log}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) GetClusterTopology(ctx context.Context) (cluster.Topology, error) {
	panic("not implemented in v88")
}

func (s *Service) Capabilities(ctx context.Context) camunda.Capabilities {
	panic("not implemented in v88")
}
