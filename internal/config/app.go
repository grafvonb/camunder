package config

import "github.com/grafvonb/camunder/internal/services/common"

type App struct {
	Backoff common.BackoffConfig `mapstructure:"backoff"`
}
