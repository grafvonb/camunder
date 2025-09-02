package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel string
)

var (
	cfg config.Config
	v   = viper.New()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "camunder",
	Short: "Camunder is a CLI tool to interact with Camunda 8",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config.Defaults(v)
		// (read file/env/flags as you already do)
		_ = v.ReadInConfig()
		v.SetEnvPrefix("CAMUNDER")
		v.AutomaticEnv()

		var err error
		cfg, err = config.LoadFrom(v)
		if err != nil {
			return err
		}

		// <- make cfg available to ALL subcommands
		cmd.SetContext(config.IntoContext(cmd.Context(), cfg))
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Example usage of loaded config:
		fmt.Printf("BaseURL=%s\n", cfg.API.BaseURL)
		fmt.Printf("Token=%q\n", cfg.API.Token)
		fmt.Printf("Timeout=%s\n", cfg.HTTP.Timeout)
		return cmd.Help()
		// return runUI(cmd, args)
	},
	SilenceUsage:  false,
	SilenceErrors: false,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)

	// User-facing flags (highest precedence)
	rootCmd.PersistentFlags().String("config", "", "Path to config file")
	rootCmd.PersistentFlags().String("base-url", "", "API base URL")
	rootCmd.PersistentFlags().String("token", "", "API bearer token")
	rootCmd.PersistentFlags().Duration("timeout", 0, "HTTP timeout (e.g. 10s, 1m)")

	// Bind flags to viper keys
	_ = v.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = v.BindPFlag("api.base_url", rootCmd.PersistentFlags().Lookup("base-url"))
	_ = v.BindPFlag("api.token", rootCmd.PersistentFlags().Lookup("token"))
	_ = v.BindPFlag("http.timeout", rootCmd.PersistentFlags().Lookup("timeout"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("CAMUNDER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.SetDefault("log-level", "info")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".camunder" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".camunder")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
