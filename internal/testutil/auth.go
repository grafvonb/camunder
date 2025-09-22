package testutil

import (
	"net/http"

	"github.com/grafvonb/camunder/internal/api/gen/clients/auth"
	"github.com/grafvonb/camunder/internal/config"
)

type tokenJSON200 = struct {
	AccessToken  string  `json:"access_token"`
	ExpiresIn    int     `json:"expires_in"`
	IdToken      *string `json:"id_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	Scope        *string `json:"scope,omitempty"`
	TokenType    string  `json:"token_type"`
}

func TestAuthJSON200Response(status int, token string, raw string) *auth.RequestTokenResponse {
	return &auth.RequestTokenResponse{
		Body: []byte(raw),
		JSON200: &tokenJSON200{
			AccessToken: token,
			TokenType:   "Bearer",
		},
		HTTPResponse: &http.Response{StatusCode: status},
	}
}

func TestConfig() *config.Config {
	return &config.Config{
		App: config.App{
			Tenant: "tenant",
		},
		Auth: config.Authentication{
			TokenURL:     "http://localhost/token",
			ClientID:     "test",
			ClientSecret: "test",
		},
		APIs: config.APIs{
			Camunda: config.API{
				BaseURL: "http://localhost/camunda",
			},
			Operate: config.API{
				BaseURL: "http://localhost/operate",
			},
			Tasklist: config.API{
				BaseURL: "http://localhost/tasklist",
			},
		},
	}
}
