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
	v = viper.New()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "camunder",
	Short: "Camunder is a CLI tool to interact with Camunda 8.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "camunder" {
			return nil
		}
		if cmd.Name() == "help" || cmd.Name() == "version" || cmd.Name() == "completion" {
			return nil
		}
		if cmd.Flags().Changed("help") {
			return nil
		}

		cfg, err := initConfig()
		if err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}
		cmd.SetContext(config.IntoContext(cmd.Context(), cfg))
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
		// return runUI(cmd, args)
	},
	SilenceUsage:  true,
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
	// User-facing flags (highest precedence)
	rootCmd.PersistentFlags().String("config", "", "path to config file")
	rootCmd.PersistentFlags().String("camunda8-base-url", "", "Camunda 8 API base URL")
	rootCmd.PersistentFlags().String("camunda8-token", "", "Camunda 8 API bearer token")
	rootCmd.PersistentFlags().String("operate-base-url", "", "Operate API base URL")
	rootCmd.PersistentFlags().String("operate-token", "", "Operate API bearer token")
	rootCmd.PersistentFlags().String("tasklist-base-url", "", "Tasklist API base URL")
	rootCmd.PersistentFlags().String("tasklist-token", "", "Tasklist API bearer token")
	rootCmd.PersistentFlags().Duration("timeout", 0, "HTTP timeout (e.g. 10s, 1m)")

	// Bind flags to viper keys
	// Resolve precedence: flags > env > config file > defaults
	_ = v.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = v.BindPFlag("camunda8_api.base_url", rootCmd.PersistentFlags().Lookup("camunda8-base-url"))
	_ = v.BindPFlag("camunda8_api.token", rootCmd.PersistentFlags().Lookup("camunda8-token"))
	_ = v.BindPFlag("operate_api.base_url", rootCmd.PersistentFlags().Lookup("operate-base-url"))
	_ = v.BindPFlag("operate_api.token", rootCmd.PersistentFlags().Lookup("operate-token"))
	_ = v.BindPFlag("tasklist_api.base_url", rootCmd.PersistentFlags().Lookup("tasklist-base-url"))
	_ = v.BindPFlag("tasklist_api.token", rootCmd.PersistentFlags().Lookup("tasklist-token"))
	_ = v.BindPFlag("http.timeout", rootCmd.PersistentFlags().Lookup("timeout"))
}

func initConfig() (config.Config, error) {
	// Define defaults
	v.SetDefault("camunda8_api.base_url", "http://localhost:8086/v2")
	v.SetDefault("operate_api.base_url", "http://localhost:8081/v1")
	v.SetDefault("tasklist_api.base_url", "http://localhost:8082/v1")
	v.SetDefault("http.timeout", "10s")

	// Environment variables
	v.SetEnvPrefix("CAMUNDER")
	v.AutomaticEnv()                                             // read in environment variables that match
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_")) // support nested env vars (e.g. CAMUNDER_API_BASE_URL)

	// Config file
	if cfgFile := v.GetString("config"); cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.camunder")
		v.SetConfigType("yaml")
		v.SetConfigName("config")
	}
	if err := v.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", v.ConfigFileUsed())
	}

	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return config.Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return cfg, nil
}
