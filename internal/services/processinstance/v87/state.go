package v87

import (
	"context"
	"fmt"
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
func (s *Service) WaitForProcessInstanceState(ctx context.Context, key int64, desiredState processinstance.State) error {
	log := s.log

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

		pi, errInDelay := s.GetProcessInstanceByKey(ctx, key)
		if errInDelay == nil {
			if pi.State.EqualsIgnoreCase(desiredState) {
				log.Info(fmt.Sprintf("process instance %d reached desired state %q", key, desiredState))
				return nil
			}
			log.Debug(fmt.Sprintf("process instance %d currently in state %q; waiting...", key, pi.State))
		} else if errInDelay != nil {
			if strings.Contains(errInDelay.Error(), "status 404") {
				log.Debug(fmt.Sprintf("process instance %d is absent (not found); waiting...", key))
			} else {
				log.Error(fmt.Sprintf("fetching state for %q failed: %v (will retry)", key, errInDelay))
			}
		}
		if cfg.MaxRetries > 0 && attempts >= cfg.MaxRetries {
			return fmt.Errorf("exceeded max_retries (%d) waiting for state %q", cfg.MaxRetries, desiredState)
		}
		select {
		case <-time.After(delay):
			delay = cfg.NextDelay(delay)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
