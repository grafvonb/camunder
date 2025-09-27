package cluster

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/camunder/internal/config"
	v87 "github.com/grafvonb/camunder/internal/services/cluster/v87"
	v88 "github.com/grafvonb/camunder/internal/services/cluster/v88"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/cluster"
)

func New(cfg *config.Config, httpClient *http.Client, log *slog.Logger) (cluster.API, error) {
	v := cfg.APIs.Version
	switch v {
	case camunda.V88:
		return v88.New(cfg, httpClient, log)
	case camunda.V87:
		return v87.New(cfg, httpClient, log)
	default:
		return nil, fmt.Errorf("%w: %q (supported: %v)", camunda.ErrUnknownAPIVersion, v, camunda.Supported())
	}
}
