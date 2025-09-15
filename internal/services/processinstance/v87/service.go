package v87

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/grafvonb/camunder/internal/api/convert"
	camundav87 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda/v87"
	operatev87 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v87"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/procesinstance"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/editors"
	"github.com/grafvonb/camunder/internal/services/auth"
)

const wrongStateMessage400 = "Process instances needs to be in one of the states [COMPLETED, CANCELED]"

type Service struct {
	cc      *camundav87.ClientWithResponses
	oc      *operatev87.ClientWithResponses
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
		APIVersion: camunda.V87,
	}
}

func (s *Service) FilterProcessInstanceWithOrphanParent(ctx context.Context, items *[]operatev87.ProcessInstance) (*[]operatev87.ProcessInstance, error) {
	if items == nil {
		return nil, nil
	}
	var result []operatev87.ProcessInstance
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

func (s *Service) GetProcessInstanceByKey(ctx context.Context, key int64) (*operatev87.ProcessInstance, error) {
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.OperateApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving operate token: %w", err)
	}
	resp, err := s.oc.GetProcessInstanceByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[operatev87.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (s *Service) GetDirectChildrenOfProcessInstance(ctx context.Context, key int64) (*[]operatev87.ProcessInstance, error) {
	filter := processinstance.SearchFilterOpts{
		ParentKey: key,
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

func (s *Service) SearchForProcessInstances(ctx context.Context, filter processinstance.SearchFilterOpts, size int32) (*operatev87.ResultsProcessInstance, error) {
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
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.OperateApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving operate token: %w", err)
	}
	resp, err := s.oc.SearchProcessInstancesWithResponse(ctx, body,
		editors.BearerTokenEditorFn[operatev87.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (s *Service) CancelProcessInstance(ctx context.Context, key int64) (*processinstance.CancelResponse, error) {
	if !s.isQuiet {
		fmt.Printf("trying to cancel process instance with key %d...\n", key)
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.CamundaApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving camunda token: %w", err)
	}
	resp, err := s.cc.CancelProcessInstanceWithResponse(ctx, strconv.Itoa(int(key)),
		camundav87.CancelProcessInstanceJSONRequestBody{},
		editors.BearerTokenEditorFn[camundav87.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	if !s.isQuiet {
		fmt.Printf("process instance with key %d was successfully cancelled\n", key)
	}
	ret := resp.ToStable()
	return &ret, nil
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key int64) (*operatev87.ChangeStatus, error) {
	if !s.isQuiet {
		fmt.Printf("trying to delete process instance with key %d...\n", key)
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.CamundaApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving camunda token: %w", err)
	}
	resp, err := s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[operatev87.RequestEditorFn](token))
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

func (s *Service) DeleteProcessInstanceWithCancel(ctx context.Context, key int64) (*operatev87.ChangeStatus, error) {
	if !s.isQuiet {
		fmt.Printf("trying to delete process instance with key %d...\n", key)
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.CamundaApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving camunda token: %w", err)
	}
	resp, err := s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[operatev87.RequestEditorFn](token))
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
		if err = s.WaitForProcessInstanceState(ctx, strconv.Itoa(int(key)), processinstance.StateCanceled.String()); err != nil {
			return nil, fmt.Errorf("waiting for canceled state failed for %d: %w", key, err)
		}

		resp, err = s.oc.DeleteProcessInstanceAndAllDependantDataByKeyWithResponse(ctx, key,
			editors.BearerTokenEditorFn[operatev87.RequestEditorFn](token))
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
