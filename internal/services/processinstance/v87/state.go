package v87

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	operatev87 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/operate/v87"
	"github.com/grafvonb/camunder/pkg/camunda/processinstance"
)

func StateOrNil(s processinstance.State) *operatev87.ProcessInstanceState {
	if s == processinstance.StateAll {
		return nil
	}
	v := operatev87.ProcessInstanceState(strings.ToUpper(s.String()))
	return &v
}

// WaitForProcessInstanceState waits until the instance reaches the desired state.
// - Respects ctx cancellation/deadline; augments with cfg.Timeout if set
// - Returns nil on success or an error on failure/timeout.
func (s *Service) WaitForProcessInstanceState(ctx context.Context, key string, desiredState string) error {
	desiredStr, err := processinstance.ParseState(desiredState)
	if err != nil {
		return err
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

	attempts := 0
	delay := cfg.InitialDelay

	for {
		if errInDelay := ctx.Err(); errInDelay != nil {
			return errInDelay
		}
		attempts++

		pi, errInDelay := s.GetProcessInstanceByKey(ctx, keyInt)
		if errInDelay == nil {
			state := ""
			if pi.State != "" {
				state = string(pi.State)
			}
			if strings.ToLower(state) == desiredStr.String() {
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
