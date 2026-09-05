package commands

import (
	"testing"

	cmdtest "github.com/chiptoe-svg/canvas-cli/commands/internal/testing"
)

// ---------------------------------------------------------------------------
// gradebook_history.go — entirely missing from new_resources_test.go
// ---------------------------------------------------------------------------

func TestGradesHistoryDayCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list graders for a specific date",
			Args: []string{"--course-id", "1", "--date", "2024-03-15"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/gradebook_history/2024-03-15": cmdtest.NewMockResponse(`[
					{"id":5,"name":"Prof. Smith"},
					{"id":6,"name":"TA Jones"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list graders - empty",
			Args: []string{"--course-id", "1", "--date", "2024-01-01"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/gradebook_history/2024-01-01": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError: false,
		},
		{
			Name:        "list graders - missing course-id",
			Args:        []string{"--date", "2024-03-15"},
			ExpectError: true,
		},
		{
			Name:        "list graders - missing date",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "list graders - API error",
			Args: []string{"--course-id", "1", "--date", "2024-03-15"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/gradebook_history/2024-03-15": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newGradesHistoryDayCmd(), tc)
		})
	}
}

func TestGradesHistorySubmissionsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list history submissions",
			Args: []string{"--course-id", "1", "--date", "2024-03-15", "--grader-id", "5", "--assignment-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/gradebook_history/2024-03-15/graders/5/assignments/10/submissions": cmdtest.NewMockResponse(`[
					{"id":101,"grade":"A"},
					{"id":102,"grade":"B"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list history submissions - empty",
			Args: []string{"--course-id", "1", "--date", "2024-03-15", "--grader-id", "5", "--assignment-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/gradebook_history/2024-03-15/graders/5/assignments/10/submissions": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No submissions found",
		},
		{
			Name:        "list history submissions - missing course-id",
			Args:        []string{"--date", "2024-03-15", "--grader-id", "5", "--assignment-id", "10"},
			ExpectError: true,
		},
		{
			Name:        "list history submissions - missing date",
			Args:        []string{"--course-id", "1", "--grader-id", "5", "--assignment-id", "10"},
			ExpectError: true,
		},
		{
			Name:        "list history submissions - missing grader-id",
			Args:        []string{"--course-id", "1", "--date", "2024-03-15", "--assignment-id", "10"},
			ExpectError: true,
		},
		{
			Name:        "list history submissions - missing assignment-id",
			Args:        []string{"--course-id", "1", "--date", "2024-03-15", "--grader-id", "5"},
			ExpectError: true,
		},
		{
			Name: "list history submissions - API error",
			Args: []string{"--course-id", "1", "--date", "2024-03-15", "--grader-id", "5", "--assignment-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/gradebook_history/2024-03-15/graders/5/assignments/10/submissions": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newGradesHistorySubmissionsCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// Additional error/edge cases for commands already partially tested
// ---------------------------------------------------------------------------

func TestCollaborationsListCmd_Extended(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list collaborations for group",
			Args: []string{"--group-id", "456"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/456/collaborations": cmdtest.NewMockResponse(`[{"id":3,"title":"Team Doc"}]`),
			},
			ExpectError: false,
		},
		{
			Name: "list collaborations - API error",
			Args: []string{"--course-id", "123"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/123/collaborations": cmdtest.NewErrorResponse(403, "unauthorized"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCollaborationsListCmd(), tc)
		})
	}
}

func TestCollaborationsMembersCmd_Extended(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name:        "invalid collaboration ID",
			Args:        []string{"notanumber"},
			ExpectError: true,
		},
		{
			Name: "members API error",
			Args: []string{"789"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/collaborations/789/members": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCollaborationsMembersCmd(), tc)
		})
	}
}
