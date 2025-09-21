package v87

import (
	"context"

	camundav87 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda/v87"
	"github.com/grafvonb/camunder/pkg/camunda/cluster"
)

type GenClusterClient interface {
	GetClusterTopologyWithResponse(ctx context.Context, reqEditors ...camundav87.RequestEditorFn) (*camundav87.GetClusterTopologyResponse, error)
}

type ClusterClient interface {
	GetClusterTopology(ctx context.Context) (cluster.Topology, error)
}

// compile-time proof of interface implementation
var _ ClusterClient = (*Service)(nil)
