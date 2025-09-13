package processinstance

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/grafvonb/camunder/internal/api/gen/clients/camunda/c87camunda"
	"github.com/grafvonb/camunder/internal/api/gen/clients/camunda/c87operate"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/editors"
	"github.com/grafvonb/camunder/internal/services/auth"
)

const wrongStateMessage400 = "Process instances needs to be in one of the states [COMPLETED, CANCELED]"

type Service struct {
	co      *c87operate.ClientWithResponses
	cc      *c87camunda.ClientWithResponses
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
	cc, err := c87camunda.NewClientWithResponses(
		cfg.APIs.Camunda8.BaseURL,
		c87camunda.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	co, err := c87operate.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		c87operate.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	s := &Service{
		co:   co,
		cc:   cc,
		auth: auth,
		cfg:  cfg,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) FilterProcessInstanceWithOrphanParent(ctx context.Context, items *[]c87operate.ProcessInstance) (*[]c87operate.ProcessInstance, error) {
	if items == nil {
		return nil, nil
	}
	var result []c87operate.ProcessInstance
	for _, it := range *items {
		if it.ParentKey == nil {
			continue
		}
		_, err := s.GetProcessInstanceByKey(ctx, *it.ParentKey)
		if err != nil && strings.Contains(err.Error(), "status 404") {
			result = append(result, it)
		} else if err != nil {
			return nil, err
		}
	}
	return &result, nil
}

func (s *Service) GetProcessInstanceByKey(ctx context.Context, key int64) (*c87operate.ProcessInstance, error) {
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.OperateApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving operate token: %w", err)
	}
	resp, err := s.co.GetProcessInstanceByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[c87operate.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (s *Service) GetDirectChildrenOfProcessInstance(ctx context.Context, key int64) (*[]c87operate.ProcessInstance, error) {
	filter := SearchFilterOpts{
		ParentKey: &key,
	}
	resp, err := s.SearchForProcessInstances(ctx, filter, 1000)
	if err != nil {
		return nil, fmt.Errorf("searching for children of process instance with key %d: %w", key, err)
	}
	if resp == nil || resp.Items == nil {
		return nil, nil
	}
	return resp.Items, nil
}

func (s *Service) SearchForProcessInstances(ctx context.Context, filter SearchFilterOpts, size int32) (*c87operate.ResultsProcessInstance, error) {
	f := c87operate.ProcessInstance{
		TenantId:          &s.cfg.App.Tenant,
		BpmnProcessId:     filter.BpmnProcessId,
		ProcessVersion:    filter.ProcessVersion,
		ProcessVersionTag: filter.ProcessVersionTag,
		State:             filter.State.Ptr(),
		ParentKey:         filter.ParentKey,
	}
	body := c87operate.SearchProcessInstancesJSONRequestBody{
		Filter: &f,
		Size:   &size,
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.OperateApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving operate token: %w", err)
	}
	resp, err := s.co.SearchProcessInstancesWithResponse(ctx, body,
		editors.BearerTokenEditorFn[c87operate.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (s *Service) CancelProcessInstance(ctx context.Context, key int64) (*c87camunda.CancelProcessInstanceResponse, error) {
	if !s.isQuiet {
		fmt.Printf("trying to cancel process instance with key %d...\n", key)
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.Camunda8ApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving camunda8 token: %w", err)
	}
	resp, err := s.cc.CancelProcessInstanceWithResponse(ctx, strconv.Itoa(int(key)),
		c87camunda.CancelProcessInstanceJSONRequestBody{},
		editors.BearerTokenEditorFn[c87camunda.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	if !s.isQuiet {
		fmt.Printf("process instance with key %d was successfully cancelled\n", key)
	}
	return resp, nil
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key int64) (*c87operate.ChangeStatus, error) {
	if !s.isQuiet {
		fmt.Printf("trying to delete process instance with key %d...\n", key)
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.Camunda8ApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving camunda8 token: %w", err)
	}
	resp, err := s.co.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[c87operate.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	if !s.isQuiet {
		fmt.Printf("process instance with key %d was successfully deleted\n", key)
	}
	return resp.JSON200, nil
}

func (s *Service) DeleteProcessInstanceWithCancel(ctx context.Context, key int64) (*c87operate.ChangeStatus, error) {
	if !s.isQuiet {
		fmt.Printf("trying to delete process instance with key %d...\n", key)
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.Camunda8ApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving camunda8 token: %w", err)
	}
	resp, err := s.co.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[c87operate.RequestEditorFn](token))
	if resp.StatusCode() == http.StatusBadRequest &&
		resp.ApplicationproblemJSON400 != nil &&
		*resp.ApplicationproblemJSON400.Message == wrongStateMessage400 {

		if !s.isQuiet {
			fmt.Printf("process instance with key %d not in state COMPLETED or CANCELED, cancelling it first...\n", key)
		}
		_, err = s.CancelProcessInstance(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("error cancelling process instance with key %d: %w", key, err)
		}

		if !s.isQuiet {
			fmt.Printf("waiting for process instance with key %d to be cancelled by workflow engine...\n", key)
		}
		if err = s.WaitForProcessInstanceState(ctx, strconv.Itoa(int(key)), &PIState{e: stateCanceled}); err != nil {
			return nil, fmt.Errorf("waiting for canceled state failed for %d: %w", key, err)
		}

		resp, err = s.co.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key,
			editors.BearerTokenEditorFn[c87operate.RequestEditorFn](token))
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	if !s.isQuiet {
		fmt.Printf("process instance with key %d was successfully deleted\n", key)
	}
	return resp.JSON200, nil
}
