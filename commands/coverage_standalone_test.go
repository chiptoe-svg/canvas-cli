package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// ---------------------------------------------------------------------------
// audit_logs.go — extend beyond the basic cases in new_resources_test.go
// ---------------------------------------------------------------------------

func TestAuditListCmd_AuthByUser(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "authentication events by user-id",
		Args: []string{"--type", "authentication", "--user-id", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/audit/authentication/users/42": cmdtest.NewMockResponse(`{"events":[
				{"id":"evt1","event_type":"login","user_id":42}
			]}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "evt1") {
				t.Error("expected event id in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newAuditListCmd(), tc)
}

func TestAuditListCmd_AuthByLogin(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "authentication events by login-id - empty",
		Args: []string{"--type", "authentication", "--login-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/audit/authentication/logins/10": cmdtest.NewMockResponse(`{"events":[]}`),
		},
		ExpectError:  false,
		ExpectOutput: "No audit events found",
	}
	cmdtest.RunCommandTest(t, newAuditListCmd(), tc)
}

func TestAuditListCmd_CourseByCourse(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course audit events by course-id",
		Args: []string{"--type", "course", "--course-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/audit/course/courses/123": cmdtest.NewMockResponse(`{"events":[{"id":"c1","event_type":"updated"}]}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newAuditListCmd(), tc)
}

func TestAuditListCmd_CourseByAccount(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course audit events by account-id",
		Args: []string{"--type", "course", "--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/audit/course/accounts/1": cmdtest.NewMockResponse(`{"events":[{"id":"c2","event_type":"created"}]}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newAuditListCmd(), tc)
}

func TestAuditListCmd_CourseNoContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "course type without context flag",
		Args:        []string{"--type", "course"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAuditListCmd(), tc)
}

func TestAuditListCmd_GradeChangeByAssignment(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "grade-change events by assignment-id",
		Args: []string{"--type", "grade-change", "--assignment-id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/audit/grade_change/assignments/456": cmdtest.NewMockResponse(`{"events":[]}`),
		},
		ExpectError:  false,
		ExpectOutput: "No audit events found",
	}
	cmdtest.RunCommandTest(t, newAuditListCmd(), tc)
}

func TestAuditListCmd_GradeChangeByGrader(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "grade-change events by grader-id",
		Args: []string{"--type", "grade-change", "--grader-id", "7"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/audit/grade_change/graders/7": cmdtest.NewMockResponse(`{"events":[]}`),
		},
		ExpectError:  false,
		ExpectOutput: "No audit events found",
	}
	cmdtest.RunCommandTest(t, newAuditListCmd(), tc)
}

func TestAuditListCmd_GradeChangeByStudent(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "grade-change events by student-id",
		Args: []string{"--type", "grade-change", "--student-id", "8"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/audit/grade_change/students/8": cmdtest.NewMockResponse(`{"events":[]}`),
		},
		ExpectError:  false,
		ExpectOutput: "No audit events found",
	}
	cmdtest.RunCommandTest(t, newAuditListCmd(), tc)
}

func TestAuditListCmd_GradeChangeAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "grade-change API error is propagated",
		Args: []string{"--type", "grade-change"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/audit/grade_change": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAuditListCmd(), tc)
}

// ---------------------------------------------------------------------------
// eportfolios.go — delete and pages are missing from new_resources_test.go
// ---------------------------------------------------------------------------

func TestEportfolioDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete eportfolio with force",
			Args: []string{"10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10": cmdtest.NewMockResponse(`{"id":10,"name":"Old Portfolio"}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete eportfolio - missing ID arg",
			Args:        []string{"--force"},
			ExpectError: true,
		},
		{
			Name:        "delete eportfolio - invalid ID",
			Args:        []string{"bad"},
			ExpectError: true,
		},
		{
			Name: "delete eportfolio - API error with force",
			Args: []string{"99", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/99": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEportfolioDeleteCmd(), tc)
		})
	}
}

func TestEportfolioPagesCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list eportfolio pages successfully",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10/pages": cmdtest.NewMockResponse(`[
					{"id":1,"name":"Introduction","position":1},
					{"id":2,"name":"Projects","position":2}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Introduction") {
					t.Error("expected 'Introduction' in output")
				}
			},
		},
		{
			Name: "list eportfolio pages - empty",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10/pages": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No pages found",
		},
		{
			Name:        "list eportfolio pages - missing ID arg",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name:        "list eportfolio pages - invalid ID",
			Args:        []string{"abc"},
			ExpectError: true,
		},
		{
			Name: "list eportfolio pages - API error",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10/pages": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEportfolioPagesCmd(), tc)
		})
	}
}

// Additional eportfolio get error case
func TestEportfolioGetCmd_Error(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name:        "get eportfolio - invalid ID string",
			Args:        []string{"notanid"},
			ExpectError: true,
		},
		{
			Name: "get eportfolio - API error",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEportfolioGetCmd(), tc)
		})
	}
}

// Additional eportfolio list error case
func TestEportfoliosListCmd_Error(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list eportfolios - API error",
		Args: []string{"--user-id", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/42/eportfolios": cmdtest.NewErrorResponse(403, "unauthorized"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newEportfoliosListCmd(), tc)
}

// ---------------------------------------------------------------------------
// epub_exports.go — get command missing from new_resources_test.go
// ---------------------------------------------------------------------------

func TestEpubExportGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get epub export successfully",
			Args: []string{"--course-id", "123", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/123/epub_exports/5": cmdtest.NewMockResponse(`{"id":5,"workflow_state":"generated"}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get epub export - missing ID arg",
			Args:        []string{"--course-id", "123"},
			ExpectError: true,
		},
		{
			Name:        "get epub export - invalid ID",
			Args:        []string{"--course-id", "123", "notanid"},
			ExpectError: true,
		},
		{
			Name:        "get epub export - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "get epub export - API error",
			Args: []string{"--course-id", "123", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/123/epub_exports/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEpubExportGetCmd(), tc)
		})
	}
}

// Additional epub export error cases
func TestEpubExportsListCmd_Error(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list epub exports - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/epub_exports": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newEpubExportsListCmd(), tc)
}

func TestEpubExportCreateCmd_Error(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create epub export - API error",
		Args: []string{"--course-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/123/epub_exports": cmdtest.NewErrorResponse(422, "unprocessable entity"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newEpubExportCreateCmd(), tc)
}

// ---------------------------------------------------------------------------
// jwts.go — refresh command missing from new_resources_test.go
// ---------------------------------------------------------------------------

func TestJWTRefreshCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "refresh JWT successfully",
			Args: []string{"--token", "eyJ.old.token"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/jwts/refresh": cmdtest.NewMockResponse(`{"token":"eyJ.new.token","expires_in":3600}`),
			},
			ExpectError:  false,
			ExpectOutput: "JWT refreshed successfully",
		},
		{
			Name:        "refresh JWT - missing token flag",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "refresh JWT - API error",
			Args: []string{"--token", "eyJ.expired.token"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/jwts/refresh": cmdtest.NewErrorResponse(401, "unauthorized"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newJWTRefreshCmd(), tc)
		})
	}
}

// Additional JWT create error case
func TestJWTCreateCmd_Error(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create JWT - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/jwts": cmdtest.NewErrorResponse(401, "unauthorized"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newJWTCreateCmd(), tc)
}

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
// progress.go — add extra error coverage beyond new_resources_test.go
// ---------------------------------------------------------------------------

func TestProgressGetCmd_Error(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get progress - API error",
		Args: []string{"42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/progress/42": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newProgressGetCmd(), tc)
}

func TestProgressCancelCmd_Extended(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name:        "cancel progress - missing ID arg",
			Args:        []string{"--force"},
			ExpectError: true,
		},
		{
			Name:        "cancel progress - invalid ID",
			Args:        []string{"abc"},
			ExpectError: true,
		},
		{
			Name: "cancel progress with force success",
			Args: []string{"42", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/progress/42": cmdtest.NewMockResponse(`{"id":42,"workflow_state":"failed"}`),
			},
			ExpectError:  false,
			ExpectOutput: "Progress job cancelled",
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newProgressCancelCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// Additional error/edge cases for commands already partially tested
// ---------------------------------------------------------------------------

func TestConferencesListCmd_Extended(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list conferences for group",
			Args: []string{"--group-id", "456"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/456/conferences": cmdtest.NewMockResponse(`[{"id":9,"title":"Group Call"}]`),
			},
			ExpectError: false,
		},
		{
			Name: "list conferences - API error",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conferences": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newConferencesListCmd(), tc)
		})
	}
}

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

func TestBrandGetCmd_Extended(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get brand variables for course",
			Args: []string{"--course-id", "123"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/123/brand_variables": cmdtest.NewMockResponse(`{"ic-brand-primary":"#FF0000"}`),
			},
			ExpectError: false,
		},
		{
			Name: "get brand variables - API error",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/brand_variables": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newBrandGetCmd(), tc)
		})
	}
}

func TestCommMessagesListCmd_Extended(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list comm messages with user-id",
			Args: []string{"--user-id", "42"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/comm_messages": cmdtest.NewMockResponse(`[{"id":3,"subject":"Grade Posted"}]`),
			},
			ExpectError: false,
		},
		{
			Name: "list comm messages - API error",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/comm_messages": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCommMessagesListCmd(), tc)
		})
	}
}

func TestHistoryListCmd_Extended(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list history - API error",
		Args: []string{"--user-id", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/42/history": cmdtest.NewErrorResponse(403, "unauthorized"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newHistoryListCmd(), tc)
}

func TestErrorReportCreateCmd_Extended(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create error report with all options",
			Args: []string{"--subject", "Login broken", "--comments", "Cannot log in", "--email", "user@example.com", "--severity", "medium"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/error_reports": cmdtest.NewMockResponse(`{"logged":true}`),
			},
			ExpectError: false,
		},
		{
			Name: "create error report - API error",
			Args: []string{"--subject", "Bug found"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/error_reports": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newErrorReportCreateCmd(), tc)
		})
	}
}
