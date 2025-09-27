package xsrf

import (
	"context"

	client "github.com/grafvonb/camunder/internal/api/gen/clients/auth/xsrf"
	"github.com/grafvonb/camunder/internal/services/auth/core"
)

type GenAuthClient interface {
	XsrfLoginPostWithResponse(ctx context.Context, appId client.XsrfLoginPostParamsAppId,
		params *client.XsrfLoginPostParams, body client.XsrfLoginPostJSONRequestBody, reqEditors ...client.RequestEditorFn) (*client.XsrfLoginPostResponse, error)
}

var _ core.Authenticator = (*Service)(nil)
