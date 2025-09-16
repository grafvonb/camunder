package processinstance

import "errors"

var (
	ErrUnknownStateFilter = errors.New("is unknown (valid: all, active, canceled, completed)")
)
