package config

import "errors"

var (
	ErrNoBaseURL = errors.New("no base_url provided in api configuration")
	ErrNoToken   = errors.New("no token provided in api configuration")
)

type Config struct {
	Config string `mapstructure:"config"`
	API    API    `mapstructure:"api"`
	HTTP   HTTP   `mapstructure:"http"`
}

type API struct {
	BaseURL string `mapstructure:"base_url"`
	Token   string `mapstructure:"token"`
}

type HTTP struct {
	Timeout string `mapstructure:"timeout"`
}

func (c Config) Validate() error {
	if c.API.BaseURL == "" {
		return ErrNoBaseURL
	}
	if c.API.Token == "" {
		return ErrNoToken
	}

	return nil
}
