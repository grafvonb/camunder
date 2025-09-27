//go:build integration_tiny

package cookie_test

import (
	"context"
	"log/slog"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth/cookie"
	"github.com/grafvonb/camunder/internal/testx"
	"github.com/stretchr/testify/require"
)

func TestCookie_Login_OK_Tiny_IT(t *testing.T) {
	srv := testx.StartAuthServerCookie(t, testx.CookieAuthOpts{
		SetCookie: true,
		ExpectUser: struct {
			Name     string
			Password string
		}{Name: "demo", Password: "demo"},
	})
	defer srv.Close()

	jar, _ := cookiejar.New(nil)
	httpClient := srv.TS.Client()
	httpClient.Jar = jar
	httpClient.Timeout = 5 * time.Second

	cfg := &config.Config{
		Auth: config.Auth{
			Mode: config.ModeCookie,
			Cookie: config.AuthCookieSession{
				BaseURL:  srv.BaseURL,
				Username: "demo",
				Password: "demo",
			},
		},
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc, err := cookie.New(cfg, httpClient, log)
	require.NoError(t, err)
	require.NotNil(t, svc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Logf("trying to authenticate aginst %s with user %q", cfg.Auth.Cookie.BaseURL, cfg.Auth.Cookie.Username)
	err = svc.Init(ctx)
	require.NoError(t, err)
	require.True(t, svc.IsAuthenticated())
	t.Log("success: got authenticated")
}
