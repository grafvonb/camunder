package config

import (
	"errors"
	"fmt"
	"strings"
)

type AuthMode string

func (m AuthMode) IsValid() bool { return m == ModeOAuth2 || m == ModeXSRF || m == ModeCookie }

const (
	ModeOAuth2 AuthMode = "oauth2"
	ModeXSRF   AuthMode = "xsrf"
	ModeCookie AuthMode = "cookie"
)

type Auth struct {
	Mode   AuthMode                    `mapstructure:"mode"`
	OAuth2 AuthOAuth2ClientCredentials `mapstructure:"oauth2"`
	XSRF   AuthXsrfSession             `mapstructure:"xsrf"`
	Cookie AuthCookieSession           `mapstructure:"cookie"`
}

func (c *Auth) Validate() error {
	var errs []error
	if !c.Mode.IsValid() {
		errs = append(errs, fmt.Errorf("mode: invalid value %q (allowed values: %q, %q)", c.Mode, ModeOAuth2, ModeXSRF))
	} else {
		switch c.Mode {
		case ModeOAuth2:
			if err := c.OAuth2.Validate(); err != nil {
				errs = append(errs, fmt.Errorf("oauth2: %w", err))
			}
		case ModeXSRF:
			if err := c.XSRF.Validate(); err != nil {
				errs = append(errs, fmt.Errorf("xsrf: %w", err))
			}
		case ModeCookie:
			if err := c.Cookie.Validate(); err != nil {
				errs = append(errs, fmt.Errorf("cookie: %w", err))
			}
		}
	}
	return errors.Join(errs...)
}

type AuthXsrfSession struct {
	BaseURL  string `mapstructure:"base_url"`
	AppId    string `mapstructure:"app_id"`
	Module   string `mapstructure:"module"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func (c *AuthXsrfSession) Validate() error {
	var errs []error

	if strings.TrimSpace(c.BaseURL) == "" {
		errs = append(errs, ErrNoXSRFBaseURL)
	}
	if strings.TrimSpace(c.AppId) == "" {
		errs = append(errs, ErrNoXSRFAppID)
	}
	if strings.TrimSpace(c.Module) == "" {
		errs = append(errs, ErrNoXSRFModule)
	}
	if strings.TrimSpace(c.User) == "" {
		errs = append(errs, ErrNoXSRFUser)
	}
	if strings.TrimSpace(c.Password) == "" {
		errs = append(errs, ErrNoXSRFPassword)
	}

	return errors.Join(errs...)
}

type AuthOAuth2ClientCredentials struct {
	TokenURL     string            `mapstructure:"token_url"`
	ClientID     string            `mapstructure:"client_id"`
	ClientSecret string            `mapstructure:"client_secret"`
	Scopes       map[string]string `mapstructure:"scopes"`
}

var allowedScopeKeys = map[string]struct{}{CamundaApiKeyConst: {}, OperateApiKeyConst: {}, TasklistApiKeyConst: {}}
var allowedScopeKeysList = []string{CamundaApiKeyConst, OperateApiKeyConst, TasklistApiKeyConst}

func (a *AuthOAuth2ClientCredentials) Validate() error {
	var errs []error

	if strings.TrimSpace(a.TokenURL) == "" {
		errs = append(errs, ErrNoTokenURL)
	}
	if strings.TrimSpace(a.ClientID) == "" {
		errs = append(errs, ErrNoClientID)
	}
	if strings.TrimSpace(a.ClientSecret) == "" {
		errs = append(errs, ErrNoClientSecret)
	}

	if len(a.Scopes) > 0 {
		for k := range a.Scopes {
			key := strings.TrimSpace(k)
			if key == "" {
				errs = append(errs, fmt.Errorf("auth.scopes contains an empty key (allowed keys: %s)",
					strings.Join(allowedScopeKeysList, ", ")))
				continue
			}
			if _, ok := allowedScopeKeys[key]; !ok {
				errs = append(errs, fmt.Errorf("auth.scopes[%s]: unsupported key (allowed keys: %s)",
					k, strings.Join(allowedScopeKeysList, ", ")))
			}
		}
	}

	return errors.Join(errs...)
}

func (a *AuthOAuth2ClientCredentials) Scope(key string) string {
	if a.Scopes == nil {
		return ""
	}
	return a.Scopes[key]
}

type AuthCookieSession struct {
	BaseURL  string `mapstructure:"base_url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (c *AuthCookieSession) Validate() error {
	var errs []error
	if strings.TrimSpace(c.BaseURL) == "" {
		errs = append(errs, errors.New("no base_url provided in cookie auth configuration"))
	}
	if strings.TrimSpace(c.Username) == "" {
		errs = append(errs, errors.New("no username provided in cookie auth configuration"))
	}
	if strings.TrimSpace(c.Password) == "" {
		errs = append(errs, errors.New("no password provided in cookie auth configuration"))
	}
	return errors.Join(errs...)
}
