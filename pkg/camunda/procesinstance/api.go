package processinstance

import (
	"context"
	"fmt"
	"strings"

	"github.com/grafvonb/camunder/pkg/camunda"
)

type API interface {
	camunda.Base
	CancelProcessInstance(ctx context.Context, key int64) (*CancelResponse, error)
	WaitForProcessInstanceState(ctx context.Context, key string, desiredState string) error
}

type CancelResponse struct {
	StatusCode int
	Status     string
}

type SearchFilterOpts struct {
	Key               int64
	BpmnProcessId     string
	ProcessVersion    int32
	ProcessVersionTag string
	State             State
	ParentKey         int64
}

// State is the process-instance state filter.
type State string

const (
	StateAll       State = "all"
	StateActive    State = "active"
	StateCompleted State = "completed"
	StateCanceled  State = "canceled"
)

func (s State) String() string { return string(s) }

// ParseState parses a string (case-insensitive) into a State.
func ParseState(in string) (State, error) {
	switch strings.ToLower(in) {
	case "all":
		return StateAll, nil
	case "active":
		return StateActive, nil
	case "canceled":
		return StateCanceled, nil
	case "completed":
		return StateCompleted, nil
	default:
		return "", fmt.Errorf("%q %w", in, ErrUnknownStateFilter)
	}
}
