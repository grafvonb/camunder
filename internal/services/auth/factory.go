package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth/core"
	"github.com/grafvonb/camunder/internal/services/auth/imx"
	"github.com/grafvonb/camunder/internal/services/auth/oauth2"
)

func BuildAuthenticator(cfg *config.Config, hc *http.Client, log *slog.Logger) (core.Authenticator, error) {
	switch cfg.Auth.Mode {
	case config.ModeOAuth2, "":
		return oauth2.New(cfg, hc, log)
	case config.ModeIMX:
		return imx.New(cfg, hc, log)
	default:
		return nil, fmt.Errorf("unknown auth mode: %s", cfg.Auth.Mode)
	}
}
