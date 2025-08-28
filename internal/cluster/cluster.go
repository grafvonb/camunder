package cluster

// Cluster describes the current topology of the cluster the gateway is part of.
type Cluster struct {
	// Brokers is a list of all brokers in the cluster.
	Brokers []Broker `json:"brokers,omitempty"`
	// ClusterSize is the number of brokers in the cluster.
	ClusterSize *int32 `json:"clusterSize,omitempty"`
	// PartitionsCount is the number of partitions spread across the cluster.
	PartitionsCount *int32 `json:"partitionsCount,omitempty"`
	// ReplicationFactor is the replication factor for this cluster.
	ReplicationFactor *int32 `json:"replicationFactor,omitempty"`
	// GatewayVersion is the version of the Zeebe gateway.
	GatewayVersion *string `json:"gatewayVersion,omitempty"`
}

// Broker describes a single broker in the cluster.
type Broker struct {
	// NodeID is the unique ID of the broker.
	NodeID int32 `json:"nodeId"`
	// Host is the host of the broker.
	Host string `json:"host"`
	// Port is the port of the broker.
	Port int32 `json:"port"`
	// Partitions is the list of partitions hosted by the broker.
	Partitions []Partition `json:"partitions,omitempty"`
	// Version is the version of the broker.
	Version string `json:"version"`
}

// Partition describes a single partition hosted by a broker.
type Partition struct {
	// PartitionID is the unique ID of the partition.
	PartitionID int32 `json:"partitionId"`
	// Role is the role of the partition.
	// Possible values: leader, follower, inactive
	Role string `json:"role"`
	// Health is the health status of the partition.
	// Possible values: healthy, unhealthy, dead
	Health string `json:"health"`
}
