package v88

import (
	"github.com/grafvonb/camunder/internal/api/convert"
	"github.com/grafvonb/camunder/pkg/camunda/cluster"
)

func (src TopologyResponse) ToStable() (cluster.Topology, error) {
	return cluster.Topology{
		Brokers:               convert.MapSlice(src.Brokers, func(b BrokerInfo) cluster.Broker { return b.ToStable() }),
		ClusterSize:           src.ClusterSize,
		GatewayVersion:        src.GatewayVersion,
		PartitionsCount:       src.PartitionsCount,
		ReplicationFactor:     src.ReplicationFactor,
		LastCompletedChangeId: src.LastCompletedChangeId,
	}, nil
}

func (src BrokerInfo) ToStable() cluster.Broker {
	return cluster.Broker{
		Host:       src.Host,
		NodeId:     src.NodeId,
		Partitions: convert.MapSlice(src.Partitions, func(p Partition) cluster.Partition { return p.ToStable() }),
		Port:       src.Port,
		Version:    src.Version,
	}
}

func (src Partition) ToStable() cluster.Partition {
	return cluster.Partition{
		Health:      cluster.PartitionHealth(src.Health),
		PartitionId: src.PartitionId,
		Role:        cluster.PartitionRole(src.Role),
	}
}
