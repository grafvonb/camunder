package v87

import (
	"net/http"
	"testing"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth"
	"github.com/stretchr/testify/require"
)

// nolint
func testConfig() *config.Config {
	return &config.Config{}
}

// nolint
func newServiceUnderTest(t *testing.T, auth auth.AuthClient) (*Service, *MockGenClusterClient) {
	t.Helper()
	m := NewMockGenClusterClient(t)
	s, err := New(testConfig(), &http.Client{}, auth, nil, WithClient(m))
	require.NoError(t, err)
	return s, m
}
