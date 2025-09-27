//go:build integration

package imx_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"log/slog"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth/imx"
	"github.com/grafvonb/camunder/internal/testx"
)

func TestIMX_Login_IT(t *testing.T) {
	if testing.Short() {
		t.Skip("short")
	}

	cfg := &config.Config{}
	cfg.Auth.IMX.BaseURL = testx.RequireEnvWithPrefix(t, "IMX_BASE_URL")
	cfg.Auth.IMX.AppId = testx.RequireEnvWithPrefix(t, "IMX_APP_ID")
	cfg.Auth.IMX.Module = testx.RequireEnvWithPrefix(t, "IMX_MODULE")
	cfg.Auth.IMX.User = testx.RequireEnvWithPrefix(t, "IMX_USER")
	cfg.Auth.IMX.Password = testx.RequireEnvWithPrefix(t, "IMX_PASSWORD")

	hc := &http.Client{Timeout: 15 * time.Second}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	svc, err := imx.New(cfg, hc, log)
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
