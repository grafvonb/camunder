//go:build integration

package xsrf_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"log/slog"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth/xsrf"
	"github.com/grafvonb/camunder/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestXSRF_Login_IT(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}

	cfg := &config.Config{
		Auth: config.Auth{
			XSRF: config.AuthXsrfSession{
				BaseURL:  testx.RequireEnvWithPrefix(t, "XSRF_BASE_URL"),
				AppId:    testx.RequireEnvWithPrefix(t, "XSRF_APP_ID"),
				Module:   testx.RequireEnvWithPrefix(t, "XSRF_MODULE"),
				User:     testx.RequireEnvWithPrefix(t, "XSRF_USER"),
				Password: testx.RequireEnvWithPrefix(t, "XSRF_PASSWORD"),
			},
		},
	}

	httpClient := &http.Client{Timeout: 15 * time.Second}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	svc, err := xsrf.New(cfg, httpClient, log)
	require.NoError(t, err)
	require.NotNil(t, svc)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = svc.Init(ctx)
	require.NoError(t, err)
	require.True(t, svc.IsAuthenticated(), "expected authenticated")
	require.NotEmpty(t, svc.Token(), "expected non-empty token")
	t.Logf("success: got xsrf token: %q...", svc.Token())
}
