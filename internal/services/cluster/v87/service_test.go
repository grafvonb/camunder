package v87

import (
	"net/http"
	"testing"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth"
	"github.com/grafvonb/camunder/internal/services/cluster/v87/mocks"
	"github.com/stretchr/testify/require"
)

// nolint
func testConfig() *config.Config {
	return &config.Config{}
}

// nolint
func newServiceUnderTest(t *testing.T, auth auth.AuthClient) (*Service, *v87mock.MockGenClusterClient) {
	t.Helper()
	m := v87mock.NewMockGenClusterClient(t)
	s, err := New(testConfig(), &http.Client{}, auth, nil, WithClient(m))
	require.NoError(t, err)
	return s, m
}
