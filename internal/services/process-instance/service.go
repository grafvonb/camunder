package processinstance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	c87camunda8v2 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda8/v2"
	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/editors"
	"github.com/grafvonb/camunder/internal/services/auth"
)

type Service struct {
	co      *c87operatev1.ClientWithResponses
	cc      *c87camunda8v2.ClientWithResponses
	auth    *auth.Service
	cfg     *config.Config
	isQuiet bool
}

func New(cfg *config.Config, httpClient *http.Client, auth *auth.Service, isQuiet bool) (*Service, error) {
	cc, err := c87camunda8v2.NewClientWithResponses(
		cfg.APIs.Camunda8.BaseURL,
		c87camunda8v2.WithHTTPClient(httpClient),
	)
	co, err := c87operatev1.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		c87operatev1.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	return &Service{
		co:      co,
		cc:      cc,
		auth:    auth,
		isQuiet: isQuiet,
		cfg:     cfg,
	}, nil
}

func (s *Service) GetProcessInstanceByKey(ctx context.Context, key int64) (*c87operatev1.ProcessInstanceItem, error) {
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.OperateApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving operate token: %w", err)
	}
	resp, err := s.co.GetProcessInstanceByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (s *Service) SearchForProcessInstances(ctx context.Context, bpmnProcessId string, state PIState) (*c87operatev1.ProcessInstanceSearchResponse, error) {
	size := int32(1000)
	body := c87operatev1.ProcessInstanceSearchRequest{
		Filter: &c87operatev1.ProcessInstanceFilter{
			BpmnProcessId: &bpmnProcessId,
			State:         state.Ptr(),
		},
		Size: &size,
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.OperateApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving operate token: %w", err)
	}
	resp, err := s.co.SearchForProcessInstancesWithResponse(ctx, body,
		editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (s *Service) CancelProcessInstance(ctx context.Context, key string) (*c87camunda8v2.CancelProcessInstanceResponse, error) {
	if !s.isQuiet {
		fmt.Printf("trying to cancel process instance with key %s...\n", key)
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.Camunda8ApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving camunda8 token: %w", err)
	}
	resp, err := s.cc.CancelProcessInstanceWithResponse(ctx, key,
		c87camunda8v2.CancelProcessInstanceJSONRequestBody{},
		editors.BearerTokenEditorFn[c87camunda8v2.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	if !s.isQuiet {
		fmt.Printf("process instance with key %s was successfully cancelled\n", key)
	}
	return resp, nil
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key string) (*c87operatev1.ProcessInstanceDeleteResponse, error) {
	if !s.isQuiet {
		fmt.Printf("trying to delete process instance with key %s...\n", key)
	}
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.Camunda8ApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving camunda8 token: %w", err)
	}
	resp, err := s.co.DeleteProcessInstanceAndDependantDataByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	if !s.isQuiet {
		fmt.Printf("process instance with key %s was successfully deleted\n", key)
	}
	return resp.JSON200, nil
}

func (s *Service) DeleteProcessInstanceWithCancel(ctx context.Context, key string) (*c87operatev1.ProcessInstanceDeleteResponse, error) {
	token, err := s.auth.RetrieveTokenForAPI(ctx, config.Camunda8ApiKeyConst)
	if err != nil {
		return nil, fmt.Errorf("error retrieving camunda8 token: %w", err)
	}
	resp, err := s.co.DeleteProcessInstanceAndDependantDataByKeyWithResponse(ctx, key,
		editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](token))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusBadRequest &&
		resp.ApplicationproblemJSON400 != nil &&
		*resp.ApplicationproblemJSON400.Message == "process instances needs to be in one of the states [COMPLETED, CANCELED]" {

		if !s.isQuiet {
			fmt.Printf("process instance with key %s not in state COMPLETED or CANCELED, cancelling it first...\n", key)
		}
		_, err = s.CancelProcessInstance(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("error cancelling process instance with key %s: %w", key, err)
		}

		// TODO implement retry with backoff
		if !s.isQuiet {
			fmt.Printf("waiting for process instance with key %s to be cancelled by workflow engine...\n", key)
		}
		time.Sleep(10 * time.Second)
		resp, err = s.co.DeleteProcessInstanceAndDependantDataByKeyWithResponse(ctx, key,
			editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](token))
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}
