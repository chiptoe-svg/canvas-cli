package commands

import (
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// ---------------------------------------------------------------------------
// course_settings.go
// ---------------------------------------------------------------------------

func TestCovCourse_CourseSettingsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get course settings successfully",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/settings": cmdtest.NewMockResponse(`{
					"allow_student_discussion_topics": true,
					"allow_student_forum_attachments": false
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get course settings - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "get course settings - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/settings": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseSettingsGetCmd(), tc)
		})
	}
}

func TestCovCourse_CourseSettingsTodoCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get course todo items",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/todo": cmdtest.NewMockResponse(`[
					{"type": "grading", "assignment": {"id": 5, "name": "Essay 1"}}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "get course todo items - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/todo": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No todo items found",
		},
		{
			Name:        "get course todo items - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "get course todo items - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/todo": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseSettingsTodoCmd(), tc)
		})
	}
}

func TestCovCourse_CourseSettingsTabsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list course tabs",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/tabs": cmdtest.NewMockResponse(`[
					{"id": "home", "label": "Home", "position": 1}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list course tabs - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/tabs": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No tabs found",
		},
		{
			Name:        "list course tabs - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list course tabs - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/tabs": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseSettingsTabsCmd(), tc)
		})
	}
}

func TestCovCourse_CourseSettingsPermissionsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get course permissions",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/permissions": cmdtest.NewMockResponse(`{
					"read_course_content": true,
					"manage_grades": false
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get course permissions - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "get course permissions - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/permissions": cmdtest.NewErrorResponse(403, "forbidden"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseSettingsPermissionsCmd(), tc)
		})
	}
}

func TestCovCourse_CourseSettingsEffectiveDueDatesCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get effective due dates",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/effective_due_dates": cmdtest.NewMockResponse(`{
					"10": {"1": {"due_at": "2024-05-01T23:59:59Z"}}
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get effective due dates - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "get effective due dates - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/effective_due_dates": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseSettingsEffectiveDueDatesCmd(), tc)
		})
	}
}

func TestCovCourse_CourseSettingsLatePolicyCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get late policy",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/late_policy": cmdtest.NewMockResponse(`{
					"late_policy": {"id": 42, "late_submission_interval": "day"}
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get late policy - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "get late policy - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/late_policy": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseSettingsLatePolicyCmd(), tc)
		})
	}
}

func TestCovCourse_CourseSettingsRecentStudentsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get recent students",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/recent_students": cmdtest.NewMockResponse(`[
					{"id": 100, "name": "Alice Smith"}
				]`),
			},
			ExpectError: false,
			ValidateOutput: func(t *testing.T, output string) {
				if !strings.Contains(output, "Alice Smith") {
					t.Error("expected 'Alice Smith' in output")
				}
			},
		},
		{
			Name: "get recent students - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/recent_students": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No recent students found",
		},
		{
			Name:        "get recent students - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "get recent students - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/recent_students": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseSettingsRecentStudentsCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// course_features.go
// ---------------------------------------------------------------------------

func TestCovCourse_CourseFeaturesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list course features",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features": cmdtest.NewMockResponse(`[
					{"feature": "new_quizzes", "feature_flag": {"state": "on"}}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list course features - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No features found",
		},
		{
			Name:        "list course features - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list course features - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseFeaturesListCmd(), tc)
		})
	}
}

func TestCovCourse_CourseFeaturesListEnabledCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list enabled course features",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features/enabled": cmdtest.NewMockResponse(`["new_quizzes", "course_pacing"]`),
			},
			ExpectError: false,
		},
		{
			Name: "list enabled course features - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features/enabled": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No enabled features found",
		},
		{
			Name:        "list enabled features - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list enabled features - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features/enabled": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseFeaturesListEnabledCmd(), tc)
		})
	}
}

func TestCovCourse_CourseFeaturesGetFlagCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get feature flag",
			Args: []string{"--course-id", "1", "--feature", "new_quizzes"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features/flags/new_quizzes": cmdtest.NewMockResponse(`{
					"feature": "new_quizzes", "state": "on"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get feature flag - missing course-id",
			Args:        []string{"--feature", "new_quizzes"},
			ExpectError: true,
		},
		{
			Name:        "get feature flag - missing feature",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "get feature flag - API error",
			Args: []string{"--course-id", "1", "--feature", "new_quizzes"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features/flags/new_quizzes": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseFeaturesGetFlagCmd(), tc)
		})
	}
}

func TestCovCourse_CourseFeaturesSetFlagCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "set feature flag",
			Args: []string{"--course-id", "1", "--feature", "new_quizzes", "--state", "on"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features/flags/new_quizzes": cmdtest.NewMockResponse(`{
					"feature": "new_quizzes", "state": "on"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "set feature flag - missing course-id",
			Args:        []string{"--feature", "new_quizzes", "--state", "on"},
			ExpectError: true,
		},
		{
			Name:        "set feature flag - missing feature",
			Args:        []string{"--course-id", "1", "--state", "on"},
			ExpectError: true,
		},
		{
			Name:        "set feature flag - missing state",
			Args:        []string{"--course-id", "1", "--feature", "new_quizzes"},
			ExpectError: true,
		},
		{
			Name: "set feature flag - API error",
			Args: []string{"--course-id", "1", "--feature", "new_quizzes", "--state", "on"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features/flags/new_quizzes": cmdtest.NewErrorResponse(422, "invalid state"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseFeaturesSetFlagCmd(), tc)
		})
	}
}

func TestCovCourse_CourseFeaturesDeleteFlagCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete feature flag",
			Args: []string{"--course-id", "1", "--feature", "new_quizzes"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features/flags/new_quizzes": cmdtest.NewMockResponse(`{
					"feature": "new_quizzes", "state": "allowed"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete feature flag - missing course-id",
			Args:        []string{"--feature", "new_quizzes"},
			ExpectError: true,
		},
		{
			Name:        "delete feature flag - missing feature",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "delete feature flag - API error",
			Args: []string{"--course-id", "1", "--feature", "new_quizzes"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/features/flags/new_quizzes": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCourseFeaturesDeleteFlagCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// content_exports.go
// ---------------------------------------------------------------------------

func TestCovCourse_ContentExportsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list content exports",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_exports": cmdtest.NewMockResponse(`[
					{"id": 1, "export_type": "common_cartridge", "workflow_state": "exported"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list content exports - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_exports": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No content exports found",
		},
		{
			Name:        "list content exports - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list content exports - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_exports": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newContentExportsListCmd(), tc)
		})
	}
}

func TestCovCourse_ContentExportsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get content export",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_exports/5": cmdtest.NewMockResponse(`{
					"id": 5, "export_type": "common_cartridge", "workflow_state": "exported"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get content export - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get content export - invalid ID",
			Args:        []string{"notanid", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get content export - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "get content export - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_exports/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newContentExportsGetCmd(), tc)
		})
	}
}

func TestCovCourse_ContentExportsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create content export",
			Args: []string{"--course-id", "1", "--export-type", "common_cartridge"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_exports": cmdtest.NewMockResponse(`{
					"id": 10, "export_type": "common_cartridge", "workflow_state": "created"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "create content export - missing course-id",
			Args:        []string{"--export-type", "common_cartridge"},
			ExpectError: true,
		},
		{
			Name:        "create content export - missing export-type",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "create content export - API error",
			Args: []string{"--course-id", "1", "--export-type", "common_cartridge"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/content_exports": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newContentExportsCreateCmd(), tc)
		})
	}
}

func TestCovCourse_EpubExportsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create epub export",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/epub_exports": cmdtest.NewMockResponse(`{
					"id": 3, "workflow_state": "created"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "create epub export - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "create epub export - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/epub_exports": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEpubExportsCreateCmd(), tc)
		})
	}
}

func TestCovCourse_EpubExportsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get epub export",
			Args: []string{"3", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/epub_exports/3": cmdtest.NewMockResponse(`{
					"id": 3, "workflow_state": "generated"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get epub export - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get epub export - invalid ID",
			Args:        []string{"notanid", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get epub export - missing course-id",
			Args:        []string{"3"},
			ExpectError: true,
		},
		{
			Name: "get epub export - API error",
			Args: []string{"3", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/epub_exports/3": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEpubExportsGetCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// blackout_dates.go
// ---------------------------------------------------------------------------

func TestCovCourse_BlackoutDatesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list blackout dates",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates": cmdtest.NewMockResponse(`[
					{"id": 1, "event_title": "Winter Break", "start_date": "2024-12-24", "end_date": "2024-12-26"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list blackout dates - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No blackout dates found",
		},
		{
			Name:        "list blackout dates - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list blackout dates - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newBlackoutDatesListCmd(), tc)
		})
	}
}

func TestCovCourse_BlackoutDatesGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get blackout date",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates/5": cmdtest.NewMockResponse(`{
					"id": 5, "event_title": "Holiday"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get blackout date - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get blackout date - invalid ID",
			Args:        []string{"abc", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get blackout date - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "get blackout date - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newBlackoutDatesGetCmd(), tc)
		})
	}
}

func TestCovCourse_BlackoutDatesCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create blackout date",
			Args: []string{
				"--course-id", "1",
				"--start-date", "2024-12-24",
				"--end-date", "2024-12-26",
				"--title", "Winter Break",
			},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates": cmdtest.NewMockResponse(`{
					"id": 1, "event_title": "Winter Break"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "create blackout date - missing course-id",
			Args:        []string{"--start-date", "2024-12-24", "--end-date", "2024-12-26", "--title", "Break"},
			ExpectError: true,
		},
		{
			Name:        "create blackout date - missing start-date",
			Args:        []string{"--course-id", "1", "--end-date", "2024-12-26", "--title", "Break"},
			ExpectError: true,
		},
		{
			Name:        "create blackout date - missing end-date",
			Args:        []string{"--course-id", "1", "--start-date", "2024-12-24", "--title", "Break"},
			ExpectError: true,
		},
		{
			Name:        "create blackout date - missing title",
			Args:        []string{"--course-id", "1", "--start-date", "2024-12-24", "--end-date", "2024-12-26"},
			ExpectError: true,
		},
		{
			Name: "create blackout date - API error",
			Args: []string{"--course-id", "1", "--start-date", "2024-12-24", "--end-date", "2024-12-26", "--title", "Break"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newBlackoutDatesCreateCmd(), tc)
		})
	}
}

func TestCovCourse_BlackoutDatesUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update blackout date",
			Args: []string{"5", "--course-id", "1", "--title", "Extended Break"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates/5": cmdtest.NewMockResponse(`{
					"id": 5, "event_title": "Extended Break"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "update blackout date - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "update blackout date - invalid ID",
			Args:        []string{"notanid", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "update blackout date - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "update blackout date - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates/5": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newBlackoutDatesUpdateCmd(), tc)
		})
	}
}

func TestCovCourse_BlackoutDatesDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete blackout date",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates/5": cmdtest.NewMockResponse(`{
					"id": 5, "event_title": "Holiday"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete blackout date - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "delete blackout date - invalid ID",
			Args:        []string{"abc", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "delete blackout date - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "delete blackout date - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/blackout_dates/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newBlackoutDatesDeleteCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// rubric_associations.go
// ---------------------------------------------------------------------------

func TestCovCourse_RubricAssociationsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update rubric association",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/rubric_associations/5": cmdtest.NewMockResponse(`{
					"id": 5, "rubric_id": 10
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "update rubric association - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "update rubric association - invalid ID",
			Args:        []string{"abc", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "update rubric association - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "update rubric association - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/rubric_associations/5": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newRubricAssociationsUpdateCmd(), tc)
		})
	}
}

func TestCovCourse_RubricAssociationsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete rubric association",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/rubric_associations/5": cmdtest.NewMockResponse(`{
					"id": 5
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete rubric association - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "delete rubric association - invalid ID",
			Args:        []string{"abc", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "delete rubric association - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "delete rubric association - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/rubric_associations/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newRubricAssociationsDeleteCmd(), tc)
		})
	}
}

func TestCovCourse_RubricAssessmentsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create rubric assessment",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/rubric_associations/5/rubric_assessments": cmdtest.NewMockResponse(`{
					"id": 20, "artifact_id": 0
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "create rubric assessment - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "create rubric assessment - invalid ID",
			Args:        []string{"abc", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "create rubric assessment - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "create rubric assessment - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/rubric_associations/5/rubric_assessments": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newRubricAssessmentsCreateCmd(), tc)
		})
	}
}

func TestCovCourse_RubricAssessmentsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete rubric assessment",
			Args: []string{"5", "--course-id", "1", "--assessment-id", "20"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/rubric_associations/5/rubric_assessments/20": cmdtest.NewMockResponse(`{
					"id": 20
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete rubric assessment - missing ID arg",
			Args:        []string{"--course-id", "1", "--assessment-id", "20"},
			ExpectError: true,
		},
		{
			Name:        "delete rubric assessment - invalid association ID",
			Args:        []string{"abc", "--course-id", "1", "--assessment-id", "20"},
			ExpectError: true,
		},
		{
			Name:        "delete rubric assessment - missing course-id",
			Args:        []string{"5", "--assessment-id", "20"},
			ExpectError: true,
		},
		{
			Name:        "delete rubric assessment - missing assessment-id",
			Args:        []string{"5", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "delete rubric assessment - API error",
			Args: []string{"5", "--course-id", "1", "--assessment-id", "20"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/rubric_associations/5/rubric_assessments/20": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newRubricAssessmentsDeleteCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// media.go
// ---------------------------------------------------------------------------

func TestCovCourse_MediaObjectsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list media objects (no context)",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_objects": cmdtest.NewMockResponse(`[
					{"media_id": "m-abc123", "title": "Intro Video"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list media objects for course",
			Args: []string{"--course-id", "123"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/123/media_objects": cmdtest.NewMockResponse(`[
					{"media_id": "m-xyz", "title": "Lecture 1"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list media objects for group",
			Args: []string{"--group-id", "456"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/456/media_objects": cmdtest.NewMockResponse(`[
					{"media_id": "m-grp", "title": "Group Video"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list media objects - empty",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_objects": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No media objects found",
		},
		{
			Name: "list media objects - API error",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_objects": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newMediaObjectsListCmd(), tc)
		})
	}
}

func TestCovCourse_MediaObjectUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update media object",
			Args: []string{"m-abc123", "--title", "New Title"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_objects/m-abc123": cmdtest.NewMockResponse(`{
					"media_id": "m-abc123", "title": "New Title"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "update media object - missing ID arg",
			Args:        []string{"--title", "New Title"},
			ExpectError: true,
		},
		{
			Name:        "update media object - missing title",
			Args:        []string{"m-abc123"},
			ExpectError: true,
		},
		{
			Name: "update media object - API error",
			Args: []string{"m-abc123", "--title", "New Title"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_objects/m-abc123": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newMediaObjectUpdateCmd(), tc)
		})
	}
}

func TestCovCourse_MediaObjectTracksCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get media object tracks",
			Args: []string{"m-abc123"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_objects/m-abc123/media_tracks": cmdtest.NewMockResponse(`[
					{"id": 1, "kind": "subtitles", "locale": "en"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "get media object tracks - empty",
			Args: []string{"m-abc123"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_objects/m-abc123/media_tracks": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No media tracks found",
		},
		{
			Name:        "get media object tracks - missing ID arg",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "get media object tracks - API error",
			Args: []string{"m-abc123"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_objects/m-abc123/media_tracks": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newMediaObjectTracksCmd(), tc)
		})
	}
}

func TestCovCourse_MediaAttachmentsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list media attachments (no context)",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_attachments": cmdtest.NewMockResponse(`[
					{"id": 1, "display_name": "lecture.mp4"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list media attachments for course",
			Args: []string{"--course-id", "123"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/123/media_attachments": cmdtest.NewMockResponse(`[
					{"id": 2, "display_name": "course_video.mp4"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list media attachments for group",
			Args: []string{"--group-id", "456"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/456/media_attachments": cmdtest.NewMockResponse(`[
					{"id": 3, "display_name": "group.mp4"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list media attachments - empty",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_attachments": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No media attachments found",
		},
		{
			Name: "list media attachments - API error",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_attachments": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newMediaAttachmentsListCmd(), tc)
		})
	}
}

func TestCovCourse_MediaAttachmentUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update media attachment",
			Args: []string{"5", "--title", "New Title"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_attachments/5": cmdtest.NewMockResponse(`{
					"id": 5, "display_name": "New Title"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "update media attachment - missing ID arg",
			Args:        []string{"--title", "New Title"},
			ExpectError: true,
		},
		{
			Name:        "update media attachment - invalid ID",
			Args:        []string{"notanid", "--title", "New Title"},
			ExpectError: true,
		},
		{
			Name:        "update media attachment - missing title",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "update media attachment - API error",
			Args: []string{"5", "--title", "New Title"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_attachments/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newMediaAttachmentUpdateCmd(), tc)
		})
	}
}

func TestCovCourse_MediaAttachmentTracksCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get media attachment tracks",
			Args: []string{"5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_attachments/5/media_tracks": cmdtest.NewMockResponse(`[
					{"id": 1, "kind": "subtitles", "locale": "en"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "get media attachment tracks - empty",
			Args: []string{"5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_attachments/5/media_tracks": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No media tracks found",
		},
		{
			Name:        "get media attachment tracks - missing ID arg",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name:        "get media attachment tracks - invalid ID",
			Args:        []string{"notanid"},
			ExpectError: true,
		},
		{
			Name: "get media attachment tracks - API error",
			Args: []string{"5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/media_attachments/5/media_tracks": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newMediaAttachmentTracksCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// course_pacing.go
// ---------------------------------------------------------------------------

func TestCovCourse_CoursePacingGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get course pace",
			Args: []string{"1", "--course-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/5/course_pacing/1": cmdtest.NewMockResponse(`{
					"id": 1, "workflow_state": "active"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get course pace - missing ID arg",
			Args:        []string{"--course-id", "5"},
			ExpectError: true,
		},
		{
			Name:        "get course pace - invalid ID",
			Args:        []string{"notanid", "--course-id", "5"},
			ExpectError: true,
		},
		{
			Name:        "get course pace - missing course-id",
			Args:        []string{"1"},
			ExpectError: true,
		},
		{
			Name: "get course pace - API error",
			Args: []string{"1", "--course-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/5/course_pacing/1": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCoursePacingGetCmd(), tc)
		})
	}
}

func TestCovCourse_CoursePacingCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create course pace",
			Args: []string{"--course-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/5/course_pacing": cmdtest.NewMockResponse(`{
					"id": 1, "workflow_state": "active"
				}`),
			},
			ExpectError: false,
		},
		{
			Name: "create course pace with options",
			Args: []string{"--course-id", "5", "--exclude-weekends", "--hard-end-dates", "--end-date", "2024-12-31"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/5/course_pacing": cmdtest.NewMockResponse(`{
					"id": 2, "workflow_state": "active"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "create course pace - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "create course pace - API error",
			Args: []string{"--course-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/5/course_pacing": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCoursePacingCreateCmd(), tc)
		})
	}
}

func TestCovCourse_CoursePacingUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update course pace",
			Args: []string{"1", "--course-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/5/course_pacing/1": cmdtest.NewMockResponse(`{
					"id": 1, "workflow_state": "active"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "update course pace - missing ID arg",
			Args:        []string{"--course-id", "5"},
			ExpectError: true,
		},
		{
			Name:        "update course pace - invalid ID",
			Args:        []string{"notanid", "--course-id", "5"},
			ExpectError: true,
		},
		{
			Name:        "update course pace - missing course-id",
			Args:        []string{"1"},
			ExpectError: true,
		},
		{
			Name: "update course pace - API error",
			Args: []string{"1", "--course-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/5/course_pacing/1": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCoursePacingUpdateCmd(), tc)
		})
	}
}

func TestCovCourse_CoursePacingDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete course pace",
			Args: []string{"1", "--course-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/5/course_pacing/1": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete course pace - missing ID arg",
			Args:        []string{"--course-id", "5"},
			ExpectError: true,
		},
		{
			Name:        "delete course pace - invalid ID",
			Args:        []string{"notanid", "--course-id", "5"},
			ExpectError: true,
		},
		{
			Name:        "delete course pace - missing course-id",
			Args:        []string{"1"},
			ExpectError: true,
		},
		{
			Name: "delete course pace - API error",
			Args: []string{"1", "--course-id", "5"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/5/course_pacing/1": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newCoursePacingDeleteCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// live_assessments.go
// ---------------------------------------------------------------------------

func TestCovCourse_LiveAssessmentsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list live assessments",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/live_assessments": cmdtest.NewMockResponse(`{
					"live_assessments": [{"id": "abc", "title": "Quiz 1"}]
				}`),
			},
			ExpectError: false,
		},
		{
			Name: "list live assessments - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/live_assessments": cmdtest.NewMockResponse(`{
					"live_assessments": []
				}`),
			},
			ExpectError:  false,
			ExpectOutput: "No live assessments found",
		},
		{
			Name:        "list live assessments - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list live assessments - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/live_assessments": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newLiveAssessmentsListCmd(), tc)
		})
	}
}

func TestCovCourse_LiveAssessmentsResultsCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list live assessment results",
			Args: []string{"--course-id", "1", "--assessment-id", "abc"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/live_assessments/abc/results": cmdtest.NewMockResponse(`{
					"results": [{"passed": true}]
				}`),
			},
			ExpectError: false,
		},
		{
			Name: "list live assessment results - empty",
			Args: []string{"--course-id", "1", "--assessment-id", "abc"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/live_assessments/abc/results": cmdtest.NewMockResponse(`{
					"results": []
				}`),
			},
			ExpectError:  false,
			ExpectOutput: "No results found",
		},
		{
			Name:        "list live assessment results - missing course-id",
			Args:        []string{"--assessment-id", "abc"},
			ExpectError: true,
		},
		{
			Name:        "list live assessment results - missing assessment-id",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "list live assessment results - API error",
			Args: []string{"--course-id", "1", "--assessment-id", "abc"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/live_assessments/abc/results": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newLiveAssessmentsResultsCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// grading_standards.go
// ---------------------------------------------------------------------------

func TestCovCourse_GradingStandardsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list grading standards",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_standards": cmdtest.NewMockResponse(`[
					{"id": 5, "title": "Letter Grades"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list grading standards - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_standards": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No grading standards found",
		},
		{
			Name:        "list grading standards - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list grading standards - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_standards": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newGradingStandardsListCmd(), tc)
		})
	}
}

func TestCovCourse_GradingStandardsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get grading standard",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_standards/5": cmdtest.NewMockResponse(`{
					"id": 5, "title": "Letter Grades"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get grading standard - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get grading standard - invalid ID",
			Args:        []string{"abc", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get grading standard - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "get grading standard - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_standards/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newGradingStandardsGetCmd(), tc)
		})
	}
}

func TestCovCourse_GradingStandardsCreateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "create grading standard",
			Args: []string{"--course-id", "1", "--title", "Letter Grades"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_standards": cmdtest.NewMockResponse(`{
					"id": 5, "title": "Letter Grades"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "create grading standard - missing course-id",
			Args:        []string{"--title", "Letter Grades"},
			ExpectError: true,
		},
		{
			Name:        "create grading standard - missing title",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name: "create grading standard - API error",
			Args: []string{"--course-id", "1", "--title", "Letter Grades"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_standards": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newGradingStandardsCreateCmd(), tc)
		})
	}
}

func TestCovCourse_GradingStandardsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete grading standard",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_standards/5": cmdtest.NewMockResponse(`{
					"id": 5
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete grading standard - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "delete grading standard - invalid ID",
			Args:        []string{"abc", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "delete grading standard - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "delete grading standard - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_standards/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newGradingStandardsDeleteCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// grading_periods.go
// ---------------------------------------------------------------------------

func TestCovCourse_GradingPeriodsListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list grading periods",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_periods": cmdtest.NewMockResponse(`{
					"grading_periods": [{"id": 5, "title": "Q1"}]
				}`),
			},
			ExpectError: false,
		},
		{
			Name: "list grading periods - empty",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_periods": cmdtest.NewMockResponse(`{
					"grading_periods": []
				}`),
			},
			ExpectError:  false,
			ExpectOutput: "No grading periods found",
		},
		{
			Name:        "list grading periods - missing course-id",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name: "list grading periods - API error",
			Args: []string{"--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_periods": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newGradingPeriodsListCmd(), tc)
		})
	}
}

func TestCovCourse_GradingPeriodsGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get grading period",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_periods/5": cmdtest.NewMockResponse(`{
					"grading_periods": [{"id": 5, "title": "Q1"}]
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get grading period - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get grading period - invalid ID",
			Args:        []string{"abc", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "get grading period - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "get grading period - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_periods/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newGradingPeriodsGetCmd(), tc)
		})
	}
}

func TestCovCourse_GradingPeriodsUpdateCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "update grading period",
			Args: []string{"5", "--course-id", "1", "--title", "Q1 Updated"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_periods/5": cmdtest.NewMockResponse(`{
					"grading_periods": [{"id": 5, "title": "Q1 Updated"}]
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "update grading period - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "update grading period - invalid ID",
			Args:        []string{"abc", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "update grading period - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "update grading period - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_periods/5": cmdtest.NewErrorResponse(422, "invalid"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newGradingPeriodsUpdateCmd(), tc)
		})
	}
}

func TestCovCourse_GradingPeriodsDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete grading period",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_periods/5": cmdtest.NewMockResponse(`{}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete grading period - missing ID arg",
			Args:        []string{"--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "delete grading period - invalid ID",
			Args:        []string{"abc", "--course-id", "1"},
			ExpectError: true,
		},
		{
			Name:        "delete grading period - missing course-id",
			Args:        []string{"5"},
			ExpectError: true,
		},
		{
			Name: "delete grading period - API error",
			Args: []string{"5", "--course-id", "1"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/1/grading_periods/5": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newGradingPeriodsDeleteCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// conferences.go
// ---------------------------------------------------------------------------

func TestCovCourse_ConferencesListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list conferences (no context)",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conferences": cmdtest.NewMockResponse(`[
					{"id": 1, "title": "Office Hours", "conference_type": "BigBlueButton"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list conferences for course",
			Args: []string{"--course-id", "123"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/courses/123/conferences": cmdtest.NewMockResponse(`[
					{"id": 2, "title": "Lecture", "conference_type": "Zoom"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list conferences for group",
			Args: []string{"--group-id", "456"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/groups/456/conferences": cmdtest.NewMockResponse(`[
					{"id": 3, "title": "Study Session"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list conferences with state filter",
			Args: []string{"--state", "live"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conferences": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No conferences found",
		},
		{
			Name: "list conferences - API error",
			Args: []string{},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/conferences": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newConferencesListCmd(), tc)
		})
	}
}

// ---------------------------------------------------------------------------
// eportfolios.go
// ---------------------------------------------------------------------------

func TestCovCourse_EportfoliosListCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list eportfolios for user",
			Args: []string{"--user-id", "42"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/42/eportfolios": cmdtest.NewMockResponse(`[
					{"id": 1, "name": "My Portfolio"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list eportfolios - empty",
			Args: []string{"--user-id", "42"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/users/42/eportfolios": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No ePortfolios found",
		},
		{
			Name:        "list eportfolios - missing user-id",
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

func TestCovCourse_EportfolioGetCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "get eportfolio",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10": cmdtest.NewMockResponse(`{
					"id": 10, "name": "My Portfolio"
				}`),
			},
			ExpectError: false,
		},
		{
			Name:        "get eportfolio - missing ID arg",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name:        "get eportfolio - invalid ID",
			Args:        []string{"notanid"},
			ExpectError: true,
		},
		{
			Name: "get eportfolio - API error",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEportfolioGetCmd(), tc)
		})
	}
}

func TestCovCourse_EportfolioDeleteCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "delete eportfolio with force",
			Args: []string{"10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10": cmdtest.NewMockResponse(`{"id": 10}`),
			},
			ExpectError: false,
		},
		{
			Name:        "delete eportfolio - missing ID arg",
			Args:        []string{"--force"},
			ExpectError: true,
		},
		{
			Name:        "delete eportfolio - invalid ID",
			Args:        []string{"bad"},
			ExpectError: true,
		},
		{
			Name: "delete eportfolio - API error",
			Args: []string{"10", "--force"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10": cmdtest.NewErrorResponse(404, "not found"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEportfolioDeleteCmd(), tc)
		})
	}
}

func TestCovCourse_EportfolioPagesCmd(t *testing.T) {
	tests := []cmdtest.CommandTestCase{
		{
			Name: "list eportfolio pages",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10/pages": cmdtest.NewMockResponse(`[
					{"id": 1, "name": "Introduction"}
				]`),
			},
			ExpectError: false,
		},
		{
			Name: "list eportfolio pages - empty",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10/pages": cmdtest.NewMockResponse(`[]`),
			},
			ExpectError:  false,
			ExpectOutput: "No pages found",
		},
		{
			Name:        "list eportfolio pages - missing ID arg",
			Args:        []string{},
			ExpectError: true,
		},
		{
			Name:        "list eportfolio pages - invalid ID",
			Args:        []string{"notanid"},
			ExpectError: true,
		},
		{
			Name: "list eportfolio pages - API error",
			Args: []string{"10"},
			MockResponses: map[string]cmdtest.MockResponse{
				"/api/v1/eportfolios/10/pages": cmdtest.NewErrorResponse(500, "server error"),
			},
			ExpectError: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			cmdtest.RunCommandTest(t, newEportfolioPagesCmd(), tc)
		})
	}
}
