package cluster

import (
	"context"

	"github.com/grafvonb/camunder/pkg/camunda"
)

type API interface {
	camunda.Base
	GetClusterTopology(ctx context.Context) (Topology, error)
}

type Topology struct {
	Brokers               []Broker
	ClusterSize           int32
	GatewayVersion        string
	PartitionsCount       int32
	ReplicationFactor     int32
	LastCompletedChangeId string
}

type Broker struct {
	Host       string
	NodeId     int32
	Partitions []Partition
	Port       int32
	Version    string
}

type Partition struct {
	Health      PartitionHealth
	PartitionId int32
	Role        PartitionRole
}

type PartitionHealth string
type PartitionRole string
