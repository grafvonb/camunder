package camunda

import "errors"

var (
	ErrCycleDetected     = errors.New("cycle detected in process instance ancestry")
	ErrUnknownAPIVersion = errors.New("unknown Camunda APIs version")
	ErrNotSupported      = errors.New("feature not supported by this version")
)
