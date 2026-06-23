package commands

import (
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// TestCovQuizzes_UpdateAllFields exercises every opt.*Set branch in runQuizzesUpdate.
// The runQuizzesUpdate function has 19 optional update fields; existing tests only set --title.
// This test sets most of them to drive coverage through all the if-block branches.
func TestCovQuizzes_UpdateAllFields(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update quiz with many fields",
		Args: []string{
			"10", "--course-id", "1",
			"--title", "Full Update",
			"--description", "Updated desc",
			"--quiz-type", "practice_quiz",
			"--assignment-group-id", "5",
			"--time-limit", "60",
			"--shuffle-answers",
			"--hide-results", "always",
			"--show-correct",
			"--scoring-policy", "keep_highest",
			"--attempts", "3",
			"--one-at-a-time",
			"--cant-go-back",
			"--access-code", "secret123",
			"--ip-filter", "192.168.1.0/24",
			"--due-at", "2024-12-01T00:00:00Z",
			"--lock-at", "2024-12-31T00:00:00Z",
			"--unlock-at", "2024-11-01T00:00:00Z",
			"--published",
			"--anonymous",
		},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/quizzes/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"title": "Full Update",
				"quiz_type": "practice_quiz",
				"published": true
			}`),
		},
		ExpectError: false,
	}
	cmd := newQuizzesUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovQuizzes_ListCourseValidationError covers the validateCourseID error path.
func TestCovQuizzes_ListCourseValidationError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list quizzes - course not found",
		Args: []string{"--course-id", "999"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/999": cmdtest.NewErrorResponse(404, "course not found"),
		},
		ExpectError: true,
	}
	cmd := newQuizzesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovQuizzes_UpdateMissingCourseID exercises the missing course-id validation path in update.
func TestCovQuizzes_UpdateMissingCourseID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "update quiz - missing course ID",
		Args:        []string{"10", "--title", "Update"},
		ExpectError: true,
	}
	cmd := newQuizzesUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovQuizzes_QuestionsDeleteNoForce exercises the "not confirmed" path without --force flag.
// The confirmation reads stdin which is wired to EOF by the test harness, returning an error.
func TestCovQuizzes_QuestionsDeleteNoForce(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete quiz question - without force (EOF stdin causes error)",
		Args: []string{"5", "--course-id", "1", "--quiz-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		// EOF on stdin causes confirmDeleteWithDetails to return an error.
		ExpectError: true,
	}
	cmd := newQuizzesQuestionsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovQuizzes_DeleteNoForce exercises the "not confirmed" path without --force for quiz delete.
func TestCovQuizzes_DeleteNoForce(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete quiz - without force (EOF stdin causes error)",
		Args: []string{"10", "--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
		},
		// EOF on stdin causes confirmDelete to return an error.
		ExpectError: true,
	}
	cmd := newQuizzesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovQuizzes_UpdateNoFields exercises missing-fields validation in update (with course+quiz but no set fields).
// This goes through cobra argument parsing and Validate() check.
func TestCovQuizzes_UpdateNoFields(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "update quiz - no update fields provided",
		Args:        []string{"10", "--course-id", "1"},
		ExpectError: true,
	}
	cmd := newQuizzesUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovQuizzes_SubmissionsListEmpty covers the empty submissions case.
func TestCovQuizzes_SubmissionsListEmpty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list quiz submissions - empty",
		Args: []string{"--course-id", "1", "--quiz-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/quizzes/10/submissions": cmdtest.NewMockResponse(`{
				"quiz_submissions": []
			}`),
		},
		ExpectError:  false,
		ExpectOutput: "No submissions found",
	}
	cmd := newQuizzesSubmissionsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
