package xsrf_test

import (
	"context"
	"log/slog"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth/xsrf"
	"github.com/grafvonb/camunder/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestXsrf_Login_OK(t *testing.T) {
	srv := testx.StartAuthServerXSRF(t, testx.XsrfAuthOpts{
		SetSessionCookie: true,
		SetXSRFToken:     true,
	})
	defer srv.Close()

	jar, _ := cookiejar.New(nil)
	httpClient := srv.TS.Client()
	httpClient.Jar = jar
	httpClient.Timeout = 5 * time.Second

	cfg := &config.Config{
		Auth: config.Auth{
			Mode: config.ModeXSRF,
			XSRF: config.AuthXsrfSession{
				BaseURL:  srv.BaseURL,
				AppId:    "app",
				Module:   "module",
				User:     "demo",
				Password: "demo",
			},
		},
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc, err := xsrf.New(cfg, httpClient, log)
	require.NoError(t, err)
	require.NotNil(t, svc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Logf("trying to authenticate aginst %s with user %q", cfg.Auth.XSRF.BaseURL, cfg.Auth.XSRF.User)
	err = svc.Init(ctx)
	require.NoError(t, err)
	require.True(t, svc.IsAuthenticated())
	t.Log("success: got authenticated")
}
