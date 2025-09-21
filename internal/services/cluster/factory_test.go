package cluster

import (
	"net/http"
	"testing"

	"log/slog"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/stretchr/testify/require"
)

func testConfig() *config.Config {
	return &config.Config{
		APIs: config.APIs{},
	}
}

func TestFactory_V87(t *testing.T) {
	cfg := testConfig()
	cfg.APIs.Version = camunda.V87
	svc, err := New(cfg, &http.Client{}, nil, slog.Default())
	require.NoError(t, err)
	require.NotNil(t, svc)
}

func TestFactory_V88(t *testing.T) {
	cfg := testConfig()
	cfg.APIs.Version = camunda.V88
	svc, err := New(cfg, &http.Client{}, nil, slog.Default())
	require.NoError(t, err)
	require.NotNil(t, svc)
}

func TestFactory_Unknown(t *testing.T) {
	cfg := testConfig()
	cfg.APIs.Version = "v0"
	svc, err := New(cfg, &http.Client{}, nil, slog.Default())
	require.Error(t, err)
	require.Nil(t, svc)
	require.Contains(t, err.Error(), "unknown Camunda APIs version")
}
