package editors

import (
	"context"
	"net/http"
)

func BearerTokenEditorFn[T ~func(context.Context, *http.Request) error](token string) T {
	return func(ctx context.Context, req *http.Request) error {
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		return nil
	}
}
