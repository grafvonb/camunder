package cmd

import (
	"fmt"
	"strings"

	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
)

func printKeys(path []int64, chain map[int64]*c87operatev1.ProcessInstanceItem) string {
	if len(path) == 0 {
		return ""
	}
	parts := make([]string, len(path))
	for i, k := range path {
		parts[i] = fmt.Sprintf("%d (%s)", k, valueOr(chain[k].BpmnProcessId, "unknown"))
	}
	return strings.Join(parts, " → ")
}
