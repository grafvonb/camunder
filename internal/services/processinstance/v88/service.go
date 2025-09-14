package v88

import (
	"context"
	"net/http"

	camundav88 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda/v88"
	operatev88 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v88"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/procesinstance"
)

type Service struct {
	cc      *camundav88.ClientWithResponses
	oc      *operatev88.ClientWithResponses
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
	cc, err := camundav88.NewClientWithResponses(
		cfg.APIs.Camunda.BaseURL,
		camundav88.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	co, err := operatev88.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		operatev88.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	s := &Service{
		oc:   co,
		cc:   cc,
		auth: auth,
		cfg:  cfg,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) Capabilities(ctx context.Context) camunda.Capabilities {
	return camunda.Capabilities{
		APIVersion: camunda.V88,
	}
}

func (s *Service) WaitForProcessInstanceState(ctx context.Context, key string, desiredState string) error {
	panic("not implemented in v88")
}

func (s *Service) CancelProcessInstance(ctx context.Context, key int64) (*processinstance.CancelResponse, error) {
	panic("not implemented in v88")
}
