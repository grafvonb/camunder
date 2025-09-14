package editors

import (
	"context"
	"net/http"

	"github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda/v87"
)

func HeaderEditor(key, val string) v87.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set(key, val)
		return nil
	}
}
