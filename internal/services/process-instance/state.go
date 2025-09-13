package processinstance

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafvonb/camunder/internal/api/gen/clients/camunda/c87operate"
)

var ErrUnknownStateFilter = errors.New("is unknown (valid: all, active, canceled, completed)")

type stateEnum int

const (
	stateAll stateEnum = iota
	stateActive
	stateCanceled
	stateCompleted
)

// PIState represents possible states of a process instance, with addition of "all" for filtering.
type PIState struct{ e stateEnum }

var (
	StateAll       = PIState{e: stateAll}
	StateActive    = PIState{e: stateActive}
	StateCompleted = PIState{e: stateCompleted}
	StateCanceled  = PIState{e: stateCanceled}
)

// Ptr converts the enum to the *ProcessInstanceFilterState
func (s PIState) Ptr() *c87operate.ProcessInstanceState {
	switch s.e {
	case stateAll:
		return nil
	case stateActive:
		v := c87operate.ACTIVE
		return &v
	case stateCompleted:
		v := c87operate.COMPLETED
		return &v
	case stateCanceled:
		v := c87operate.CANCELED
		return &v
	default:
		return nil
	}
}

func (s PIState) String() string {
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

func PIStateFilterFromString(s string) (PIState, error) {
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
		return PIState{}, fmt.Errorf("%q %w", s, ErrUnknownStateFilter)
	}
}

// WaitForProcessInstanceState waits until the instance reaches the desired state.
// - Respects ctx cancellation/deadline; augments with cfg.Timeout if set
// - Returns nil on success or an error on failure/timeout.
func (s *Service) WaitForProcessInstanceState(ctx context.Context, key string, desiredState *PIState) error {
	if desiredState == nil {
		return errors.New("desired state must be provided")
	}

	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse provided instance key %q to int64: %w", key, err)
	}

	cfg := s.cfg.App.Backoff
	if cfg.Timeout > 0 {
		deadline := time.Now().Add(cfg.Timeout)
		if dl, ok := ctx.Deadline(); !ok || deadline.Before(dl) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithDeadline(ctx, deadline)
			defer cancel()
		}
	}

	desiredStr := strings.ToUpper(desiredState.String())
	attempts := 0
	delay := cfg.InitialDelay

	for {
		if errInDelay := ctx.Err(); errInDelay != nil {
			return errInDelay
		}
		attempts++

		pi, errInDelay := s.GetProcessInstanceByKey(ctx, keyInt)
		if errInDelay == nil && pi != nil {
			state := ""
			if pi.State != nil {
				state = string(*pi.State)
			}
			if strings.ToUpper(state) == desiredStr {
				if !s.isQuiet {
					fmt.Printf("process instance %q reached desired state %q\n", key, state)
				}
				return nil
			}
			if !s.isQuiet {
				fmt.Printf("process instance %q currently in state %q; waiting...\n", key, state)
			}
		} else if errInDelay != nil {
			if !s.isQuiet {
				fmt.Printf("fetching state for %q failed: %v (will retry)\n", key, errInDelay)
			}
		}

		if cfg.MaxRetries > 0 && attempts >= cfg.MaxRetries {
			return fmt.Errorf("exceeded max_retries (%d) waiting for state %q", cfg.MaxRetries, desiredStr)
		}
		select {
		case <-time.After(delay):
			delay = cfg.NextDelay(delay)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
