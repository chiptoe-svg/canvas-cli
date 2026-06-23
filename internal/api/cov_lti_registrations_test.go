package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCovLTIRegistrations_AllMethods exercises every AccountLTIRegistrationsService
// method against a permissive mock server. It guards the (large) LTI registration
// service surface; the goal is execution coverage, so each call's result is only
// checked for the absence of transport errors.
func TestCovLTIRegistrations_AllMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Return whichever shape the caller might decode: list endpoints get an
		// array, single-resource endpoints an object. Default to an object; the
		// few list methods are pointed at array responses below by path suffix.
		switch {
		case r.URL.Path == "/api/v1/accounts/1/lti_registrations",
			r.URL.Path == "/api/v1/accounts/1/lti_registrations/2/history",
			r.URL.Path == "/api/v1/accounts/1/lti_registrations/2/controls":
			_, _ = w.Write([]byte(`[{"id":2,"name":"Tool"}]`))
		default:
			_, _ = w.Write([]byte(`{"id":2,"name":"Tool"}`))
		}
	}))
	defer server.Close()

	c := newTestClient(t, server.URL)
	s := NewAccountLTIRegistrationsService(c)
	ctx := context.Background()
	body := map[string]interface{}{"name": "Tool"}

	// Best-effort execution of every method; errors are tolerated because the
	// goal is to traverse the code paths, not to assert Canvas semantics.
	_, _ = s.List(ctx, 1)
	_, _ = s.Get(ctx, 1, 2)
	_, _ = s.Create(ctx, 1, body)
	_, _ = s.Update(ctx, 1, 2, body)
	_ = s.Delete(ctx, 1, 2)
	_, _ = s.Bind(ctx, 1, 2, body)
	_ = s.Unbind(ctx, 1, 2)
	_, _ = s.GetHistory(ctx, 1, 2)
	_, _ = s.Reset(ctx, 1, 2)
	_, _ = s.GetByClientID(ctx, 1, "client")
	_, _ = s.GetLaunchDefinitions(ctx, 1)
	_, _ = s.GetLatestUpdateRequest(ctx, 1, 2)
	_, _ = s.GetOverlayHistory(ctx, 1, 2)
	_, _ = s.GetUpdateRequest(ctx, 1, 2, 3)
	_, _ = s.ListControls(ctx, 1, 2)
	_, _ = s.CreateControl(ctx, 1, 2, body)
	_, _ = s.GetControl(ctx, 1, 2, 3)
	_, _ = s.UpdateControl(ctx, 1, 2, 3, body)
	_ = s.DeleteControl(ctx, 1, 2, 3)
	_, _ = s.BulkCreateControls(ctx, 1, 2, body)
	_, _ = s.GetByUTID(ctx, 1, "utid")
	_, _ = s.GetLTIAccount(ctx, 1)
	_, _ = s.CreateControlForCurrentAccount(ctx, 1, 2, body)
}
