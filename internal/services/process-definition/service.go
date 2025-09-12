package processdefinition

import (
	"context"
	"fmt"
	"net/http"

	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/editors"
	"github.com/grafvonb/camunder/internal/services/auth"
)

type Service struct {
	c       *c87operatev1.ClientWithResponses
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
	c, err := c87operatev1.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		c87operatev1.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	s := &Service{
		c:    c,
		cfg:  cfg,
		auth: auth,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) GetProcessDefinitionByKey(ctx context.Context, key int64) (*c87operatev1.ProcessDefinitionItem, error) {
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.OperateApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving operate token: %w", err)
	}
	resp, err := s.c.GetProcessDefinitionByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (s *Service) SearchForProcessDefinitions(ctx context.Context, filter SearchFilterOpts, size int32) (*c87operatev1.ProcessDefinitionSearchResponse, error) {
	body := c87operatev1.ProcessDefinitionSearchRequest{
		Filter: &c87operatev1.ProcessDefinitionFilter{
			BpmnProcessId: filter.BpmnProcessId,
			Version:       filter.Version,
			VersionTag:    filter.VersionTag,
		},
		Size: &size,
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.OperateApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving operate token: %w", err)
	}
	resp, err := s.c.SearchForProcessDefinitionsWithResponse(ctx, body,
		editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}
