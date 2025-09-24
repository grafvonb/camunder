package v88

import (
	"context"
	"log/slog"
	"net/http"

	camundav88 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda/v88"
	operatev88 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v88"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/processinstance"
)

// nolint
type Service struct {
	cc  *camundav88.ClientWithResponses
	oc  *operatev88.ClientWithResponses
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	panic("not implemented in v88")
}

func (s *Service) Capabilities(ctx context.Context) camunda.Capabilities {
	panic("not implemented in v88")
}

func (s *Service) GetProcessInstanceByKey(ctx context.Context, key int64) (processinstance.ProcessInstance, error) {
	panic("not implemented in v88")
}

func (s *Service) WaitForProcessInstanceState(ctx context.Context, key int64, desiredState processinstance.State) error {
	panic("not implemented in v88")
}

func (s *Service) CancelProcessInstance(ctx context.Context, key int64) (processinstance.CancelResponse, error) {
	panic("not implemented in v88")
}

func (s *Service) SearchForProcessInstances(ctx context.Context, filter processinstance.SearchFilterOpts, size int32) (processinstance.ProcessInstances, error) {
	panic("not implemented in v88")
}

func (s *Service) GetDirectChildrenOfProcessInstance(ctx context.Context, key int64) (processinstance.ProcessInstances, error) {
	panic("not implemented in v88")
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key int64) (processinstance.ChangeStatus, error) {
	panic("not implemented in v88")
}

func (s *Service) FilterProcessInstanceWithOrphanParent(ctx context.Context, items []processinstance.ProcessInstance) ([]processinstance.ProcessInstance, error) {
	panic("not implemented in v88")
}

func (s *Service) DeleteProcessInstanceWithCancel(ctx context.Context, key int64) (processinstance.ChangeStatus, error) {
	panic("not implemented in v88")
}
