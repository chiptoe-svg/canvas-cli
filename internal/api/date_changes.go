package api

import (
	"context"
	"fmt"
	"time"
)

// The three availability timestamps quizzes and assignments share.
// Canvas calls them unlock_at (available from), due_at and lock_at (closed).
const (
	DateFieldUnlockAt = "unlock_at"
	DateFieldDueAt    = "due_at"
	DateFieldLockAt   = "lock_at"
)

// DateChanges maps a date field to its new value. A nil value clears the
// date (sent as JSON null); a field that is absent is left untouched, so a
// change to due_at never disturbs unlock_at or lock_at.
type DateChanges map[string]*time.Time

// Payload builds the request body Canvas expects: {wrapper: {field: value}}
// where value is RFC 3339 in UTC or null. It is exported so a --dry-run can
// print exactly what would be sent.
func (d DateChanges) Payload(wrapper string) (map[string]interface{}, error) {
	fields := make(map[string]interface{}, len(d))
	for field, value := range d {
		switch field {
		case DateFieldUnlockAt, DateFieldDueAt, DateFieldLockAt:
		default:
			return nil, fmt.Errorf("unknown date field %q (expected %s, %s or %s)", field, DateFieldUnlockAt, DateFieldDueAt, DateFieldLockAt)
		}
		if value == nil {
			fields[field] = nil
			continue
		}
		fields[field] = value.UTC().Format(time.RFC3339)
	}
	return map[string]interface{}{wrapper: fields}, nil
}

// UpdateDates sets or clears a quiz's unlock_at / due_at / lock_at and
// returns the quiz as Canvas reports it after the write.
//
// Canvas API: PUT /api/v1/courses/:course_id/quizzes/:id with
// quiz[unlock_at], quiz[due_at], quiz[lock_at] (DateTime; null clears) —
// https://canvas.instructure.com/doc/api/quizzes.html#method.quizzes/quizzes_api.update
func (s *QuizzesService) UpdateDates(ctx context.Context, courseID, quizID int64, changes DateChanges) (*Quiz, error) {
	body, err := changes.Payload("quiz")
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/api/v1/courses/%d/quizzes/%d", courseID, quizID)
	var quiz Quiz
	if err := s.client.PutJSON(ctx, path, body, &quiz); err != nil {
		return nil, err
	}
	return &quiz, nil
}

// UpdateDates sets or clears an assignment's unlock_at / due_at / lock_at
// and returns the assignment as Canvas reports it after the write. Only the
// base dates change; overrides are untouched.
//
// Canvas API: PUT /api/v1/courses/:course_id/assignments/:id with
// assignment[unlock_at], assignment[due_at], assignment[lock_at] (DateTime;
// null clears) —
// https://canvas.instructure.com/doc/api/assignments.html#method.assignments_api.update
func (s *AssignmentsService) UpdateDates(ctx context.Context, courseID, assignmentID int64, changes DateChanges) (*Assignment, error) {
	body, err := changes.Payload("assignment")
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d", courseID, assignmentID)
	var result Assignment
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}
	return NormalizeAssignment(&result), nil
}
