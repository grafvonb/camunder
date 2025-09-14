package camunda

import "context"

type Base interface {
	Capabilities(ctx context.Context) Capabilities
}

type Capabilities struct {
	APIVersion APIVersion
}
