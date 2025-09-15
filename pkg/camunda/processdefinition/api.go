package processdefinition

import (
	"context"

	"github.com/grafvonb/camunder/pkg/camunda"
)

type API interface {
	camunda.Base
	GetProcessDefinitionByKey(ctx context.Context, key int64) (*ProcessDefinition, error)
	SearchProcessDefinitions(ctx context.Context, filter SearchFilterOpts, size int32) (*ResultsProcessDefinition, error)
}

type ProcessDefinition struct {
	BpmnProcessId string
	Key           int64
	Name          string
	TenantId      string
	Version       int32
	VersionTag    string
}

type SearchFilterOpts struct {
	Key           int64
	BpmnProcessId string
	Version       int32
	VersionTag    string
}

type ResultsProcessDefinition struct {
	Total int32
	Items []ProcessDefinition
}
