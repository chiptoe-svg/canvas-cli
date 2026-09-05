package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/chiptoe-svg/canvas-cli/commands/internal/testing"
)

func TestAppointmentGroupListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list appointment groups successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups": cmdtest.NewMockResponse(`[
					{"id":1,"title":"Office Hours","workflow_state":"active","appointments_count":3},
					{"id":2,"title":"Project Review","workflow_state":"active","appointments_count":2}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Office Hours") {
					t.Error("expected 'Office Hours' in output")
				}
			},
		},
		{
			Name: "list appointment groups with scope",
			Args: []string{"--scope", "manageable"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups": cmdtest.NewMockResponse(`[{"id":1,"title":"My Group","workflow_state":"active"}]`),
			},
			ExpectError: false,
		},
		{
			Name: "list appointment groups - empty",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No appointment groups found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAppointmentGroupListCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAppointmentGroupGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get appointment group successfully",
			Args: []string{"543"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups/543": cmdtest.NewMockResponse(`{
					"id":543,
					"title":"Final Presentation",
					"description":"Meet with professor",
					"workflow_state":"active",
					"context_codes":["course_123"]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Final Presentation") {
					t.Error("expected 'Final Presentation' in output")
				}
			},
		},
		{
			Name:        "get appointment group - missing ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAppointmentGroupGetCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAppointmentGroupCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create appointment group successfully",
			Args: []string{"--context", "course_123", "--title", "Office Hours"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups": cmdtest.NewMockResponse(`{
					"id":10,
					"title":"Office Hours",
					"workflow_state":"pending",
					"context_codes":["course_123"]
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Office Hours") {
					t.Error("expected 'Office Hours' in output")
				}
			},
		},
		{
			Name:        "create appointment group - missing title",
			Args:        []string{"--context", "course_123"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAppointmentGroupCreateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAppointmentGroupUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update appointment group successfully",
			Args: []string{"543", "--title", "Updated Group"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups/543": cmdtest.NewMockResponse(`{
					"id":543,
					"title":"Updated Group",
					"workflow_state":"active"
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Updated Group") {
					t.Error("expected 'Updated Group' in output")
				}
			},
		},
		{
			Name:        "update appointment group - missing ID",
			Args:        []string{"--title", "x"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAppointmentGroupUpdateCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAppointmentGroupDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete appointment group with force",
			Args: []string{"543", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups/543": cmdtest.NewMockResponse(`{"id":543,"title":"Deleted Group","workflow_state":"deleted"}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete appointment group - missing ID",
			Args:        []string{"--force"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAppointmentGroupDeleteCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAppointmentGroupUsersCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list appointment group users successfully",
			Args: []string{"543"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups/543/users": cmdtest.NewMockResponse(`[
					{"id":42,"name":"Alice"},
					{"id":43,"name":"Bob"}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Alice") {
					t.Error("expected 'Alice' in output")
				}
			},
		},
		{
			Name: "list appointment group users - empty",
			Args: []string{"543"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups/543/users": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No users found",
		},
		{
			Name:        "list appointment group users - missing ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAppointmentGroupUsersCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAppointmentGroupGroupsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list appointment group student groups successfully",
			Args: []string{"543"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups/543/groups": cmdtest.NewMockResponse(`[{"id":5,"name":"Team Alpha"}]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Team Alpha") {
					t.Error("expected 'Team Alpha' in output")
				}
			},
		},
		{
			Name:        "list appointment group student groups - missing ID",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAppointmentGroupGroupsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestAppointmentGroupNextCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get next appointment successfully",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups/next_appointment": cmdtest.NewMockResponse(`[{
					"id":99,
					"title":"Next Slot",
					"context_code":"course_123",
					"workflow_state":"active"
				}]`),
			},
			ExpectError: false,
		},
		{
			Name: "get next appointment - none available",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/appointment_groups/next_appointment": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No upcoming appointments available",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newAppointmentGroupNextCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
