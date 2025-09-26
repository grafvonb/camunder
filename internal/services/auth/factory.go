package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth/cookie"
	"github.com/grafvonb/camunder/internal/services/auth/core"
	"github.com/grafvonb/camunder/internal/services/auth/imx"
	"github.com/grafvonb/camunder/internal/services/auth/oauth2"
)

func BuildAuthenticator(cfg *config.Config, httpClient *http.Client, log *slog.Logger) (core.Authenticator, error) {
	switch cfg.Auth.Mode {
	case config.ModeOAuth2, "":
		return oauth2.New(cfg, httpClient, log)
	case config.ModeIMX:
		return imx.New(cfg, httpClient, log)
	case config.ModeCookie:
		return cookie.New(cfg, httpClient, log)
	default:
		return nil, fmt.Errorf("unknown auth mode: %s", cfg.Auth.Mode)
	}
}
