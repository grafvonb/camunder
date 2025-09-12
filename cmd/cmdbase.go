package cmd

import (
	"context"
	"fmt"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth"
	"github.com/grafvonb/camunder/internal/services/httpc"
)

type Services struct {
	Config *config.Config
	HTTP   *httpc.Service
	Auth   *auth.Service
}

func NewFromContext(ctx context.Context) (*Services, error) {
	cfg, err := config.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("config from context: %w", err)
	}
	httpSvc, err := httpc.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("http service from context: %w", err)
	}
	authSvc, err := auth.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth service from context: %w", err)
	}
	return &Services{
		Config: cfg,
		HTTP:   httpSvc,
		Auth:   authSvc,
	}, nil
}
