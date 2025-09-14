package camunda

import "errors"

var ErrNotSupported = errors.New("feature not supported by this version")

type Capabilities struct {
}
