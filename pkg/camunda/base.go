package camunda

import "context"

// Base is an optional “common” port you can embed in entity services.
type Base interface {
	Capabilities(ctx context.Context) Capabilities
}
