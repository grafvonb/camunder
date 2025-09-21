package v87

import (
	"context"

	camundav88 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda/v88"
	"github.com/grafvonb/camunder/pkg/camunda/cluster"
)

type GenClusterClient interface {
	GetClusterTopologyWithResponse(ctx context.Context, reqEditors ...camundav88.RequestEditorFn) (*camundav88.GetClusterTopologyResponse, error)
}

type ClusterClient interface {
	GetClusterTopology(ctx context.Context) (cluster.Topology, error)
}

// compile-time proof of interface implementation
var _ ClusterClient = (*Service)(nil)
