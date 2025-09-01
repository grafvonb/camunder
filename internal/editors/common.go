package editors

import (
	"context"
	"net/http"

	c87camunda8v2 "github.com/grafvonb/camunder/internal/api/gen/clients/camunda/camunda8/v2"
)

func HeaderEditor(key, val string) c87camunda8v2.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set(key, val)
		return nil
	}
}
