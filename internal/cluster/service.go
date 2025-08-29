package cluster

import (
	"context"
	"io"
	"net/http"

	c87camunda8v2 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda8/v2"
)

type Service struct {
	c *c87camunda8v2.Client
}

func New(c *c87camunda8v2.Client) *Service {
	return &Service{
		c: c,
	}
}

func (s *Service) GetTopology(ctx context.Context, cluster *Cluster) error {
	resp, err := s.c.GetClusterTypology(ctx, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)

	}

	return nil
}
