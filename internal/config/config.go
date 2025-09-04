package config

import (
	"errors"
	"fmt"
)

var (
	ErrNoBaseURL = errors.New("no base_url provided in api configuration")
	ErrNoToken   = errors.New("no token provided in api configuration")
)

type Config struct {
	Config      string `mapstructure:"config"`
	Camunda8API API    `mapstructure:"camunda8_api"`
	OperateAPI  API    `mapstructure:"operate_api"`
	TasklistAPI API    `mapstructure:"tasklist_api"`
	HTTP        HTTP   `mapstructure:"http"`
}

func (c Config) Validate() error {
	if err := c.Camunda8API.Validate(); err != nil {
		return fmt.Errorf("camunda8 API: %w", err)
	}
	if err := c.OperateAPI.Validate(); err != nil {
		return fmt.Errorf("operate API: %w", err)
	}
	if err := c.TasklistAPI.Validate(); err != nil {
		return fmt.Errorf("tasklist API: %w", err)
	}
	return nil
}

type API struct {
	BaseURL string `mapstructure:"base_url"`
	Token   string `mapstructure:"token"`
}

func (a API) Validate() error {
	if a.BaseURL == "" {
		return ErrNoBaseURL
	}
	if a.Token == "" {
		return ErrNoToken
	}

	return nil
}

type HTTP struct {
	Timeout string `mapstructure:"timeout"`
}
