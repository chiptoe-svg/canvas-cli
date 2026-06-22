package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// ---- Quiz Submissions: Create ----

func TestQuizzesSubmissionsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "start quiz submission successfully",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/submissions": cmdtest.NewMockResponse(`{
					"quiz_submissions": [{"id": 99, "quiz_id": 10, "workflow_state": "untaken"}]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "submission started") {
					t.Error("Expected 'submission started' in output")
				}
			},
		},
		{
			Name:        "start quiz submission - missing course ID",
			Args:        []string{"--quiz-id", "10"},
			ExpectError: true,
		},
		{
			Name:        "start quiz submission - missing quiz ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "start quiz submission - API error",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                        courseMock,
				"/api/v1/courses/1/quizzes/10/submissions": cmdtest.NewErrorResponse(403, "unauthorized"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesSubmissionsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

// ---- Quiz Question Groups ----

func TestQuizzesGroupsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get quiz group successfully",
			Args: []string{"3", "--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/groups/3": cmdtest.NewMockResponse(`{
					"quiz_groups": [{"id": 3, "quiz_id": 10, "name": "Chapter Questions"}]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Chapter Questions") {
					t.Error("Expected 'Chapter Questions' in output")
				}
			},
		},
		{
			Name:        "get quiz group - missing group ID",
			Args:        []string{"--course-id", "1", "--quiz-id", "10"},
			ExpectError: true,
		},
		{
			Name: "get quiz group - API error",
			Args: []string{"3", "--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                        courseMock,
				"/api/v1/courses/1/quizzes/10/groups/3":   cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesGroupsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesGroupsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create quiz group successfully",
			Args: []string{"--course-id", "1", "--quiz-id", "10", "--name", "Random Group", "--pick-count", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/groups": cmdtest.NewMockResponse(`{
					"quiz_groups": [{"id": 7, "quiz_id": 10, "name": "Random Group"}]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "group created") {
					t.Error("Expected 'group created' in output")
				}
			},
		},
		{
			Name:        "create quiz group - missing course ID",
			Args:        []string{"--quiz-id", "10", "--name", "G"},
			ExpectError: true,
		},
		{
			Name: "create quiz group - API error",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                   courseMock,
				"/api/v1/courses/1/quizzes/10/groups": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesGroupsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesGroupsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update quiz group successfully",
			Args: []string{"3", "--course-id", "1", "--quiz-id", "10", "--pick-count", "8"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/groups/3": cmdtest.NewMockResponse(`{
					"quiz_groups": [{"id": 3, "quiz_id": 10, "pick_count": 8}]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "group updated") {
					t.Error("Expected 'group updated' in output")
				}
			},
		},
		{
			Name:        "update quiz group - no fields",
			Args:        []string{"3", "--course-id", "1", "--quiz-id", "10"},
			ExpectError: true,
		},
		{
			Name: "update quiz group - API error",
			Args: []string{"3", "--course-id", "1", "--quiz-id", "10", "--name", "New Name"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                      courseMock,
				"/api/v1/courses/1/quizzes/10/groups/3":  cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesGroupsUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesGroupsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete quiz group with force",
			Args: []string{"3", "--course-id", "1", "--quiz-id", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                     courseMock,
				"/api/v1/courses/1/quizzes/10/groups/3": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete quiz group - missing group ID",
			Args:        []string{"--course-id", "1", "--quiz-id", "10", "--force"},
			ExpectError: true,
		},
		{
			Name: "delete quiz group - API error",
			Args: []string{"3", "--course-id", "1", "--quiz-id", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                     courseMock,
				"/api/v1/courses/1/quizzes/10/groups/3": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesGroupsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

// ---- Quiz Reports ----

func TestQuizzesReportsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list quiz reports successfully",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/reports": cmdtest.NewMockResponse(`[
					{"id": 5, "quiz_id": 10, "report_type": "student_analysis"}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "student_analysis") {
					t.Error("Expected 'student_analysis' in output")
				}
			},
		},
		{
			Name:        "list quiz reports - missing quiz ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "list quiz reports - API error",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                    courseMock,
				"/api/v1/courses/1/quizzes/10/reports": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesReportsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesReportsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get quiz report successfully",
			Args: []string{"5", "--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                       courseMock,
				"/api/v1/courses/1/quizzes/10/reports/5": cmdtest.NewMockResponse(`{
					"id": 5, "quiz_id": 10, "report_type": "item_analysis"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get quiz report - missing report ID",
			Args:        []string{"--course-id", "1", "--quiz-id", "10"},
			ExpectError: true,
		},
		{
			Name: "get quiz report - API error",
			Args: []string{"5", "--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                      courseMock,
				"/api/v1/courses/1/quizzes/10/reports/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesReportsGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesReportsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create quiz report successfully",
			Args: []string{"--course-id", "1", "--quiz-id", "10", "--report-type", "student_analysis"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/reports": cmdtest.NewMockResponse(`{
					"id": 8, "quiz_id": 10, "report_type": "student_analysis"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "report created") {
					t.Error("Expected 'report created' in output")
				}
			},
		},
		{
			Name:        "create quiz report - missing report type",
			Args:        []string{"--course-id", "1", "--quiz-id", "10"},
			ExpectError: true,
		},
		{
			Name: "create quiz report - API error",
			Args: []string{"--course-id", "1", "--quiz-id", "10", "--report-type", "student_analysis"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                    courseMock,
				"/api/v1/courses/1/quizzes/10/reports": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesReportsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesReportsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete quiz report with force",
			Args: []string{"5", "--course-id", "1", "--quiz-id", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                      courseMock,
				"/api/v1/courses/1/quizzes/10/reports/5": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete quiz report - missing report ID",
			Args:        []string{"--course-id", "1", "--quiz-id", "10", "--force"},
			ExpectError: true,
		},
		{
			Name: "delete quiz report - API error",
			Args: []string{"5", "--course-id", "1", "--quiz-id", "10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                      courseMock,
				"/api/v1/courses/1/quizzes/10/reports/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesReportsDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

// ---- Quiz Statistics ----

func TestQuizzesStatisticsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list quiz statistics successfully",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/statistics": cmdtest.NewMockResponse(`{
					"quiz_statistics": [{"id": 1, "quiz_id": 10}]
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "list quiz statistics - missing quiz ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "list quiz statistics - API error",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                        courseMock,
				"/api/v1/courses/1/quizzes/10/statistics":  cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesStatisticsListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

// ---- Quiz Extensions ----

func TestQuizzesExtensionsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create quiz extension successfully",
			Args: []string{"--course-id", "1", "--quiz-id", "10", "--user-id", "42", "--extra-time", "30"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/10/extensions": cmdtest.NewMockResponse(`{
					"quiz_extensions": [{"user_id": 42, "quiz_id": 10, "extra_time": 30}]
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "create quiz extension - missing user ID",
			Args:        []string{"--course-id", "1", "--quiz-id", "10", "--extra-time", "30"},
			ExpectError: true,
		},
		{
			Name: "create quiz extension - API error",
			Args: []string{"--course-id", "1", "--quiz-id", "10", "--user-id", "42"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                        courseMock,
				"/api/v1/courses/1/quizzes/10/extensions":  cmdtest.NewErrorResponse(422, "invalid user"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesExtensionsCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

// ---- Quiz IP Filters ----

func TestQuizzesIPFiltersListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list quiz IP filters successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/ip_filters": cmdtest.NewMockResponse(`{
					"quiz_ip_filters": [{"name": "Campus", "filter": "10.0.0.0/8"}]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Campus") {
					t.Error("Expected 'Campus' in output")
				}
			},
		},
		{
			Name:        "list quiz IP filters - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list quiz IP filters - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                    courseMock,
				"/api/v1/courses/1/quizzes/ip_filters": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesIPFiltersListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

// ---- Quiz Assignment Overrides ----

func TestQuizzesAssignmentOverridesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list quiz assignment overrides successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/assignment_overrides": cmdtest.NewMockResponse(`{
					"quiz_assignment_overrides": [{"quiz_id": "10"}]
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "list quiz assignment overrides - missing course ID",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list quiz assignment overrides - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                              courseMock,
				"/api/v1/courses/1/quizzes/assignment_overrides": cmdtest.NewErrorResponse(500, "error"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesAssignmentOverridesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestQuizzesAssignmentOverridesSetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "set quiz assignment overrides successfully",
			Args: []string{"--course-id", "1", "--quiz-id", "10", "--section-id", "5", "--due-at", "2026-08-01T00:00:00Z"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": courseMock,
				"/api/v1/courses/1/quizzes/assignment_overrides": cmdtest.NewMockResponse(`{
					"quiz_assignment_overrides": [{"quiz_id": "10"}]
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "set quiz assignment overrides - missing quiz ID",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "set quiz assignment overrides - API error",
			Args: []string{"--course-id", "1", "--quiz-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1":                              courseMock,
				"/api/v1/courses/1/quizzes/assignment_overrides": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newQuizzesAssignmentOverridesSetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
