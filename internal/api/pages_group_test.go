package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// handleVersionDetection is already defined in testhelpers_test.go / client_test.go;
// do not redeclare here.

func TestPagesService_GroupList(t *testing.T) {
	tests := []struct {
		name    string
		groupID int64
		wantLen int
		wantErr bool
	}{
		{name: "valid group", groupID: 42, wantLen: 2, wantErr: false},
		{name: "zero group id", groupID: 0, wantLen: 0, wantErr: true},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/pages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Page{
			{PageID: 1, Title: "Intro", URL: "intro"},
			{PageID: 2, Title: "Rules", URL: "rules"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages, err := svc.GroupList(ctx, tt.groupID, nil)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("GroupList: %v", err)
			}
			if len(pages) != tt.wantLen {
				t.Errorf("expected %d pages, got %d", tt.wantLen, len(pages))
			}
		})
	}
}

func TestPagesService_GroupGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/pages/intro" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page{PageID: 1, Title: "Intro", URL: "intro"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)
	ctx := context.Background()

	page, err := svc.GroupGet(ctx, 42, "intro")
	if err != nil {
		t.Fatalf("GroupGet: %v", err)
	}
	if page.Title != "Intro" {
		t.Errorf("expected title 'Intro', got %q", page.Title)
	}
}

func TestPagesService_GroupGet_InvalidArgs(t *testing.T) {
	svc := NewPagesService(&Client{})
	ctx := context.Background()

	if _, err := svc.GroupGet(ctx, 0, "intro"); err == nil {
		t.Error("expected error for zero group_id")
	}
	if _, err := svc.GroupGet(ctx, 1, ""); err == nil {
		t.Error("expected error for empty url_or_id")
	}
}

func TestPagesService_GroupGetFrontPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/front_page" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page{PageID: 3, Title: "Home", FrontPage: true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)

	page, err := svc.GroupGetFrontPage(context.Background(), 42)
	if err != nil {
		t.Fatalf("GroupGetFrontPage: %v", err)
	}
	if !page.FrontPage {
		t.Error("expected FrontPage=true")
	}
}

func TestPagesService_GroupCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/42/pages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page{PageID: 10, Title: "New Page", URL: "new-page"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)

	page, err := svc.GroupCreate(context.Background(), 42, &CreatePageParams{Title: "New Page"})
	if err != nil {
		t.Fatalf("GroupCreate: %v", err)
	}
	if page.Title != "New Page" {
		t.Errorf("expected title 'New Page', got %q", page.Title)
	}
}

func TestPagesService_GroupCreate_NilParams(t *testing.T) {
	svc := NewPagesService(&Client{})
	if _, err := svc.GroupCreate(context.Background(), 1, nil); err == nil {
		t.Error("expected error for nil params")
	}
}

func TestPagesService_GroupUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/42/pages/intro" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page{PageID: 1, Title: "Updated", URL: "intro"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)

	title := "Updated"
	page, err := svc.GroupUpdate(context.Background(), 42, "intro", &UpdatePageParams{Title: &title})
	if err != nil {
		t.Fatalf("GroupUpdate: %v", err)
	}
	if page.Title != "Updated" {
		t.Errorf("expected title 'Updated', got %q", page.Title)
	}
}

func TestPagesService_GroupUpdateFrontPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/42/front_page" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Page{PageID: 1, Title: "Home", FrontPage: true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)

	title := "Home"
	page, err := svc.GroupUpdateFrontPage(context.Background(), 42, &UpdatePageParams{Title: &title})
	if err != nil {
		t.Fatalf("GroupUpdateFrontPage: %v", err)
	}
	if page.Title != "Home" {
		t.Errorf("expected title 'Home', got %q", page.Title)
	}
}

func TestPagesService_GroupDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/42/pages/intro" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)

	if err := svc.GroupDelete(context.Background(), 42, "intro"); err != nil {
		t.Fatalf("GroupDelete: %v", err)
	}
}

func TestPagesService_GroupDelete_InvalidArgs(t *testing.T) {
	svc := NewPagesService(&Client{})
	ctx := context.Background()

	if err := svc.GroupDelete(ctx, 0, "intro"); err == nil {
		t.Error("expected error for zero group_id")
	}
	if err := svc.GroupDelete(ctx, 1, ""); err == nil {
		t.Error("expected error for empty url_or_id")
	}
}

func TestPagesService_GroupListRevisions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/pages/intro/revisions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]PageRevision{{RevisionID: 1}, {RevisionID: 2}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)

	revs, err := svc.GroupListRevisions(context.Background(), 42, "intro")
	if err != nil {
		t.Fatalf("GroupListRevisions: %v", err)
	}
	if len(revs) != 2 {
		t.Errorf("expected 2 revisions, got %d", len(revs))
	}
}

func TestPagesService_GroupGetRevision(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/pages/intro/revisions/3" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PageRevision{RevisionID: 3})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)

	rev, err := svc.GroupGetRevision(context.Background(), 42, "intro", 3, false)
	if err != nil {
		t.Fatalf("GroupGetRevision: %v", err)
	}
	if rev.RevisionID != 3 {
		t.Errorf("expected RevisionID 3, got %d", rev.RevisionID)
	}
}

func TestPagesService_GroupGetLatestRevision(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/pages/intro/revisions/latest" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PageRevision{RevisionID: 5, Latest: true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)

	rev, err := svc.GroupGetLatestRevision(context.Background(), 42, "intro", false)
	if err != nil {
		t.Fatalf("GroupGetLatestRevision: %v", err)
	}
	if !rev.Latest {
		t.Error("expected Latest=true")
	}
}

func TestPagesService_GroupRevertToRevision(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/groups/42/pages/intro/revisions/2" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PageRevision{RevisionID: 2})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewPagesService(client)

	rev, err := svc.GroupRevertToRevision(context.Background(), 42, "intro", 2)
	if err != nil {
		t.Fatalf("GroupRevertToRevision: %v", err)
	}
	if rev.RevisionID != 2 {
		t.Errorf("expected RevisionID 2, got %d", rev.RevisionID)
	}
}

func TestPagesContextPath(t *testing.T) {
	tests := []struct {
		courseID int64
		groupID  int64
		want     string
		wantErr  bool
	}{
		{courseID: 1, groupID: 0, want: "courses/1", wantErr: false},
		{courseID: 0, groupID: 2, want: "groups/2", wantErr: false},
		{courseID: 0, groupID: 0, want: "", wantErr: true},
		{courseID: 1, groupID: 2, want: "", wantErr: true},
	}

	for _, tt := range tests {
		seg, err := pagesContextPath(tt.courseID, tt.groupID)
		if tt.wantErr {
			if err == nil {
				t.Errorf("pagesContextPath(%d,%d): expected error", tt.courseID, tt.groupID)
			}
			continue
		}
		if err != nil {
			t.Errorf("pagesContextPath(%d,%d): unexpected error: %v", tt.courseID, tt.groupID, err)
			continue
		}
		if seg != tt.want {
			t.Errorf("pagesContextPath(%d,%d) = %q, want %q", tt.courseID, tt.groupID, seg, tt.want)
		}
	}
}
