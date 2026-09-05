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
