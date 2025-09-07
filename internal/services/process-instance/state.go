package processinstance

import (
	"errors"
	"fmt"

	c87operatev1 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v1"
)

var ErrUnknownStateFilter = errors.New("is unknown (valid: all, active, canceled, completed)")

type stateEnum int

const (
	stateAll stateEnum = iota
	stateActive
	stateCanceled
	stateCompleted
)

// PIStateFilter represents possible states for filtering process instances
type PIStateFilter struct{ e stateEnum }

var (
	StateAll       = PIStateFilter{e: stateAll}
	StateActive    = PIStateFilter{e: stateActive}
	StateCompleted = PIStateFilter{e: stateCompleted}
	StateCanceled  = PIStateFilter{e: stateCanceled}
)

// Ptr converts the enum to the *ProcessInstanceFilterState
func (s PIStateFilter) Ptr() *c87operatev1.ProcessInstanceFilterState {
	switch s.e {
	case stateAll:
		return nil
	case stateActive:
		v := c87operatev1.ProcessInstanceFilterStateACTIVE
		return &v
	case stateCompleted:
		v := c87operatev1.ProcessInstanceFilterStateCOMPLETED
		return &v
	case stateCanceled:
		v := c87operatev1.ProcessInstanceFilterStateCANCELED
		return &v
	default:
		return nil
	}
}

func (s PIStateFilter) String() string {
	switch s.e {
	case stateAll:
		return "all"
	case stateActive:
		return "active"
	case stateCanceled:
		return "canceled"
	case stateCompleted:
		return "completed"
	default:
		return "unknown"
	}
}

func PIStateFilterFromString(s string) (PIStateFilter, error) {
	switch s {
	case "all":
		return StateAll, nil
	case "active":
		return StateActive, nil
	case "canceled":
		return StateCanceled, nil
	case "completed":
		return StateCompleted, nil
	default:
		return PIStateFilter{}, fmt.Errorf("%q %w", s, ErrUnknownStateFilter)
	}
}
