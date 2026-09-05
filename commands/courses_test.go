package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/chiptoe-svg/canvas-cli/commands/internal/testing"
)

func TestCoursesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list courses successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "Test Course",
						"course_code": "TEST101",
						"workflow_state": "available",
						"uuid": "test-uuid"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Test Course") {
					t.Error("Expected 'Test Course' in output")
				}
				if !strings.Contains(output, "TEST101") {
					t.Error("Expected 'TEST101' in output")
				}
			},
		},
		{
			Name: "list courses with enrollment filter",
			Args: []string{"--enrollment-type", "teacher"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses": cmdtest.NewMockResponse(`[
					{
						"id": 1,
						"name": "Teacher Course",
						"course_code": "TEACH101",
						"workflow_state": "available"
					}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Teacher Course") {
					t.Error("Expected 'Teacher Course' in output")
				}
			},
		},
		{
			Name: "list courses - empty response",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No courses found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newCoursesListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestCoursesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get course successfully",
			Args: []string{"1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": cmdtest.NewMockResponse(`{
					"id": 1,
					"name": "Test Course",
					"course_code": "TEST101",
					"workflow_state": "available",
					"start_at": "2024-01-01T00:00:00Z",
					"end_at": "2024-12-31T23:59:59Z",
					"default_view": "modules"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Test Course") {
					t.Error("Expected 'Test Course' in output")
				}
				if !strings.Contains(output, "TEST101") {
					t.Error("Expected 'TEST101' in output")
				}
			},
		},
		{
			Name: "get course - not found",
			Args: []string{"999"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/999": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
		{
			Name:        "get course - invalid ID",
			Args:        []string{"invalid"},
			ExpectError: true,
		},
		{
			Name:        "get course - missing ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newCoursesGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestCoursesUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update course name",
			Args: []string{"1", "--name", "Updated Course"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1": cmdtest.NewMockResponse(`{
					"id": 1,
					"name": "Updated Course",
					"course_code": "TEST101",
					"workflow_state": "available"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Course") {
					t.Error("Expected 'Updated Course' in output")
				}
			},
		},
		{
			Name:        "update course - missing ID",
			Args:        []string{"--name", "Updated"},
			ExpectError: true,
		},
		{
			Name:        "update course - no changes",
			Args:        []string{"1"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newCoursesUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
