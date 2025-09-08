package config

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Camunda8ApiKeyConst = "camunda8_api"
	OperateApiKeyConst  = "operate_api"
	TasklistApiKeyConst = "tasklist_api"
)

var ValidAPIKeys = []string{
	Camunda8ApiKeyConst,
	OperateApiKeyConst,
	TasklistApiKeyConst,
}

type APIs struct {
	Camunda8 API `mapstructure:"camunda8_api"`
	Operate  API `mapstructure:"operate_api"`
	Tasklist API `mapstructure:"tasklist_api"`
}

func (a *APIs) Validate() error {
	var errs []error
	if err := a.Camunda8.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("camunda8: %w", err))
	}
	if err := a.Operate.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("operate: %w", err))
	}
	if err := a.Tasklist.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("tasklist: %w", err))
	}
	return errors.Join(errs...)
}

type API struct {
	Key     string `mapstructure:"key"`
	BaseURL string `mapstructure:"base_url"`
}

func (a *API) Validate() error {
	if strings.TrimSpace(a.BaseURL) == "" {
		return ErrNoBaseURL
	}
	return nil
}
