package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grafvonb/camunder/pkg/camunda"
)

const (
	CamundaApiKeyConst  = "camunda_api"
	OperateApiKeyConst  = "operate_api"
	TasklistApiKeyConst = "tasklist_api"
)

var ValidAPIKeys = []string{
	CamundaApiKeyConst,
	OperateApiKeyConst,
	TasklistApiKeyConst,
}

type APIs struct {
	Version  camunda.APIVersion `mapstructure:"version"`
	Camunda  API                `mapstructure:"camunda_api"`
	Operate  API                `mapstructure:"operate_api"`
	Tasklist API                `mapstructure:"tasklist_api"`
}

func (a *APIs) Validate() error {
	var errs []error
	switch a.Version {
	case "":
		a.Version = camunda.Current
	default:
		v, err := camunda.Normalize(string(a.Version))
		if err != nil {
			errs = append(errs, fmt.Errorf("version: %w", err))
		} else {
			a.Version = v
		}
	}
	if err := a.Camunda.Validate(); err != nil {
		errs = append(errs, fmt.Errorf("camunda: %w", err))
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
