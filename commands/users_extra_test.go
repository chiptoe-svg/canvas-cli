package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestUsersSearchCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "search users - API error",
		// search-term is a positional argument
		Args: []string{"John"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/search/recipients": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newUsersSearchCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersSearchCmd_EmptyResults(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "search users - no results",
		Args: []string{"nonexistent"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/search/recipients": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No users found") {
				t.Error("Expected 'No users found' in output")
			}
		},
	}
	cmd := newUsersSearchCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersMeCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get current user - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/profile": cmdtest.NewErrorResponse(401, "unauthorized"),
		},
		ExpectError: true,
	}
	cmd := newUsersMeCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestUsersListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list users - API error",
		Args: []string{"--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newUsersListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
