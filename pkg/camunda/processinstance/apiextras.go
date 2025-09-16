package processinstance

import (
	"context"
)

type Walker interface {
	Ancestry(ctx context.Context, startKey int64) (rootKey int64, path []int64, chain map[int64]ProcessInstance, err error)
	Descendants(ctx context.Context, rootKey int64) (desc []int64, edges map[int64][]int64, chain map[int64]ProcessInstance, err error)
	Family(ctx context.Context, startKey int64) (fam []int64, edges map[int64][]int64, chain map[int64]ProcessInstance, err error)
}

func AsWalker(api API) (Walker, bool) {
	w, ok := api.(Walker)
	return w, ok
}
