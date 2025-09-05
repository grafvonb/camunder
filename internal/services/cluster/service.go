package cluster

import (
	"context"
	"fmt"
	"net/http"

	c87camunda8v2 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda8/v2"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/editors"
)

type Service struct {
	c *c87camunda8v2.ClientWithResponses
}

func New(cfg config.Config, httpClient *http.Client) (*Service, error) {
	c, err := c87camunda8v2.NewClientWithResponses(
		cfg.Camunda8API.BaseURL,
		c87camunda8v2.WithHTTPClient(httpClient),
		c87camunda8v2.WithRequestEditorFn(editors.BearerTokenEditorFn[c87camunda8v2.RequestEditorFn](cfg.Camunda8API.Token)),
	)
	if err != nil {
		return nil, err
	}
	return &Service{c: c}, nil
}

func (s *Service) GetClusterTopology(ctx context.Context) (*c87camunda8v2.Topology, error) {
	resp, err := s.c.GetClusterTopologyWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}
