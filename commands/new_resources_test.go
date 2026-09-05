package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// --- collaborations ---

func TestCollaborationsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list course collaborations",
			Args: []string{"--course-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/5/collaborations": cmdtest.NewMockResponse(`[{"id":1,"title":"Group Doc"}]`),
			},
			ExpectError: false,
		},
		{
			Name:        "missing course-id or group-id",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCollaborationsListCmd(), tc)
		})
	}
}

func TestCollaborationsMembersCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list collaboration members",
			Args: []string{"11"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/collaborations/11/members": cmdtest.NewMockResponse(`[{"id":100,"type":"user","name":"Alice"}]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Alice") {
					t.Error("expected 'Alice' in output")
				}
			},
		},
		{
			Name:        "missing collaboration id",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCollaborationsMembersCmd(), tc)
		})
	}
}
