package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

func TestFoldersListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list folders for course",
			Args: []string{"--course-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/10/folders": cmdtest.NewMockResponse(`[
					{"id": 1, "name": "Root", "context_type": "Course"},
					{"id": 2, "name": "Lectures", "context_type": "Course"}
				]`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Lectures") {
					t.Error("expected 'Lectures' in output")
				}
			},
		},
		{
			Name: "list folders for group",
			Args: []string{"--group-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/5/folders": cmdtest.NewMockResponse(`[
					{"id": 3, "name": "Group Root"}
				]`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Group Root") {
					t.Error("expected 'Group Root' in output")
				}
			},
		},
		{
			Name: "list folders for user",
			Args: []string{"--user-id", "99"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/99/folders": cmdtest.NewMockResponse(`[
					{"id": 4, "name": "My Files"}
				]`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "My Files") {
					t.Error("expected 'My Files' in output")
				}
			},
		},
		{
			Name:        "no context specified",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list sub-folders",
			Args: []string{"--folder-id", "7"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/folders/7/folders": cmdtest.NewMockResponse(`[
					{"id": 8, "name": "Week 1"},
					{"id": 9, "name": "Week 2"}
				]`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Week 1") {
					t.Error("expected 'Week 1' in output")
				}
			},
		},
		{
			Name: "empty result",
			Args: []string{"--course-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/10/folders": cmdtest.NewMockResponse(`[]`),
			},
			ExpectOutput: "No folders found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newFoldersListCmd(), tc)
		})
	}
}

func TestFoldersGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get folder by id",
			Args: []string{"--folder-id", "55"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/folders/55": cmdtest.NewMockResponse(`{"id": 55, "name": "Lectures"}`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Lectures") {
					t.Error("expected 'Lectures' in output")
				}
			},
		},
		{
			Name:        "missing folder-id",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newFoldersGetCmd(), tc)
		})
	}
}

func TestFoldersResolvePathCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "resolve path for course",
			Args: []string{"--course-id", "10", "--path", "lectures/week1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/10/folders/by_path/lectures/week1": cmdtest.NewMockResponse(`[
					{"id": 1, "name": "Root"},
					{"id": 2, "name": "lectures"},
					{"id": 3, "name": "week1"}
				]`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "week1") {
					t.Error("expected 'week1' in output")
				}
			},
		},
		{
			Name:        "no context",
			Args:        []string{"--path", "some/path"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newFoldersResolvePathCmd(), tc)
		})
	}
}

func TestFoldersCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create folder in course",
			Args: []string{"--course-id", "10", "--name", "Assignments"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/10/folders": cmdtest.NewMockResponse(`{"id": 20, "name": "Assignments"}`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Assignments") {
					t.Error("expected 'Assignments' in output")
				}
			},
		},
		{
			Name: "create folder in group",
			Args: []string{"--group-id", "5", "--name", "Resources"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/5/folders": cmdtest.NewMockResponse(`{"id": 21, "name": "Resources"}`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Resources") {
					t.Error("expected 'Resources' in output")
				}
			},
		},
		{
			Name: "create sub-folder",
			Args: []string{"--parent-folder-id", "7", "--name", "Week 1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/folders/7/folders": cmdtest.NewMockResponse(`{"id": 22, "name": "Week 1"}`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Week 1") {
					t.Error("expected 'Week 1' in output")
				}
			},
		},
		{
			Name:        "missing name",
			Args:        []string{"--course-id", "10"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newFoldersCreateCmd(), tc)
		})
	}
}

func TestFoldersDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete folder with force",
			Args: []string{"--folder-id", "55", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/folders/55": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "missing folder-id",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newFoldersDeleteCmd(), tc)
		})
	}
}

func TestFoldersMediaCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get media folder for course",
			Args: []string{"--course-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/10/folders/media": cmdtest.NewMockResponse(`{"id": 100, "name": "Uploaded Media"}`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Uploaded Media") {
					t.Error("expected 'Uploaded Media' in output")
				}
			},
		},
		{
			Name:        "no context specified",
			Args:        []string{},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newFoldersMediaCmd(), tc)
		})
	}
}

func TestFoldersCopyCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "copy folder",
			Args: []string{"--dest-folder-id", "20", "--source-folder-id", "10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/folders/20/copy_folder": cmdtest.NewMockResponse(`{"id": 30, "name": "CopiedFolder"}`),
			},
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "CopiedFolder") {
					t.Error("expected 'CopiedFolder' in output")
				}
			},
		},
		{
			Name:        "missing dest-folder-id",
			Args:        []string{"--source-folder-id", "10"},
			ExpectError: true,
		},
		{
			Name:        "missing source-folder-id",
			Args:        []string{"--dest-folder-id", "20"},
			ExpectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newFoldersCopyCmd(), tc)
		})
	}
}
