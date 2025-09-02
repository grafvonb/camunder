package config

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type API struct {
	BaseURL string `mapstructure:"base_url"`
	Token   string `mapstructure:"token"`
}

type HTTP struct {
	Timeout int `mapstructure:"timeout"`
}

type C87Config struct {
	Camunda8API API  `mapstructure:"camunda8_api"`
	OperateAPI  API  `mapstructure:"operate_api"`
	TasklistAPI API  `mapstructure:"tasklist_api"`
	HTTP        HTTP `mapstructure:"http"`
}

func Load(cmd *cobra.Command) (C87Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.camunder")
	v.SetEnvPrefix("CAMUNDER")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	v.SetDefault("camunda8_api.base_url", "http://localhost:8086/v2")
	v.SetDefault("http.timeout", "10")

	// config file is optional
	err := v.ReadInConfig()
	if err != nil {
		cmd.PrintErrln(err)
	} else {
		cmd.Println("config file loaded: " + v.ConfigFileUsed())
	}

	var cfg C87Config
	if err := v.Unmarshal(&cfg); err != nil {
		return C87Config{}, err
	}

	return cfg, nil
}
