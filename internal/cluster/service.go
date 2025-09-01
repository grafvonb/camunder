package cluster

import (
	"context"
	"net/http"

	c87camunda8v2 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda8/v2"
)

type Service struct {
	c *c87camunda8v2.ClientWithResponses
}

func New(baseUrl string, httpClient *http.Client, token string) (*Service, error) {
	authEditor := func(ctx context.Context, req *http.Request) error {
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		return nil
	}

	c, err := c87camunda8v2.NewClientWithResponses(
		baseUrl,
		c87camunda8v2.WithHTTPClient(httpClient),
		c87camunda8v2.WithRequestEditorFn(authEditor),
	)
	if err != nil {
		return nil, err
	}
	return &Service{c: c}, nil
}

func (s *Service) GetTopology(ctx context.Context) (*Cluster, error) {
	resp, err := s.c.GetClusterTypologyWithResponse(ctx, nil)
	if err != nil {
		return nil, err
	}

	return *resp.JSON200, nil
}
