package systemservice

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/camunder/internal/api/gen/clients/uamoim"
	"github.com/grafvonb/camunder/internal/config"
)

type Service struct {
	c   *uamoim.ClientWithResponses
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	c, err := uamoim.NewClientWithResponses(
		cfg.APIs.Camunda.BaseURL,
		uamoim.WithHTTPClient(httpClient),
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

func (s *Service) GetSystems(ctx context.Context) (*[]uamoim.SystemSummary, error) {
	resp, err := s.c.ListSystemsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (s *Service) GetSystemByKey(ctx context.Context, key string) (*uamoim.SystemDetail, error) {
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}
	resp, err := s.c.GetSystemWithResponse(ctx, key)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}
