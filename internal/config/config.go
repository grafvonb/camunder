package config

import (
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type API struct {
	BaseURL string `mapstructure:"base_url"`
	Token   string `mapstructure:"token"`
}

type HTTP struct {
	Timeout time.Duration `mapstructure:"timeout"`
}

type Config struct {
	API  API  `mapstructure:"api"`
	HTTP HTTP `mapstructure:"http"`
}

// Defaults sets sensible defaults on the provided viper instance.
func Defaults(v *viper.Viper) {
	v.SetDefault("api.base_url", "http://localhost:8086/v2")
	v.SetDefault("http.timeout", "10s")
}

// LoadFrom unmarshals config values from the given viper into a Config struct.
func LoadFrom(v *viper.Viper) (Config, error) {
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}

// HTTPClient returns a ready-to-use http.Client based on the config.
func (c Config) HTTPClient() *http.Client {
	return &http.Client{Timeout: c.HTTP.Timeout}
}
