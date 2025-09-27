package testx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type AuthServerCookie struct {
	TS      *httptest.Server
	BaseURL string
}

type CookieAuthOpts struct {
	LoginPath  string // default "/api/login"
	CookieName string // default "sessionid"
	SetCookie  bool   // default true
	StatusCode int    // default 200
	ExpectUser struct {
		Name     string // default "demo"
		Password string // default "demo"
	}
}

func StartAuthServerCookie(t testing.TB, opts CookieAuthOpts) *AuthServerCookie {
	t.Helper()
	if opts.LoginPath == "" {
		opts.LoginPath = "/api/login"
	}
	if opts.CookieName == "" {
		opts.CookieName = "sessionid"
	}
	if opts.StatusCode == 0 {
		opts.StatusCode = http.StatusOK
	}
	if opts.ExpectUser.Name == "" {
		opts.ExpectUser.Name = "demo"
	}
	if opts.ExpectUser.Password == "" {
		opts.ExpectUser.Password = "demo"
	}

	mux := http.NewServeMux()
	mux.HandleFunc(opts.LoginPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query()
		user, pass := q.Get("username"), q.Get("password")
		if user != opts.ExpectUser.Name || pass != opts.ExpectUser.Password {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": false})
			return
		}

		if opts.SetCookie {
			http.SetCookie(w, &http.Cookie{
				Name:     opts.CookieName,
				Value:    "ok",
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Secure:   true,
				Expires:  time.Now().Add(1 * time.Hour),
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(opts.StatusCode)
		json.NewEncoder(w).Encode(map[string]any{"ok": opts.StatusCode == http.StatusOK})
	})

	t.Log("trying to start AuthServerCookie...")
	ts := httptest.NewTLSServer(mux)
	t.Log("AuthServerCookie started")
	return &AuthServerCookie{TS: ts, BaseURL: ts.URL}
}

func (s *AuthServerCookie) Close() { s.TS.Close() }
