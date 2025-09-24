package testutil

import (
	"net/http"

	"github.com/grafvonb/camunder/internal/api/gen/clients/auth/oauth2"
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

func TestAuthJSON200Response(status int, token string, raw string) *oauth2.RequestTokenResponse {
	return &oauth2.RequestTokenResponse{
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
		Auth: config.Auth{
			OAuth2: config.AuthOAuth2ClientCredentials{
				TokenURL:     "http://localhost/token",
				ClientID:     "test",
				ClientSecret: "test",
			},
			IMX: config.AuthImxSession{
				BaseURL:  "http://localhost/imx",
				AppId:    "test",
				Module:   "test",
				User:     "test",
				Password: "test",
			},
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
