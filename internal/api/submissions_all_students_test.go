package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// gridServer is a fake /students/submissions endpoint that records every
// request URL and answers with one submission per assignment_ids[] entry (or
// a fixed page set when no assignment filter is given).
type gridServer struct {
	*httptest.Server
	mu   sync.Mutex
	urls []string
}

func newGridServer(t *testing.T) *gridServer {
	t.Helper()
	gs := &gridServer{}
	gs.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			handleVersionDetection(w)
			return
		}
		if r.URL.Path != "/api/v1/courses/7/students/submissions" {
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		gs.mu.Lock()
		gs.urls = append(gs.urls, r.URL.RequestURI())
		gs.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		q := r.URL.Query()
		ids := q["assignment_ids[]"]
		if len(ids) == 0 {
			// Two-page unfiltered grid, joined by a Link header.
			if q.Get("page") == "2" {
				fmt.Fprint(w, `[{"assignment_id":3,"user_id":10,"workflow_state":"graded"}]`)
				return
			}
			next := gs.URL + "/api/v1/courses/7/students/submissions?" + q.Encode() + "&page=2"
			w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next"`, next))
			fmt.Fprint(w, `[{"assignment_id":1,"user_id":10,"workflow_state":"unsubmitted"},{"assignment_id":2,"user_id":10,"workflow_state":"submitted"}]`)
			return
		}
		var rows []map[string]interface{}
		for _, id := range ids {
			n, _ := strconv.ParseInt(id, 10, 64)
			rows = append(rows, map[string]interface{}{"assignment_id": n, "user_id": 10, "workflow_state": "unsubmitted"})
		}
		_ = json.NewEncoder(w).Encode(rows)
	}))
	return gs
}

func TestSubmissionsService_ListForAllStudents_UnfilteredGridFollowsPagination(t *testing.T) {
	gs := newGridServer(t)
	defer gs.Close()
	svc := NewSubmissionsService(newTestClient(t, gs.URL))

	subs, err := svc.ListForAllStudents(context.Background(), 7, nil, &ListSubmissionsOptions{Include: []string{"assignment"}, PerPage: 100})
	if err != nil {
		t.Fatalf("ListForAllStudents: %v", err)
	}
	if len(subs) != 3 {
		t.Fatalf("got %d submissions, want 3 across two pages", len(subs))
	}
	if len(gs.urls) != 2 {
		t.Fatalf("requests = %v, want the first page plus the Link-header page", gs.urls)
	}
	first := gs.urls[0]
	for _, want := range []string{"student_ids%5B%5D=all", "include%5B%5D=assignment", "per_page=100"} {
		if !strings.Contains(first, want) {
			t.Errorf("first request %q lacks %s", first, want)
		}
	}
	if strings.Contains(first, "assignment_ids") {
		t.Errorf("unfiltered grid must not send assignment_ids[]: %s", first)
	}
}

func TestSubmissionsService_ListForAllStudents_ChunksLongAssignmentLists(t *testing.T) {
	gs := newGridServer(t)
	defer gs.Close()
	svc := NewSubmissionsService(newTestClient(t, gs.URL))

	var ids []int64
	for i := int64(0); i < 250; i++ {
		ids = append(ids, 1000000+i)
	}
	subs, err := svc.ListForAllStudents(context.Background(), 7, ids, &ListSubmissionsOptions{Include: []string{"assignment"}, PerPage: 100})
	if err != nil {
		t.Fatalf("ListForAllStudents: %v", err)
	}
	if len(subs) != len(ids) {
		t.Fatalf("got %d submissions, want one per assignment (%d)", len(subs), len(ids))
	}
	if len(gs.urls) < 2 {
		t.Fatalf("expected the 250 ids to be split across several requests, got %d", len(gs.urls))
	}
	seen := map[string]bool{}
	for _, u := range gs.urls {
		if len(u) > maxSubmissionsGridURLLen {
			t.Errorf("request URL is %d chars, over the %d limit", len(u), maxSubmissionsGridURLLen)
		}
		parsed, err := url.Parse(u)
		if err != nil {
			t.Fatal(err)
		}
		q := parsed.Query()
		if q.Get("student_ids[]") != "all" || q.Get("include[]") != "assignment" || q.Get("per_page") != "100" {
			t.Errorf("chunk lost the base query: %s", u)
		}
		for _, id := range q["assignment_ids[]"] {
			if seen[id] {
				t.Errorf("assignment %s requested twice", id)
			}
			seen[id] = true
		}
	}
	if len(seen) != len(ids) {
		t.Errorf("chunks covered %d ids, want %d", len(seen), len(ids))
	}
}

func TestChunkAssignmentIDs_SingleChunkWhenShort(t *testing.T) {
	base := url.Values{}
	base.Add("student_ids[]", "all")
	urls := chunkAssignmentIDs("/api/v1/courses/7/students/submissions", base, []int64{1, 2, 3})
	if len(urls) != 1 {
		t.Fatalf("got %d urls, want 1: %v", len(urls), urls)
	}
	if !strings.Contains(urls[0], "assignment_ids%5B%5D=1&assignment_ids%5B%5D=2&assignment_ids%5B%5D=3") {
		t.Errorf("url = %s", urls[0])
	}
}
