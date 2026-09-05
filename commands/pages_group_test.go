package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/chiptoe-svg/canvas-cli/commands/internal/testing"
)

func TestPagesListCmd_GroupContext(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list pages for group",
			Args: []string{"--group-id", "7"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/7/pages": cmdtest.NewMockResponse(`[
					{"page_id": 10, "url": "group-home", "title": "Group Home", "published": true}
				]`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Group Home") {
					t.Error("expected 'Group Home' in output")
				}
			},
		},
		{
			Name:        "both course-id and group-id specified",
			Args:        []string{"--course-id", "1", "--group-id", "7"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newPagesListCmd(), tc)
		})
	}
}

func TestPagesGetCmd_GroupContext(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get page by group context",
			Args: []string{"--group-id", "7", "group-home"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/7/pages/group-home": cmdtest.NewMockResponse(`{
					"page_id": 10, "url": "group-home", "title": "Group Home", "body": "<p>Hi</p>", "published": true
				}`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Group Home") {
					t.Error("expected 'Group Home' in output")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newPagesGetCmd(), tc)
		})
	}
}

func TestPagesFrontCmd_GroupContext(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get front page for group",
			Args: []string{"--group-id", "7"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/7/front_page": cmdtest.NewMockResponse(`{
					"page_id": 10, "url": "group-home", "title": "Group Home", "published": true
				}`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Group Home") {
					t.Error("expected 'Group Home' in output")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newPagesFrontCmd(), tc)
		})
	}
}

func TestPagesCreateCmd_GroupContext(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create page for group",
			Args: []string{"--group-id", "7", "--title", "Group Resource"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/7/pages": cmdtest.NewMockResponse(`{
					"page_id": 11, "url": "group-resource", "title": "Group Resource", "published": false
				}`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Group Resource") {
					t.Error("expected 'Group Resource' in output")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newPagesCreateCmd(), tc)
		})
	}
}

func TestPagesDeleteCmd_GroupContext(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete page from group with force",
			Args: []string{"--group-id", "7", "group-page", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/7/pages/group-page": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newPagesDeleteCmd(), tc)
		})
	}
}

func TestPagesRevisionsCmd_GroupContext(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list revisions for group page",
			Args: []string{"--group-id", "7", "group-page"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/7/pages/group-page/revisions": cmdtest.NewMockResponse(`[
					{"revision_id": 1, "url": "group-page", "updated_at": "2024-01-01T00:00:00Z"}
				]`),
			},
			ExpectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newPagesRevisionsCmd(), tc)
		})
	}
}
