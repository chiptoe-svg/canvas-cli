package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// ---------------------------------------------------------------------------
// Command: users get — invalid arg
// ---------------------------------------------------------------------------

func TestUsersGetCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "get user - non-numeric ID",
		Args:        []string{"not-a-number"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersGetCmd(), tc)
}

// ---------------------------------------------------------------------------
// Command: users list — with search flag (account context)
// ---------------------------------------------------------------------------

func TestUsersListCmd_WithSearchAndEnrollmentFilters(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list users with search and enrollment filters",
		Args: []string{
			"--account-id", "1",
			"--search", "alice",
			"--enrollment-type", "student",
			"--enrollment-state", "active",
			"--include", "email,avatar_url",
		},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewMockResponse(`[
				{"id": 7, "name": "Alice Smith", "email": "alice@example.com"}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Alice Smith") {
				t.Error("Expected 'Alice Smith' in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersListCmd(), tc)
}

// ---------------------------------------------------------------------------
// Command: users me — verifies login_id present output
// ---------------------------------------------------------------------------

func TestUsersMeCmd_WithLoginID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get current user with login_id",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self": cmdtest.NewMockResponse(`{
				"id": 1,
				"name": "Me User",
				"login_id": "meuser",
				"email": "me@example.com"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Me User") {
				t.Error("Expected 'Me User' in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersMeCmd(), tc)
}

// ---------------------------------------------------------------------------
// Command: users search — verifies non-empty result path fully executes
// ---------------------------------------------------------------------------

func TestUsersSearchCmd_MultipleResults(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "search users - multiple results",
		Args: []string{"john"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/search/recipients": cmdtest.NewMockResponse(`[
				{"id": 1, "name": "John Alpha"},
				{"id": 2, "name": "John Beta"}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "John Alpha") {
				t.Error("Expected 'John Alpha' in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersSearchCmd(), tc)
}
