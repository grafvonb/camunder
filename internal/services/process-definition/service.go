package processdefinition

import (
	"context"
	"fmt"
	"net/http"

	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/editors"
)

type Service struct {
	c *c87operatev1.ClientWithResponses
}

func New(cfg config.Config, httpClient *http.Client) (*Service, error) {
	c, err := c87operatev1.NewClientWithResponses(
		cfg.OperateAPI.BaseURL,
		c87operatev1.WithHTTPClient(httpClient),
		c87operatev1.WithRequestEditorFn(editors.BearerTokenEditorFn[c87operatev1.RequestEditorFn](cfg.OperateAPI.Token)),
	)
	if err != nil {
		return nil, err
	}
	return &Service{c: c}, nil
}

func (s *Service) SearchForProcessDefinitions(ctx context.Context) (*c87operatev1.ProcessDefinitionSearchResponse, error) {
	size := int32(1000)
	body := c87operatev1.ProcessDefinitionSearchRequest{
		Filter: &c87operatev1.ProcessDefinitionFilter{},
		Size:   &size,
	}
	resp, err := s.c.SearchForProcessDefinitionsWithResponse(ctx, body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	return resp.JSON200, nil
}
