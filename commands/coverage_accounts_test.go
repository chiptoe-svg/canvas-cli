package commands

// coverage_accounts_test.go — additional tests for accounts-domain commands.
// Focuses on error paths and branch coverage for groups, accounts, and
// related run* helpers that were below 80% in the pre-expansion baseline.

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// ---------------------------------------------------------------------------
// Accounts command error paths
// ---------------------------------------------------------------------------

func TestAccountsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list accounts - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newAccountsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAccountsListCmd_WithInclude(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list accounts with include flag",
		Args: []string{"--include", "lti_guid"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts": cmdtest.NewMockResponse(`[
				{
					"id": 1,
					"name": "Root Account",
					"workflow_state": "active",
					"lti_guid": "abc123"
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Root Account") {
				t.Error("Expected 'Root Account' in output")
			}
		},
	}
	cmd := newAccountsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAccountsListCmd_LongName(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list accounts with long name truncation",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts": cmdtest.NewMockResponse(`[
				{
					"id": 1,
					"name": "This Is A Very Long Account Name That Exceeds The Display Width Limit",
					"workflow_state": "active",
					"parent_account_id": 0,
					"root_account_id": 0
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Total: 1 account(s)") {
				t.Error("Expected account count in output")
			}
		},
	}
	cmd := newAccountsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAccountsListCmd_WithParentID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list accounts - account with parent ID",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts": cmdtest.NewMockResponse(`[
				{
					"id": 2,
					"name": "Sub Account",
					"workflow_state": "active",
					"parent_account_id": 1,
					"root_account_id": 1
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Sub Account") {
				t.Error("Expected 'Sub Account' in output")
			}
		},
	}
	cmd := newAccountsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAccountsGetCmd_WithSISID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get account with SIS account ID",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1": cmdtest.NewMockResponse(`{
				"id": 1,
				"name": "My Account",
				"workflow_state": "active",
				"sis_account_id": "SIS-001",
				"uuid": "uuid-123",
				"parent_account_id": 0,
				"root_account_id": 0,
				"default_time_zone": "UTC"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "SIS-001") {
				t.Error("Expected SIS account ID in output")
			}
		},
	}
	cmd := newAccountsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAccountsSubAccountsCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list sub-accounts - API error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/sub_accounts": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmd := newAccountsSubAccountsCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAccountsSubAccountsCmd_LongName(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list sub-accounts with long name truncation",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/sub_accounts": cmdtest.NewMockResponse(`[
				{
					"id": 2,
					"name": "This Is A Very Long Sub Account Name That Exceeds Display Width",
					"workflow_state": "active",
					"parent_account_id": 1,
					"root_account_id": 1
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Total: 1 sub-account(s)") {
				t.Error("Expected sub-account count in output")
			}
		},
	}
	cmd := newAccountsSubAccountsCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAccountsGetCmd_InvalidIDFormat(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "get account - non-numeric ID",
		Args:        []string{"abc"},
		ExpectError: true,
	}
	cmd := newAccountsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAccountsSubAccountsCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "sub accounts - non-numeric ID",
		Args:        []string{"xyz"},
		ExpectError: true,
	}
	cmd := newAccountsSubAccountsCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestAccountsSubAccountsCmd_MissingArg(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "sub accounts - missing required arg",
		Args:        []string{},
		ExpectError: true,
	}
	cmd := newAccountsSubAccountsCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// ---------------------------------------------------------------------------
// Groups command additional coverage — API error paths
// ---------------------------------------------------------------------------

func TestGroupsListCmd_UserContextError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups - user context API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/groups": cmdtest.NewErrorResponse(401, "unauthorized"),
		},
		ExpectError: true,
	}
	cmd := newGroupsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsListCmd_AccountContextError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups by account - API error",
		Args: []string{"--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/groups": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmd := newGroupsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsListCmd_WithIncludes(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups with include-users and include-permissions",
		Args: []string{"--course-id", "1", "--include-users", "--include-permissions"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/groups": cmdtest.NewMockResponse(`[
				{
					"id": 1,
					"name": "Group With Users",
					"members_count": 3
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Group With Users") {
				t.Error("Expected 'Group With Users' in output")
			}
		},
	}
	cmd := newGroupsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsGetCmd_WithIncludes(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get group with include-users and include-permissions",
		Args: []string{"10", "--include-users", "--include-permissions"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"name": "Group With Permissions",
				"members_count": 3
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Group With Permissions") {
				t.Error("Expected 'Group With Permissions' in output")
			}
		},
	}
	cmd := newGroupsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsDeleteCmd_Force_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete group - force - API error",
		Args: []string{"10", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10": cmdtest.NewErrorResponse(404, "group not found"),
		},
		ExpectError: true,
	}
	cmd := newGroupsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsMembersListCmd_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list group members - empty",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10/users": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
	}
	cmd := newGroupsMembersListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsCategoriesListCmd_CourseAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list group categories - course API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/group_categories": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsCategoriesListCmd_AccountAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list group categories by account - API error",
		Args: []string{"--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/group_categories": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsCategoriesCreateCmd_CourseAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create group category - course API error",
		Args: []string{"--course-id", "1", "--name", "Fail Category"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/group_categories": cmdtest.NewErrorResponse(422, "unprocessable"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsCategoriesCreateCmd_AccountAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create group category in account - API error",
		Args: []string{"--account-id", "1", "--name", "Fail Category"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/group_categories": cmdtest.NewErrorResponse(422, "unprocessable"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsCategoriesDeleteCmd_Force(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete group category - force",
		Args: []string{"5", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5": cmdtest.NewMockResponse(`{
				"id": 5,
				"name": "Deleted Category"
			}`),
		},
		ExpectError: false,
	}
	cmd := newGroupsCategoriesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsCategoriesDeleteCmd_ForceAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete group category - force - API error",
		Args: []string{"5", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5": cmdtest.NewErrorResponse(404, "category not found"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestGroupsCategoriesGroupsCmd_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups in category - empty",
		Args: []string{"5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5/groups": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
	}
	cmd := newGroupsCategoriesGroupsCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// ---------------------------------------------------------------------------
// Courses command additional coverage
// ---------------------------------------------------------------------------

func TestCoursesListCmd_AccountContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list courses in account context",
		Args: []string{"--account", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/courses": cmdtest.NewMockResponse(`[
				{
					"id": 10,
					"name": "Account Course",
					"course_code": "AC101",
					"workflow_state": "available"
				}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Account Course") {
				t.Error("Expected 'Account Course' in output")
			}
		},
	}
	cmd := newCoursesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestCoursesListCmd_AccountContextError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list courses in account context - API error",
		Args: []string{"--account", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/courses": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmd := newCoursesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestCoursesListCmd_UserContextError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list courses in user context - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newCoursesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestCoursesDeleteCmd_Force(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete course - conclude (force)",
		Args: []string{"1", "--event", "conclude", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": cmdtest.NewMockResponse(`{"delete": "true"}`),
		},
		ExpectError: false,
	}
	cmd := newCoursesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestCoursesDeleteCmd_PermanentForce(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete course - permanent delete (force)",
		Args: []string{"1", "--event", "delete", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": cmdtest.NewMockResponse(`{"delete": "true"}`),
		},
		ExpectError: false,
	}
	cmd := newCoursesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestCoursesDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete course - API error",
		Args: []string{"1", "--event", "conclude", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": cmdtest.NewErrorResponse(404, "course not found"),
		},
		ExpectError: true,
	}
	cmd := newCoursesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
