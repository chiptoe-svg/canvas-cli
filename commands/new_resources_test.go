package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// --- progress ---

func TestProgressGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get progress job",
			Args: []string{"42"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/progress/42": cmdtest.NewMockResponse(`{"id":42,"workflow_state":"running","completion":50}`),
			},
			ExpectError: false,
		},
		{
			Name:        "missing id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name:        "invalid id",
			Args:        []string{"abc"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newProgressGetCmd(), tc)
		})
	}
}

func TestProgressCancelCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "cancel progress job with force",
			Args: []string{"42", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/progress/42/cancel": cmdtest.NewMockResponse(`{"id":42,"workflow_state":"failed"}`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newProgressCancelCmd(), tc)
		})
	}
}

// --- jwts ---

func TestJWTCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create jwt",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/jwts": cmdtest.NewMockResponse(`{"token":"eyJ.test"}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "eyJ.test") {
					t.Error("expected token in output")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newJWTCreateCmd(), tc)
		})
	}
}

// --- error-reports ---

func TestErrorReportCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create error report",
			Args: []string{"--subject", "Bug report"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/error_reports": cmdtest.NewMockResponse(`{"logged":true}`),
			},
			ExpectError: false,
		},
		{
			Name:        "missing subject",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newErrorReportCreateCmd(), tc)
		})
	}
}

// --- comm-messages ---

func TestCommMessagesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list comm messages",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/comm_messages": cmdtest.NewMockResponse(`[{"id":1,"subject":"Welcome","workflow_state":"sent"}]`),
			},
			ExpectError: false,
		},
		{
			Name: "empty comm messages",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/comm_messages": cmdtest.NewMockResponse(`[]`),
			},
			ExpectOutput: "No communication messages found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCommMessagesListCmd(), tc)
		})
	}
}

// --- conferences ---

func TestConferencesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list all conferences",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conferences": cmdtest.NewMockResponse(`[{"id":1,"title":"Lecture 1"}]`),
			},
			ExpectError: false,
		},
		{
			Name: "list course conferences",
			Args: []string{"--course-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/10/conferences": cmdtest.NewMockResponse(`[{"id":5,"title":"Course Conf"}]`),
			},
			ExpectError: false,
		},
		{
			Name:         "empty conferences",
			Args:         []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conferences": cmdtest.NewMockResponse(`[]`),
			},
			ExpectOutput: "No conferences found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newConferencesListCmd(), tc)
		})
	}
}

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

// --- eportfolios ---

func TestEportfoliosListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list eportfolios",
			Args: []string{"--user-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/1/eportfolios": cmdtest.NewMockResponse(`[{"id":10,"name":"My Portfolio"}]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "My Portfolio") {
					t.Error("expected 'My Portfolio' in output")
				}
			},
		},
		{
			Name:        "missing user-id",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEportfoliosListCmd(), tc)
		})
	}
}

func TestEportfolioGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get eportfolio",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10": cmdtest.NewMockResponse(`{"id":10,"name":"My Portfolio"}`),
			},
			ExpectError: false,
		},
		{
			Name:        "missing id",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEportfolioGetCmd(), tc)
		})
	}
}

// --- epub-exports ---

func TestEpubExportsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list epub exports",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/epub_exports": cmdtest.NewMockResponse(`[{"id":1,"workflow_state":"generated"}]`),
			},
			ExpectError: false,
		},
		{
			Name:         "empty list",
			Args:         []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/epub_exports": cmdtest.NewMockResponse(`[]`),
			},
			ExpectOutput: "No ePub exports found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEpubExportsListCmd(), tc)
		})
	}
}

func TestEpubExportCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create epub export",
			Args: []string{"--course-id", "42"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/42/epub_exports": cmdtest.NewMockResponse(`{"id":5,"workflow_state":"created"}`),
			},
			ExpectError: false,
		},
		{
			Name:        "missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEpubExportCreateCmd(), tc)
		})
	}
}

// --- account-calendars ---

func TestAccountCalendarsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list account calendars",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/account_calendars": cmdtest.NewMockResponse(`[{"id":1,"name":"Root","visible":true}]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Root") {
					t.Error("expected 'Root' in output")
				}
			},
		},
		{
			Name: "list calendars for account",
			Args: []string{"--account-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/5/account_calendars": cmdtest.NewMockResponse(`[{"id":10,"name":"Sub"}]`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newAccountCalendarsListCmd(), tc)
		})
	}
}

// --- brand ---

func TestBrandGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get global brand variables",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/brand_variables": cmdtest.NewMockResponse(`{"ic-brand-primary":"#E66000"}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "#E66000") {
					t.Error("expected brand color in output")
				}
			},
		},
		{
			Name: "get account brand variables",
			Args: []string{"--account-id", "3"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/accounts/3/brand_variables": cmdtest.NewMockResponse(`{"ic-brand-primary":"#0770A3"}`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newBrandGetCmd(), tc)
		})
	}
}

// --- media ---

func TestMediaObjectsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list media objects",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_objects": cmdtest.NewMockResponse(`[{"media_id":"m-abc","title":"Video 1"}]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Video 1") {
					t.Error("expected 'Video 1' in output")
				}
			},
		},
		{
			Name: "list course media objects",
			Args: []string{"--course-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/10/media_objects": cmdtest.NewMockResponse(`[{"media_id":"m-xyz"}]`),
			},
			ExpectError: false,
		},
		{
			Name:         "empty list",
			Args:         []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_objects": cmdtest.NewMockResponse(`[]`),
			},
			ExpectOutput: "No media objects found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newMediaObjectsListCmd(), tc)
		})
	}
}

func TestMediaAttachmentsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list media attachments",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_attachments": cmdtest.NewMockResponse(`[{"id":5,"filename":"video.mp4"}]`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newMediaAttachmentsListCmd(), tc)
		})
	}
}

// --- history ---

func TestHistoryListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list user history",
			Args: []string{"--user-id", "42"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/42/history": cmdtest.NewMockResponse(`[{"asset_code":"assignment_1","asset_name":"Midterm"}]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Midterm") {
					t.Error("expected 'Midterm' in output")
				}
			},
		},
		{
			Name:        "missing user-id",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newHistoryListCmd(), tc)
		})
	}
}

// --- audit ---

func TestAuditListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "grade-change global",
			Args: []string{"--type", "grade-change"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/audit/grade_change": cmdtest.NewMockResponse(`[{"id":"gc1","event_type":"grade_change"}]`),
			},
			ExpectError: false,
		},
		{
			Name: "grade-change for course",
			Args: []string{"--type", "grade-change", "--course-id", "20"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/audit/grade_change/courses/20": cmdtest.NewMockResponse(`[{"id":"gc2"}]`),
			},
			ExpectError: false,
		},
		{
			Name: "authentication for account",
			Args: []string{"--type", "authentication", "--account-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/audit/authentication/accounts/1": cmdtest.NewMockResponse(`[{"id":"e1","event_type":"login"}]`),
			},
			ExpectError: false,
		},
		{
			Name:        "authentication without context",
			Args:        []string{"--type", "authentication"},
			ExpectError: true,
		},
		{
			Name:        "unknown audit type",
			Args:        []string{"--type", "unknown"},
			ExpectError: true,
		},
		{
			Name:         "empty results",
			Args:         []string{"--type", "grade-change"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/audit/grade_change": cmdtest.NewMockResponse(`[]`),
			},
			ExpectOutput: "No audit events found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newAuditListCmd(), tc)
		})
	}
}
