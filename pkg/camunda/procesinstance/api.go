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

type ProcessInstance struct {
	BpmnProcessId             string               `json:"bpmnProcessId,omitempty"`
	EndDate                   string               `json:"endDate,omitempty"`
	Incident                  bool                 `json:"incident,omitempty"`
	Key                       int64                `json:"key,omitempty"`
	ParentFlowNodeInstanceKey int64                `json:"parentFlowNodeInstanceKey,omitempty"`
	ParentKey                 int64                `json:"parentKey,omitempty"`
	ParentProcessInstanceKey  *ProcessInstance     `json:"parentProcessInstanceKey,omitempty"`
	ProcessDefinitionKey      int64                `json:"processDefinitionKey,omitempty"`
	ProcessVersion            int32                `json:"processVersion,omitempty"`
	ProcessVersionTag         string               `json:"processVersionTag,omitempty"`
	StartDate                 string               `json:"startDate,omitempty"`
	State                     ProcessInstanceState `json:"state,omitempty"`
	TenantId                  string               `json:"tenantId,omitempty"`
}

// ProcessInstanceState defines model for ProcessInstance.State.
type ProcessInstanceState string

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
