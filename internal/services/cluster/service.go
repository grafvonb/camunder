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

func New(cfg *config.Config, httpClient *http.Client, auth *auth.Service, isQuiet bool) (*Service, error) {
	c, err := c87camunda8v2.NewClientWithResponses(
		cfg.APIs.Camunda8.BaseURL,
		c87camunda8v2.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	return &Service{
		c:       c,
		cfg:     cfg,
		auth:    auth,
		isQuiet: isQuiet,
	}, nil
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
