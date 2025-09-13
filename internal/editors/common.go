package editors

import (
	"context"
	"net/http"

	"github.com/grafvonb/camunder/internal/api/gen/clients/camunda/c87camunda"
)

func HeaderEditor(key, val string) c87camunda.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set(key, val)
		return nil
	}
}
