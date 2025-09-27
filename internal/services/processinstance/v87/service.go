package v87

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/grafvonb/camunder/internal/api/convert"
	camundav87 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda/v87"
	operatev87 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v87"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/processinstance"

	"github.com/grafvonb/camunder/internal/config"
)

const wrongStateMessage400 = "Process instances needs to be in one of the states [COMPLETED, CANCELED]"

type Service struct {
	cc  *camundav87.ClientWithResponses
	oc  *operatev87.ClientWithResponses
	cfg *config.Config
	log *slog.Logger
}

type Option func(*Service)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	cc, err := camundav87.NewClientWithResponses(
		cfg.APIs.Camunda.BaseURL,
		camundav87.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	co, err := operatev87.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		operatev87.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	s := &Service{oc: co, cc: cc, cfg: cfg, log: log}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) Capabilities(ctx context.Context) camunda.Capabilities {
	return camunda.Capabilities{
		APIVersion: camunda.V87,
	}
}

func (s *Service) FilterProcessInstanceWithOrphanParent(ctx context.Context, items []processinstance.ProcessInstance) ([]processinstance.ProcessInstance, error) {
	if items == nil {
		return nil, nil
	}
	var result []processinstance.ProcessInstance
	for _, it := range items {
		if it.ParentKey == 0 {
			continue
		}
		_, err := s.GetProcessInstanceByKey(ctx, it.ParentKey)
		if err != nil && strings.Contains(err.Error(), "status 404") {
			result = append(result, it)
		} else if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Service) GetProcessInstanceByKey(ctx context.Context, key int64) (processinstance.ProcessInstance, error) {
	resp, err := s.oc.GetProcessInstanceByKeyWithResponse(ctx, key)
	if err != nil {
		return processinstance.ProcessInstance{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return processinstance.ProcessInstance{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	ret := resp.JSON200.ToStable()
	return ret, nil
}

func (s *Service) GetDirectChildrenOfProcessInstance(ctx context.Context, key int64) (processinstance.ProcessInstances, error) {
	filter := processinstance.SearchFilterOpts{
		ParentKey: key,
	}
	resp, err := s.SearchForProcessInstances(ctx, filter, 1000)
	if err != nil {
		return processinstance.ProcessInstances{}, fmt.Errorf("searching for children of process instance with key %d: %w", key, err)
	}
	return resp, nil
}

func (s *Service) SearchForProcessInstances(ctx context.Context, filter processinstance.SearchFilterOpts, size int32) (processinstance.ProcessInstances, error) {
	st := StateOrNil(filter.State)
	f := operatev87.ProcessInstance{
		TenantId:          &s.cfg.App.Tenant,
		BpmnProcessId:     &filter.BpmnProcessId,
		ProcessVersion:    convert.PtrIfNonZero(filter.ProcessVersion),
		ProcessVersionTag: &filter.ProcessVersionTag,
		State:             st,
		ParentKey:         convert.PtrIfNonZero(filter.ParentKey),
	}
	body := operatev87.SearchProcessInstancesJSONRequestBody{
		Filter: &f,
		Size:   &size,
	}
	resp, err := s.oc.SearchProcessInstancesWithResponse(ctx, body)
	if err != nil {
		return processinstance.ProcessInstances{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return processinstance.ProcessInstances{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200.ToStable(), nil
}

func (s *Service) CancelProcessInstance(ctx context.Context, key int64) (processinstance.CancelResponse, error) {
	s.log.Debug(fmt.Sprintf("trying to cancel process instance with key %d...", key))
	resp, err := s.cc.CancelProcessInstanceWithResponse(ctx, strconv.Itoa(int(key)),
		camundav87.CancelProcessInstanceJSONRequestBody{})
	if err != nil {
		return processinstance.CancelResponse{}, err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return processinstance.CancelResponse{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	s.log.Info(fmt.Sprintf("process instance with key %d was successfully cancelled", key))
	return resp.ToStable(), nil
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key int64) (processinstance.ChangeStatus, error) {
	s.log.Debug(fmt.Sprintf("trying to delete process instance with key %d...", key))
	resp, err := s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key)
	if err != nil {
		return processinstance.ChangeStatus{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return processinstance.ChangeStatus{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	s.log.Info(fmt.Sprintf("process instance with key %d was successfully deleted", key))
	ret := resp.JSON200.ToStable()
	return ret, nil
}

func (s *Service) DeleteProcessInstanceWithCancel(ctx context.Context, key int64) (processinstance.ChangeStatus, error) {
	s.log.Debug(fmt.Sprintf("trying to delete process instance with key %d...", key))
	resp, err := s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key)
	if resp.StatusCode() == http.StatusBadRequest &&
		resp.ApplicationproblemJSON400 != nil &&
		*resp.ApplicationproblemJSON400.Message == wrongStateMessage400 {
		s.log.Info(fmt.Sprintf("process instance with key %d not in state COMPLETED or CANCELED, cancelling it first...", key))
		_, err = s.CancelProcessInstance(ctx, key)
		if err != nil {
			return processinstance.ChangeStatus{}, fmt.Errorf("error cancelling process instance with key %d: %w", key, err)
		}
		s.log.Info(fmt.Sprintf("waiting for process instance with key %d to be cancelled by workflow engine...", key))
		if err = s.WaitForProcessInstanceState(ctx, key, processinstance.StateCanceled); err != nil {
			return processinstance.ChangeStatus{}, fmt.Errorf("waiting for canceled state failed for %d: %w", key, err)
		}
		resp, err = s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key)
	}
	if err != nil {
		return processinstance.ChangeStatus{}, err
	}
	if resp.StatusCode() != http.StatusOK {
		return processinstance.ChangeStatus{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	s.log.Info(fmt.Sprintf("process instance with key %d was successfully deleted", key))
	ret := resp.JSON200.ToStable()
	return ret, nil
}
