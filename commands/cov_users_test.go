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
