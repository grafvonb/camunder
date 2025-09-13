package walk

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/grafvonb/camunder/internal/api/gen/clients/camunda/c87camunda"
	"github.com/grafvonb/camunder/internal/api/gen/clients/camunda/c87operate"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth"
	processinstance "github.com/grafvonb/camunder/internal/services/process-instance"
)

var (
	ErrCycleDetected = errors.New("cycle detected in process instance ancestry")
)

type Service struct {
	co      *c87operate.ClientWithResponses
	cc      *c87camunda.ClientWithResponses
	auth    *auth.Service
	cfg     *config.Config
	piSvc   *processinstance.Service
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
	piSvc, err := processinstance.New(cfg, httpClient, auth)
	if err != nil {
		return nil, fmt.Errorf("init process instance service: %w", err)
	}
	s := &Service{
		co:    co,
		cc:    cc,
		auth:  auth,
		cfg:   cfg,
		piSvc: piSvc,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) Ancestry(ctx context.Context, startKey int64) (rootKey int64, path []int64, chain map[int64]*c87operate.ProcessInstance, err error) {
	// visited keeps track of visited nodes to detect cycles
	// well-know pattern to have fast lookups, no duplicates, clear semantic and low memory usage with visited[cur] = struct{}{} below
	visited := make(map[int64]struct{})
	chain = make(map[int64]*c87operate.ProcessInstance)

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

		it, getErr := s.piSvc.GetProcessInstanceByKey(ctx, cur)
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

func (s *Service) Descendants(ctx context.Context, rootKey int64) (desc []int64, edges map[int64][]int64, chain map[int64]*c87operate.ProcessInstance, err error) {
	visited := make(map[int64]struct{})
	edges = make(map[int64][]int64)
	chain = make(map[int64]*c87operate.ProcessInstance)

	// depth-first search (DFS) to explore the tree
	var dfs func(int64) error
	dfs = func(parent int64) error {
		// check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if _, seen := visited[parent]; seen {
			// already expanded this subtree
			return nil
		}
		visited[parent] = struct{}{}

		desc = append(desc, parent)
		if _, ok := chain[parent]; !ok {
			it, getErr := s.piSvc.GetProcessInstanceByKey(ctx, parent)
			if getErr != nil {
				return fmt.Errorf("get %d: %w", parent, getErr)
			}
			chain[parent] = it
		}

		children, e := s.piSvc.GetDirectChildrenOfProcessInstance(ctx, parent)
		if e != nil {
			return fmt.Errorf("list children of %d: %w", parent, e)
		}

		// keep an entry even if no children (useful for tree rendering)
		if _, ok := edges[parent]; !ok {
			edges[parent] = nil
		}

		items := *children
		for i := range items {
			it := &items[i]
			k := *it.Key

			edges[parent] = append(edges[parent], k)
			chain[k] = it

			if dfsErr := dfs(k); dfsErr != nil {
				return dfsErr
			}
		}
		return nil
	}

	if err = dfs(rootKey); err != nil {
		return nil, nil, nil, err
	}
	return desc, edges, chain, nil
}

func (s *Service) Family(ctx context.Context, startKey int64) (fam []int64, edges map[int64][]int64, chain map[int64]*c87operate.ProcessInstance, err error) {
	rootKey, _, _, err := s.Ancestry(ctx, startKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("ancestry fetch: %w", err)
	}
	return s.Descendants(ctx, rootKey)
}
