package cluster

import (
	"context"

	"github.com/grafvonb/camunder/internal/httpcore"
)

type Service struct {
	c *httpcore.Client
}

func New(c *httpcore.Client) *Service {
	return &Service{
		c: c,
	}
}

func (s *Service) GetTopology(ctx context.Context) (*Cluster, error) {
	var topology Cluster
	err := s.c.Get(ctx, "/v2/topology", &topology)
	if err != nil {
		return nil, err
	}
	return &topology, nil
}
