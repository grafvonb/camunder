package processdefinition

import (
	"context"

	"github.com/grafvonb/camunder/pkg/camunda"
)

type API interface {
	camunda.Base
	GetProcessDefinitionByKey(ctx context.Context, key int64) (ProcessDefinition, error)
	SearchProcessDefinitions(ctx context.Context, filter SearchFilterOpts, size int32) (ProcessDefinitions, error)
}

type ProcessDefinition struct {
	BpmnProcessId string `json:"bpmnProcessId,omitempty"`
	Key           int64  `json:"key,omitempty"`
	Name          string `json:"name,omitempty"`
	TenantId      string `json:"tenantId,omitempty"`
	Version       int32  `json:"version,omitempty"`
	VersionTag    string `json:"versionTag,omitempty"`
}

type SearchFilterOpts struct {
	Key           int64  `json:"key,omitempty"`
	BpmnProcessId string `json:"bpmnProcessId,omitempty"`
	Version       int32  `json:"version,omitempty"`
	VersionTag    string `json:"versionTag,omitempty"`
}

type ProcessDefinitions struct {
	Total int32               `json:"total,omitempty"`
	Items []ProcessDefinition `json:"items,omitempty"`
}
