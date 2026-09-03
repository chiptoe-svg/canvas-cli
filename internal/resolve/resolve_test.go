package resolve

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var roster = []Student{
	{ID: 10, Name: "Ada Lovelace", SortableName: "Lovelace, Ada", ShortName: "Ada", LoginID: "ada@example.edu", SISUserID: "S001"},
	{ID: 11, Name: "Ann Smith", SortableName: "Smith, Ann", LoginID: "asmith"},
	{ID: 12, Name: "Anne Smithers", SortableName: "Smithers, Anne", LoginID: "asmithers"},
	{ID: 13, Name: "Bo Li", SortableName: "Li, Bo"},
	{ID: 14, Name: "Bo Lin", SortableName: "Lin, Bo"},
}

func TestFindStudent(t *testing.T) {
	cases := []struct {
		query  string
		wantID int64
		errSub string
	}{
		{"10", 10, ""},
		{" 13 ", 13, ""},
		{"ada lovelace", 10, ""},
		{"Lovelace, Ada", 10, ""},
		{"ADA", 10, ""},             // short name, exact
		{"ada@example.edu", 10, ""}, // login id
		{"S001", 10, ""},            // SIS id
		{"lovel", 10, ""},           // unique substring
		{"Ann Smith", 11, ""},       // exact beats the substring in "Anne Smithers"
		{"Bo Li", 13, ""},           // exact beats "Bo Lin"
		{"smith", 0, `"smith" matches 2 students`},
		{"bo", 0, `"bo" matches 2 students`},
		{"99", 0, `no student matches "99" (5 searched)`},
		{"zzz", 0, `no student matches "zzz"`},
		{"", 0, "student is required"},
	}
	for _, c := range cases {
		t.Run(c.query, func(t *testing.T) {
			s, err := FindStudent(c.query, roster)
			if c.errSub != "" {
				if err == nil || !strings.Contains(err.Error(), c.errSub) {
					t.Fatalf("err = %v, want %q", err, c.errSub)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if s.ID != c.wantID {
				t.Errorf("got %d, want %d", s.ID, c.wantID)
			}
		})
	}
	// Ambiguity lists the candidates, sorted, with sortable name and login.
	_, err := FindStudent("smith", roster)
	var amb *AmbiguousError
	if !errors.As(err, &amb) || len(amb.Candidates) != 2 || amb.Candidates[0] != "11 Ann Smith (Smith, Ann) <asmith>" || amb.Candidates[1] != "12 Anne Smithers (Smithers, Anne) <asmithers>" {
		t.Errorf("ambiguous = %v", err)
	}
	var none *NoMatchError
	if _, err := FindStudent("zzz", roster); !errors.As(err, &none) || none.Total != 5 {
		t.Errorf("no match = %v", err)
	}
}

var items = []Assignment{
	{ID: 456, Name: "Quiz 3", Kind: KindQuiz, QuizID: 77},
	{ID: 457, Name: "Essay 1", Kind: KindAssignment},
	{ID: 458, Name: "Essay 10", Kind: KindAssignment},
	{ID: 459, Name: "Week 1 discussion", Kind: KindDiscussion},
	{ID: 77, Name: "Seventy-seven", Kind: KindAssignment},
}

func TestFindAssignment(t *testing.T) {
	cases := []struct {
		query  string
		wantID int64
		errSub string
	}{
		{"456", 456, ""},
		{"77", 77, ""}, // an assignment id wins over a quiz id
		{"quiz 3", 456, ""},
		{"QUIZ", 456, ""},    // unique substring
		{"Essay 1", 457, ""}, // exact beats the substring in "Essay 10"
		{"essay 10", 458, ""},
		{"discussion", 459, ""},
		{"essay", 0, `"essay" matches 2 assignments`},
		{"999", 0, `no assignment matches "999" (5 searched)`},
		{"nope", 0, "no assignment matches"},
		{"", 0, "assignment is required"},
	}
	for _, c := range cases {
		t.Run(c.query, func(t *testing.T) {
			a, err := FindAssignment(c.query, items)
			if c.errSub != "" {
				if err == nil || !strings.Contains(err.Error(), c.errSub) {
					t.Fatalf("err = %v, want %q", err, c.errSub)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if a.ID != c.wantID {
				t.Errorf("got %d, want %d", a.ID, c.wantID)
			}
		})
	}
	// A quiz id resolves to its assignment when no assignment has that id.
	a, err := FindAssignment("77", items[:4])
	if err != nil || a.ID != 456 {
		t.Errorf("quiz id: %v, %v", a, err)
	}
	if got := DescribeAssignment(&items[0]); got != "456 Quiz 3 [quiz 77]" {
		t.Errorf("DescribeAssignment = %q", got)
	}
	if got := DescribeAssignment(&items[3]); got != "459 Week 1 discussion [discussion]" {
		t.Errorf("DescribeAssignment = %q", got)
	}
	if got := DescribeStudent(&roster[3]); got != "13 Bo Li (Li, Bo)" {
		t.Errorf("DescribeStudent = %q", got)
	}
}

func TestAmbiguousError_Cap(t *testing.T) {
	var many []string
	for i := 0; i < 14; i++ {
		many = append(many, fmt.Sprintf("%d Student %02d", i, i))
	}
	msg := (&AmbiguousError{What: "student", Query: "student", Candidates: many}).Error()
	if !strings.Contains(msg, "matches 14 students") || !strings.Contains(msg, "… and 4 more") || strings.Contains(msg, "Student 12") {
		t.Errorf("message = %s", msg)
	}
}

func TestListStudentsAndAssignments(t *testing.T) {
	var uris []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/accounts" {
			_, _ = w.Write([]byte(`[]`))
			return
		}
		uris = append(uris, r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/courses/1/users":
			fmt.Fprint(w, `[{"id":10,"name":"Ada Lovelace","sortable_name":"Lovelace, Ada","short_name":"Ada","login_id":"ada","sis_user_id":"S1"},{"id":13,"name":"Test Student","sortable_name":"Student, Test"}]`)
		case "/api/v1/courses/1/assignments":
			fmt.Fprint(w, `[{"id":456,"name":"Quiz 3","quiz_id":77,"is_quiz_assignment":true,"submission_types":["online_quiz"]},{"id":457,"name":"Essay 1","submission_types":["online_upload"]},{"id":459,"name":"Discuss","submission_types":["discussion_topic"]}]`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client, err := api.NewClient(api.ClientConfig{BaseURL: server.URL, Token: "t", RequestsPerSec: 1000})
	if err != nil {
		t.Fatal(err)
	}

	students, err := ListStudents(context.Background(), client, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(students) != 1 || students[0] != (Student{ID: 10, Name: "Ada Lovelace", SortableName: "Lovelace, Ada", ShortName: "Ada", LoginID: "ada", SISUserID: "S1"}) {
		t.Errorf("students = %+v", students)
	}
	if len(uris) != 1 || !strings.Contains(uris[0], "enrollment_type%5B%5D=student") || !strings.Contains(uris[0], "enrollment_state%5B%5D=active") || !strings.Contains(uris[0], "per_page=100") {
		t.Errorf("roster request = %v", uris)
	}

	items, err := ListAssignments(context.Background(), client, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 || items[0] != (Assignment{ID: 456, Name: "Quiz 3", Kind: KindQuiz, QuizID: 77}) || items[1].Kind != KindAssignment || items[2].Kind != KindDiscussion {
		t.Errorf("items = %+v", items)
	}

	if _, err := ListStudents(context.Background(), client, 2); err == nil || !strings.Contains(err.Error(), "failed to list students of course 2") {
		t.Errorf("error wrapping: %v", err)
	}
	if _, err := ListAssignments(context.Background(), client, 2); err == nil || !strings.Contains(err.Error(), "failed to list assignments of course 2") {
		t.Errorf("error wrapping: %v", err)
	}
}
