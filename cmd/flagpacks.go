package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultBackoffStrategy   = "exponential"
	defaultBackoffMultiplier = 2.0
)

var (
	defaultBackoffInitialDelay = 500 * time.Millisecond
	defaultBackoffMaxDelay     = 8 * time.Second
	defaultBackoffMaxRetries   = 0 // 0 = unlimited
	defaultBackoffTimeout      = 2 * time.Minute
)

func AddApiCommandsFlagsAndBindings(cmd *cobra.Command, v *viper.Viper) {
	// fs := cmd.PersistentFlags()

	// fs.String("api-url", "", "Base URL of the Camunda 8 API (e.g. https://operate.camunda.io/your-cluster-id)")
	// _ = v.BindPFlag("app.api.url", fs.Lookup("api-url"))
	// v.SetDefault("app.api.url", "")
}

func AddBackoffFlagsAndBindings(cmd *cobra.Command, v *viper.Viper) {
	fs := cmd.PersistentFlags()

	fs.String("backoff-strategy", defaultBackoffStrategy, "Backoff strategy: fixed|exponential")
	fs.Duration("backoff-initial-delay", defaultBackoffInitialDelay, "Initial delay between retries")
	fs.Duration("backoff-max-delay", defaultBackoffMaxDelay, "Maximum delay between retries")
	fs.Int("backoff-max-retries", defaultBackoffMaxRetries, "Max retry attempts (0 = unlimited)")
	fs.Float64("backoff-multiplier", defaultBackoffMultiplier, "Exponential multiplier (>1)")
	fs.Duration("backoff-timeout", defaultBackoffTimeout, "Overall timeout for the retry loop")

	_ = v.BindPFlag("app.backoff.strategy", fs.Lookup("backoff-strategy"))
	_ = v.BindPFlag("app.backoff.initial_delay", fs.Lookup("backoff-initial-delay"))
	_ = v.BindPFlag("app.backoff.max_delay", fs.Lookup("backoff-max-delay"))
	_ = v.BindPFlag("app.backoff.max_retries", fs.Lookup("backoff-max-retries"))
	_ = v.BindPFlag("app.backoff.multiplier", fs.Lookup("backoff-multiplier"))
	_ = v.BindPFlag("app.backoff.timeout", fs.Lookup("backoff-timeout"))

	v.SetDefault("app.backoff.strategy", defaultBackoffStrategy)
	v.SetDefault("app.backoff.initial_delay", defaultBackoffInitialDelay)
	v.SetDefault("app.backoff.max_delay", defaultBackoffMaxDelay)
	v.SetDefault("app.backoff.max_retries", defaultBackoffMaxRetries)
	v.SetDefault("app.backoff.multiplier", defaultBackoffMultiplier)
	v.SetDefault("app.backoff.timeout", defaultBackoffTimeout)
}
