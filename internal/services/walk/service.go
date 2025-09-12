package walk

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	c87camunda8v2 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda8/v2"
	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/editors"
	"github.com/grafvonb/camunder/internal/services/auth"
)

var (
	ErrCycleDetected = errors.New("cycle detected in process instance ancestry")
)

type Service struct {
	co      *c87operatev1.ClientWithResponses
	cc      *c87camunda8v2.ClientWithResponses
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
	cc, err := c87camunda8v2.NewClientWithResponses(
		cfg.APIs.Camunda8.BaseURL,
		c87camunda8v2.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, err
	}
	co, err := c87operatev1.NewClientWithResponses(
		cfg.APIs.Operate.BaseURL,
		c87operatev1.WithHTTPClient(httpClient),
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

func (s *Service) Ancestry(ctx context.Context, startKey int64) (rootKey int64, path []int64, chain map[int64]*c87operatev1.ProcessInstanceItem, err error) {
	// visited keeps track of visited nodes to detect cycles
	// well-know pattern to have fast lookups, no duplicates, clear semantic and low memory usage with visited[cur] = struct{}{} below
	visited := make(map[int64]struct{})
	chain = make(map[int64]*c87operatev1.ProcessInstanceItem)

	cur := startKey
	for {
		// check for context cancellation
		select {
		case <-ctx.Done():
			return 0, nil, chain, ctx.Err()
		default:
		}

		if _, seen := visited[cur]; seen {
			return 0, nil, chain, fmt.Errorf("%w for this key %d", ErrCycleDetected, cur)
		}
		visited[cur] = struct{}{}

		it, getErr := s.GetProcessInstanceByKey(ctx, cur)
		if getErr != nil {
			return 0, nil, chain, fmt.Errorf("get %d: %w", cur, getErr)
		}
		chain[cur] = it
		path = append(path, cur)

		// no parent => cur is root
		if it.ParentKey == nil || *it.ParentKey == 0 {
			rootKey = cur
			return
		}

		cur = *it.ParentKey
	}
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
