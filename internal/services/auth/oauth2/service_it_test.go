//go:build integration

package oauth2_test

import (
	"context"
	"os"
	"testing"
	"time"

	"log/slog"
	"net/http"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth/oauth2"
)

func requireEnv(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Skipf("missing %s; skipping integration test", key)
	}
	return v
}

func TestOAuth2_TokenAndEditor_IT(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}

	tokenURL := requireEnv(t, "OAUTH_TOKEN_URL")
	clientID := requireEnv(t, "OAUTH_CLIENT_ID")
	clientSecret := requireEnv(t, "OAUTH_CLIENT_SECRET")
	// scope := os.Getenv("OAUTH_SCOPE")   // optional
	target := os.Getenv("OAUTH_TARGET") // optional; defaults to req host

	cfg := &config.Config{}
	cfg.Auth.OAuth2.TokenURL = tokenURL
	cfg.Auth.OAuth2.ClientID = clientID
	cfg.Auth.OAuth2.ClientSecret = clientSecret

	httpClient := &http.Client{Timeout: 15 * time.Second}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	svc, err := oauth2.New(cfg, httpClient, log)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tok, err := svc.RetrieveTokenForAPI(ctx, target)
	if err != nil {
		t.Fatalf("RetrieveTokenForAPI: %v", err)
	}
	if tok == "" {
		t.Fatalf("empty access token")
	}
	t.Logf("got token: %q...", tok[:30])

	// Editor adds header on non-token URL
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com/", nil)
	_ = svc.Editor()(ctx, req)
	if got := req.Header.Get("Authorization"); got == "" {
		t.Fatalf("Authorization header not set")
	}

	// Editor must NOT add header on token URL
	req2, _ := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, nil)
	_ = svc.Editor()(ctx, req2)
	if req2.Header.Get("Authorization") != "" {
		t.Fatalf("editor must skip token URL")
	}
}
