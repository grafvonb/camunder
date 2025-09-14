package walkv87

import (
	"context"
	"fmt"
	"net/http"

	camundav87 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda/v87"
	operatev87 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v87"
	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth"
	"github.com/grafvonb/camunder/internal/services/processinstance/v87"
	"github.com/grafvonb/camunder/pkg/camunda"
)

type Service struct {
	oc      *operatev87.ClientWithResponses
	cc      *camundav87.ClientWithResponses
	auth    *auth.Service
	cfg     *config.Config
	piSvc   *v87.Service
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
	piSvc, err := v87.New(cfg, httpClient, auth)
	if err != nil {
		return nil, fmt.Errorf("init process instance service: %w", err)
	}
	s := &Service{
		oc:    co,
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

func (s *Service) Ancestry(ctx context.Context, startKey int64) (rootKey int64, path []int64, chain map[int64]*operatev87.ProcessInstance, err error) {
	// visited keeps track of visited nodes to detect cycles
	// well-know pattern to have fast lookups, no duplicates, clear semantic and low memory usage with visited[cur] = struct{}{} below
	visited := make(map[int64]struct{})
	chain = make(map[int64]*operatev87.ProcessInstance)

	cur := startKey
	for {
		// check for context cancellation
		select {
		case <-ctx.Done():
			return 0, nil, chain, ctx.Err()
		default:
		}

		if _, seen := visited[cur]; seen {
			return 0, nil, chain, fmt.Errorf("%w for this key %d", camunda.ErrCycleDetected, cur)
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

func (s *Service) Descendants(ctx context.Context, rootKey int64) (desc []int64, edges map[int64][]int64, chain map[int64]*operatev87.ProcessInstance, err error) {
	visited := make(map[int64]struct{})
	edges = make(map[int64][]int64)
	chain = make(map[int64]*operatev87.ProcessInstance)

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

func (s *Service) Family(ctx context.Context, startKey int64) (fam []int64, edges map[int64][]int64, chain map[int64]*operatev87.ProcessInstance, err error) {
	rootKey, _, _, err := s.Ancestry(ctx, startKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("ancestry fetch: %w", err)
	}
	return s.Descendants(ctx, rootKey)
}
