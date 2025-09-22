package auth

import (
	"context"
	"io"
	"net/url"
	"testing"

	gen "github.com/grafvonb/camunder/internal/api/gen/clients/auth"
	"github.com/grafvonb/camunder/internal/config"

	"github.com/grafvonb/camunder/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func testConfig() *config.Config {
	return &config.Config{
		Auth: config.Authentication{
			TokenURL:     "http://localhost:8080/auth/realms/camunda-platform/protocol/openid-connect/token",
			ClientID:     "test",
			ClientSecret: "test",
		},
	}
}

func newServiceUnderTest(t *testing.T) (*Service, *MockGenAuthClient) {
	t.Helper()
	m := NewMockGenAuthClient(t)
	s, err := New(testConfig(), nil, nil, WithClient(m))
	require.NoError(t, err)
	return s, m
}

func TestRetrieveTokenForAPI_SuccessAndCaches(t *testing.T) {
	s, m := newServiceUnderTest(t)
	ctx := context.Background()

	m.EXPECT().
		RequestTokenWithBodyWithResponse(mock.Anything, formContentType, mock.Anything).
		Run(func(ctx context.Context, contentType string, body io.Reader, _ ...gen.RequestEditorFn) {
			b, _ := io.ReadAll(body)
			v, _ := url.ParseQuery(string(b))
			assert.Equal(t, "client_credentials", v.Get("grant_type"))
			assert.Equal(t, "test", v.Get("client_id"))
			assert.Equal(t, "test", v.Get("client_secret"))
		}).
		Return(testutil.TestAuthJSON200Response(200, "token", `{"access_token":"token1","token_type":"Bearer"}`), nil).
		Once()

	tok, err := s.RetrieveTokenForAPI(ctx, "camunda")
	require.NoError(t, err)
	assert.Equal(t, "token", tok)

	// second call should hit the cache, no new request
	tok, err = s.RetrieveTokenForAPI(ctx, "camunda")
	require.NoError(t, err)
	assert.Equal(t, "token", tok)
}

func TestRetrieveTokenForAPI_HTTPErrorStatus(t *testing.T) {
	s, m := newServiceUnderTest(t)
	ctx := context.Background()

	m.EXPECT().
		RequestTokenWithBodyWithResponse(mock.Anything, formContentType, mock.Anything).
		Return(testutil.TestAuthJSON200Response(400, "", `{"error":"invalid_client"}`), nil).
		Once()

	_, err := s.RetrieveTokenForAPI(ctx, "camunda")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "token request failed: status=400 body={\"error\":\"invalid_client\"}")
}

func TestRetrieveTokenForAPI_MissingToken(t *testing.T) {
	s, m := newServiceUnderTest(t)
	ctx := context.Background()

	m.EXPECT().
		RequestTokenWithBodyWithResponse(mock.Anything, formContentType, mock.Anything).
		Return(testutil.TestAuthJSON200Response(200, "", `{"access_token":"","token_type":"Bearer"}`), nil).
		Once()

	_, err := s.RetrieveTokenForAPI(ctx, "camunda")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing access token")
}

func TestRetrieveTokenForAPI_CleanCache(t *testing.T) {
	s, m := newServiceUnderTest(t)
	ctx := context.Background()

	m.EXPECT().
		RequestTokenWithBodyWithResponse(mock.Anything, formContentType, mock.Anything).
		Return(testutil.TestAuthJSON200Response(200, "token", `{"access_token":"token","token_type":"Bearer"}`), nil).
		Twice()

	tok, err := s.RetrieveTokenForAPI(ctx, "camunda")
	require.NoError(t, err)
	assert.Equal(t, "token", tok)

	s.ClearCache()

	tok, err = s.RetrieveTokenForAPI(ctx, "camunda")
	require.NoError(t, err)
	assert.Equal(t, "token", tok)
}

func TestFromContext_NoValue(t *testing.T) {
	svc, err := FromContext(context.Background())
	require.Nil(t, svc)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNoAuthServiceInContext)
}

func TestFromContext_WrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), ctxAuthServiceKey{}, "wrong type")
	svc, err := FromContext(ctx)
	require.Nil(t, svc)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidServiceInContext)
}
