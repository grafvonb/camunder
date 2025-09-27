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
)

func TestXSRF_Login_IT(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}

	cfg := &config.Config{}
	cfg.Auth.XSRF.BaseURL = testx.RequireEnvWithPrefix(t, "XSRF_BASE_URL")
	cfg.Auth.XSRF.AppId = testx.RequireEnvWithPrefix(t, "XSRF_APP_ID")
	cfg.Auth.XSRF.Module = testx.RequireEnvWithPrefix(t, "XSRF_MODULE")
	cfg.Auth.XSRF.User = testx.RequireEnvWithPrefix(t, "XSRF_USER")
	cfg.Auth.XSRF.Password = testx.RequireEnvWithPrefix(t, "XSRF_PASSWORD")

	hc := &http.Client{Timeout: 15 * time.Second}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	svc, err := xsrf.New(cfg, hc, log)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := svc.Init(ctx); err != nil {
		t.Fatalf("Init: %v", err)
	}

	if !svc.IsAuthenticated() {
		t.Fatal("expected authenticated")
	}
	t.Logf("got xsrf token: %q...", svc.Token())
}
