package camunda

import "context"

type API interface {
	GetClusterTopology(ctx context.Context) (*Topology, error)
	Capabilities() Capabilities
}

type Capabilities struct {
}

type Topology struct {
	// Brokers A list of brokers that are part of this cluster.
	Brokers *[]Broker `json:"brokers,omitempty"`
	// ClusterSize The number of brokers in the cluster.
	ClusterSize *int32 `json:"clusterSize,omitempty"`
	// GatewayVersion The version of the Zeebe Gateway.
	GatewayVersion *string `json:"gatewayVersion,omitempty"`
	// PartitionsCount The number of partitions are spread across the cluster.
	PartitionsCount *int32 `json:"partitionsCount,omitempty"`
	// ReplicationFactor The configured replication factor for this cluster.
	ReplicationFactor *int32 `json:"replicationFactor,omitempty"`
}

type Broker struct {
	// Host The hostname for reaching the broker.
	Host *string `json:"host,omitempty"`
	// NodeId The unique (within a cluster) node ID for the broker.
	NodeId *int32 `json:"nodeId,omitempty"`
	// Partitions A list of partitions managed or replicated on this broker.
	Partitions *[]Partition `json:"partitions,omitempty"`
	// Port The port for reaching the broker.
	Port *int32 `json:"port,omitempty"`
	// Version The broker version.
	Version *string `json:"version,omitempty"`
}

type Partition struct {
	// Health Describes the current health of the partition.
	Health *PartitionHealth `json:"health,omitempty"`
	// PartitionId The unique ID of this partition.
	PartitionId *int32 `json:"partitionId,omitempty"`
	// Role Describes the Raft role of the broker for a given partition.
	Role *PartitionRole `json:"role,omitempty"`
}

type PartitionHealth string

type PartitionRole string
