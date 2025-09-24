package oauth2

import (
	"context"
	"io"

	client "github.com/grafvonb/camunder/internal/api/gen/clients/auth/oauth2"
	"github.com/grafvonb/camunder/internal/services/auth/core"
)

const formContentType = "application/x-www-form-urlencoded"

type GenAuthClient interface {
	RequestTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...client.RequestEditorFn) (*client.RequestTokenResponse, error)
}

var _ core.Authenticator = (*Service)(nil)
