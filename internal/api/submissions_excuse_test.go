package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetExcused(t *testing.T) {
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
		_, _ = w.Write([]byte(`{"id":1,"user_id":10,"assignment_id":456,"excused":true,"workflow_state":"graded"}`))
	}))
	defer server.Close()
	svc := NewSubmissionsService(newTestClient(t, server.URL))

	sub, err := svc.SetExcused(context.Background(), 1, 456, 10, true)
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != "PUT" || gotPath != "/api/v1/courses/1/assignments/456/submissions/10" {
		t.Errorf("request = %s %s", gotMethod, gotPath)
	}
	if gotBody != `{"submission":{"excuse":true}}` {
		t.Errorf("body = %s", gotBody)
	}
	if !sub.ExcusedTLN || sub.UserID != 10 {
		t.Errorf("decoded = %+v", sub)
	}

	if _, err := svc.SetExcused(context.Background(), 1, 456, 10, false); err != nil {
		t.Fatal(err)
	}
	if gotBody != `{"submission":{"excuse":false}}` {
		t.Errorf("an explicit false must be sent: %s", gotBody)
	}
}
