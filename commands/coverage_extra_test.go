package commands

// coverage_extra_test.go — targeted tests to push commands package coverage above 80%.
// All test names here are unique — they cover branches not reached by existing test files.

import (
	"testing"

	cmdtest "github.com/chiptoe-svg/canvas-cli/commands/internal/testing"
	"github.com/chiptoe-svg/canvas-cli/internal/diagnostics"
)

// ---------------------------------------------------------------------------
// doctor.go — getStatusIcon (33% → covers Warning, Skipped, default branches)
// ---------------------------------------------------------------------------

func TestGetStatusIcon(t *testing.T) {
	cases := []struct {
		status diagnostics.CheckStatus
		want   string
	}{
		{diagnostics.StatusPass, "✓"},
		{diagnostics.StatusFail, "✗"},
		{diagnostics.StatusWarning, "⚠"},
		{diagnostics.StatusSkipped, "○"},
		{diagnostics.CheckStatus("unknown"), "?"},
	}

	for _, tc := range cases {
		got := getStatusIcon(tc.status)
		if got != tc.want {
			t.Errorf("getStatusIcon(%q) = %q, want %q", tc.status, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// enrollments.go — runEnrollmentsConclude (50% → covers conclude/deactivate tasks)
// The existing TestEnrollmentsConcludeCmd only covers the delete+force path.
// ---------------------------------------------------------------------------

func TestEnrollmentsConcludeCmd_ConcludeTask(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "conclude enrollment (default task)",
		Args: []string{"--course-id", "1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/enrollments/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"enrollment_state": "completed"
			}`),
		},
		ExpectError: false,
	}
	cmd := newEnrollmentsConcludeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestEnrollmentsConcludeCmd_DeactivateTask(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "deactivate enrollment",
		Args: []string{"--course-id", "1", "--task", "deactivate", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/enrollments/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"enrollment_state": "inactive"
			}`),
		},
		ExpectError: false,
	}
	cmd := newEnrollmentsConcludeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestEnrollmentsConcludeCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "conclude enrollment - API error",
		Args: []string{"--course-id", "1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                courseMock,
			"/api/v1/courses/1/enrollments/10": cmdtest.NewErrorResponse(404, "enrollment not found"),
		},
		ExpectError: true,
	}
	cmd := newEnrollmentsConcludeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// ---------------------------------------------------------------------------
// analytics.go — runAnalyticsUser (branches for type=assignments, communication)
// and runAnalyticsDepartment error and account context paths.
// The analytics_extra_test.go covers activity/assignments/communication
// but we need the API error path for user analytics.
// ---------------------------------------------------------------------------

func TestAnalyticsUserCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "analytics user - API error",
		Args: []string{"100", "--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/analytics/users/100/activity": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmd := newAnalyticsUserCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// ---------------------------------------------------------------------------
// announcements.go — API error paths (not covered by existing announcements_test.go)
// ---------------------------------------------------------------------------

func TestAnnouncementsGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get announcement - API error",
		Args: []string{"--course-id", "1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/discussion_topics/5": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmd := newAnnouncementsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAnnouncementsUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update announcement - API error",
		Args: []string{"--course-id", "1", "5", "--title", "Updated Title"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/discussion_topics/5": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newAnnouncementsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAnnouncementsCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create announcement - API error",
		Args: []string{"--course-id", "1", "--title", "My Announcement"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/discussion_topics": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newAnnouncementsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAnnouncementsDeleteCmd_Force_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete announcement - force - API error",
		Args: []string{"--course-id", "1", "5", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                     courseMock,
			"/api/v1/courses/1/discussion_topics/5": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmd := newAnnouncementsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAnnouncementsDeleteCmd_Force_Success(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete announcement - force - success",
		Args: []string{"--course-id", "1", "5", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                     courseMock,
			"/api/v1/courses/1/discussion_topics/5": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: false,
	}
	cmd := newAnnouncementsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// ---------------------------------------------------------------------------
// outcomes.go — runOutcomesGroupsList API error path (59.1%)
// ---------------------------------------------------------------------------

func TestOutcomesGroupsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list outcome groups - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                courseMock,
			"/api/v1/courses/1/outcome_groups": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newOutcomesGroupsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// ---------------------------------------------------------------------------
// sections.go — runSectionsUpdate (69.6%) - success and error paths
// Note: sections update uses /api/v1/sections/:id (not /courses/:id/sections/:id)
// ---------------------------------------------------------------------------

func TestSectionsUpdateCmd_Success(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update section - success",
		Args: []string{"10", "--name", "Updated Section"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/sections/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"name": "Updated Section",
				"course_id": 1
			}`),
		},
		ExpectError: false,
	}
	cmd := newSectionsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestSectionsUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update section - API error",
		Args: []string{"10", "--name", "Updated Section"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/sections/10": cmdtest.NewErrorResponse(404, "section not found"),
		},
		ExpectError: true,
	}
	cmd := newSectionsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// ---------------------------------------------------------------------------
// pages.go — runPagesDelete (58.8%) - force success, error, and missing arg paths
// ---------------------------------------------------------------------------

func TestPagesDeleteCmd_Force(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete page successfully (force)",
		Args: []string{"--course-id", "1", "home-page", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                 courseMock,
			"/api/v1/courses/1/pages/home-page": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: false,
	}
	cmd := newPagesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestPagesDeleteCmd_ForceAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete page - force - API error",
		Args: []string{"--course-id", "1", "home-page", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                 courseMock,
			"/api/v1/courses/1/pages/home-page": cmdtest.NewErrorResponse(404, "page not found"),
		},
		ExpectError: true,
	}
	cmd := newPagesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestPagesDeleteCmd_MissingCourseID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "delete page - missing course ID",
		Args:        []string{"home-page", "--force"},
		ExpectError: true,
	}
	cmd := newPagesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestPagesDeleteCmd_MissingPageID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "delete page - missing page URL/ID",
		Args:        []string{"--course-id", "1", "--force"},
		ExpectError: true,
	}
	cmd := newPagesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// ---------------------------------------------------------------------------
// assignment_groups.go — runAssignmentGroupsDelete (69.6%)
// ---------------------------------------------------------------------------

func TestAssignmentGroupsDeleteCmd_Force_Success(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete assignment group - force - success",
		Args: []string{"5", "--course-id", "1", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/assignment_groups/5": cmdtest.NewMockResponse(`{
				"id": 5,
				"name": "Old Assignment Group"
			}`),
		},
		ExpectError: false,
	}
	cmd := newAssignmentGroupsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAssignmentGroupsDeleteCmd_Force_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete assignment group - force - API error",
		Args: []string{"5", "--course-id", "1", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1":                     courseMock,
			"/api/v1/courses/1/assignment_groups/5": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmd := newAssignmentGroupsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// ---------------------------------------------------------------------------
// skills.go — newSkillsPathCmd RunE (50%)
// ---------------------------------------------------------------------------

func TestSkillsPathCmd_ViaConstructor(t *testing.T) {
	cmd := newSkillsPathCmd()
	// RunE just calls runSkillsPath which prints the path — just verify it exists
	if cmd.RunE == nil && cmd.Run == nil {
		t.Error("expected RunE or Run to be set")
	}
}
