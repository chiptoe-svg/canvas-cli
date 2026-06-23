package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAccountsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/10" {
			t.Errorf("expected /api/v1/accounts/10, got %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		acct, ok := body["account"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected 'account' key in body")
		}
		if acct["name"] != "Updated Account" {
			t.Errorf("expected name 'Updated Account', got %v", acct["name"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Account{ID: 10, Name: "Updated Account"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountsService(client)

	params := &UpdateAccountParams{Name: "Updated Account"}
	account, err := svc.Update(context.Background(), 10, params)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if account.ID != 10 {
		t.Errorf("expected ID 10, got %d", account.ID)
	}
	if account.Name != "Updated Account" {
		t.Errorf("expected name 'Updated Account', got %s", account.Name)
	}
}

func TestAccountsService_Update_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"errors":[{"message":"forbidden"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountsService(client)
	_, err := svc.Update(context.Background(), 10, &UpdateAccountParams{Name: "X"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAccountsService_CreateSubAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/5/sub_accounts" {
			t.Errorf("expected /api/v1/accounts/5/sub_accounts, got %s", r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		acct, ok := body["account"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected 'account' key")
		}
		if acct["name"] != "Sub Account" {
			t.Errorf("expected name 'Sub Account', got %v", acct["name"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Account{ID: 20, Name: "Sub Account"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountsService(client)
	params := &CreateSubAccountParams{Name: "Sub Account"}
	account, err := svc.CreateSubAccount(context.Background(), 5, params)
	if err != nil {
		t.Fatalf("CreateSubAccount: %v", err)
	}
	if account.ID != 20 {
		t.Errorf("expected ID 20, got %d", account.ID)
	}
}

func TestAccountsService_CreateSubAccount_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"errors":[{"message":"invalid"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountsService(client)
	_, err := svc.CreateSubAccount(context.Background(), 5, &CreateSubAccountParams{Name: "X"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsService_DeleteUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/5/users/99" {
			t.Errorf("expected /api/v1/accounts/5/users/99, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{ID: 99, Name: "Removed User"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountsService(client)
	user, err := svc.DeleteUser(context.Background(), 5, 99)
	if err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	if user.ID != 99 {
		t.Errorf("expected ID 99, got %d", user.ID)
	}
}

func TestAccountsService_DeleteUser_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountsService(client)
	_, err := svc.DeleteUser(context.Background(), 5, 99)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAccountsService_UpdateCourses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/accounts/3/courses" {
			t.Errorf("expected /api/v1/accounts/3/courses, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"completed": true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountsService(client)
	result, err := svc.UpdateCourses(context.Background(), 3, map[string]interface{}{"event": "offer"})
	if err != nil {
		t.Fatalf("UpdateCourses: %v", err)
	}
	if result["completed"] != true {
		t.Errorf("expected completed=true, got %v", result["completed"])
	}
}

func TestAccountsService_GetAdminSelf(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/admins/self") {
			t.Errorf("expected path ending in /admins/self, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": float64(1), "role": "AccountAdmin"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAccountsService(client)
	result, err := svc.GetAdminSelf(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetAdminSelf: %v", err)
	}
	if result["role"] != "AccountAdmin" {
		t.Errorf("expected role=AccountAdmin, got %v", result["role"])
	}
}

// TestAccountsService_Update_WithStorageQuota asserts that DefaultStorageQuotaMb
// is sent correctly when set (including value 0), and omitted when nil.
// The fix changed the field from int (omitempty silently drops 0) to *int.
func TestAccountsService_Update_WithStorageQuota(t *testing.T) {
	var receivedQuota interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		acct, _ := body["account"].(map[string]interface{})
		receivedQuota = acct["default_storage_quota_mb"]
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Account{ID: 1})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewAccountsService(client)

	// With a non-zero quota: field must be present.
	quota := 500
	_, err := svc.Update(context.Background(), 1, &UpdateAccountParams{DefaultStorageQuotaMb: &quota})
	if err != nil {
		t.Fatalf("Update with quota: %v", err)
	}
	if receivedQuota == nil {
		t.Error("expected default_storage_quota_mb in body when set to non-zero")
	}

	// With nil quota: field must be absent (omitempty on *int drops nil).
	_, err = svc.Update(context.Background(), 1, &UpdateAccountParams{Name: "X"})
	if err != nil {
		t.Fatalf("Update without quota: %v", err)
	}
	if receivedQuota != nil {
		t.Errorf("expected default_storage_quota_mb absent when nil, got %v", receivedQuota)
	}
}

// TestAccountsService_CreateSubAccount_WithStorageQuota mirrors the Update test
// for CreateSubAccountParams.DefaultStorageQuotaMb.
func TestAccountsService_CreateSubAccount_WithStorageQuota(t *testing.T) {
	var receivedQuota interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		acct, _ := body["account"].(map[string]interface{})
		receivedQuota = acct["default_storage_quota_mb"]
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Account{ID: 2})
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)
	svc := NewAccountsService(client)

	quota := 1000
	_, err := svc.CreateSubAccount(context.Background(), 1, &CreateSubAccountParams{
		Name:                  "Child",
		DefaultStorageQuotaMb: &quota,
	})
	if err != nil {
		t.Fatalf("CreateSubAccount with quota: %v", err)
	}
	if receivedQuota == nil {
		t.Error("expected default_storage_quota_mb in body when set")
	}
}
