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
)

type Service struct {
	co      *c87operatev1.ClientWithResponses
	cc      *c87camunda8v2.ClientWithResponses
	isQuiet bool
}

func New(cfg config.Config, httpClient *http.Client, isQuiet bool) (*Service, error) {
	cc, err := c87camunda8v2.NewClientWithResponses(
		cfg.Camunda8API.BaseURL,
		c87camunda8v2.WithHTTPClient(httpClient),
		c87camunda8v2.WithRequestEditorFn(editors.BearerTokenEditorFn[c87camunda8v2.RequestEditorFn](cfg.Camunda8API.Token)),
	)
	co, err := c87operatev1.NewClientWithResponses(
		cfg.OperateAPI.BaseURL,
		c87operatev1.WithHTTPClient(httpClient),
		c87operatev1.WithRequestEditorFn(editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](cfg.OperateAPI.Token)),
	)
	if err != nil {
		return nil, err
	}
	return &Service{
		co:      co,
		cc:      cc,
		isQuiet: isQuiet,
	}, nil
}

func (s *Service) SearchForProcessInstances(ctx context.Context, bpmnProcessId string) (*c87operatev1.ProcessInstanceSearchResponse, error) {
	size := int32(1000)
	body := c87operatev1.ProcessInstanceSearchRequest{
		Filter: &c87operatev1.ProcessInstanceFilter{
			BpmnProcessId: &bpmnProcessId,
		},
		Size: &size,
	}
	resp, err := s.co.SearchForProcessInstancesWithResponse(ctx, body)
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
		fmt.Printf("Trying to cancel process instance with key %s...\n", key)
	}
	resp, err := s.cc.CancelProcessInstanceWithResponse(ctx, key, c87camunda8v2.CancelProcessInstanceJSONRequestBody{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	if !s.isQuiet {
		fmt.Printf("Process instance with key %s was successfully cancelled\n", key)
	}
	return resp, nil
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key string) (*c87operatev1.ProcessInstanceDeleteResponse, error) {
	if !s.isQuiet {
		fmt.Printf("Trying to delete process instance with key %s...\n", key)
	}
	resp, err := s.co.DeleteProcessInstanceAndDependantDataByKeyWithResponse(ctx, key)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	if !s.isQuiet {
		fmt.Printf("Process instance with key %s was successfully deleted\n", key)
	}
	return resp.JSON200, nil
}

func (s *Service) DeleteProcessInstanceWithCancel(ctx context.Context, key string) (*c87operatev1.ProcessInstanceDeleteResponse, error) {
	resp, err := s.co.DeleteProcessInstanceAndDependantDataByKeyWithResponse(ctx, key)
	if resp.StatusCode() == http.StatusBadRequest &&
		resp.ApplicationproblemJSON400 != nil &&
		*resp.ApplicationproblemJSON400.Message == "Process instances needs to be in one of the states [COMPLETED, CANCELED]" {

		if !s.isQuiet {
			fmt.Printf("Process instance with key %s not in state COMPLETED or CANCELED, cancelling it first...\n", key)
		}
		_, err = s.CancelProcessInstance(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("error cancelling process instance with key %s: %w", key, err)
		}

		// TODO implement retry with backoff
		if !s.isQuiet {
			fmt.Printf("Waiting for process instance with key %s to be cancelled by workflow engine...\n", key)
		}
		time.Sleep(10 * time.Second)
		resp, err = s.co.DeleteProcessInstanceAndDependantDataByKeyWithResponse(ctx, key)
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}
