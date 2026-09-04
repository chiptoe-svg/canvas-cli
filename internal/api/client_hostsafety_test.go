package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A request path is appended to the configured base URL. Anything that can
// change the host after concatenation — "@evil.example/x" turns the base host
// into userinfo, "//evil.example/x" is a protocol-relative URL — would send
// the bearer token to a foreign host. Such paths must be refused before any
// request is made; this matters most for the raw `canvas api` command, which
// is exposed to AI agents over MCP.
func TestClient_RefusesPathsThatEscapeBaseHost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		t.Errorf("no request should reach the server for a rejected path, got %s %s", r.Method, r.URL)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	for _, path := range []string{
		"@evil.example/steal",
		"//evil.example/steal",
		"evil.example/steal",
		"https://evil.example/steal",
		"",
	} {
		t.Run(path, func(t *testing.T) {
			resp, err := client.doRequest(context.Background(), http.MethodGet, path, nil)
			if err == nil {
				resp.Body.Close()
				t.Fatalf("path %q must be rejected", path)
			}
			if !strings.Contains(err.Error(), "path") {
				t.Errorf("error should explain the path rule, got: %v", err)
			}
		})
	}
}

func TestClient_AcceptsOrdinaryPaths(t *testing.T) {
	var hit string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		hit = r.URL.RequestURI()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	resp, err := client.doRequest(context.Background(), http.MethodGet, "/api/v1/courses/1?include[]=term", nil)
	if err != nil {
		t.Fatalf("ordinary path rejected: %v", err)
	}
	resp.Body.Close()
	if hit != "/api/v1/courses/1?include[]=term" {
		t.Errorf("unexpected request URI %q", hit)
	}
}
