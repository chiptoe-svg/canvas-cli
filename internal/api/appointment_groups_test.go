package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppointmentGroupsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/appointment_groups" {
			t.Errorf("expected path /api/v1/appointment_groups, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"id":1,"title":"Office Hours","workflow_state":"active","appointments_count":3},
			{"id":2,"title":"Project Review","workflow_state":"active","appointments_count":2}
		]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAppointmentGroupsService(client)

	groups, err := svc.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Title != "Office Hours" {
		t.Errorf("expected 'Office Hours', got %s", groups[0].Title)
	}
}

func TestAppointmentGroupsService_ListWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Query().Get("scope") != "manageable" {
			t.Errorf("expected scope=manageable, got %s", r.URL.Query().Get("scope"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":1,"title":"Office Hours","workflow_state":"active"}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAppointmentGroupsService(client)

	groups, err := svc.List(context.Background(), &ListAppointmentGroupsOptions{
		Scope:   "manageable",
		Include: []string{"appointments", "participant_count"},
	})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
}

func TestAppointmentGroupsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/appointment_groups/1" {
			t.Errorf("expected path /api/v1/appointment_groups/1, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id":1,
			"title":"Office Hours",
			"description":"Weekly office hours",
			"workflow_state":"active",
			"appointments_count":3,
			"context_codes":["course_123"],
			"participant_type":"User"
		}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAppointmentGroupsService(client)

	group, err := svc.Get(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if group.ID != 1 {
		t.Errorf("expected group ID 1, got %d", group.ID)
	}
	if group.Description != "Weekly office hours" {
		t.Errorf("expected description 'Weekly office hours', got %s", group.Description)
	}
}

func TestAppointmentGroupsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/appointment_groups" {
			t.Errorf("expected path /api/v1/appointment_groups, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		ag, ok := body["appointment_group"].(map[string]interface{})
		if !ok {
			t.Fatal("expected appointment_group in body")
		}
		if ag["title"] != "Final Presentation" {
			t.Errorf("expected title 'Final Presentation', got %v", ag["title"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id":10,
			"title":"Final Presentation",
			"workflow_state":"pending",
			"context_codes":["course_123"],
			"appointments_count":0
		}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAppointmentGroupsService(client)

	group, err := svc.Create(context.Background(), &CreateAppointmentGroupParams{
		ContextCodes:                  []string{"course_123"},
		Title:                         "Final Presentation",
		ParticipantsPerAppointment:    1,
		MinAppointmentsPerParticipant: 1,
		MaxAppointmentsPerParticipant: 1,
		NewAppointments: [][2]string{
			{"2024-12-01T09:00:00Z", "2024-12-01T10:00:00Z"},
		},
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if group.Title != "Final Presentation" {
		t.Errorf("expected 'Final Presentation', got %s", group.Title)
	}
}

func TestAppointmentGroupsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/appointment_groups/1" {
			t.Errorf("expected path /api/v1/appointment_groups/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id":1,
			"title":"Updated Group",
			"workflow_state":"active",
			"context_codes":["course_123"]
		}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAppointmentGroupsService(client)

	group, err := svc.Update(context.Background(), 1, &UpdateAppointmentGroupParams{
		ContextCodes: []string{"course_123"},
		Title:        "Updated Group",
		Publish:      true,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if group.Title != "Updated Group" {
		t.Errorf("expected 'Updated Group', got %s", group.Title)
	}
}

func TestAppointmentGroupsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/appointment_groups/1" {
			t.Errorf("expected path /api/v1/appointment_groups/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Query().Get("cancel_reason") != "test reason" {
			t.Errorf("expected cancel_reason 'test reason', got %s", r.URL.Query().Get("cancel_reason"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":1,"title":"Deleted Group","workflow_state":"deleted"}`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAppointmentGroupsService(client)

	group, err := svc.Delete(context.Background(), 1, "test reason")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if group == nil {
		t.Fatal("expected non-nil group response")
	}
}

func TestAppointmentGroupsService_ListUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/appointment_groups/1/users" {
			t.Errorf("expected path /api/v1/appointment_groups/1/users, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":42,"name":"Alice"},{"id":43,"name":"Bob"}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAppointmentGroupsService(client)

	users, err := svc.ListUsers(context.Background(), 1, "")
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestAppointmentGroupsService_ListGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/appointment_groups/1/groups" {
			t.Errorf("expected path /api/v1/appointment_groups/1/groups, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":5,"name":"Team Alpha"}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAppointmentGroupsService(client)

	groups, err := svc.ListGroups(context.Background(), 1, "registered")
	if err != nil {
		t.Fatalf("ListGroups failed: %v", err)
	}
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "Team Alpha" {
		t.Errorf("expected 'Team Alpha', got %s", groups[0].Name)
	}
}

func TestAppointmentGroupsService_NextAppointment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/appointment_groups/next_appointment" {
			t.Errorf("expected path /api/v1/appointment_groups/next_appointment, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":99,"title":"Next Slot","context_code":"course_123","workflow_state":"active"}]`))
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewAppointmentGroupsService(client)

	events, err := svc.NextAppointment(context.Background(), []int64{1, 2})
	if err != nil {
		t.Fatalf("NextAppointment failed: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].ID != 99 {
		t.Errorf("expected event ID 99, got %d", events[0].ID)
	}
}

func TestNewAppointmentGroupsService(t *testing.T) {
	client := &Client{}
	svc := NewAppointmentGroupsService(client)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
	if svc.client != client {
		t.Error("expected service client to match input client")
	}
}
