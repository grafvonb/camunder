package auth

import (
	"context"
	"io"

	client "github.com/grafvonb/camunder/internal/api/gen/clients/auth"
)

const formContentType = "application/x-www-form-urlencoded"

type GenAuthClient interface {
	RequestTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...client.RequestEditorFn) (*client.RequestTokenResponse, error)
}

type AuthClient interface {
	RetrieveTokenForAPI(ctx context.Context, target string) (string, error)
}

var _ AuthClient = (*Service)(nil)
