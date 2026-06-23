package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// =============================================================================
// users_extra.go — users profile / settings / page-views / logins / courses /
//                  missing-submissions / activity-stream / todo / upcoming-events
//                  / merge / split
// =============================================================================

func TestCovUsers_Profile_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users profile - happy path",
		Args: []string{"123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/profile": cmdtest.NewMockResponse(`{"id":123,"name":"Alice","short_name":"Alice","bio":"","avatar_url":""}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Alice") {
				t.Errorf("expected 'Alice' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersProfileCmd(), tc)
}

func TestCovUsers_Profile_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users profile - API error",
		Args: []string{"123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/profile": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersProfileCmd(), tc)
}

func TestCovUsers_Profile_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "users profile - invalid ID",
		Args:        []string{"not-a-number"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersProfileCmd(), tc)
}

func TestCovUsers_Settings_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users settings - happy path",
		Args: []string{"42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/42/settings": cmdtest.NewMockResponse(`{"manual_mark_as_read":false,"collapse_global_nav":false}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUsersSettingsCmd(), tc)
}

func TestCovUsers_Settings_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users settings - API error",
		Args: []string{"42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/42/settings": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersSettingsCmd(), tc)
}

func TestCovUsers_Settings_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "users settings - invalid ID",
		Args:        []string{"abc"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersSettingsCmd(), tc)
}

func TestCovUsers_UpdateSettings_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users update-settings - happy path",
		Args: []string{"50", "--manual-mark-as-read", "--collapse-global-nav"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/50/settings": cmdtest.NewMockResponse(`{"manual_mark_as_read":true,"collapse_global_nav":true}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUsersUpdateSettingsCmd(), tc)
}

func TestCovUsers_UpdateSettings_NoFlags(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users update-settings - no optional flags",
		Args: []string{"50"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/50/settings": cmdtest.NewMockResponse(`{"manual_mark_as_read":false,"collapse_global_nav":false}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUsersUpdateSettingsCmd(), tc)
}

func TestCovUsers_UpdateSettings_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users update-settings - API error",
		Args: []string{"50"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/50/settings": cmdtest.NewErrorResponse(500, "internal server error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersUpdateSettingsCmd(), tc)
}

func TestCovUsers_UpdateSettings_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "users update-settings - invalid ID",
		Args:        []string{"bad"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersUpdateSettingsCmd(), tc)
}

func TestCovUsers_PageViews_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users page-views - happy path",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/10/page_views": cmdtest.NewMockResponse(`[{"id":"abc","url":"https://example.com","http_method":"GET"}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUsersPageViewsCmd(), tc)
}

func TestCovUsers_PageViews_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users page-views - empty result",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/10/page_views": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No page views found") {
				t.Errorf("expected 'No page views found' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersPageViewsCmd(), tc)
}

func TestCovUsers_PageViews_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users page-views - API error",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/10/page_views": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersPageViewsCmd(), tc)
}

func TestCovUsers_PageViews_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "users page-views - invalid ID",
		Args:        []string{"x"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersPageViewsCmd(), tc)
}

func TestCovUsers_Logins_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users logins - happy path",
		Args: []string{"20"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/20/logins": cmdtest.NewMockResponse(`[{"id":1,"unique_id":"alice@example.com"}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUsersLoginsCmd(), tc)
}

func TestCovUsers_Logins_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users logins - empty result",
		Args: []string{"20"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/20/logins": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No logins found") {
				t.Errorf("expected 'No logins found' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersLoginsCmd(), tc)
}

func TestCovUsers_Logins_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users logins - API error",
		Args: []string{"20"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/20/logins": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersLoginsCmd(), tc)
}

func TestCovUsers_Logins_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "users logins - invalid ID",
		Args:        []string{"!"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersLoginsCmd(), tc)
}

func TestCovUsers_Courses_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users courses - happy path",
		Args: []string{"30"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/30/courses": cmdtest.NewMockResponse(`[{"id":1,"name":"Intro to Go"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Intro to Go") {
				t.Errorf("expected course name in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersCoursesCmd(), tc)
}

func TestCovUsers_Courses_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users courses - empty result",
		Args: []string{"30"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/30/courses": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No courses found") {
				t.Errorf("expected 'No courses found' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersCoursesCmd(), tc)
}

func TestCovUsers_Courses_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users courses - API error",
		Args: []string{"30"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/30/courses": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersCoursesCmd(), tc)
}

func TestCovUsers_Courses_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "users courses - invalid ID",
		Args:        []string{"nope"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersCoursesCmd(), tc)
}

func TestCovUsers_MissingSubmissions_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users missing-submissions - happy path",
		Args: []string{"40"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/40/missing_submissions": cmdtest.NewMockResponse(`[{"id":1,"name":"Assignment 1"}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUsersMissingSubmissionsCmd(), tc)
}

func TestCovUsers_MissingSubmissions_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users missing-submissions - empty",
		Args: []string{"40"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/40/missing_submissions": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No missing submissions found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersMissingSubmissionsCmd(), tc)
}

func TestCovUsers_MissingSubmissions_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users missing-submissions - API error",
		Args: []string{"40"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/40/missing_submissions": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersMissingSubmissionsCmd(), tc)
}

func TestCovUsers_MissingSubmissions_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "users missing-submissions - invalid ID",
		Args:        []string{"bad"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersMissingSubmissionsCmd(), tc)
}

func TestCovUsers_ActivityStream_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users activity-stream - happy path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/activity_stream": cmdtest.NewMockResponse(`[{"id":1,"type":"DiscussionTopic","title":"Hello"}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUsersActivityStreamCmd(), tc)
}

func TestCovUsers_ActivityStream_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users activity-stream - empty",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/activity_stream": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No activity stream items found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersActivityStreamCmd(), tc)
}

func TestCovUsers_ActivityStream_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users activity-stream - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/activity_stream": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersActivityStreamCmd(), tc)
}

func TestCovUsers_Todo_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users todo - happy path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/todo": cmdtest.NewMockResponse(`[{"type":"submitting","assignment":{"id":1,"name":"Essay"}}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUsersTodoCmd(), tc)
}

func TestCovUsers_Todo_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users todo - empty",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/todo": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No todo items found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersTodoCmd(), tc)
}

func TestCovUsers_Todo_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users todo - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/todo": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersTodoCmd(), tc)
}

func TestCovUsers_UpcomingEvents_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users upcoming-events - happy path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/upcoming_events": cmdtest.NewMockResponse(`[{"id":1,"title":"Midterm"}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUsersUpcomingEventsCmd(), tc)
}

func TestCovUsers_UpcomingEvents_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users upcoming-events - empty",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/upcoming_events": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No upcoming events found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersUpcomingEventsCmd(), tc)
}

func TestCovUsers_UpcomingEvents_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users upcoming-events - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/upcoming_events": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersUpcomingEventsCmd(), tc)
}

func TestCovUsers_Merge_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users merge - happy path",
		Args: []string{"123", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/merge_into/456": cmdtest.NewMockResponse(`{"id":456,"name":"Merged User"}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "merged") {
				t.Errorf("expected 'merged' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersMergeCmd(), tc)
}

func TestCovUsers_Merge_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users merge - API error",
		Args: []string{"123", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/merge_into/456": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersMergeCmd(), tc)
}

func TestCovUsers_Merge_InvalidSourceID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "users merge - invalid source ID",
		Args:        []string{"bad", "456"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersMergeCmd(), tc)
}

func TestCovUsers_Merge_InvalidDestID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "users merge - invalid dest ID",
		Args:        []string{"123", "bad"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersMergeCmd(), tc)
}

func TestCovUsers_Split_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users split - happy path",
		Args: []string{"999"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/999/split": cmdtest.NewMockResponse(`[{"id":100,"name":"User A"},{"id":101,"name":"User B"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "split") {
				t.Errorf("expected 'split' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersSplitCmd(), tc)
}

func TestCovUsers_Split_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users split - empty response",
		Args: []string{"999"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/999/split": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUsersSplitCmd(), tc)
}

func TestCovUsers_Split_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "users split - API error",
		Args: []string{"999"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/999/split": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersSplitCmd(), tc)
}

func TestCovUsers_Split_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "users split - invalid ID",
		Args:        []string{"xyz"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersSplitCmd(), tc)
}

// =============================================================================
// favorites.go
// =============================================================================

func TestCovFavorites_CoursesList_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites courses list - happy path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/courses": cmdtest.NewMockResponse(`[{"id":1,"name":"Fav Course"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Fav Course") {
				t.Errorf("expected 'Fav Course' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesListCmd(), tc)
}

func TestCovFavorites_CoursesList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites courses list - empty",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/courses": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No favorite courses found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesListCmd(), tc)
}

func TestCovFavorites_CoursesList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites courses list - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/courses": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesListCmd(), tc)
}

func TestCovFavorites_CoursesAdd_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites courses add - happy path",
		Args: []string{"123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/courses/123": cmdtest.NewMockResponse(`{"id":123,"name":"Added Course"}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesAddCmd(), tc)
}

func TestCovFavorites_CoursesAdd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites courses add - API error",
		Args: []string{"123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/courses/123": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesAddCmd(), tc)
}

func TestCovFavorites_CoursesAdd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "favorites courses add - invalid ID",
		Args:        []string{"bad"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesAddCmd(), tc)
}

func TestCovFavorites_CoursesRemove_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites courses remove - happy path",
		Args: []string{"123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/courses/123": cmdtest.NewMockResponse(`{"id":123,"name":"Removed Course"}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "removed") {
				t.Errorf("expected 'removed' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesRemoveCmd(), tc)
}

func TestCovFavorites_CoursesRemove_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites courses remove - API error",
		Args: []string{"123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/courses/123": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesRemoveCmd(), tc)
}

func TestCovFavorites_CoursesRemove_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "favorites courses remove - invalid ID",
		Args:        []string{"nope"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesRemoveCmd(), tc)
}

func TestCovFavorites_CoursesReset_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites courses reset - happy path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/courses": cmdtest.NewMockResponse(`[{"id":1,"name":"Default Course"}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesResetCmd(), tc)
}

func TestCovFavorites_CoursesReset_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites courses reset - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/courses": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesCoursesResetCmd(), tc)
}

func TestCovFavorites_GroupsList_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites groups list - happy path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/groups": cmdtest.NewMockResponse(`[{"id":10,"name":"My Group"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "My Group") {
				t.Errorf("expected 'My Group' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsListCmd(), tc)
}

func TestCovFavorites_GroupsList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites groups list - empty",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/groups": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No favorite groups found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsListCmd(), tc)
}

func TestCovFavorites_GroupsList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites groups list - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/groups": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsListCmd(), tc)
}

func TestCovFavorites_GroupsAdd_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites groups add - happy path",
		Args: []string{"456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/groups/456": cmdtest.NewMockResponse(`{"id":456,"name":"Added Group"}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsAddCmd(), tc)
}

func TestCovFavorites_GroupsAdd_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites groups add - API error",
		Args: []string{"456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/groups/456": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsAddCmd(), tc)
}

func TestCovFavorites_GroupsAdd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "favorites groups add - invalid ID",
		Args:        []string{"bad"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsAddCmd(), tc)
}

func TestCovFavorites_GroupsRemove_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites groups remove - happy path",
		Args: []string{"456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/groups/456": cmdtest.NewMockResponse(`{"id":456,"name":"Removed Group"}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "removed") {
				t.Errorf("expected 'removed' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsRemoveCmd(), tc)
}

func TestCovFavorites_GroupsRemove_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites groups remove - API error",
		Args: []string{"456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/groups/456": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsRemoveCmd(), tc)
}

func TestCovFavorites_GroupsRemove_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "favorites groups remove - invalid ID",
		Args:        []string{"!"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsRemoveCmd(), tc)
}

func TestCovFavorites_GroupsReset_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites groups reset - happy path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/groups": cmdtest.NewMockResponse(`[{"id":1,"name":"Default Group"}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsResetCmd(), tc)
}

func TestCovFavorites_GroupsReset_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "favorites groups reset - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/favorites/groups": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newFavoritesGroupsResetCmd(), tc)
}

// =============================================================================
// bookmarks.go
// =============================================================================

func TestCovBookmarks_List_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks list - happy path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks": cmdtest.NewMockResponse(`[{"id":1,"name":"My Bookmark","url":"https://example.com"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "My Bookmark") {
				t.Errorf("expected bookmark in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newBookmarksListCmd(), tc)
}

func TestCovBookmarks_List_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks list - empty",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No bookmarks found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newBookmarksListCmd(), tc)
}

func TestCovBookmarks_List_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks list - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newBookmarksListCmd(), tc)
}

func TestCovBookmarks_Get_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks get - happy path",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks/1": cmdtest.NewMockResponse(`{"id":1,"name":"My Bookmark","url":"https://example.com"}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "My Bookmark") {
				t.Errorf("expected bookmark in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newBookmarksGetCmd(), tc)
}

func TestCovBookmarks_Get_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks get - API error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks/1": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newBookmarksGetCmd(), tc)
}

func TestCovBookmarks_Get_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "bookmarks get - invalid ID",
		Args:        []string{"bad"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newBookmarksGetCmd(), tc)
}

func TestCovBookmarks_Create_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks create - happy path",
		Args: []string{"--name", "New Bookmark", "--url", "https://example.com/page"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks": cmdtest.NewMockResponse(`{"id":5,"name":"New Bookmark","url":"https://example.com/page"}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "New Bookmark") {
				t.Errorf("expected bookmark in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newBookmarksCreateCmd(), tc)
}

func TestCovBookmarks_Create_WithPosition(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks create - with position",
		Args: []string{"--name", "Positioned", "--url", "https://x.com", "--position", "3"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks": cmdtest.NewMockResponse(`{"id":6,"name":"Positioned","url":"https://x.com"}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newBookmarksCreateCmd(), tc)
}

func TestCovBookmarks_Create_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks create - API error",
		Args: []string{"--name", "New Bookmark", "--url", "https://example.com"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks": cmdtest.NewErrorResponse(422, "unprocessable"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newBookmarksCreateCmd(), tc)
}

func TestCovBookmarks_Update_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks update - happy path",
		Args: []string{"1", "--name", "Updated Bookmark"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks/1": cmdtest.NewMockResponse(`{"id":1,"name":"Updated Bookmark","url":"https://example.com"}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newBookmarksUpdateCmd(), tc)
}

func TestCovBookmarks_Update_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks update - API error",
		Args: []string{"1", "--name", "Updated"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks/1": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newBookmarksUpdateCmd(), tc)
}

func TestCovBookmarks_Update_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "bookmarks update - invalid ID",
		Args:        []string{"bad", "--name", "Updated"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newBookmarksUpdateCmd(), tc)
}

func TestCovBookmarks_Delete_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks delete - happy path",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks/1": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "deleted") {
				t.Errorf("expected 'deleted' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newBookmarksDeleteCmd(), tc)
}

func TestCovBookmarks_Delete_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bookmarks delete - API error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/bookmarks/1": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newBookmarksDeleteCmd(), tc)
}

func TestCovBookmarks_Delete_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "bookmarks delete - invalid ID",
		Args:        []string{"bad"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newBookmarksDeleteCmd(), tc)
}

// =============================================================================
// course_nicknames.go
// =============================================================================

func TestCovCourseNicknames_List_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames list - happy path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames": cmdtest.NewMockResponse(`[{"course_id":1,"name":"CS101","nickname":"My CS"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "My CS") {
				t.Errorf("expected nickname in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesListCmd(), tc)
}

func TestCovCourseNicknames_List_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames list - empty",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No course nicknames found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesListCmd(), tc)
}

func TestCovCourseNicknames_List_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames list - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesListCmd(), tc)
}

func TestCovCourseNicknames_Get_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames get - happy path",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames/10": cmdtest.NewMockResponse(`{"course_id":10,"name":"CS101","nickname":"My CS"}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "My CS") {
				t.Errorf("expected nickname in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesGetCmd(), tc)
}

func TestCovCourseNicknames_Get_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames get - API error",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames/10": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesGetCmd(), tc)
}

func TestCovCourseNicknames_Get_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "course-nicknames get - invalid ID",
		Args:        []string{"bad"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesGetCmd(), tc)
}

func TestCovCourseNicknames_Set_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames set - happy path",
		Args: []string{"10", "--nickname", "Fav Course"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames/10": cmdtest.NewMockResponse(`{"course_id":10,"name":"CS101","nickname":"Fav Course"}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesSetCmd(), tc)
}

func TestCovCourseNicknames_Set_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames set - API error",
		Args: []string{"10", "--nickname", "Fav"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames/10": cmdtest.NewErrorResponse(422, "invalid"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesSetCmd(), tc)
}

func TestCovCourseNicknames_Set_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "course-nicknames set - invalid ID",
		Args:        []string{"bad", "--nickname", "Name"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesSetCmd(), tc)
}

func TestCovCourseNicknames_Delete_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames delete - happy path",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames/10": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "deleted") {
				t.Errorf("expected 'deleted' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesDeleteCmd(), tc)
}

func TestCovCourseNicknames_Delete_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames delete - API error",
		Args: []string{"10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames/10": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesDeleteCmd(), tc)
}

func TestCovCourseNicknames_Delete_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "course-nicknames delete - invalid ID",
		Args:        []string{"bad"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesDeleteCmd(), tc)
}

func TestCovCourseNicknames_DeleteAll_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames delete-all - happy path",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "deleted") {
				t.Errorf("expected 'deleted' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesDeleteAllCmd(), tc)
}

func TestCovCourseNicknames_DeleteAll_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "course-nicknames delete-all - API error",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self/course_nicknames": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCourseNicknamesDeleteAllCmd(), tc)
}

// =============================================================================
// observees.go
// =============================================================================

func TestCovObservees_List_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "observees list - happy path",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/observees": cmdtest.NewMockResponse(`[{"id":1,"name":"Student One"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Student One") {
				t.Errorf("expected observee in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newObserveesListCmd(), tc)
}

func TestCovObservees_List_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "observees list - empty",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/observees": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No observees found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newObserveesListCmd(), tc)
}

func TestCovObservees_List_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "observees list - API error",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/observees": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newObserveesListCmd(), tc)
}

func TestCovObservees_Get_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "observees get - happy path",
		Args: []string{"--user-id", "123", "--observee-id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/observees": cmdtest.NewMockResponse(`{"id":456,"name":"Student Two"}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newObserveesGetCmd(), tc)
}

func TestCovObservees_Get_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "observees get - API error",
		Args: []string{"--user-id", "123", "--observee-id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/observees": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newObserveesGetCmd(), tc)
}

func TestCovObservees_Remove_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "observees remove - happy path",
		Args: []string{"--user-id", "123", "--observee-id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/observees/456": cmdtest.NewMockResponse(`{"id":456,"name":"Student Two"}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "removed") {
				t.Errorf("expected 'removed' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newObserveesRemoveCmd(), tc)
}

func TestCovObservees_Remove_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "observees remove - API error",
		Args: []string{"--user-id", "123", "--observee-id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/observees/456": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newObserveesRemoveCmd(), tc)
}

func TestCovObservers_List_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "observees observers list - happy path",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/observers": cmdtest.NewMockResponse(`[{"id":10,"name":"Parent One"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Parent One") {
				t.Errorf("expected observer in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newObserversListCmd(), tc)
}

func TestCovObservers_List_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "observees observers list - empty",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/observers": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No observers found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newObserversListCmd(), tc)
}

func TestCovObservers_List_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "observees observers list - API error",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/observers": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newObserversListCmd(), tc)
}

// =============================================================================
// user_features.go
// =============================================================================

func TestCovUserFeatures_List_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features list - happy path",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features": cmdtest.NewMockResponse(`[{"feature":"new_gradebook","applies_to":"User","feature_flag":{"feature":"new_gradebook","state":"on"}}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUserFeaturesListCmd(), tc)
}

func TestCovUserFeatures_List_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features list - empty",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No features found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUserFeaturesListCmd(), tc)
}

func TestCovUserFeatures_List_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features list - API error",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUserFeaturesListCmd(), tc)
}

func TestCovUserFeatures_ListEnabled_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features list-enabled - happy path",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features/enabled": cmdtest.NewMockResponse(`[{"feature":"new_gradebook","applies_to":"User","feature_flag":{"feature":"new_gradebook","state":"on"}}]`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUserFeaturesListEnabledCmd(), tc)
}

func TestCovUserFeatures_ListEnabled_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features list-enabled - API error",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features/enabled": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUserFeaturesListEnabledCmd(), tc)
}

func TestCovUserFeatures_GetFlag_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features get-flag - happy path",
		Args: []string{"--user-id", "123", "--feature", "new_gradebook"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features/flags/new_gradebook": cmdtest.NewMockResponse(`{"feature":"new_gradebook","state":"on"}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUserFeaturesGetFlagCmd(), tc)
}

func TestCovUserFeatures_GetFlag_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features get-flag - API error",
		Args: []string{"--user-id", "123", "--feature", "new_gradebook"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features/flags/new_gradebook": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUserFeaturesGetFlagCmd(), tc)
}

func TestCovUserFeatures_SetFlag_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features set-flag - happy path",
		Args: []string{"--user-id", "123", "--feature", "new_gradebook", "--state", "on"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features/flags/new_gradebook": cmdtest.NewMockResponse(`{"feature":"new_gradebook","state":"on"}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newUserFeaturesSetFlagCmd(), tc)
}

func TestCovUserFeatures_SetFlag_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features set-flag - API error",
		Args: []string{"--user-id", "123", "--feature", "new_gradebook", "--state", "on"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features/flags/new_gradebook": cmdtest.NewErrorResponse(422, "invalid state"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUserFeaturesSetFlagCmd(), tc)
}

func TestCovUserFeatures_DeleteFlag_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features delete-flag - happy path",
		Args: []string{"--user-id", "123", "--feature", "new_gradebook"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features/flags/new_gradebook": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "deleted") {
				t.Errorf("expected 'deleted' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newUserFeaturesDeleteFlagCmd(), tc)
}

func TestCovUserFeatures_DeleteFlag_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "user-features delete-flag - API error",
		Args: []string{"--user-id", "123", "--feature", "new_gradebook"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/features/flags/new_gradebook": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUserFeaturesDeleteFlagCmd(), tc)
}

// =============================================================================
// comm_channels.go
// =============================================================================

func TestCovCommChannels_List_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "comm-channels list - happy path",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/communication_channels": cmdtest.NewMockResponse(`[{"id":1,"address":"user@example.com","type":"email"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "user@example.com") {
				t.Errorf("expected channel in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newCommChannelsListCmd(), tc)
}

func TestCovCommChannels_List_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "comm-channels list - empty",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/communication_channels": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No communication channels found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newCommChannelsListCmd(), tc)
}

func TestCovCommChannels_List_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "comm-channels list - API error",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/communication_channels": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCommChannelsListCmd(), tc)
}

func TestCovCommChannels_Create_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "comm-channels create - happy path",
		Args: []string{"--user-id", "123", "--address", "new@example.com", "--type", "email"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/communication_channels": cmdtest.NewMockResponse(`{"id":5,"address":"new@example.com","type":"email"}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "new@example.com") {
				t.Errorf("expected channel address in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newCommChannelsCreateCmd(), tc)
}

func TestCovCommChannels_Create_WithSkip(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "comm-channels create - with skip-confirmation",
		Args: []string{"--user-id", "123", "--address", "+15551234567", "--type", "sms", "--skip-confirmation"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/communication_channels": cmdtest.NewMockResponse(`{"id":6,"address":"+15551234567","type":"sms"}`),
		},
		ExpectError: false,
	}
	cmdtest.RunCommandTest(t, newCommChannelsCreateCmd(), tc)
}

func TestCovCommChannels_Create_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "comm-channels create - API error",
		Args: []string{"--user-id", "123", "--address", "bad@example.com", "--type", "email"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/communication_channels": cmdtest.NewErrorResponse(422, "unprocessable"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCommChannelsCreateCmd(), tc)
}

func TestCovCommChannels_Delete_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "comm-channels delete - happy path",
		Args: []string{"--user-id", "123", "--channel-id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/communication_channels/456": cmdtest.NewMockResponse(`{"id":456,"address":"user@example.com","type":"email"}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "deleted") {
				t.Errorf("expected 'deleted' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newCommChannelsDeleteCmd(), tc)
}

func TestCovCommChannels_Delete_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "comm-channels delete - API error",
		Args: []string{"--user-id", "123", "--channel-id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/communication_channels/456": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCommChannelsDeleteCmd(), tc)
}

// =============================================================================
// content_shares.go
// =============================================================================

func TestCovContentShares_ListSent_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "content-shares list-sent - happy path",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/content_shares/sent": cmdtest.NewMockResponse(`[{"id":1,"name":"My Assignment","content_type":"Assignment"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "My Assignment") {
				t.Errorf("expected share in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newContentSharesListSentCmd(), tc)
}

func TestCovContentShares_ListSent_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "content-shares list-sent - empty",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/content_shares/sent": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No sent content shares found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newContentSharesListSentCmd(), tc)
}

func TestCovContentShares_ListSent_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "content-shares list-sent - API error",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/content_shares/sent": cmdtest.NewErrorResponse(403, "forbidden"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newContentSharesListSentCmd(), tc)
}

func TestCovContentShares_ListReceived_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "content-shares list-received - happy path",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/content_shares/received": cmdtest.NewMockResponse(`[{"id":2,"name":"Received Quiz","content_type":"Quiz"}]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Received Quiz") {
				t.Errorf("expected share in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newContentSharesListReceivedCmd(), tc)
}

func TestCovContentShares_ListReceived_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "content-shares list-received - empty",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/content_shares/received": cmdtest.NewMockResponse(`[]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "No received content shares found") {
				t.Errorf("expected empty message in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newContentSharesListReceivedCmd(), tc)
}

func TestCovContentShares_ListReceived_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "content-shares list-received - API error",
		Args: []string{"--user-id", "123"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/content_shares/received": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newContentSharesListReceivedCmd(), tc)
}

func TestCovContentShares_Get_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "content-shares get - happy path",
		Args: []string{"--user-id", "123", "--id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/content_shares/456": cmdtest.NewMockResponse(`{"id":456,"name":"My Share","content_type":"Assignment"}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "My Share") {
				t.Errorf("expected share in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newContentSharesGetCmd(), tc)
}

func TestCovContentShares_Get_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "content-shares get - API error",
		Args: []string{"--user-id", "123", "--id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/content_shares/456": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newContentSharesGetCmd(), tc)
}

func TestCovContentShares_Delete_Happy(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "content-shares delete - happy path",
		Args: []string{"--user-id", "123", "--id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/content_shares/456": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "deleted") {
				t.Errorf("expected 'deleted' in output, got: %q", output)
			}
		},
	}
	cmdtest.RunCommandTest(t, newContentSharesDeleteCmd(), tc)
}

func TestCovContentShares_Delete_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "content-shares delete - API error",
		Args: []string{"--user-id", "123", "--id", "456"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/123/content_shares/456": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newContentSharesDeleteCmd(), tc)
}
