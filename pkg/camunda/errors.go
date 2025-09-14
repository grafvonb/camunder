package camunda

import "errors"

var (
	ErrCycleDetected = errors.New("cycle detected in process instance ancestry")
)
