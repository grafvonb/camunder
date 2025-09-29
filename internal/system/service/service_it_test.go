package systemservice_test

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/grafvonb/camunder/internal/config"
	systemservice "github.com/grafvonb/camunder/internal/system/service"
	"github.com/grafvonb/camunder/internal/testx"
	"github.com/stretchr/testify/require"
)

func Test_GetSystems_OK_IT(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}

	svc := getService(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	systems, err := svc.GetSystems(ctx)
	require.NoError(t, err)
	require.NotNil(t, systems)
	require.Greater(t, len(*systems), 0)
	t.Logf("success: got %d systems", len(*systems))
}

func TestService_GetSystemByKey(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}

	svc := getService(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	key := testx.RequireEnvWithPrefix(t, "SYSTEM_KEY")
	sys, err := svc.GetSystemByKey(ctx, key)
	require.NoError(t, err)
	require.NotNil(t, sys)

	require.NotEmpty(t, sys.DisplayName)

}

func getService(t *testing.T) *systemservice.Service {
	t.Helper()
	cfg, httpClient, log := getSetup(t)
	svc, err := systemservice.New(cfg, httpClient, log)
	require.NoError(t, err)
	require.NotNil(t, svc)
	return svc
}

func getSetup(t *testing.T) (*config.Config, *http.Client, *slog.Logger) {
	t.Helper()
	cfg := getConfig(t)
	httpClient := &http.Client{Timeout: 15 * time.Second}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return cfg, httpClient, log
}

func getConfig(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		APIs: config.APIs{
			Camunda: config.API{
				BaseURL: testx.RequireEnv(t, "BASE_URL"),
			},
		},
	}
}
