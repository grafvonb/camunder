package camunda

import (
	"fmt"
	"strings"
)

type APIVersion string

const (
	V87 APIVersion = "8.7"
	V88 APIVersion = "8.8"
)

const Current = V87

func Normalize(s string) (APIVersion, error) {
	v := strings.TrimSpace(strings.ToLower(s))
	switch v {
	case "8.7", "87", "v87", "v8.7":
		return V87, nil
	case "8.8", "88", "v88", "v8.8":
		return V88, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnknownAPIVersion, v)
	}
}

func Supported() []APIVersion {
	return []APIVersion{V87, V88}
}
