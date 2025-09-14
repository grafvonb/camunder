package cmd

import (
	"fmt"
	"strings"

	"github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v87"
)

type Chain map[int64]*v87.ProcessInstance
type KeysPath []int64

type Label func(*v87.ProcessInstance) string

func (p KeysPath) KeysOnly(c Chain) string {
	return p.join(c, func(it *v87.ProcessInstance) string {
		return fmt.Sprint(*it.Key)
	}, "\n")
}

func (p KeysPath) StandardLine(c Chain) string {
	return p.join(c, func(it *v87.ProcessInstance) string {
		key := valueOr(it.Key, int64(0))
		tenant := valueOr(it.TenantId, "")
		bpmnID := valueOr(it.BpmnProcessId, "")
		version := valueOr(it.ProcessVersion, int32(0))
		versionTag := valueOr(it.ProcessVersionTag, "")
		state := valueOr(it.State, "")
		start := valueOr(it.StartDate, "")
		end := valueOr(it.EndDate, "")
		parent := valueOr(it.ParentKey, int64(0))
		incident := valueOr(it.Incident, false)

		var pTag, eTag, vTag string
		if parent > 0 {
			pTag = fmt.Sprintf(" p:%d", parent)
		} else {
			pTag = " p:<root>"
		}
		if end != "" {
			eTag = fmt.Sprintf(" e:%s", end)
		}
		if versionTag != "" {
			vTag = "/" + versionTag
		}

		return fmt.Sprintf(
			"%-16d %s %s v%d%s %s s:%s%s%s i:%t",
			key, tenant, bpmnID, version, vTag, state, start, eTag, pTag, incident,
		)
	}, "\n")
}

func (p KeysPath) PrettyLine(c Chain) string {
	return p.join(c, func(it *v87.ProcessInstance) string {
		return fmt.Sprintf("%d (%s)", *it.Key, valueOr(it.BpmnProcessId, "undefined"))
	}, " → ")
}

func (p KeysPath) join(c Chain, label Label, sep string) string {
	if len(p) == 0 {
		return ""
	}
	if label == nil {
		label = func(it *v87.ProcessInstance) string {
			return fmt.Sprintf("%d (%s)", *it.Key, valueOr(it.BpmnProcessId, "undefined"))
		}
	}
	out := make([]string, 0, len(p))
	for _, k := range p {
		if it := c[k]; it != nil {
			out = append(out, label(it))
		} else {
			out = append(out, fmt.Sprint(k)) // fallback if chain missing
		}
	}
	return strings.Join(out, sep)
}
