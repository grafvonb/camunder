package v87

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grafvonb/camunder/internal/api/convert"
	operatev88 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v88"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/editors"
	"github.com/grafvonb/camunder/internal/services/auth"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/processdefinition"
)

type Service struct {
	c       *operatev88.ClientWithResponses
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
	c, err := operatev88.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		operatev88.WithHTTPClient(httpClient),
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

func (s *Service) Capabilities(ctx context.Context) camunda.Capabilities {
	panic("not implemented for v87")
}

func (s *Service) GetProcessDefinitionByKey(ctx context.Context, key int64) (processdefinition.ProcessDefinition, error) {
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.OperateApiKeyConst)
	if err != nil {
		return processdefinition.ProcessDefinition{}, fmt.Errorf("error retrieving operate token: %w", err)
	}
	resp, err := s.c.GetProcessDefinitionByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[operatev88.RequestEditorFn](token))
	if err != nil {
		return processdefinition.ProcessDefinition{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return processdefinition.ProcessDefinition{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	ret := resp.JSON200.ToStable()
	return ret, nil
}

func (s *Service) SearchProcessDefinitions(ctx context.Context, filter processdefinition.SearchFilterOpts, size int32) (processdefinition.ProcessDefinitions, error) {
	body := operatev88.QueryProcessDefinition{
		Filter: &operatev88.ProcessDefinition{
			BpmnProcessId: &filter.BpmnProcessId,
			Version:       convert.PtrIfNonZero(filter.Version),
			VersionTag:    &filter.VersionTag,
		},
		Size: &size,
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.OperateApiKeyConst)
	if err != nil {
		return processdefinition.ProcessDefinitions{}, fmt.Errorf("error retrieving operate token: %w", err)
	}
	resp, err := s.c.SearchProcessDefinitionsWithResponse(ctx, body,
		editors.BearerTokenEditorFn[operatev88.RequestEditorFn](token))
	if err != nil {
		return processdefinition.ProcessDefinitions{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return processdefinition.ProcessDefinitions{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200.ToStable(), nil
}
