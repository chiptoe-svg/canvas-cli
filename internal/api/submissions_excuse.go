package api

import (
	"context"
	"fmt"
)

// SetExcused excuses (or un-excuses) a student from an assignment and
// returns the submission as Canvas reports it after the write. Unlike
// Grade, an explicit false is sent, so the call can clear an excusal.
//
// Canvas API: PUT /api/v1/courses/:course_id/assignments/:assignment_id/submissions/:user_id
// with submission[excuse] (boolean: "Sets the 'excused' status of an
// assignment") —
// https://canvas.instructure.com/doc/api/submissions.html#method.submissions_api.update
func (s *SubmissionsService) SetExcused(ctx context.Context, courseID, assignmentID, userID int64, excused bool) (*Submission, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/assignments/%d/submissions/%d", courseID, assignmentID, userID)
	body := map[string]interface{}{
		"submission": map[string]interface{}{"excuse": excused},
	}
	var result Submission
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}
	return NormalizeSubmission(&result), nil
}
