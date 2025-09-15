package v87

import (
	"github.com/grafvonb/camunder/internal/api/convert"
	"github.com/grafvonb/camunder/pkg/camunda/processdefinition"
)

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

func (src *ResultsProcessDefinition) ToStable() *processdefinition.ResultsProcessDefinition {
	return &processdefinition.ResultsProcessDefinition{
		Total: int32(convert.Deref(src.Total, 0)),
		Items: convert.MapSlice(*src.Items, func(i ProcessDefinition) processdefinition.ProcessDefinition {
			return i.ToStable()
		}),
	}
}
