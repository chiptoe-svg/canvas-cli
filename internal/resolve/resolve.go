// Package resolve turns the names people use — "ada lovelace", "Quiz 3",
// "lineup" — into Canvas ids, refusing anything that does not name exactly
// one student or one assignment. Commands share these resolvers so that
// "excuse Ada from Quiz 3" and "reschedule the lineup assignment" pick
// their targets by the same rules.
//
// Matching rules, in order, for a query q:
//  1. all digits: the object whose id is q (for assignments, a quiz id
//     resolves to the quiz's assignment);
//  2. a case-insensitive exact match on a name (for students: name,
//     sortable name, short name, login id or SIS id);
//  3. a case-insensitive substring of a name, which must match exactly one
//     object.
//
// Zero matches and more than one match are errors that list the
// candidates, so the caller can show the user what to type instead.
package resolve

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// Student is a course member as the resolvers see them.
type Student struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	SortableName string `json:"sortable_name,omitempty"`
	ShortName    string `json:"short_name,omitempty"`
	LoginID      string `json:"login_id,omitempty"`
	SISUserID    string `json:"sis_user_id,omitempty"`
}

// Assignment is a gradable item: a plain assignment, or the assignment
// behind a quiz or graded discussion.
type Assignment struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Kind   string `json:"type"` // assignment | quiz | discussion
	QuizID int64  `json:"quiz_id,omitempty"`
}

// Kinds of Assignment.
const (
	KindAssignment = "assignment"
	KindQuiz       = "quiz"
	KindDiscussion = "discussion"
)

// MaxCandidates caps the candidate list in error messages.
const MaxCandidates = 10

// NoMatchError is returned when nothing matches the query.
type NoMatchError struct {
	What  string // "student" | "assignment"
	Query string
	Total int // how many objects were searched
}

func (e *NoMatchError) Error() string {
	return fmt.Sprintf("no %s matches %q (%d searched); give an id or the exact name", e.What, e.Query, e.Total)
}

// AmbiguousError is returned when the query matches more than one object.
type AmbiguousError struct {
	What       string
	Query      string
	Candidates []string // rendered "id name" lines, sorted
}

func (e *AmbiguousError) Error() string {
	shown := e.Candidates
	more := ""
	if len(shown) > MaxCandidates {
		more = fmt.Sprintf("\n  … and %d more", len(shown)-MaxCandidates)
		shown = shown[:MaxCandidates]
	}
	return fmt.Sprintf("%q matches %d %ss; give an id or a more specific name:\n  %s%s", e.Query, len(e.Candidates), e.What, strings.Join(shown, "\n  "), more)
}

// FindStudent picks exactly one student for query. The result names the
// student a write is applied to, so only an exact match counts: the Canvas
// id, or the full name, sortable name, short name, login id or SIS id
// (case-insensitive). A query that only matches part of a name is refused
// with the candidates listed — "lee" must not quietly become whichever Lee
// happens to be enrolled. A numeric query that is one student's Canvas id
// and another's SIS or login id is refused as ambiguous.
func FindStudent(query string, students []Student) (*Student, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, fmt.Errorf("student is required")
	}
	lq := strings.ToLower(q)
	var exact []int
	if id, err := strconv.ParseInt(q, 10, 64); err == nil {
		for i := range students {
			s := &students[i]
			if s.ID == id || equalFold(lq, s.LoginID, s.SISUserID) {
				exact = append(exact, i)
			}
		}
	} else {
		for i := range students {
			s := &students[i]
			if equalFold(lq, s.Name, s.SortableName, s.ShortName, s.LoginID, s.SISUserID) {
				exact = append(exact, i)
			}
		}
	}
	switch len(exact) {
	case 1:
		return &students[exact[0]], nil
	case 0:
		var partial []string
		for i := range students {
			s := &students[i]
			if containsFold(lq, s.Name, s.SortableName) {
				partial = append(partial, describeStudent(s))
			}
		}
		if len(partial) > 0 {
			sort.Strings(partial)
			return nil, &PartialMatchError{What: "student", Query: query, Candidates: partial}
		}
		return nil, &NoMatchError{What: "student", Query: query, Total: len(students)}
	}
	var lines []string
	for _, i := range exact {
		lines = append(lines, describeStudent(&students[i]))
	}
	sort.Strings(lines)
	return nil, &AmbiguousError{What: "student", Query: query, Candidates: lines}
}

// PartialMatchError: the query matched part of one or more names but none
// exactly. For a write that is not enough; the caller must give the full
// name or an id.
type PartialMatchError struct {
	What       string
	Query      string
	Candidates []string
}

func (e *PartialMatchError) Error() string {
	shown, more := e.Candidates, ""
	if len(shown) > MaxCandidates {
		more = fmt.Sprintf("\n  … and %d more", len(shown)-MaxCandidates)
		shown = shown[:MaxCandidates]
	}
	return fmt.Sprintf("no %s is named exactly %q; give the full name, login or id of the one you mean:\n  %s%s", e.What, e.Query, strings.Join(shown, "\n  "), more)
}

// FindAssignment picks exactly one assignment for query. A numeric query
// matches an assignment id or a quiz id.
func FindAssignment(query string, items []Assignment) (*Assignment, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, fmt.Errorf("assignment is required")
	}
	if id, err := strconv.ParseInt(q, 10, 64); err == nil {
		for i := range items {
			if items[i].ID == id {
				return &items[i], nil
			}
		}
		for i := range items {
			if items[i].QuizID != 0 && items[i].QuizID == id {
				return &items[i], nil
			}
		}
		return nil, &NoMatchError{What: "assignment", Query: query, Total: len(items)}
	}
	lq := strings.ToLower(q)
	var exact, partial []int
	for i := range items {
		switch {
		case equalFold(lq, items[i].Name):
			exact = append(exact, i)
		case containsFold(lq, items[i].Name):
			partial = append(partial, i)
		}
	}
	pick := exact
	if len(pick) == 0 {
		pick = partial
	}
	switch len(pick) {
	case 0:
		return nil, &NoMatchError{What: "assignment", Query: query, Total: len(items)}
	case 1:
		return &items[pick[0]], nil
	}
	var lines []string
	for _, i := range pick {
		lines = append(lines, DescribeAssignment(&items[i]))
	}
	sort.Strings(lines)
	return nil, &AmbiguousError{What: "assignment", Query: query, Candidates: lines}
}

func describeStudent(s *Student) string {
	line := fmt.Sprintf("%d %s", s.ID, s.Name)
	if s.SortableName != "" && s.SortableName != s.Name {
		line += " (" + s.SortableName + ")"
	}
	if s.LoginID != "" {
		line += " <" + s.LoginID + ">"
	}
	return line
}

// DescribeStudent renders a student as "id Name (Sortable, Name) <login>".
func DescribeStudent(s *Student) string { return describeStudent(s) }

// DescribeAssignment renders an assignment as "id Name [quiz 77]".
func DescribeAssignment(a *Assignment) string {
	line := fmt.Sprintf("%d %s", a.ID, a.Name)
	switch {
	case a.QuizID != 0:
		line += fmt.Sprintf(" [quiz %d]", a.QuizID)
	case a.Kind != "" && a.Kind != KindAssignment:
		line += " [" + a.Kind + "]"
	}
	return line
}

func equalFold(lq string, names ...string) bool {
	for _, n := range names {
		if n != "" && strings.ToLower(n) == lq {
			return true
		}
	}
	return false
}

func containsFold(lq string, names ...string) bool {
	for _, n := range names {
		if n != "" && strings.Contains(strings.ToLower(n), lq) {
			return true
		}
	}
	return false
}

// ListStudents reads the course's active student roster.
//
// Canvas API: GET /api/v1/courses/:course_id/users?enrollment_type[]=student
// &enrollment_state[]=active — the Test Student is not included unless
// include[]=test_student is asked for, and is excluded by name anyway.
// https://canvas.instructure.com/doc/api/courses.html#method.courses.users
func ListStudents(ctx context.Context, client *api.Client, courseID int64) ([]Student, error) {
	users, err := api.NewCoursesService(client).ListUsers(ctx, courseID, &api.ListCourseUsersOptions{
		EnrollmentType:  []string{"student"},
		EnrollmentState: []string{"active"},
		PerPage:         100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list students of course %d: %w", courseID, err)
	}
	students := make([]Student, 0, len(users))
	for i := range users {
		u := &users[i]
		if u.Name == "Test Student" {
			continue
		}
		students = append(students, Student{ID: u.ID, Name: u.Name, SortableName: u.SortableName, ShortName: u.ShortName, LoginID: u.LoginID, SISUserID: u.SisUserID})
	}
	return students, nil
}

// ListAssignments reads the course's assignments. Graded quizzes and
// discussions are included through their assignments, so a quiz name or
// quiz id resolves to the assignment Canvas grades it under. Practice
// quizzes and surveys have no assignment and are not listed.
//
// Canvas API: GET /api/v1/courses/:course_id/assignments — each quiz-backed
// assignment carries quiz_id and is_quiz_assignment.
// https://canvas.instructure.com/doc/api/assignments.html#method.assignments_api.index
func ListAssignments(ctx context.Context, client *api.Client, courseID int64) ([]Assignment, error) {
	list, err := api.NewAssignmentsService(client).List(ctx, courseID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list assignments of course %d: %w", courseID, err)
	}
	items := make([]Assignment, 0, len(list))
	for i := range list {
		a := &list[i]
		items = append(items, Assignment{ID: a.ID, Name: a.Name, Kind: assignmentKind(a), QuizID: a.QuizID})
	}
	return items, nil
}

func assignmentKind(a *api.Assignment) string {
	if a.IsQuizAssignment || a.QuizID > 0 {
		return KindQuiz
	}
	for _, st := range a.SubmissionTypes {
		switch st {
		case "online_quiz":
			return KindQuiz
		case "discussion_topic":
			return KindDiscussion
		}
	}
	return KindAssignment
}
