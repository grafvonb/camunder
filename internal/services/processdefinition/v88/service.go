package v87

import (
	"context"
	"log/slog"
	"net/http"

	operatev88 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v88"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/processdefinition"
)

// nolint
type Service struct {
	c   *operatev88.ClientWithResponses
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	panic("not implemented for v88")
}

func (s *Service) Capabilities(ctx context.Context) camunda.Capabilities {
	panic("not implemented for v88")
}

func (s *Service) GetProcessDefinitionByKey(ctx context.Context, key int64) (processdefinition.ProcessDefinition, error) {
	panic("not implemented for v88")
}

func (s *Service) SearchProcessDefinitions(ctx context.Context, filter processdefinition.SearchFilterOpts, size int32) (processdefinition.ProcessDefinitions, error) {
	panic("not implemented for v88")
}
