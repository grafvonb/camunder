package processinstance

import (
	"context"
	"fmt"
	"net/http"

	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/grafvonb/camunder/internal/editors"
)

type Service struct {
	c *c87operatev1.ClientWithResponses
}

func New(baseUrl string, httpClient *http.Client, token string) (*Service, error) {
	c, err := c87operatev1.NewClientWithResponses(
		baseUrl,
		c87operatev1.WithHTTPClient(httpClient),
		c87operatev1.WithRequestEditorFn(editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](token)),
	)
	if err != nil {
		return nil, err
	}
	return &Service{c: c}, nil
}

func (s *Service) DeleteProcessInstance(ctx context.Context, key string) (*c87operatev1.ProcessInstanceDeleteResponse, error) {
	resp, err := s.c.DeleteProcessInstanceAndDependantDataByKeyWithResponse(ctx, key)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusNoContent {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}
