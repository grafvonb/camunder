package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth"
	"github.com/grafvonb/camunder/internal/services/httpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagQuiet      bool // quiet mode, suppress output, use exit code only
	flagShowConfig bool // show effective config and exit
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "camunder",
	Short: "Camunder is a CLI tool to interact with Camunda 8.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for root/help-like commands
		if cmd.Name() == "help" || cmd.Name() == "version" || cmd.Name() == "completion" {
			return nil
		}
		if cmd.Flags().Changed("help") {
			return nil
		}

		v := viper.New()
		if err := initViper(v, cmd); err != nil {
			return err
		}
		cfg, err := retrieveConfig(v, !flagShowConfig)
		if err != nil {
			return err
		}
		cmd.SetContext(cfg.ToContext(cmd.Context()))
		if flagShowConfig {
			cfgpath := v.ConfigFileUsed()
			if cfgpath == "" {
				cfgpath = "(none)"
			}
			cmd.Println("config loaded: "+cfgpath, v.ConfigFileUsed())
			cmd.Println(cfg.String())
			os.Exit(0)
		}
		httpSvc, err := httpc.New(cfg, flagQuiet)
		if err != nil {
			return fmt.Errorf("create http service: %w", err)
		}
		cmd.SetContext(httpSvc.ToContext(cmd.Context()))
		authSvc, err := auth.New(cfg, httpSvc.Client(), flagQuiet)
		if err != nil {
			return fmt.Errorf("create auth service: %w", err)
		}
		if err := authSvc.Warmup(cmd.Context()); err != nil {
			cmd.PrintErrf("warming up auth service: %v\n", err)
			return err
		}
		cmd.SetContext(authSvc.ToContext(cmd.Context()))

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
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	pf := rootCmd.PersistentFlags()

	pf.String("config", "", "path to config file")

	pf.String("tenant", "", "default tenant ID")

	pf.String("auth-token-url", "", "auth token URL")
	pf.String("auth-client-id", "", "auth client ID")
	pf.String("auth-client-secret", "", "auth client secret")
	pf.StringToString("auth-scopes", nil, "auth scopes as key=value (repeatable or comma-separated)")

	pf.String("http-timeout", "", "HTTP timeout (Go duration, e.g. 30s)")

	pf.String("camunda8-base-url", "", "Camunda8 API base URL")
	pf.String("operate-base-url", "", "Operate API base URL")
	pf.String("tasklist-base-url", "", "Tasklist API base URL")

	pf.BoolVar(&flagQuiet, "quiet", false, "suppress output, use exit code only")

	// TODO show-config flag should be in a "config" subcommand
	pf.BoolVar(&flagShowConfig, "show-config", false, "print effective config (secrets redacted)")

	// TODO add --dry-run flag to commands that perform actions
}

func initViper(v *viper.Viper, cmd *cobra.Command) error {
	// Resolve precedence: flags > env > config file > defaults

	// Bind scalar flags directly to config keys
	_ = v.BindPFlag("app.tenant", cmd.Flags().Lookup("tenant"))
	_ = v.BindPFlag("config", cmd.Flags().Lookup("config"))
	_ = v.BindPFlag("auth.token_url", cmd.Flags().Lookup("auth-token-url"))
	_ = v.BindPFlag("auth.client_id", cmd.Flags().Lookup("auth-client-id"))
	_ = v.BindPFlag("auth.client_secret", cmd.Flags().Lookup("auth-client-secret"))
	_ = v.BindPFlag("http.timeout", cmd.Flags().Lookup("http-timeout"))

	_ = v.BindPFlag("apis.camunda8_api.base_url", cmd.Flags().Lookup("camunda8-base-url"))
	_ = v.BindPFlag("apis.operate_api.base_url", cmd.Flags().Lookup("operate-base-url"))
	_ = v.BindPFlag("apis.tasklist_api.base_url", cmd.Flags().Lookup("tasklist-base-url"))

	// Bind map flag to a tmp key so we can merge later
	_ = v.BindPFlag("tmp.auth_scopes", cmd.Flags().Lookup("auth-scopes"))

	// Force hardcoded keys
	v.Set("apis.camunda8_api.key", config.Camunda8ApiKeyConst)
	v.Set("apis.operate_api.key", config.OperateApiKeyConst)
	v.Set("apis.tasklist_api.key", config.TasklistApiKeyConst)

	// Defaults
	v.SetDefault("http.timeout", "30s")

	// Config file discovery
	if cfgFile := v.GetString("config"); cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")

		// Search config paths (in order):
		// Look in the current dir (./config.yaml)
		// Then $XDG_CONFIG_HOME/camunder/config.yaml
		// Then $HOME/.config/camunder/config.yaml
		// Finally fallback to $HOME/.camunder/config.yaml
		v.AddConfigPath(".")
		if xdg, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok && xdg != "" {
			v.AddConfigPath(filepath.Join(xdg, "camunder"))
		} else if home, err := os.UserHomeDir(); err == nil {
			v.AddConfigPath(filepath.Join(home, ".config", "camunder"))
		}
		if home, err := os.UserHomeDir(); err == nil {
			v.AddConfigPath(filepath.Join(home, ".camunder"))
		}
	}

	// ENV: CAMUNDER_AUTH_CLIENT_ID, etc.
	v.SetEnvPrefix("CAMUNDER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config (ignore "not found")
	if err := v.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		if !errors.As(err, &nf) {
			return fmt.Errorf("read config: %w", err)
		}
	}
	return nil
}

func retrieveConfig(v *viper.Viper, validate bool) (*config.Config, error) {
	var cfg config.Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if tmpScopes := v.GetStringMapString("tmp.auth_scopes"); len(tmpScopes) > 0 {
		if cfg.Auth.Scopes == nil {
			cfg.Auth.Scopes = make(map[string]string, len(tmpScopes))
		}
		for k, scope := range tmpScopes {
			cfg.Auth.Scopes[strings.TrimSpace(k)] = strings.TrimSpace(scope)
		}
	}

	if validate {
		if err := cfg.Validate(); err != nil {
			return nil, fmt.Errorf("validate config\n%w", err)
		}
	}
	return &cfg, nil
}
