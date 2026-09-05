package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/chiptoe-svg/canvas-cli/commands/internal/testing"
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
