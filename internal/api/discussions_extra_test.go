package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Group-context variants of existing operations ---

func TestDiscussionsService_ListContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/discussion_topics" {
			t.Errorf("expected path /api/v1/groups/42/discussion_topics, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]DiscussionTopic{{ID: 10, Title: "Group Topic"}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	topics, err := svc.ListContext(context.Background(), "groups", 42, nil)
	if err != nil {
		t.Fatalf("ListContext: %v", err)
	}
	if len(topics) != 1 || topics[0].ID != 10 {
		t.Errorf("unexpected topics: %+v", topics)
	}
}

func TestDiscussionsService_GetContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/discussion_topics/7" {
			t.Errorf("expected path /api/v1/groups/42/discussion_topics/7, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionTopic{ID: 7, Title: "Group Discussion"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	topic, err := svc.GetContext(context.Background(), "groups", 42, 7, nil)
	if err != nil {
		t.Fatalf("GetContext: %v", err)
	}
	if topic.ID != 7 {
		t.Errorf("expected id 7, got %d", topic.ID)
	}
}

func TestDiscussionsService_CreateContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/discussion_topics" {
			t.Errorf("expected path /api/v1/groups/42/discussion_topics, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["title"] != "Group Discussion" {
			t.Errorf("unexpected title: %v", body["title"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionTopic{ID: 8, Title: "Group Discussion"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	topic, err := svc.CreateContext(context.Background(), "groups", 42, &CreateDiscussionParams{
		Title: "Group Discussion",
	})
	if err != nil {
		t.Fatalf("CreateContext: %v", err)
	}
	if topic.ID != 8 {
		t.Errorf("expected id 8, got %d", topic.ID)
	}
}

func TestDiscussionsService_UpdateContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/discussion_topics/7" {
			t.Errorf("expected path /api/v1/groups/42/discussion_topics/7, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionTopic{ID: 7, Title: "Updated"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)
	title := "Updated"

	topic, err := svc.UpdateContext(context.Background(), "groups", 42, 7, &UpdateDiscussionParams{Title: &title})
	if err != nil {
		t.Fatalf("UpdateContext: %v", err)
	}
	if topic.Title != "Updated" {
		t.Errorf("unexpected title: %s", topic.Title)
	}
}

func TestDiscussionsService_DeleteContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/42/discussion_topics/7" {
			t.Errorf("expected path /api/v1/groups/42/discussion_topics/7, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.DeleteContext(context.Background(), "groups", 42, 7); err != nil {
		t.Fatalf("DeleteContext: %v", err)
	}
}

// --- GetView ---

func TestDiscussionsService_GetView_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/view" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionView{
			UnreadEntries: []int64{1, 2},
			View:          []DiscussionEntry{{ID: 1, Message: "hello"}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	view, err := svc.GetView(context.Background(), "courses", 10, 20)
	if err != nil {
		t.Fatalf("GetView: %v", err)
	}
	if len(view.UnreadEntries) != 2 {
		t.Errorf("expected 2 unread entries, got %d", len(view.UnreadEntries))
	}
	if len(view.View) != 1 || view.View[0].Message != "hello" {
		t.Errorf("unexpected view entries: %+v", view.View)
	}
}

func TestDiscussionsService_GetView_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/view" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionView{})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	_, err := svc.GetView(context.Background(), "groups", 5, 9)
	if err != nil {
		t.Fatalf("GetView group: %v", err)
	}
}

// --- Duplicate ---

func TestDiscussionsService_Duplicate_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/duplicate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionTopic{ID: 21, Title: "Copy of topic"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	topic, err := svc.Duplicate(context.Background(), "courses", 10, 20)
	if err != nil {
		t.Fatalf("Duplicate: %v", err)
	}
	if topic.ID != 21 {
		t.Errorf("expected id 21, got %d", topic.ID)
	}
}

func TestDiscussionsService_Duplicate_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/duplicate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionTopic{ID: 10})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	_, err := svc.Duplicate(context.Background(), "groups", 5, 9)
	if err != nil {
		t.Fatalf("Duplicate group: %v", err)
	}
}

// --- Reorder ---

func TestDiscussionsService_Reorder_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/reorder" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.Reorder(context.Background(), "courses", 10, []int64{3, 1, 2}); err != nil {
		t.Fatalf("Reorder: %v", err)
	}
}

func TestDiscussionsService_Reorder_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/reorder" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.Reorder(context.Background(), "groups", 5, []int64{1, 2}); err != nil {
		t.Fatalf("Reorder group: %v", err)
	}
}

// --- UpdateEntry / DeleteEntry ---

func TestDiscussionsService_UpdateEntry_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/entries/30" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["message"] != "updated text" {
			t.Errorf("unexpected message: %v", body["message"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionEntry{ID: 30, Message: "updated text"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	entry, err := svc.UpdateEntry(context.Background(), "courses", 10, 20, 30, "updated text")
	if err != nil {
		t.Fatalf("UpdateEntry: %v", err)
	}
	if entry.Message != "updated text" {
		t.Errorf("unexpected message: %s", entry.Message)
	}
}

func TestDiscussionsService_UpdateEntry_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/entries/11" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionEntry{ID: 11})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	_, err := svc.UpdateEntry(context.Background(), "groups", 5, 9, 11, "msg")
	if err != nil {
		t.Fatalf("UpdateEntry group: %v", err)
	}
}

func TestDiscussionsService_DeleteEntry_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/entries/30" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.DeleteEntry(context.Background(), "courses", 10, 20, 30); err != nil {
		t.Fatalf("DeleteEntry: %v", err)
	}
}

func TestDiscussionsService_DeleteEntry_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/entries/11" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.DeleteEntry(context.Background(), "groups", 5, 9, 11); err != nil {
		t.Fatalf("DeleteEntry group: %v", err)
	}
}

// --- ListReplies ---

func TestDiscussionsService_ListReplies_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/entries/30/replies" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]DiscussionEntry{
			{ID: 100, Message: "reply 1"},
			{ID: 101, Message: "reply 2"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	replies, err := svc.ListReplies(context.Background(), "courses", 10, 20, 30)
	if err != nil {
		t.Fatalf("ListReplies: %v", err)
	}
	if len(replies) != 2 {
		t.Errorf("expected 2 replies, got %d", len(replies))
	}
}

func TestDiscussionsService_ListReplies_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/entries/11/replies" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]DiscussionEntry{{ID: 200}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	replies, err := svc.ListReplies(context.Background(), "groups", 5, 9, 11)
	if err != nil {
		t.Fatalf("ListReplies group: %v", err)
	}
	if len(replies) != 1 {
		t.Errorf("expected 1 reply, got %d", len(replies))
	}
}

// --- GetEntryList ---

func TestDiscussionsService_GetEntryList_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/entry_list" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		// Verify ids[] query parameters
		ids := r.URL.Query()["ids[]"]
		if len(ids) != 2 {
			t.Errorf("expected 2 ids, got %d: %v", len(ids), ids)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]DiscussionEntry{
			{ID: 1, Message: "entry 1"},
			{ID: 2, Message: "entry 2"},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	entries, err := svc.GetEntryList(context.Background(), "courses", 10, 20, []int64{1, 2})
	if err != nil {
		t.Fatalf("GetEntryList: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestDiscussionsService_GetEntryList_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/entry_list" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]DiscussionEntry{{ID: 5}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	_, err := svc.GetEntryList(context.Background(), "groups", 5, 9, nil)
	if err != nil {
		t.Fatalf("GetEntryList group: %v", err)
	}
}

// --- PostEntryContext group variant ---

func TestDiscussionsService_PostEntryContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/entries" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionEntry{ID: 50, Message: "hello group"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	entry, err := svc.PostEntryContext(context.Background(), "groups", 5, 9, "hello group")
	if err != nil {
		t.Fatalf("PostEntryContext group: %v", err)
	}
	if entry.Message != "hello group" {
		t.Errorf("unexpected message: %s", entry.Message)
	}
}

// --- PostReplyContext group variant ---

func TestDiscussionsService_PostReplyContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/entries/11/replies" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		parentID := int64(11)
		json.NewEncoder(w).Encode(DiscussionEntry{ID: 60, ParentID: &parentID})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	entry, err := svc.PostReplyContext(context.Background(), "groups", 5, 9, 11, "group reply")
	if err != nil {
		t.Fatalf("PostReplyContext group: %v", err)
	}
	if entry.ParentID == nil || *entry.ParentID != 11 {
		t.Errorf("expected parent_id 11, got %v", entry.ParentID)
	}
}

// --- ListEntriesContext group variant ---

func TestDiscussionsService_ListEntriesContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/entries" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]DiscussionEntry{{ID: 1}, {ID: 2}})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	entries, err := svc.ListEntriesContext(context.Background(), "groups", 5, 9)
	if err != nil {
		t.Fatalf("ListEntriesContext group: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

// --- MarkAllTopicsRead ---

func TestDiscussionsService_MarkAllTopicsRead_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/read_all" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkAllTopicsRead(context.Background(), "courses", 10); err != nil {
		t.Fatalf("MarkAllTopicsRead: %v", err)
	}
}

func TestDiscussionsService_MarkAllTopicsRead_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/read_all" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkAllTopicsRead(context.Background(), "groups", 5); err != nil {
		t.Fatalf("MarkAllTopicsRead group: %v", err)
	}
}

// --- MarkAllEntriesRead / MarkAllEntriesUnread ---

func TestDiscussionsService_MarkAllEntriesRead_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/read_all" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkAllEntriesRead(context.Background(), "courses", 10, 20); err != nil {
		t.Fatalf("MarkAllEntriesRead: %v", err)
	}
}

func TestDiscussionsService_MarkAllEntriesRead_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/read_all" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkAllEntriesRead(context.Background(), "groups", 5, 9); err != nil {
		t.Fatalf("MarkAllEntriesRead group: %v", err)
	}
}

func TestDiscussionsService_MarkAllEntriesUnread_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/read_all" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkAllEntriesUnread(context.Background(), "courses", 10, 20); err != nil {
		t.Fatalf("MarkAllEntriesUnread: %v", err)
	}
}

func TestDiscussionsService_MarkAllEntriesUnread_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/read_all" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkAllEntriesUnread(context.Background(), "groups", 5, 9); err != nil {
		t.Fatalf("MarkAllEntriesUnread group: %v", err)
	}
}

// --- MarkEntryRead / MarkEntryUnread ---

func TestDiscussionsService_MarkEntryRead_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/entries/30/read" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkEntryRead(context.Background(), "courses", 10, 20, 30); err != nil {
		t.Fatalf("MarkEntryRead: %v", err)
	}
}

func TestDiscussionsService_MarkEntryRead_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/entries/11/read" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkEntryRead(context.Background(), "groups", 5, 9, 11); err != nil {
		t.Fatalf("MarkEntryRead group: %v", err)
	}
}

func TestDiscussionsService_MarkEntryUnread_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/entries/30/read" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkEntryUnread(context.Background(), "courses", 10, 20, 30); err != nil {
		t.Fatalf("MarkEntryUnread: %v", err)
	}
}

func TestDiscussionsService_MarkEntryUnread_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/entries/11/read" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkEntryUnread(context.Background(), "groups", 5, 9, 11); err != nil {
		t.Fatalf("MarkEntryUnread group: %v", err)
	}
}

// --- MarkTopicReadContext / MarkTopicUnreadContext group variants ---

func TestDiscussionsService_MarkTopicReadContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/read" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkTopicReadContext(context.Background(), "groups", 5, 9); err != nil {
		t.Fatalf("MarkTopicReadContext group: %v", err)
	}
}

func TestDiscussionsService_MarkTopicUnreadContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/read" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.MarkTopicUnreadContext(context.Background(), "groups", 5, 9); err != nil {
		t.Fatalf("MarkTopicUnreadContext group: %v", err)
	}
}

// --- RateEntry ---

func TestDiscussionsService_RateEntry_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/entries/30/rating" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		// JSON numbers decode to float64
		if rating, ok := body["rating"].(float64); !ok || rating != 1 {
			t.Errorf("expected rating=1, got %v", body["rating"])
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.RateEntry(context.Background(), "courses", 10, 20, 30, 1); err != nil {
		t.Fatalf("RateEntry: %v", err)
	}
}

func TestDiscussionsService_RateEntry_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/entries/11/rating" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.RateEntry(context.Background(), "groups", 5, 9, 11, 0); err != nil {
		t.Fatalf("RateEntry group: %v", err)
	}
}

// --- SubscribeContext / UnsubscribeContext group variants ---

func TestDiscussionsService_SubscribeContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/subscribed" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.SubscribeContext(context.Background(), "groups", 5, 9); err != nil {
		t.Fatalf("SubscribeContext group: %v", err)
	}
}

func TestDiscussionsService_UnsubscribeContext_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/subscribed" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.UnsubscribeContext(context.Background(), "groups", 5, 9); err != nil {
		t.Fatalf("UnsubscribeContext group: %v", err)
	}
}

// --- GetSummary / CreateSummary / DisableSummary / SummaryFeedback ---

func TestDiscussionsService_GetSummary_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/summaries" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionSummary{ID: 1, Text: "Summary text"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	summary, err := svc.GetSummary(context.Background(), "courses", 10, 20)
	if err != nil {
		t.Fatalf("GetSummary: %v", err)
	}
	if summary.Text != "Summary text" {
		t.Errorf("unexpected text: %s", summary.Text)
	}
}

func TestDiscussionsService_GetSummary_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/summaries" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionSummary{ID: 2})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	_, err := svc.GetSummary(context.Background(), "groups", 5, 9)
	if err != nil {
		t.Fatalf("GetSummary group: %v", err)
	}
}

func TestDiscussionsService_CreateSummary_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/summaries" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["userInput"] != "Focus on key themes" {
			t.Errorf("unexpected userInput: %v", body["userInput"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionSummary{ID: 3, Text: "Generated summary"})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	summary, err := svc.CreateSummary(context.Background(), "courses", 10, 20, "Focus on key themes")
	if err != nil {
		t.Fatalf("CreateSummary: %v", err)
	}
	if summary.Text != "Generated summary" {
		t.Errorf("unexpected text: %s", summary.Text)
	}
}

func TestDiscussionsService_CreateSummary_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/summaries" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DiscussionSummary{ID: 4})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	// No userInput — body should be nil
	_, err := svc.CreateSummary(context.Background(), "groups", 5, 9, "")
	if err != nil {
		t.Fatalf("CreateSummary group: %v", err)
	}
}

func TestDiscussionsService_DisableSummary_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/summaries/disable" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.DisableSummary(context.Background(), "courses", 10, 20); err != nil {
		t.Fatalf("DisableSummary: %v", err)
	}
}

func TestDiscussionsService_DisableSummary_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/summaries/disable" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	if err := svc.DisableSummary(context.Background(), "groups", 5, 9); err != nil {
		t.Fatalf("DisableSummary group: %v", err)
	}
}

func TestDiscussionsService_SummaryFeedback_Course(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/10/discussion_topics/20/summaries/5/feedback" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["_action"] != "like" {
			t.Errorf("unexpected _action: %v", body["_action"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SummaryFeedbackResult{Liked: true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	result, err := svc.SummaryFeedback(context.Background(), "courses", 10, 20, 5, "like")
	if err != nil {
		t.Fatalf("SummaryFeedback: %v", err)
	}
	if !result.Liked {
		t.Error("expected liked=true")
	}
}

func TestDiscussionsService_SummaryFeedback_Group(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/groups/5/discussion_topics/9/summaries/3/feedback" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SummaryFeedbackResult{Disliked: true})
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)
	svc := NewDiscussionsService(client)

	result, err := svc.SummaryFeedback(context.Background(), "groups", 5, 9, 3, "dislike")
	if err != nil {
		t.Fatalf("SummaryFeedback group: %v", err)
	}
	if !result.Disliked {
		t.Error("expected disliked=true")
	}
}
