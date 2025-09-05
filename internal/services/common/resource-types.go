package common

import (
	"fmt"
	"strings"
)

type ResourceTypes map[string]string

func (s ResourceTypes) String() string {
	var sb strings.Builder
	for k, v := range s {
		sb.WriteString(fmt.Sprintf("%s (%s), ", v, k))
	}
	result := sb.String()
	return strings.TrimSuffix(result, ", ")
}

func (s ResourceTypes) PrettyString() string {
	return fmt.Sprintf("Supported resource types are: %s", s.String())
}
