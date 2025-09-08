package config

import "github.com/grafvonb/camunder/internal/services/common"

type App struct {
	Tenant  string               `mapstructure:"tenant"`
	Backoff common.BackoffConfig `mapstructure:"backoff"`
}
