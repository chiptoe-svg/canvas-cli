package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFilesService_ListGroupFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/files" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Attachment{
			{ID: 1, DisplayName: "file1.pdf"},
			{ID: 2, DisplayName: "file2.txt"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFilesService(client)

	files, err := svc.ListGroupFiles(context.Background(), 5, nil)
	if err != nil {
		t.Fatalf("ListGroupFiles: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestFilesService_GetGroupQuota(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/files/quota" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(QuotaInfo{Quota: 1048576, QuotaUsed: 512000})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFilesService(client)

	quota, err := svc.GetGroupQuota(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetGroupQuota: %v", err)
	}
	if quota.Quota != 1048576 {
		t.Errorf("expected quota 1048576, got %d", quota.Quota)
	}
}

func TestFilesService_ResetVerifier(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/files/42/reset_verifier" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Attachment{ID: 42, DisplayName: "file.pdf"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFilesService(client)

	file, err := svc.ResetVerifier(context.Background(), 42)
	if err != nil {
		t.Fatalf("ResetVerifier: %v", err)
	}
	if file.ID != 42 {
		t.Errorf("expected ID 42, got %d", file.ID)
	}
}

func TestFilesService_ResetVerifier_InvalidID(t *testing.T) {
	svc := NewFilesService(&Client{})
	_, err := svc.ResetVerifier(context.Background(), 0)
	if err == nil {
		t.Error("expected error for zero file_id")
	}
}

func TestFilesService_CopyFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/folders/10/copy_file" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Attachment{ID: 99, DisplayName: "copy.pdf"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFilesService(client)

	file, err := svc.CopyFile(context.Background(), 10, &CopyFileParams{SourceFileID: 5})
	if err != nil {
		t.Fatalf("CopyFile: %v", err)
	}
	if file.ID != 99 {
		t.Errorf("expected ID 99, got %d", file.ID)
	}
}

func TestFilesService_CopyFile_NilParams(t *testing.T) {
	svc := NewFilesService(&Client{})
	_, err := svc.CopyFile(context.Background(), 1, nil)
	if err == nil {
		t.Error("expected error for nil params")
	}
}

func TestFilesService_CopyFile_InvalidSourceID(t *testing.T) {
	svc := NewFilesService(&Client{})
	_, err := svc.CopyFile(context.Background(), 1, &CopyFileParams{SourceFileID: 0})
	if err == nil {
		t.Error("expected error for zero source_file_id")
	}
}

func TestFilesUsageRightsContextPath(t *testing.T) {
	tests := []struct {
		name     string
		courseID int64
		groupID  int64
		userID   int64
		want     string
		wantErr  bool
	}{
		{name: "course", courseID: 1, want: "courses/1"},
		{name: "group", groupID: 2, want: "groups/2"},
		{name: "user", userID: 3, want: "users/3"},
		{name: "none", wantErr: true},
		{name: "course+group", courseID: 1, groupID: 2, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seg, err := filesUsageRightsContextPath(tt.courseID, tt.groupID, tt.userID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got %q", seg)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if seg != tt.want {
				t.Errorf("got %q, want %q", seg, tt.want)
			}
		})
	}
}

func TestFilesService_SetUsageRights(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/usage_rights" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UsageRights{UseJustification: "own_copyright"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFilesService(client)

	rights, err := svc.SetUsageRights(context.Background(), 10, 0, 0, &SetUsageRightsParams{
		FileIDs:          []int64{1, 2},
		UseJustification: "own_copyright",
	})
	if err != nil {
		t.Fatalf("SetUsageRights: %v", err)
	}
	if rights.UseJustification != "own_copyright" {
		t.Errorf("expected use_justification 'own_copyright', got %q", rights.UseJustification)
	}
}

func TestFilesService_SetUsageRights_NilParams(t *testing.T) {
	svc := NewFilesService(&Client{})
	_, err := svc.SetUsageRights(context.Background(), 1, 0, 0, nil)
	if err == nil {
		t.Error("expected error for nil params")
	}
}

func TestFilesService_SetUsageRights_NoFiles(t *testing.T) {
	svc := NewFilesService(&Client{})
	_, err := svc.SetUsageRights(context.Background(), 1, 0, 0, &SetUsageRightsParams{
		UseJustification: "own_copyright",
	})
	if err == nil {
		t.Error("expected error when no file_ids or folder_ids provided")
	}
}

func TestFilesService_RemoveUsageRights(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/5/usage_rights" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFilesService(client)

	err := svc.RemoveUsageRights(context.Background(), 0, 5, 0, &RemoveUsageRightsParams{
		FileIDs: []int64{10},
	})
	if err != nil {
		t.Fatalf("RemoveUsageRights: %v", err)
	}
}

func TestFilesService_RemoveUsageRights_NilParams(t *testing.T) {
	svc := NewFilesService(&Client{})
	if err := svc.RemoveUsageRights(context.Background(), 1, 0, 0, nil); err == nil {
		t.Error("expected error for nil params")
	}
}

func TestFilesService_ListLicenses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/99/content_licenses" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ContentLicense{
			{ID: "public_domain", Name: "Public Domain"},
			{ID: "cc_by", Name: "CC Attribution"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFilesService(client)

	licenses, err := svc.ListLicenses(context.Background(), 0, 0, 99)
	if err != nil {
		t.Fatalf("ListLicenses: %v", err)
	}
	if len(licenses) != 2 {
		t.Errorf("expected 2 licenses, got %d", len(licenses))
	}
}
