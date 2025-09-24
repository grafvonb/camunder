package imx

import (
	"context"

	client "github.com/grafvonb/camunder/internal/api/gen/clients/auth/imx"
	"github.com/grafvonb/camunder/internal/services/auth/core"
)

type GenAuthClient interface {
	ImxLoginPostWithResponse(ctx context.Context, appId client.ImxLoginPostParamsAppId,
		params *client.ImxLoginPostParams, body client.ImxLoginPostJSONRequestBody, reqEditors ...client.RequestEditorFn) (*client.ImxLoginPostResponse, error)
}

var _ core.Authenticator = (*Service)(nil)
