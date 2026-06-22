package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewFoldersService(t *testing.T) {
	client := &Client{}
	svc := NewFoldersService(client)
	if svc == nil {
		t.Fatal("NewFoldersService returned nil")
	}
	if svc.client != client {
		t.Error("NewFoldersService did not set client correctly")
	}
}

func TestFoldersContextPath(t *testing.T) {
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
		{name: "all three", courseID: 1, groupID: 2, userID: 3, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seg, err := foldersContextPath(tt.courseID, tt.groupID, tt.userID)
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

func TestFoldersService_ListContextFolders_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/folders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Folder{
			{ID: 1, Name: "Root", ContextType: "Course"},
			{ID: 2, Name: "Assignments", ContextType: "Course"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folders, err := svc.ListContextFolders(context.Background(), 10, 0, 0, nil)
	if err != nil {
		t.Fatalf("ListContextFolders: %v", err)
	}
	if len(folders) != 2 {
		t.Errorf("expected 2 folders, got %d", len(folders))
	}
}

func TestFoldersService_ListContextFolders_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/folders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Folder{{ID: 3, Name: "Group Root"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folders, err := svc.ListContextFolders(context.Background(), 0, 5, 0, nil)
	if err != nil {
		t.Fatalf("ListContextFolders: %v", err)
	}
	if len(folders) != 1 {
		t.Errorf("expected 1 folder, got %d", len(folders))
	}
}

func TestFoldersService_ListContextFolders_User(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/users/99/folders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Folder{{ID: 4, Name: "My Files"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folders, err := svc.ListContextFolders(context.Background(), 0, 0, 99, nil)
	if err != nil {
		t.Fatalf("ListContextFolders: %v", err)
	}
	if len(folders) != 1 {
		t.Errorf("expected 1 folder, got %d", len(folders))
	}
}

func TestFoldersService_ListContextFolders_InvalidContext(t *testing.T) {
	svc := NewFoldersService(&Client{})
	_, err := svc.ListContextFolders(context.Background(), 0, 0, 0, nil)
	if err == nil {
		t.Error("expected error for no context")
	}
}

func TestFoldersService_ListFolderSubFolders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/folders/7/folders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Folder{{ID: 8, Name: "Sub"}, {ID: 9, Name: "Sub2"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folders, err := svc.ListFolderSubFolders(context.Background(), 7, nil)
	if err != nil {
		t.Fatalf("ListFolderSubFolders: %v", err)
	}
	if len(folders) != 2 {
		t.Errorf("expected 2 sub-folders, got %d", len(folders))
	}
}

func TestFoldersService_ListFolderSubFolders_InvalidID(t *testing.T) {
	svc := NewFoldersService(&Client{})
	_, err := svc.ListFolderSubFolders(context.Background(), 0, nil)
	if err == nil {
		t.Error("expected error for zero folder_id")
	}
}

func TestFoldersService_GetFolder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/folders/55" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Folder{ID: 55, Name: "Lectures"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folder, err := svc.GetFolder(context.Background(), 55)
	if err != nil {
		t.Fatalf("GetFolder: %v", err)
	}
	if folder.ID != 55 {
		t.Errorf("expected ID 55, got %d", folder.ID)
	}
	if folder.Name != "Lectures" {
		t.Errorf("expected name 'Lectures', got %q", folder.Name)
	}
}

func TestFoldersService_GetFolder_InvalidID(t *testing.T) {
	svc := NewFoldersService(&Client{})
	_, err := svc.GetFolder(context.Background(), 0)
	if err == nil {
		t.Error("expected error for zero folder_id")
	}
}

func TestFoldersService_ResolvePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/folders/by_path/lectures/week1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Folder{
			{ID: 1, Name: "Root"},
			{ID: 2, Name: "lectures"},
			{ID: 3, Name: "week1"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folders, err := svc.ResolvePath(context.Background(), 10, 0, 0, "lectures/week1")
	if err != nil {
		t.Fatalf("ResolvePath: %v", err)
	}
	if len(folders) != 3 {
		t.Errorf("expected 3 folders in path, got %d", len(folders))
	}
}

func TestFoldersService_ResolvePath_NoPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/folders/by_path" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Folder{{ID: 1, Name: "Root"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folders, err := svc.ResolvePath(context.Background(), 0, 5, 0, "")
	if err != nil {
		t.Fatalf("ResolvePath empty path: %v", err)
	}
	if len(folders) != 1 {
		t.Errorf("expected 1 folder, got %d", len(folders))
	}
}

func TestFoldersService_CreateContextFolder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/courses/10/folders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Folder{ID: 20, Name: "New Folder"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folder, err := svc.CreateContextFolder(context.Background(), 10, 0, 0, &CreateFolderParams{Name: "New Folder"})
	if err != nil {
		t.Fatalf("CreateContextFolder: %v", err)
	}
	if folder.Name != "New Folder" {
		t.Errorf("expected name 'New Folder', got %q", folder.Name)
	}
}

func TestFoldersService_CreateContextFolder_NilParams(t *testing.T) {
	svc := NewFoldersService(&Client{})
	_, err := svc.CreateContextFolder(context.Background(), 1, 0, 0, nil)
	if err == nil {
		t.Error("expected error for nil params")
	}
}

func TestFoldersService_CreateSubFolder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/folders/7/folders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Folder{ID: 21, Name: "Sub"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folder, err := svc.CreateSubFolder(context.Background(), 7, &CreateFolderParams{Name: "Sub"})
	if err != nil {
		t.Fatalf("CreateSubFolder: %v", err)
	}
	if folder.Name != "Sub" {
		t.Errorf("expected name 'Sub', got %q", folder.Name)
	}
}

func TestFoldersService_UpdateFolder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/folders/55" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Folder{ID: 55, Name: "Renamed"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	newName := "Renamed"
	folder, err := svc.UpdateFolder(context.Background(), 55, &UpdateFolderParams{Name: &newName})
	if err != nil {
		t.Fatalf("UpdateFolder: %v", err)
	}
	if folder.Name != "Renamed" {
		t.Errorf("expected name 'Renamed', got %q", folder.Name)
	}
}

func TestFoldersService_UpdateFolder_NilParams(t *testing.T) {
	svc := NewFoldersService(&Client{})
	_, err := svc.UpdateFolder(context.Background(), 1, nil)
	if err == nil {
		t.Error("expected error for nil params")
	}
}

func TestFoldersService_DeleteFolder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/folders/55" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	if err := svc.DeleteFolder(context.Background(), 55, false); err != nil {
		t.Fatalf("DeleteFolder: %v", err)
	}
}

func TestFoldersService_DeleteFolder_Force(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("force") != "true" {
			t.Error("expected force=true query parameter")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	if err := svc.DeleteFolder(context.Background(), 55, true); err != nil {
		t.Fatalf("DeleteFolder force: %v", err)
	}
}

func TestFoldersService_DeleteFolder_InvalidID(t *testing.T) {
	svc := NewFoldersService(&Client{})
	if err := svc.DeleteFolder(context.Background(), 0, false); err == nil {
		t.Error("expected error for zero folder_id")
	}
}

func TestFoldersService_GetMediaFolder_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/folders/media" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Folder{ID: 100, Name: "Uploaded Media"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folder, err := svc.GetMediaFolder(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("GetMediaFolder: %v", err)
	}
	if folder.Name != "Uploaded Media" {
		t.Errorf("expected name 'Uploaded Media', got %q", folder.Name)
	}
}

func TestFoldersService_GetMediaFolder_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/folders/media" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Folder{ID: 101, Name: "Group Media"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folder, err := svc.GetMediaFolder(context.Background(), 0, 5)
	if err != nil {
		t.Fatalf("GetMediaFolder group: %v", err)
	}
	if folder.ID != 101 {
		t.Errorf("expected ID 101, got %d", folder.ID)
	}
}

func TestFoldersService_GetMediaFolder_NoContext(t *testing.T) {
	svc := NewFoldersService(&Client{})
	_, err := svc.GetMediaFolder(context.Background(), 0, 0)
	if err == nil {
		t.Error("expected error for no context")
	}
}

func TestFoldersService_CopyFolder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/folders/20/copy_folder" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Folder{ID: 30, Name: "CopiedFolder"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewFoldersService(client)

	folder, err := svc.CopyFolder(context.Background(), 20, &CopyFolderParams{SourceFolderID: 10})
	if err != nil {
		t.Fatalf("CopyFolder: %v", err)
	}
	if folder.ID != 30 {
		t.Errorf("expected ID 30, got %d", folder.ID)
	}
}

func TestFoldersService_CopyFolder_NilParams(t *testing.T) {
	svc := NewFoldersService(&Client{})
	_, err := svc.CopyFolder(context.Background(), 1, nil)
	if err == nil {
		t.Error("expected error for nil params")
	}
}
