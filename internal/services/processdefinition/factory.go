package processdefinition

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grafvonb/camunder/internal/config"
	"github.com/grafvonb/camunder/internal/services/auth"
	v87 "github.com/grafvonb/camunder/internal/services/processdefinition/v87"
	v88 "github.com/grafvonb/camunder/internal/services/processdefinition/v88"
	"github.com/grafvonb/camunder/pkg/camunda"
	"github.com/grafvonb/camunder/pkg/camunda/processdefinition"
)

func New(cfg *config.Config, httpClient *http.Client, auth *auth.Service, log *slog.Logger) (processdefinition.API, error) {
	v := cfg.APIs.Version
	switch v {
	case camunda.V88:
		return v88.New(cfg, httpClient, auth, log)
	case camunda.V87:
		return v87.New(cfg, httpClient, auth, log)
	default:
		return nil, fmt.Errorf("%w: %q (supported: %v)", camunda.ErrUnknownAPIVersion, v, camunda.Supported())
	}
}
