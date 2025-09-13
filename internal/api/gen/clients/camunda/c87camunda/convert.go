package c87camunda

import (
	"github.com/grafvonb/camunder/internal/api/convert"
	"github.com/grafvonb/camunder/pkg/camunda"
)

func (src TopologyResponse) ToStable() (camunda.Topology, error) {
	br, err := convert.MapNullableSlice(src.Brokers, func(b BrokerInfo) camunda.Broker { return b.ToStable() })
	if err != nil {
		return camunda.Topology{}, err
	}
	cs, err := convert.MapNullable(src.ClusterSize, func(v int32) int32 { return v })
	if err != nil {
		return camunda.Topology{}, err
	}
	gv, err := convert.MapNullable(src.GatewayVersion, func(s string) string { return s })
	if err != nil {
		return camunda.Topology{}, err
	}
	pc, err := convert.MapNullable(src.PartitionsCount, func(v int32) int32 { return v })
	if err != nil {
		return camunda.Topology{}, err
	}
	rf, err := convert.MapNullable(src.ReplicationFactor, func(v int32) int32 { return v })
	if err != nil {
		return camunda.Topology{}, err
	}

	return camunda.Topology{
		Brokers:           br, // *([]camunda.Broker) or nil
		ClusterSize:       cs, // *int32 or nil
		GatewayVersion:    gv, // *string or nil
		PartitionsCount:   pc, // *int32 or nil
		ReplicationFactor: rf, // *int32 or nil
	}, nil
}

func (src BrokerInfo) ToStable() camunda.Broker {
	return camunda.Broker{
		Host:       convert.CopyPtr(src.Host),
		NodeId:     convert.CopyPtr(src.NodeId),
		Partitions: convert.MapPtrSlice(src.Partitions, func(p Partition) camunda.Partition { return p.ToStable() }),
		Port:       convert.CopyPtr(src.Port),
		Version:    convert.CopyPtr(src.Version),
	}
}

func (src Partition) ToStable() camunda.Partition {
	return camunda.Partition{
		Health:      convert.MapPtr(src.Health, func(h PartitionHealth) camunda.PartitionHealth { return camunda.PartitionHealth(h) }),
		PartitionId: convert.CopyPtr(src.PartitionId),
		Role:        convert.MapPtr(src.Role, func(r PartitionRole) camunda.PartitionRole { return camunda.PartitionRole(r) }),
	}
}
