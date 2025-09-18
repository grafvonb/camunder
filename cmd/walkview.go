package cmd

import (
	"fmt"
	"strings"

	"github.com/grafvonb/camunder/pkg/camunda/processinstance"
)

type Chain map[int64]processinstance.ProcessInstance
type KeysPath []int64

type Label func(processinstance.ProcessInstance) string

func (p KeysPath) KeysOnly(c Chain) string {
	return p.join(c, func(item processinstance.ProcessInstance) string {
		return fmt.Sprint(item.Key)
	}, "\n")
}

func (p KeysPath) StandardLine(c Chain) string {
	return p.join(c, func(item processinstance.ProcessInstance) string {
		var pTag, eTag, vTag string
		if item.ProcessVersion > 0 {
			pTag = fmt.Sprintf(" p:%d", item.ParentKey)
		} else {
			pTag = " p:<root>"
		}
		if item.EndDate != "" {
			eTag = fmt.Sprintf(" e:%s", item.EndDate)
		}
		if item.ProcessVersionTag != "" {
			vTag = "/" + item.ProcessVersionTag
		}

		return fmt.Sprintf(
			"%-16d %s %s v%s%s %s s:%s%s%s i:%t",
			item.Key, item.TenantId, item.BpmnProcessId, version, vTag, item.State, item.StartDate, eTag, pTag, item.Incident,
		)
	}, "\n")
}

func (p KeysPath) PrettyLine(c Chain) string {
	return p.join(c, func(it processinstance.ProcessInstance) string {
		return fmt.Sprintf("%d (%s)", it.Key, it.BpmnProcessId)
	}, " → ")
}

func (p KeysPath) join(c Chain, label Label, sep string) string {
	if len(p) == 0 {
		return ""
	}
	if label == nil {
		label = func(it processinstance.ProcessInstance) string {
			return fmt.Sprintf("%d (%s)", it.Key, it.BpmnProcessId)
		}
	}
	out := make([]string, 0, len(p))
	for _, k := range p {
		out = append(out, label(c[k]))
	}
	return strings.Join(out, sep)
}
