package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCourseNicknamesService_List(t *testing.T) {
	want := []CourseNickname{{CourseID: 10, Name: "Math 101", Nickname: "Math"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/course_nicknames" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseNicknamesService(client)
	got, err := svc.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d nicknames, want %d", len(got), len(want))
	}
	if got[0].CourseID != want[0].CourseID {
		t.Errorf("got CourseID %d, want %d", got[0].CourseID, want[0].CourseID)
	}
	if got[0].Nickname != want[0].Nickname {
		t.Errorf("got Nickname %q, want %q", got[0].Nickname, want[0].Nickname)
	}
}

func TestCourseNicknamesService_Get(t *testing.T) {
	want := &CourseNickname{CourseID: 42, Name: "Biology", Nickname: "Bio"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/course_nicknames/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseNicknamesService(client)
	got, err := svc.Get(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if got.CourseID != want.CourseID {
		t.Errorf("got CourseID %d, want %d", got.CourseID, want.CourseID)
	}
	if got.Nickname != want.Nickname {
		t.Errorf("got Nickname %q, want %q", got.Nickname, want.Nickname)
	}
}

func TestCourseNicknamesService_Set(t *testing.T) {
	want := &CourseNickname{CourseID: 99, Name: "Chemistry", Nickname: "Chem"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/course_nicknames/99" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseNicknamesService(client)
	params := SetCourseNicknameParams{Nickname: "Chem"}
	got, err := svc.Set(context.Background(), 99, params)
	if err != nil {
		t.Fatal(err)
	}
	if got.CourseID != want.CourseID {
		t.Errorf("got CourseID %d, want %d", got.CourseID, want.CourseID)
	}
	if got.Nickname != want.Nickname {
		t.Errorf("got Nickname %q, want %q", got.Nickname, want.Nickname)
	}
}

func TestCourseNicknamesService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/course_nicknames/55" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseNicknamesService(client)
	err := svc.Delete(context.Background(), 55)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCourseNicknamesService_DeleteAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/self/course_nicknames" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}")) // #nosec G104
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewCourseNicknamesService(client)
	err := svc.DeleteAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCourseNicknamesService(t *testing.T) {
	client := &Client{}
	svc := NewCourseNicknamesService(client)
	if svc == nil {
		t.Fatal("NewCourseNicknamesService returned nil")
	}
	if svc.client != client {
		t.Error("NewCourseNicknamesService did not set client correctly")
	}
}
