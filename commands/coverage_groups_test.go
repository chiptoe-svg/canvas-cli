package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/chiptoe-svg/canvas-cli/commands/internal/testing"
)

// groupCategoryAPIErr is a reusable 404 mock for category-not-found scenarios.
var groupCategoryAPIErr = cmdtest.NewErrorResponse(404, "category not found")

// ---------------------------------------------------------------------------
// groups list – error path and user-context branch
// ---------------------------------------------------------------------------

func TestGroupsListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups - API error",
		Args: []string{"--course-id", "99"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/99/groups": cmdtest.NewErrorResponse(500, "internal error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsListCmd(), tc)
}

func TestGroupsListCmd_UserContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups for current user (no course/account)",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/groups": cmdtest.NewMockResponse(`[
				{"id": 7, "name": "My Group", "members_count": 2}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "My Group") {
				t.Error("expected 'My Group' in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newGroupsListCmd(), tc)
}

func TestGroupsListCmd_EmptyUserGroups(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:          "list groups for current user - empty",
		Args:          []string{},
		MockResponses: map[string]cmdtest.MockResponse{"/api/v1/users/self/groups": cmdtest.NewMockResponse(`[]`)},
		ExpectError:   false,
		ExpectOutput:  "No groups found",
	}
	cmdtest.RunCommandTest(t, newGroupsListCmd(), tc)
}

func TestGroupsListCmd_IncludeFlags(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups with include flags",
		Args: []string{"--course-id", "1", "--include-users", "--include-permissions"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/groups": cmdtest.NewMockResponse(`[{"id": 1, "name": "G1"}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newGroupsListCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups get – error paths
// ---------------------------------------------------------------------------

func TestGroupsGetCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "get group - non-numeric ID",
		Args:        []string{"abc"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsGetCmd(), tc)
}

func TestGroupsGetCmd_IncludeFlags(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get group with include flags",
		Args: []string{"10", "--include-users", "--include-permissions"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10": cmdtest.NewMockResponse(`{"id": 10, "name": "Team A"}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newGroupsGetCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups create – error path
// ---------------------------------------------------------------------------

func TestGroupsCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create group - API error",
		Args: []string{"--category-id", "5", "--name", "Fail Group"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5/groups": cmdtest.NewErrorResponse(422, "unprocessable entity"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCreateCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups update – error paths and additional update opts
// ---------------------------------------------------------------------------

func TestGroupsUpdateCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "update group - non-numeric ID",
		Args:        []string{"abc", "--name", "New"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsUpdateCmd(), tc)
}

func TestGroupsUpdateCmd_NoFieldsProvided(t *testing.T) {
	// Providing the required group ID but no --name/--description/etc
	// The options Validate() should return an error.
	tc := cmdtest.CommandTestCase{
		Name:        "update group - no fields set",
		Args:        []string{"10"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsUpdateCmd(), tc)
}

func TestGroupsUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update group - API error",
		Args: []string{"10", "--name", "Fail"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10": cmdtest.NewErrorResponse(500, "internal error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsUpdateCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups delete – additional paths
// ---------------------------------------------------------------------------

func TestGroupsDeleteCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "delete group - non-numeric ID",
		Args:        []string{"abc", "--force"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsDeleteCmd(), tc)
}

func TestGroupsDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete group - API error",
		Args: []string{"99", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/99": cmdtest.NewErrorResponse(404, "group not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsDeleteCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups members list – error paths
// ---------------------------------------------------------------------------

func TestGroupsMembersListCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "list group members - non-numeric ID",
		Args:        []string{"abc"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsMembersListCmd(), tc)
}

func TestGroupsMembersListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list group members - API error",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10/users": cmdtest.NewErrorResponse(404, "group not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsMembersListCmd(), tc)
}

func TestGroupsMembersListCmd_EmptyMembers(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list group members - empty",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10/users": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError:  false,
		ExpectOutput: "No members found in group",
	}
	cmdtest.RunCommandTest(t, newGroupsMembersListCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups members add – error paths
// ---------------------------------------------------------------------------

func TestGroupsMembersAddCmd_InvalidGroupID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "add member - non-numeric group ID",
		Args:        []string{"abc", "--user-id", "5"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsMembersAddCmd(), tc)
}

func TestGroupsMembersAddCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "add member - API error",
		Args: []string{"10", "--user-id", "99"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10/memberships": cmdtest.NewErrorResponse(422, "already a member"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsMembersAddCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups members remove – error paths
// ---------------------------------------------------------------------------

func TestGroupsMembersRemoveCmd_InvalidGroupID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "remove member - non-numeric group ID",
		Args:        []string{"abc", "--membership-id", "5"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsMembersRemoveCmd(), tc)
}

func TestGroupsMembersRemoveCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "remove member - API error",
		Args: []string{"10", "--membership-id", "55"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10/memberships/55": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsMembersRemoveCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups categories list – error paths and account branch
// ---------------------------------------------------------------------------

func TestGroupsCategoriesListCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list group categories - API error",
		Args: []string{"--course-id", "99"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/99/group_categories": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesListCmd(), tc)
}

func TestGroupsCategoriesListCmd_NoContextError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "list group categories - no course or account",
		Args:        []string{},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesListCmd(), tc)
}

func TestGroupsCategoriesListCmd_EmptyCategories(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:          "list group categories - empty course",
		Args:          []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{"/api/v1/courses/1/group_categories": cmdtest.NewMockResponse(`[]`)},
		ExpectError:   false,
		ExpectOutput:  "No group categories found",
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesListCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups categories get – error paths
// ---------------------------------------------------------------------------

func TestGroupsCategoriesGetCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "get group category - non-numeric ID",
		Args:        []string{"abc"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesGetCmd(), tc)
}

func TestGroupsCategoriesGetCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get group category - API error",
		Args: []string{"99"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/99": groupCategoryAPIErr,
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesGetCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups categories create – account-context branch
// ---------------------------------------------------------------------------

func TestGroupsCategoriesCreateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create group category - API error",
		Args: []string{"--course-id", "1", "--name", "Fail Category"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/group_categories": cmdtest.NewErrorResponse(422, "invalid"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesCreateCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups categories update – error paths
// ---------------------------------------------------------------------------

func TestGroupsCategoriesUpdateCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "update group category - non-numeric ID",
		Args:        []string{"abc", "--name", "New"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesUpdateCmd(), tc)
}

func TestGroupsCategoriesUpdateCmd_NoFieldsProvided(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "update group category - no fields set",
		Args:        []string{"5"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesUpdateCmd(), tc)
}

func TestGroupsCategoriesUpdateCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update group category - API error",
		Args: []string{"5", "--name", "Fail"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesUpdateCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups categories delete – error paths
// ---------------------------------------------------------------------------

func TestGroupsCategoriesDeleteCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "delete group category - non-numeric ID",
		Args:        []string{"abc", "--force"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesDeleteCmd(), tc)
}

func TestGroupsCategoriesDeleteCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete group category - API error",
		Args: []string{"99", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/99": groupCategoryAPIErr,
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesDeleteCmd(), tc)
}

// ---------------------------------------------------------------------------
// groups categories groups – error paths
// ---------------------------------------------------------------------------

func TestGroupsCategoriesGroupsCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "list groups in category - non-numeric ID",
		Args:        []string{"abc"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesGroupsCmd(), tc)
}

func TestGroupsCategoriesGroupsCmd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups in category - API error",
		Args: []string{"99"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/99/groups": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesGroupsCmd(), tc)
}

func TestGroupsCategoriesGroupsCmd_EmptyGroups(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:          "list groups in category - empty",
		Args:          []string{"5"},
		MockResponses: map[string]cmdtest.MockResponse{"/api/v1/group_categories/5/groups": cmdtest.NewMockResponse(`[]`)},
		ExpectError:   false,
		ExpectOutput:  "No groups found in category",
	}
	cmdtest.RunCommandTest(t, newGroupsCategoriesGroupsCmd(), tc)
}
