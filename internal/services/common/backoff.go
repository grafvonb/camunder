package common

import (
	"errors"
	"fmt"
	"time"
)

type BackoffStrategy string

const (
	BackoffFixed       BackoffStrategy = "fixed"
	BackoffExponential BackoffStrategy = "exponential"
)

type BackoffConfig struct {
	Strategy     BackoffStrategy `mapstructure:"strategy"`
	InitialDelay time.Duration   `mapstructure:"initial_delay"`
	MaxDelay     time.Duration   `mapstructure:"max_delay"`
	MaxRetries   int             `mapstructure:"max_retries"`
	Multiplier   float64         `mapstructure:"multiplier"` // Used for exponential backoff strategy
	Timeout      time.Duration   `mapstructure:"timeout"`
}

func (c BackoffConfig) NextDelay(prev time.Duration) time.Duration {
	if prev <= 0 {
		prev = c.InitialDelay
	}
	var delay time.Duration
	switch c.Strategy {
	case BackoffExponential:
		delay = time.Duration(float64(prev) * c.Multiplier)
	default: // fixed strategy
		delay = prev
	}
	if c.MaxDelay > 0 && delay > c.MaxDelay {
		delay = c.MaxDelay
	}
	return delay
}

func (c BackoffConfig) Validate() error {
	var errs []error

	if c.Strategy != BackoffFixed && c.Strategy != BackoffExponential {
		errs = append(errs, fmt.Errorf("strategy must be either %q or %q", BackoffFixed, BackoffExponential))
	}
	if c.InitialDelay <= 0 {
		errs = append(errs, errors.New("initial_delay must be a positive duration"))
	}
	if c.MaxDelay < 0 {
		errs = append(errs, errors.New("max_delay must be non-negative"))
	}
	if c.MaxDelay > 0 && c.MaxDelay < c.InitialDelay {
		errs = append(errs, errors.New("max_delay must be greater than or equal to initial_delay"))
	}
	if c.MaxRetries < 0 {
		errs = append(errs, errors.New("max_retries must be non-negative"))
	}
	if c.Strategy == BackoffExponential && c.Multiplier <= 1 {
		errs = append(errs, errors.New("multiplier must be greater than 1 for exponential backoff"))
	}
	if c.Timeout <= 0 {
		errs = append(errs, errors.New("timeout must be a positive duration"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
