package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestFilesListCmd_GroupContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list files for group",
		Args: []string{"--group-id", "9"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/9/files": cmdtest.NewMockResponse(`[
				{"id": 7, "display_name": "Group_Doc.pdf", "filename": "group_doc.pdf", "size": 4096}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Group_Doc.pdf") {
				t.Error("expected 'Group_Doc.pdf' in output")
			}
		},
	}
	cmd := newFilesListCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestFilesQuotaCmd_GroupContext(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get quota for group",
		Args: []string{"--group-id", "9"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/groups/9/files/quota": cmdtest.NewMockResponse(`{
				"quota": 524288000,
				"quota_used": 1048576
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Storage Quota") {
				t.Error("expected 'Storage Quota' in output")
			}
		},
	}
	cmd := newFilesQuotaCmd()
	cmdtest.RunCommandTest(t, cmd, tc)
}

func TestFilesResetVerifierCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "reset verifier successfully",
			Args: []string{"42"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/files/42/reset_verifier": cmdtest.NewMockResponse(`{
					"id": 42, "display_name": "file.pdf", "filename": "file.pdf", "size": 2048
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "reset") {
					t.Error("expected 'reset' in output")
				}
			},
		},
		{
			Name:        "missing file-id argument",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newFilesResetVerifierCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestFilesCopyCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "copy file to folder successfully",
			Args: []string{"--dest-folder-id", "10", "--source-file-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/folders/10/copy_file": cmdtest.NewMockResponse(`{
					"id": 99, "display_name": "copied.pdf", "filename": "copied.pdf", "size": 1024
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "copied") {
					t.Error("expected 'copied' in output")
				}
			},
		},
		{
			Name:        "missing dest-folder-id",
			Args:        []string{"--source-file-id", "5"},
			ExpectError: true,
		},
		{
			Name:        "missing source-file-id",
			Args:        []string{"--dest-folder-id", "10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newFilesCopyCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestFilesSetUsageRightsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "set usage rights for course",
			Args: []string{
				"--course-id", "10",
				"--file-ids", "5",
				"--use-justification", "public_domain",
			},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/10/usage_rights": cmdtest.NewMockResponse(`{
					"use_justification": "public_domain",
					"license": "",
					"legal_copyright": ""
				}`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "public_domain") {
					t.Error("expected 'public_domain' in output")
				}
			},
		},
		{
			Name:        "missing context",
			Args:        []string{"--file-ids", "5", "--use-justification", "public_domain"},
			ExpectError: true,
		},
		{
			Name:        "missing use-justification",
			Args:        []string{"--course-id", "10", "--file-ids", "5"},
			ExpectError: true,
		},
		{
			Name:        "missing file-ids and folder-ids",
			Args:        []string{"--course-id", "10", "--use-justification", "public_domain"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newFilesUsageRightsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestFilesRemoveUsageRightsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "remove usage rights for course",
			Args: []string{"--course-id", "10", "--file-ids", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/10/usage_rights": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "missing context",
			Args:        []string{"--file-ids", "5"},
			ExpectError: true,
		},
		{
			Name:        "missing file-ids and folder-ids",
			Args:        []string{"--course-id", "10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newFilesRemoveUsageRightsCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}

func TestFilesLicensesCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list licenses for course",
			Args: []string{"--course-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/10/content_licenses": cmdtest.NewMockResponse(`[
					{"id": "public_domain", "name": "Public Domain"},
					{"id": "cc_by", "name": "CC Attribution"}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "public_domain") {
					t.Error("expected 'public_domain' in output")
				}
			},
		},
		{
			Name: "list licenses for group",
			Args: []string{"--group-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/5/content_licenses": cmdtest.NewMockResponse(`[
					{"id": "cc_by_sa", "name": "CC Attribution Share-Alike"}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "cc_by_sa") {
					t.Error("expected 'cc_by_sa' in output")
				}
			},
		},
		{
			Name:        "no context",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := newFilesLicensesCmd()
			cmdtest.RunCommandTest(t, cmd, tc)
		})
	}
}
