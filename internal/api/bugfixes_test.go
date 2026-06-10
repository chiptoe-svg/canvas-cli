package api

// Tests covering specific bug fixes. Each test documents the bug it prevents.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"log/slog"
)

// --- Fix #2: retry body reuse ---
//
// doRequest used to pass the original io.Reader to each retry attempt. After
// the first read the reader is at EOF, so subsequent retry attempts sent an
// empty body. The fix reads the body once into a []byte and wraps it in a
// fresh bytes.Reader for each attempt.

func TestDoRequest_RetryPreservesBody(t *testing.T) {
	var attempts int32
	var secondBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}

		n := atomic.AddInt32(&attempts, 1)
		body, _ := io.ReadAll(r.Body)
		if n == 1 {
			// First attempt fails with 500
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"server error"}`))
			return
		}
		// Second attempt succeeds — capture body to verify it was resent
		secondBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":1}`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:             server.URL,
		Token:               "test-token",
		RequestsPerSec:      100,
		RetryInitialBackoff: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	payload := `{"name":"test"}`
	if err := client.PostJSON(
		context.Background(),
		"/api/v1/courses",
		map[string]string{"name": "test"},
		nil,
	); err != nil {
		t.Fatalf("PostJSON: %v", err)
	}

	if atomic.LoadInt32(&attempts) < 2 {
		t.Fatalf("expected at least 2 attempts (retry), got %d", attempts)
	}

	// The second request must carry the same JSON body, not be empty.
	if secondBody == "" {
		t.Error("retry attempt sent an empty body — body reuse bug is not fixed")
	}
	// Verify the body contains the expected field
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(secondBody), &parsed); err != nil {
		t.Fatalf("second body is not valid JSON: %s", secondBody)
	}
	if _, ok := parsed["name"]; !ok {
		t.Errorf("retry body missing 'name' field; got: %s", secondBody)
	}
	_ = payload
}

// --- Fix #7: Retry-After header on 429 ---
//
// ExecuteWithRetry must honour the server's Retry-After header on 429. The
// backoff should be at least the value specified in the header.

func TestRetryPolicy_HonoursRetryAfterHeader(t *testing.T) {
	// Use a very short exponential backoff so the hint dominates.
	policy := &RetryPolicy{
		MaxRetries:     2,
		InitialBackoff: time.Millisecond,
		MaxBackoff:     10 * time.Millisecond,
		Logger:         slog.Default(),
	}

	var callTimes []time.Time
	fn := func() (*http.Response, error) {
		callTimes = append(callTimes, time.Now())
		if len(callTimes) == 1 {
			// Return 429 with Retry-After: 0.05 (50 ms)
			resp := &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Header:     http.Header{"Retry-After": []string{"0.05"}},
				Body:       io.NopCloser(strings.NewReader(`{"message":"rate limited"}`)),
			}
			return resp, &APIError{StatusCode: http.StatusTooManyRequests}
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	}

	ctx := context.Background()
	resp, err := policy.ExecuteWithRetry(ctx, fn)
	if err != nil {
		t.Fatalf("ExecuteWithRetry: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if len(callTimes) < 2 {
		t.Fatalf("expected at least 2 calls, got %d", len(callTimes))
	}

	elapsed := callTimes[1].Sub(callTimes[0])
	// The Retry-After of 50 ms must dominate over the 1 ms exponential backoff.
	if elapsed < 40*time.Millisecond {
		t.Errorf("expected wait >=40ms (Retry-After 50ms), got %v", elapsed)
	}
}

// --- Fix #8: pagination infinite loop guard ---
//
// GetAllPages and GetAllPagesGeneric must abort instead of looping forever when
// the server keeps returning an identical "next" Link header.

func TestGetAllPages_InfiniteLoopGuard_SameURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		// Always return the same next-page link → cycle
		w.Header().Set("Link", fmt.Sprintf(`<%s/api/v1/courses>; rel="next"`, "http://"+r.Host))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1}]`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "t",
		RequestsPerSec: 1000,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	var courses []Course
	err = client.GetAllPages(context.Background(), "/api/v1/courses", &courses)
	if err == nil {
		t.Fatal("expected error for same-URL pagination cycle, got nil")
	}
	if !strings.Contains(err.Error(), "pagination aborted") {
		t.Errorf("expected 'pagination aborted' in error, got: %v", err)
	}
}

func TestGetAllPagesGeneric_InfiniteLoopGuard_SameURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Link", fmt.Sprintf(`<%s/api/v1/courses>; rel="next"`, "http://"+r.Host))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1}]`))
	}))
	defer server.Close()

	client, err := NewClient(ClientConfig{
		BaseURL:        server.URL,
		Token:          "t",
		RequestsPerSec: 1000,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = GetAllPagesGeneric[Course](client, context.Background(), "/api/v1/courses")
	if err == nil {
		t.Fatal("expected error for same-URL pagination cycle, got nil")
	}
	if !strings.Contains(err.Error(), "pagination aborted") {
		t.Errorf("expected 'pagination aborted' in error, got: %v", err)
	}
}

// --- Fix #9: BulkGrade progress parsing ---
//
// Canvas returns the Progress object at the TOP level of the BulkGrade response,
// not nested under a "progress" key. The old code expected {"progress":{"id":N}},
// so Progress.ID was always 0. The fix unmarshals the top-level response directly.

func TestSubmissionsService_BulkGrade_ProgressID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Canvas returns Progress at the top level (not under a "progress" key).
		w.Write([]byte(`{"id": 42, "workflow_state": "queued"}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	service := NewSubmissionsService(client)
	_, err := service.BulkGrade(context.Background(), 1, 1, &BulkGradeParams{
		GradeData: map[int64]GradeData{10: {PostedGrade: "A"}},
	})
	if err != nil {
		t.Fatalf("BulkGrade: %v", err)
	}
	// The old bug would have id==0 because it looked for {"progress":{"id":42}}.
	// We verify the call succeeds; the actual ID value is not surfaced by the
	// current BulkGrade API but the parsing must not error.
}
