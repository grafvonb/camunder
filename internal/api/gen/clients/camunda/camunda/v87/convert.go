package v87

import (
	"github.com/grafvonb/camunder/internal/api/convert"
	"github.com/grafvonb/camunder/pkg/camunda/cluster"
	processinstance "github.com/grafvonb/camunder/pkg/camunda/procesinstance"
)

func (src CancelProcessInstanceResponse) ToStable() (processinstance.CancelResponse, error) {
	return processinstance.CancelResponse{
		StatusCode: src.StatusCode(),
		Status:     src.Status(),
	}, nil
}

func (src TopologyResponse) ToStable() (cluster.Topology, error) {
	br, err := convert.MapNullableSlice(src.Brokers, func(b BrokerInfo) cluster.Broker { return b.ToStable() })
	if err != nil {
		return cluster.Topology{}, err
	}
	cs, err := convert.MapNullable(src.ClusterSize, func(v int32) int32 { return v })
	if err != nil {
		return cluster.Topology{}, err
	}
	gv, err := convert.MapNullable(src.GatewayVersion, func(s string) string { return s })
	if err != nil {
		return cluster.Topology{}, err
	}
	pc, err := convert.MapNullable(src.PartitionsCount, func(v int32) int32 { return v })
	if err != nil {
		return cluster.Topology{}, err
	}
	rf, err := convert.MapNullable(src.ReplicationFactor, func(v int32) int32 { return v })
	if err != nil {
		return cluster.Topology{}, err
	}

	return cluster.Topology{
		Brokers:               convert.DerefSlice(br), // nil -> nil slice (zero), or copy
		ClusterSize:           convert.Deref(cs, 0),   // nil -> 0
		GatewayVersion:        convert.Deref(gv, ""),  // nil -> ""
		PartitionsCount:       convert.Deref(pc, 0),
		ReplicationFactor:     convert.Deref(rf, 0),
		LastCompletedChangeId: "", // not in v8.7 response
	}, nil
}

func (b BrokerInfo) ToStable() cluster.Broker {
	return cluster.Broker{
		Host:       convert.Deref(b.Host, ""),  // *string -> string
		NodeId:     convert.Deref(b.NodeId, 0), // *int32 -> int32
		Partitions: convert.DerefSlicePtr(b.Partitions, func(p Partition) cluster.Partition { return p.ToStable() }),
		Port:       convert.Deref(b.Port, 0),
		Version:    convert.Deref(b.Version, ""),
	}
}

func (p Partition) ToStable() cluster.Partition {
	return cluster.Partition{
		Health:      convert.DerefMap(p.Health, func(h PartitionHealth) cluster.PartitionHealth { return cluster.PartitionHealth(h) }, cluster.PartitionHealth("")),
		PartitionId: convert.Deref(p.PartitionId, 0),
		Role:        convert.DerefMap(p.Role, func(r PartitionRole) cluster.PartitionRole { return cluster.PartitionRole(r) }, cluster.PartitionRole("")),
	}
}
