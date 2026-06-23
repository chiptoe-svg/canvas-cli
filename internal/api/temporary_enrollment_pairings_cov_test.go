package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountTemporaryEnrollmentPairingsService_New(t *testing.T) {
	client := newTestClient(t, "http://localhost")
	svc := NewAccountTemporaryEnrollmentPairingsService(client)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestAccountTemporaryEnrollmentPairingsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/temporary_enrollment_pairings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]TemporaryEnrollmentPairing{
			{ID: 10, WorkflowState: "active"},
			{ID: 11, WorkflowState: "active"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountTemporaryEnrollmentPairingsService(client)
	pairings, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(pairings) != 2 {
		t.Fatalf("expected 2 pairings, got %d", len(pairings))
	}
	if pairings[0].ID != 10 {
		t.Errorf("expected ID 10, got %d", pairings[0].ID)
	}
}

func TestAccountTemporaryEnrollmentPairingsService_List_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountTemporaryEnrollmentPairingsService(client)
	_, err := svc.List(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountTemporaryEnrollmentPairingsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/temporary_enrollment_pairings/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TemporaryEnrollmentPairing{ID: 10, WorkflowState: "active", RootAccountID: 1})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountTemporaryEnrollmentPairingsService(client)
	pairing, err := svc.Get(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if pairing.ID != 10 {
		t.Errorf("expected ID 10, got %d", pairing.ID)
	}
	if pairing.WorkflowState != "active" {
		t.Errorf("expected workflow_state=active, got %s", pairing.WorkflowState)
	}
}

func TestAccountTemporaryEnrollmentPairingsService_Get_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountTemporaryEnrollmentPairingsService(client)
	_, err := svc.Get(context.Background(), 1, 10)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountTemporaryEnrollmentPairingsService_GetNew(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/temporary_enrollment_pairings/new" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TemporaryEnrollmentPairing{ID: 0, WorkflowState: "pending"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountTemporaryEnrollmentPairingsService(client)
	pairing, err := svc.GetNew(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetNew: %v", err)
	}
	if pairing.WorkflowState != "pending" {
		t.Errorf("expected pending, got %s", pairing.WorkflowState)
	}
}

func TestAccountTemporaryEnrollmentPairingsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/temporary_enrollment_pairings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body TemporaryEnrollmentPairingParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.StartingEnrollmentState != "active" {
			t.Errorf("expected starting_enrollment_state=active, got %s", body.StartingEnrollmentState)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TemporaryEnrollmentPairing{ID: 15, WorkflowState: "active", StartingEnrollmentState: "active"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountTemporaryEnrollmentPairingsService(client)
	params := &TemporaryEnrollmentPairingParams{StartingEnrollmentState: "active"}
	pairing, err := svc.Create(context.Background(), 1, params)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if pairing.ID != 15 {
		t.Errorf("expected ID 15, got %d", pairing.ID)
	}
}

func TestAccountTemporaryEnrollmentPairingsService_Delete(t *testing.T) {
	deleted := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/1/temporary_enrollment_pairings/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		deleted = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountTemporaryEnrollmentPairingsService(client)
	err := svc.Delete(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !deleted {
		t.Error("expected DELETE request to be made")
	}
}

func TestAccountTemporaryEnrollmentPairingsService_Delete_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountTemporaryEnrollmentPairingsService(client)
	err := svc.Delete(context.Background(), 1, 10)
	if err == nil {
		t.Fatal("expected error")
	}
}
