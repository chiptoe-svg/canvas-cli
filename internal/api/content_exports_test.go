package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentExportsService_List(t *testing.T) {
	exports := []ContentExport{
		{ID: 1, ExportType: "common_cartridge", WorkflowState: "created"},
		{ID: 2, ExportType: "zip", WorkflowState: "exported"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/10/content_exports" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(exports); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentExportsService(client)
	got, err := svc.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 exports, got %d", len(got))
	}
}

func TestContentExportsService_Get(t *testing.T) {
	want := ContentExport{ID: 5, ExportType: "qti", WorkflowState: "exported"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/10/content_exports/5" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentExportsService(client)
	got, err := svc.Get(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != 5 {
		t.Errorf("want ID 5, got %d", got.ID)
	}
}

func TestContentExportsService_Create(t *testing.T) {
	want := ContentExport{ID: 7, ExportType: "common_cartridge", WorkflowState: "created"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/10/content_exports" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentExportsService(client)
	got, err := svc.Create(context.Background(), 10, CreateContentExportParams{ExportType: "common_cartridge"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.ID != 7 {
		t.Errorf("want ID 7, got %d", got.ID)
	}
}

func TestContentExportsService_CreateEpub(t *testing.T) {
	want := EpubExport{ID: 3, WorkflowState: "created"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/courses/10/epub_exports" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentExportsService(client)
	got, err := svc.CreateEpub(context.Background(), 10)
	if err != nil {
		t.Fatalf("CreateEpub: %v", err)
	}
	if got.ID != 3 {
		t.Errorf("want ID 3, got %d", got.ID)
	}
}

func TestContentExportsService_GetEpub(t *testing.T) {
	want := EpubExport{ID: 3, WorkflowState: "exported"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/courses/10/epub_exports/3" {
			http.NotFound(w, r)
			return
		}
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewContentExportsService(client)
	got, err := svc.GetEpub(context.Background(), 10, 3)
	if err != nil {
		t.Fatalf("GetEpub: %v", err)
	}
	if got.WorkflowState != "exported" {
		t.Errorf("want exported, got %s", got.WorkflowState)
	}
}
