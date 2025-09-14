package config

import (
	"errors"
	"fmt"
	"strings"
)

type Authentication struct {
	TokenURL     string            `mapstructure:"token_url"`
	ClientID     string            `mapstructure:"client_id"`
	ClientSecret string            `mapstructure:"client_secret"`
	Scopes       map[string]string `mapstructure:"scopes"`
}

var allowedScopeKeys = map[string]struct{}{CamundaApiKeyConst: {}, OperateApiKeyConst: {}, TasklistApiKeyConst: {}}
var allowedScopeKeysList = []string{CamundaApiKeyConst, OperateApiKeyConst, TasklistApiKeyConst}

func (a *Authentication) Validate() error {
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

func (a *Authentication) Scope(key string) string {
	if a.Scopes == nil {
		return ""
	}
	return a.Scopes[key]
}
