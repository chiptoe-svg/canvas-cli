package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/chiptoe-svg/canvas-cli/commands/internal/testing"
)

// ---------- groups.go coverage boosters ----------

// TestCovGroups_ListUserContext exercises the default (user context) branch of runGroupsList.
func TestCovGroups_ListUserContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups for current user",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/groups": cmdtest.NewMockResponse(`[
				{"id": 7, "name": "Self Group", "members_count": 2}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Self Group") {
				t.Errorf("Expected 'Self Group' in output, got: %s", output)
			}
		},
	}
	cmd := newGroupsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_ListCourseAPIError covers the error branch of runGroupsList when API call fails.
func TestCovGroups_ListCourseAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups for course - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/groups": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newGroupsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_ListAccountAPIError covers the error branch for account-scoped list.
func TestCovGroups_ListAccountAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups for account - API error",
		Args: []string{"--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/groups": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newGroupsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_ListWithIncludes exercises the include-users/include-permissions branches.
func TestCovGroups_ListWithIncludes(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups with includes",
		Args: []string{"--course-id", "1", "--include-users", "--include-permissions"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/groups": cmdtest.NewMockResponse(`[
				{"id": 1, "name": "Inclusive Group", "members_count": 3}
			]`),
		},
		ExpectError: false,
	}
	cmd := newGroupsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_GetWithIncludes exercises include flags on runGroupsGet.
func TestCovGroups_GetWithIncludes(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get group with includes",
		Args: []string{"10", "--include-users", "--include-permissions"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"name": "Full Group",
				"members_count": 5
			}`),
		},
		ExpectError: false,
	}
	cmd := newGroupsGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CreateAPIError covers error branch in runGroupsCreate.
func TestCovGroups_CreateAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create group - API error",
		Args: []string{"--category-id", "5", "--name", "Bad Group"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5/groups": cmdtest.NewErrorResponse(422, "invalid group"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_UpdateAPIError covers error branch in runGroupsUpdate.
func TestCovGroups_UpdateAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update group - API error",
		Args: []string{"10", "--name", "Bad Update"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10": cmdtest.NewErrorResponse(404, "group not found"),
		},
		ExpectError: true,
	}
	cmd := newGroupsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_UpdateMultipleFields exercises more of the update field branches.
func TestCovGroups_UpdateMultipleFields(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update group with multiple fields",
		Args: []string{"10", "--name", "Updated", "--description", "New desc", "--public", "--join-level", "invitation_only", "--storage-quota-mb", "100", "--sis-group-id", "SIS123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"name": "Updated"
			}`),
		},
		ExpectError: false,
	}
	cmd := newGroupsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_DeleteAPIError covers the API error path in runGroupsDelete.
func TestCovGroups_DeleteAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete group - API error",
		Args: []string{"10", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10": cmdtest.NewErrorResponse(404, "group not found"),
		},
		ExpectError: true,
	}
	cmd := newGroupsDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_MembersListAPIError covers error path in runGroupsMembersList.
func TestCovGroups_MembersListAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list group members - API error",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10/users": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newGroupsMembersListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_MembersAddAPIError covers error path in runGroupsMembersAdd.
func TestCovGroups_MembersAddAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "add member - API error",
		Args: []string{"10", "--user-id", "100"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10/memberships": cmdtest.NewErrorResponse(422, "already a member"),
		},
		ExpectError: true,
	}
	cmd := newGroupsMembersAddCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_MembersRemoveAPIError covers error path in runGroupsMembersRemove.
func TestCovGroups_MembersRemoveAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "remove member - API error",
		Args: []string{"10", "--membership-id", "55"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10/memberships/55": cmdtest.NewErrorResponse(404, "membership not found"),
		},
		ExpectError: true,
	}
	cmd := newGroupsMembersRemoveCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CategoriesListAccountContext covers the account branch of runGroupsCategoriesList.
func TestCovGroups_CategoriesListAccountContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list group categories by account",
		Args: []string{"--account-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/group_categories": cmdtest.NewMockResponse(`[
				{"id": 3, "name": "Account Category", "self_signup": "enabled"}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Account Category") {
				t.Errorf("Expected 'Account Category' in output")
			}
		},
	}
	cmd := newGroupsCategoriesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CategoriesListAPIError covers the error path in runGroupsCategoriesList.
func TestCovGroups_CategoriesListAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list group categories - API error",
		Args: []string{"--course-id", "1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/group_categories": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CategoriesGetAPIError covers the error path in runGroupsCategoriesGet.
func TestCovGroups_CategoriesGetAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get group category - API error",
		Args: []string{"99"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/99": cmdtest.NewErrorResponse(404, "category not found"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesGetCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CategoriesCreateAccountContext covers the account branch of runGroupsCategoriesCreate.
func TestCovGroups_CategoriesCreateAccountContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create group category in account",
		Args: []string{"--account-id", "1", "--name", "Account Category"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/group_categories": cmdtest.NewMockResponse(`{
				"id": 20,
				"name": "Account Category",
				"account_id": 1
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Account Category") {
				t.Errorf("Expected 'Account Category' in output")
			}
		},
	}
	cmd := newGroupsCategoriesCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CategoriesCreateAPIError covers the error path in runGroupsCategoriesCreate.
func TestCovGroups_CategoriesCreateAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create group category - API error",
		Args: []string{"--course-id", "1", "--name", "Bad Category"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1/group_categories": cmdtest.NewErrorResponse(422, "invalid category"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesCreateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CategoriesUpdateAPIError covers the error path in runGroupsCategoriesUpdate.
func TestCovGroups_CategoriesUpdateAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update group category - API error",
		Args: []string{"5", "--name", "Bad Update"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5": cmdtest.NewErrorResponse(404, "category not found"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CategoriesUpdateMultipleFields exercises more field branches in runGroupsCategoriesUpdate.
func TestCovGroups_CategoriesUpdateMultipleFields(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update group category with multiple fields",
		Args: []string{"5", "--name", "New Name", "--self-signup", "restricted", "--auto-leader", "first", "--group-limit", "5", "--sis-category-id", "SIS456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5": cmdtest.NewMockResponse(`{
				"id": 5,
				"name": "New Name"
			}`),
		},
		ExpectError: false,
	}
	cmd := newGroupsCategoriesUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CategoriesDeleteAPIError covers the API error path in runGroupsCategoriesDelete.
func TestCovGroups_CategoriesDeleteAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete group category - API error",
		Args: []string{"5", "--force"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5": cmdtest.NewErrorResponse(404, "category not found"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesDeleteCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CategoriesGroupsAPIError covers the error path in runGroupsCategoriesGroups.
func TestCovGroups_CategoriesGroupsAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups in category - API error",
		Args: []string{"5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5/groups": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmd := newGroupsCategoriesGroupsCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_CategoriesGroupsEmpty covers the empty result path.
func TestCovGroups_CategoriesGroupsEmpty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups in category - empty",
		Args: []string{"5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/group_categories/5/groups": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError:  false,
		ExpectOutput: "No groups found in category",
	}
	cmd := newGroupsCategoriesGroupsCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovGroups_MembersListEmpty covers the empty members path.
func TestCovGroups_MembersListEmpty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list group members - empty",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/10/users": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError:  false,
		ExpectOutput: "No members found in group",
	}
	cmd := newGroupsMembersListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// ---------- discussions.go coverage boosters ----------

// TestCovDisc_UpdateAllFields exercises all the opt.*Set branches in runDiscussionsUpdate.
func TestCovDisc_UpdateAllFields(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update discussion with all optional fields",
		Args: []string{
			"--course-id", "1", "10",
			"--title", "Full Update",
			"--message", "Updated message",
			"--type", "threaded",
			"--published",
			"--delayed-post-at", "2024-12-01T00:00:00Z",
			"--allow-rating",
			"--lock-at", "2024-12-31T00:00:00Z",
			"--require-initial-post",
			"--pinned",
			"--locked",
		},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/discussion_topics/10": cmdtest.NewMockResponse(`{
				"id": 10,
				"title": "Full Update",
				"message": "Updated message",
				"pinned": true,
				"locked": true
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Full Update") {
				t.Errorf("Expected 'Full Update' in output")
			}
		},
	}
	cmd := newDiscussionsUpdateCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovDisc_PostWithPositionalMessage exercises taking message from positional arg.
func TestCovDisc_PostWithPositionalMessage(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "post entry with positional message",
		Args: []string{"--course-id", "1", "10", "My positional message"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/courses/1": courseMock,
			"/api/v1/courses/1/discussion_topics/10/entries": cmdtest.NewMockResponse(`{
				"id": 55,
				"message": "My positional message",
				"user_id": 100
			}`),
		},
		ExpectError: false,
	}
	cmd := newDiscussionsPostCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

// TestCovDisc_ListUserGroups exercises default user-context list.
func TestCovGroups_ListUserAPIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list groups for user - API error",
		Args: []string{"--user-id", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/5/groups": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmd := newGroupsListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}
