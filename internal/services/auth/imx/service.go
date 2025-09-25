package imx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"

	imxapi "github.com/grafvonb/camunder/internal/api/gen/clients/auth/imx"
	"github.com/grafvonb/camunder/internal/config"
	authcore "github.com/grafvonb/camunder/internal/services/auth/core"
)

type Service struct {
	c         GenAuthClient
	cfg       *config.Config
	http      *http.Client
	log       *slog.Logger
	baseURL   *url.URL
	xsrfToken string

	mu sync.Mutex
}

type Option func(*Service)

func WithClient(c GenAuthClient) Option {
	return func(s *Service) { s.c = c }
}

func WithHTTPClient(h *http.Client) Option {
	return func(s *Service) { s.http = h }
}

func New(cfg *config.Config, hc *http.Client, log *slog.Logger, opts ...Option) (*Service, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}
	if hc == nil {
		hc = http.DefaultClient
	}
	if hc.Jar == nil {
		jar, _ := cookiejar.New(nil)
		hc.Jar = jar
	}

	baseURL, err := url.Parse(cfg.Auth.IMX.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse imx base url: %w", err)
	}
	c, err := imxapi.NewClientWithResponses(cfg.Auth.IMX.BaseURL, imxapi.WithHTTPClient(hc))
	if err != nil {
		return nil, fmt.Errorf("init imx auth client: %w", err)
	}
	s := &Service{
		c:       c,
		cfg:     cfg,
		http:    hc,
		log:     log,
		baseURL: baseURL,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Service) Name() string { return "imx" }

func (s *Service) IsAuthenticated() bool {
	return s.xsrfToken != ""
}

func (s *Service) Init(ctx context.Context) error {
	if s.xsrfToken != "" {
		return nil
	}
	appID := imxapi.ImxLoginPostParamsAppId(s.cfg.Auth.IMX.AppId)
	body := imxapi.ImxLoginPostJSONRequestBody{
		"Module":   s.cfg.Auth.IMX.Module,
		"User":     s.cfg.Auth.IMX.User,
		"Password": s.cfg.Auth.IMX.Password,
	}
	nox := true
	resp, err := s.c.ImxLoginPostWithResponse(ctx, appID, &imxapi.ImxLoginPostParams{Noxsrf: &nox}, body)
	if err != nil {
		return fmt.Errorf("imx login request: %w", err)
	}
	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusMultipleChoices {
		return fmt.Errorf("imx login failed: status=%d body=%s", resp.StatusCode(), string(resp.Body))
	}
	if s.http.Jar == nil {
		return errors.New("http client has no cookie jar")
	}
	for _, c := range s.http.Jar.Cookies(s.baseURL) {
		if c.Name == "XSRF-TOKEN" && c.Value != "" {
			s.xsrfToken = c.Value
			break
		}
	}
	if s.xsrfToken == "" {
		return errors.New("imx login missing XSRF-TOKEN cookie")
	}
	return nil
}

func (s *Service) Editor() authcore.RequestEditor {
	return func(ctx context.Context, req *http.Request) error {
		sameHost := strings.EqualFold(req.URL.Host, s.baseURL.Host)
		isLogin := strings.Contains(req.URL.Path, "/imx/login/")
		if sameHost && !isLogin && s.xsrfToken == "" {
			return errors.New("imx: not authenticated; call Init first")
		}
		req.Header.Set("Accept", "application/json")
		if s.xsrfToken != "" {
			req.Header.Set("X-XSRF-TOKEN", s.xsrfToken)
		}
		return nil
	}
}

func (s *Service) ClearCache() {
	s.mu.Lock()
	s.xsrfToken = ""
	if s.http.Jar != nil {
		s.http.Jar.SetCookies(s.baseURL, nil)
	}
	s.mu.Unlock()
}
