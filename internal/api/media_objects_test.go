package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMediaObjectsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/media_objects" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		objs := []MediaObject{
			{MediaID: "m-abc123", Title: "Intro Video", MediaType: "video"},
		}
		if err := json.NewEncoder(w).Encode(objs); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewMediaObjectsService(client)

	objs, err := svc.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(objs) != 1 || objs[0].MediaID != "m-abc123" {
		t.Errorf("unexpected result: %+v", objs)
	}
}

func TestMediaObjectsService_ListForCourse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/media_objects" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		objs := []MediaObject{{MediaID: "m-xyz"}}
		if err := json.NewEncoder(w).Encode(objs); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewMediaObjectsService(client)

	objs, err := svc.ListForCourse(context.Background(), 10, nil)
	if err != nil {
		t.Fatalf("ListForCourse failed: %v", err)
	}
	if len(objs) != 1 {
		t.Errorf("expected 1 object, got %d", len(objs))
	}
}

func TestMediaObjectsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/media_objects/m-abc123" || r.Method != http.MethodPut {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		obj := MediaObject{MediaID: "m-abc123", UserEnteredTitle: "New Title"}
		if err := json.NewEncoder(w).Encode(obj); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewMediaObjectsService(client)

	obj, err := svc.Update(context.Background(), "m-abc123", "New Title")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if obj.UserEnteredTitle != "New Title" {
		t.Errorf("unexpected title: %s", obj.UserEnteredTitle)
	}
}

func TestMediaObjectsService_GetMediaTracks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/media_objects/m-abc123/media_tracks" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		tracks := []MediaTrack{{Locale: "en", Kind: "subtitles"}}
		if err := json.NewEncoder(w).Encode(tracks); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewMediaObjectsService(client)

	tracks, err := svc.GetMediaTracks(context.Background(), "m-abc123")
	if err != nil {
		t.Fatalf("GetMediaTracks failed: %v", err)
	}
	if len(tracks) != 1 || tracks[0].Locale != "en" {
		t.Errorf("unexpected tracks: %+v", tracks)
	}
}

func TestMediaObjectsService_ListAttachments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/media_attachments" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		atts := []MediaAttachment{{ID: 5, Filename: "video.mp4"}}
		if err := json.NewEncoder(w).Encode(atts); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewMediaObjectsService(client)

	atts, err := svc.ListAttachments(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListAttachments failed: %v", err)
	}
	if len(atts) != 1 || atts[0].ID != 5 {
		t.Errorf("unexpected result: %+v", atts)
	}
}

func TestMediaObjectsService_UpdateAttachment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/media_attachments/5" || r.Method != http.MethodPut {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		att := MediaAttachment{ID: 5, DisplayName: "Updated"}
		if err := json.NewEncoder(w).Encode(att); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewMediaObjectsService(client)

	att, err := svc.UpdateAttachment(context.Background(), 5, "Updated")
	if err != nil {
		t.Fatalf("UpdateAttachment failed: %v", err)
	}
	if att.DisplayName != "Updated" {
		t.Errorf("unexpected display name: %s", att.DisplayName)
	}
}

func TestMediaObjectsService_GetAttachmentMediaTracks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/media_attachments/5/media_tracks" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		tracks := []MediaTrack{{Locale: "fr", Kind: "captions"}}
		if err := json.NewEncoder(w).Encode(tracks); err != nil {
			t.Errorf("encode error: %v", err)
		}
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewMediaObjectsService(client)

	tracks, err := svc.GetAttachmentMediaTracks(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetAttachmentMediaTracks failed: %v", err)
	}
	if len(tracks) != 1 || tracks[0].Locale != "fr" {
		t.Errorf("unexpected tracks: %+v", tracks)
	}
}
