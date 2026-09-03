package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDateChanges_Payload(t *testing.T) {
	due := time.Date(2026, 9, 9, 16, 50, 0, 0, time.FixedZone("EDT", -4*3600))
	body, err := DateChanges{DateFieldDueAt: &due, DateFieldLockAt: nil}.Payload("quiz")
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(body)
	if got := string(raw); got != `{"quiz":{"due_at":"2026-09-09T20:50:00Z","lock_at":null}}` {
		t.Errorf("payload = %s", got)
	}
	if _, err := (DateChanges{"show_correct_answers_at": &due}).Payload("quiz"); err == nil || !strings.Contains(err.Error(), "unknown date field") {
		t.Errorf("unknown field accepted: %v", err)
	}
	empty, _ := DateChanges{}.Payload("assignment")
	raw, _ = json.Marshal(empty)
	if string(raw) != `{"assignment":{}}` {
		t.Errorf("empty payload = %s", raw)
	}
}

func TestUpdateDates(t *testing.T) {
	var gotMethod, gotPath, gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		gotMethod, gotPath = r.Method, r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.Contains(r.URL.Path, "/quizzes/"):
			_, _ = w.Write([]byte(`{"id":456,"title":"Attendance","unlock_at":"2026-09-09T20:00:00Z","due_at":"2026-09-09T20:50:00Z","lock_at":null}`))
		default:
			_, _ = w.Write([]byte(`{"id":789,"name":"Essay","due_at":"2026-09-09T20:50:00Z","lock_at":null}`))
		}
	}))
	defer server.Close()
	client := newTestClient(t, server.URL)

	unlock := time.Date(2026, 9, 9, 20, 0, 0, 0, time.UTC)
	due := time.Date(2026, 9, 9, 20, 50, 0, 0, time.UTC)

	t.Run("quiz", func(t *testing.T) {
		quiz, err := NewQuizzesService(client).UpdateDates(context.Background(), 1, 456, DateChanges{
			DateFieldUnlockAt: &unlock, DateFieldDueAt: &due, DateFieldLockAt: nil,
		})
		if err != nil {
			t.Fatal(err)
		}
		if gotMethod != "PUT" || gotPath != "/api/v1/courses/1/quizzes/456" {
			t.Errorf("request = %s %s", gotMethod, gotPath)
		}
		var body map[string]map[string]interface{}
		if err := json.Unmarshal([]byte(gotBody), &body); err != nil {
			t.Fatalf("body not json: %s", gotBody)
		}
		q := body["quiz"]
		if q["unlock_at"] != "2026-09-09T20:00:00Z" || q["due_at"] != "2026-09-09T20:50:00Z" {
			t.Errorf("body = %s", gotBody)
		}
		if v, ok := q["lock_at"]; !ok || v != nil {
			t.Errorf("lock_at should be an explicit null: %s", gotBody)
		}
		if quiz.DueAt == nil || !quiz.DueAt.Equal(due) || quiz.LockAt != nil {
			t.Errorf("decoded quiz = %+v", quiz)
		}
	})

	t.Run("assignment", func(t *testing.T) {
		a, err := NewAssignmentsService(client).UpdateDates(context.Background(), 1, 789, DateChanges{DateFieldDueAt: &due})
		if err != nil {
			t.Fatal(err)
		}
		if gotMethod != "PUT" || gotPath != "/api/v1/courses/1/assignments/789" {
			t.Errorf("request = %s %s", gotMethod, gotPath)
		}
		if gotBody != `{"assignment":{"due_at":"2026-09-09T20:50:00Z"}}` {
			t.Errorf("body = %s", gotBody)
		}
		if !a.DueAt.Equal(due) || !a.LockAt.IsZero() {
			t.Errorf("decoded assignment = due %v lock %v", a.DueAt, a.LockAt)
		}
	})

	t.Run("unknown field is refused before any request", func(t *testing.T) {
		gotPath = ""
		if _, err := NewAssignmentsService(client).UpdateDates(context.Background(), 1, 789, DateChanges{"nope": nil}); err == nil {
			t.Fatal("want an error")
		}
		if gotPath != "" {
			t.Errorf("a request was sent: %s", gotPath)
		}
	})
}
