package processinstance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	c87camunda8v2 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda8/v2"
	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/grafvonb/camunder/internal/editors"
)

type Service struct {
	co *c87operatev1.ClientWithResponses
	cc *c87camunda8v2.ClientWithResponses
}

func New(operateBaseUrl string, camunda8BaseUrl string, httpClient *http.Client, token string) (*Service, error) {
	co, err := c87operatev1.NewClientWithResponses(
		operateBaseUrl,
		c87operatev1.WithHTTPClient(httpClient),
		c87operatev1.WithRequestEditorFn(editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](token)),
	)
	cc, err := c87camunda8v2.NewClientWithResponses(
		camunda8BaseUrl,
		c87camunda8v2.WithHTTPClient(httpClient),
		c87camunda8v2.WithRequestEditorFn(editors.BearerTokenEditorFn[c87camunda8v2.RequestEditorFn](token)),
	)
	if err != nil {
		return nil, err
	}
	return &Service{
		co: co,
		cc: cc,
	}, nil
}

func (s *Service) CancelProcessInstance(ctx context.Context, key string) (*c87camunda8v2.CancelProcessInstanceResponse, error) {
	resp, err := s.cc.CancelProcessInstanceWithResponse(ctx, key, c87camunda8v2.CancelProcessInstanceJSONRequestBody{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp, nil
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key string) (*c87operatev1.ProcessInstanceDeleteResponse, error) {
	resp, err := s.co.DeleteProcessInstanceAndDependantDataByKeyWithResponse(ctx, key)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}

func (s *Service) DeleteProcessInstanceWithCancel(ctx context.Context, key string) (*c87operatev1.ProcessInstanceDeleteResponse, error) {
	resp, err := s.co.DeleteProcessInstanceAndDependantDataByKeyWithResponse(ctx, key)
	if resp.StatusCode() == http.StatusBadRequest &&
		resp.ApplicationproblemJSON400 != nil &&
		*resp.ApplicationproblemJSON400.Message == "Process instances needs to be in one of the states [COMPLETED, CANCELED]" {

		_, err = s.CancelProcessInstance(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("error cancelling process instance with key %s: %w", key, err)
		}

		// TODO implement retry with backoff
		time.Sleep(15 * time.Second)
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
