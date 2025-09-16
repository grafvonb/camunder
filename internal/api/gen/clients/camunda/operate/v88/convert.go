package v88

import (
	"github.com/grafvonb/camunder/internal/api/convert"
	"github.com/grafvonb/camunder/pkg/camunda/processdefinition"
	"github.com/grafvonb/camunder/pkg/camunda/processinstance"
)

func (src ProcessInstance) ToStable() processinstance.ProcessInstance {
	return processinstance.ProcessInstance{
		BpmnProcessId:             convert.Deref(src.BpmnProcessId, ""),
		Key:                       convert.Deref(src.Key, 0),
		EndDate:                   convert.Deref(src.EndDate, ""),
		Incident:                  convert.Deref(src.Incident, false),
		ParentFlowNodeInstanceKey: convert.Deref(src.ParentFlowNodeInstanceKey, 0),
		ParentKey:                 convert.Deref(src.ParentKey, 0),
		ParentProcessInstanceKey:  convert.Deref(src.ParentProcessInstanceKey.Key, 0),
		ProcessDefinitionKey:      convert.Deref(src.ProcessDefinitionKey, 0),
		ProcessVersion:            convert.Deref(src.ProcessVersion, 0),
		ProcessVersionTag:         convert.Deref(src.ProcessVersionTag, ""),
		StartDate:                 convert.Deref(src.StartDate, ""),
		State:                     convert.DerefMap(src.State, func(s ProcessInstanceState) processinstance.State { return processinstance.State(s) }, ""),
		TenantId:                  convert.Deref(src.TenantId, ""),
	}
}

func (src *ResultsProcessInstance) ToStable() *processinstance.ProcessInstances {
	return &processinstance.ProcessInstances{
		Total: int32(convert.Deref(src.Total, 0)),
		Items: convert.MapSlice(*src.Items, func(i ProcessInstance) processinstance.ProcessInstance {
			return i.ToStable()
		}),
	}
}

func (src ProcessDefinition) ToStable() processdefinition.ProcessDefinition {
	return processdefinition.ProcessDefinition{
		BpmnProcessId: convert.Deref(src.BpmnProcessId, ""),
		Key:           convert.Deref(src.Key, 0),
		Name:          convert.Deref(src.Name, ""),
		TenantId:      convert.Deref(src.TenantId, ""),
		Version:       convert.Deref(src.Version, 0),
		VersionTag:    convert.Deref(src.VersionTag, ""),
	}
}

func (src *ResultsProcessDefinition) ToStable() processdefinition.ProcessDefinitions {
	return processdefinition.ProcessDefinitions{
		Total: int32(convert.Deref(src.Total, 0)),
		Items: convert.MapSlice(*src.Items, func(i ProcessDefinition) processdefinition.ProcessDefinition {
			return i.ToStable()
		}),
	}
}
