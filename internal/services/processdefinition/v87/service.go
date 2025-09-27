package v87

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/camunder/internal/api/convert"
	operatev87 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v87"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/processdefinition"
)

type Service struct {
	c   *operatev87.ClientWithResponses
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	c, err := operatev87.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		operatev87.WithHTTPClient(httpClient),
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

func (s *Service) Capabilities(ctx context.Context) camunda.Capabilities {
	panic("not implemented for v88")
}

func (s *Service) GetProcessDefinitionByKey(ctx context.Context, key int64) (processdefinition.ProcessDefinition, error) {
	resp, err := s.c.GetProcessDefinitionByKeyWithResponse(ctx, key)
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
	body := operatev87.QueryProcessDefinition{
		Filter: &operatev87.ProcessDefinition{
			BpmnProcessId: &filter.BpmnProcessId,
			Version:       convert.PtrIfNonZero(filter.Version),
			VersionTag:    &filter.VersionTag,
		},
		Size: &size,
	}
	resp, err := s.c.SearchProcessDefinitionsWithResponse(ctx, body)
	if err != nil {
		return processdefinition.ProcessDefinitions{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return processdefinition.ProcessDefinitions{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200.ToStable(), nil
}
