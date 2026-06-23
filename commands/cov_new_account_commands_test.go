package commands

// cov_new_account_commands_test.go — coverage tests for the Wave-4 account-admin
// command files that were added in feature/spec-compliance:
//   auth_providers, account_features, grading_period_sets, account_reports,
//   enrollment_terms, account_notifications, account_content_migrations,
//   temporary_enrollment_pairings, account_logins, developer_keys,
//   csp_settings, account_analytics
//
// Each run* function is driven through its happy path (200 OK mock) and at
// least one error path (500 error mock).  Prefix TestCovAcct_ prevents
// collisions with existing test functions.

import (
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
)

// ---------------------------------------------------------------------------
// auth_providers
// ---------------------------------------------------------------------------

func TestCovAcct_AuthProvidersList_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list auth providers - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers": cmdtest.NewMockResponse(`[
				{"id":10,"auth_type":"saml","position":1,"workflow_state":"active"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAuthProvidersListCmd(), tc)
}

func TestCovAcct_AuthProvidersList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list auth providers - api error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAuthProvidersListCmd(), tc)
}

func TestCovAcct_AuthProvidersList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list auth providers - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newAuthProvidersListCmd(), tc)
}

func TestCovAcct_AuthProvidersGet_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get auth provider - ok",
		Args: []string{"1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers/10": cmdtest.NewMockResponse(`{
				"id":10,"auth_type":"saml","position":1,"workflow_state":"active","jit_provisioning":true
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAuthProvidersGetCmd(), tc)
}

func TestCovAcct_AuthProvidersGet_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get auth provider - error",
		Args: []string{"1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers/10": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAuthProvidersGetCmd(), tc)
}

func TestCovAcct_AuthProvidersCreate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create auth provider - ok",
		Args: []string{"1", "--auth-type", "saml"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers": cmdtest.NewMockResponse(`{
				"id":20,"auth_type":"saml","position":1,"workflow_state":"active"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAuthProvidersCreateCmd(), tc)
}

func TestCovAcct_AuthProvidersCreate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create auth provider - error",
		Args: []string{"1", "--auth-type", "saml"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers": cmdtest.NewErrorResponse(422, "invalid"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAuthProvidersCreateCmd(), tc)
}

func TestCovAcct_AuthProvidersDelete_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete auth provider - ok",
		Args: []string{"1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers/10": cmdtest.NewMockResponse(`{}`),
		},
	}
	cmdtest.RunCommandTest(t, newAuthProvidersDeleteCmd(), tc)
}

func TestCovAcct_AuthProvidersDelete_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete auth provider - error",
		Args: []string{"1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers/10": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAuthProvidersDeleteCmd(), tc)
}

func TestCovAcct_AuthProvidersRestore_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "restore auth provider - ok",
		Args: []string{"1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers/10": cmdtest.NewMockResponse(`{
				"id":10,"auth_type":"saml","position":1,"workflow_state":"active"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAuthProvidersRestoreCmd(), tc)
}

func TestCovAcct_AuthProvidersRestore_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "restore auth provider - error",
		Args: []string{"1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers/10": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAuthProvidersRestoreCmd(), tc)
}

func TestCovAcct_AuthProvidersForcePasswordReset_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "force password reset - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers/force_password_reset": cmdtest.NewMockResponse(`{}`),
		},
	}
	cmdtest.RunCommandTest(t, newAuthProvidersForcePasswordResetCmd(), tc)
}

func TestCovAcct_AuthProvidersForcePasswordReset_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "force password reset - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/authentication_providers/force_password_reset": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAuthProvidersForcePasswordResetCmd(), tc)
}

func TestCovAcct_AuthProvidersSSOSettings_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "sso settings - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/sso_settings": cmdtest.NewMockResponse(`{
				"login_handle_name":"Email","change_password_url":"","auth_discovery_url":"","unknown_user_url":""
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAuthProvidersSSOSettingsCmd(), tc)
}

func TestCovAcct_AuthProvidersSSOSettings_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "sso settings - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/sso_settings": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAuthProvidersSSOSettingsCmd(), tc)
}

// SSO settings with optional fields populated
func TestCovAcct_AuthProvidersSSOSettings_WithURLs(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "sso settings - with optional urls",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/sso_settings": cmdtest.NewMockResponse(`{
				"login_handle_name":"Email",
				"change_password_url":"https://example.com/pw",
				"auth_discovery_url":"https://example.com/auth",
				"unknown_user_url":"https://example.com/unknown"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAuthProvidersSSOSettingsCmd(), tc)
}

// ---------------------------------------------------------------------------
// account_features
// ---------------------------------------------------------------------------

func TestCovAcct_AccountFeaturesList_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account features - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features": cmdtest.NewMockResponse(`[
				{"feature":"analytics_2","display_name":"Analytics 2","applies_to":"Account"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesListCmd(), tc)
}

func TestCovAcct_AccountFeaturesList_LongDisplayName(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account features - long display name",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features": cmdtest.NewMockResponse(`[
				{"feature":"x","display_name":"A Very Long Feature Display Name That Exceeds The Limit","applies_to":"Account"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesListCmd(), tc)
}

func TestCovAcct_AccountFeaturesList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account features - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesListCmd(), tc)
}

func TestCovAcct_AccountFeaturesListEnabled_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list enabled account features - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features/enabled": cmdtest.NewMockResponse(`["analytics_2","course_pacing"]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesListEnabledCmd(), tc)
}

func TestCovAcct_AccountFeaturesListEnabled_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list enabled account features - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features/enabled": cmdtest.NewErrorResponse(500, "server error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesListEnabledCmd(), tc)
}

func TestCovAcct_AccountFeaturesGetFlag_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get feature flag - ok",
		Args: []string{"1", "analytics_2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features/flags/analytics_2": cmdtest.NewMockResponse(`{
				"feature":"analytics_2","state":"on","transition_locked":false
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesGetFlagCmd(), tc)
}

func TestCovAcct_AccountFeaturesGetFlag_WithLockedAt(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get feature flag - locked at set",
		Args: []string{"1", "analytics_2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features/flags/analytics_2": cmdtest.NewMockResponse(`{
				"feature":"analytics_2","state":"on","transition_locked":true,"locked_at":"2026-01-01T00:00:00Z"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesGetFlagCmd(), tc)
}

func TestCovAcct_AccountFeaturesGetFlag_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get feature flag - error",
		Args: []string{"1", "analytics_2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features/flags/analytics_2": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesGetFlagCmd(), tc)
}

func TestCovAcct_AccountFeaturesSetFlag_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "set feature flag - ok",
		Args: []string{"1", "analytics_2", "--state", "on"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features/flags/analytics_2": cmdtest.NewMockResponse(`{
				"feature":"analytics_2","state":"on","transition_locked":false
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesSetFlagCmd(), tc)
}

func TestCovAcct_AccountFeaturesSetFlag_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "set feature flag - error",
		Args: []string{"1", "analytics_2", "--state", "on"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features/flags/analytics_2": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesSetFlagCmd(), tc)
}

func TestCovAcct_AccountFeaturesDeleteFlag_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete feature flag - ok",
		Args: []string{"1", "analytics_2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features/flags/analytics_2": cmdtest.NewMockResponse(`{}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesDeleteFlagCmd(), tc)
}

func TestCovAcct_AccountFeaturesDeleteFlag_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete feature flag - error",
		Args: []string{"1", "analytics_2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/features/flags/analytics_2": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesDeleteFlagCmd(), tc)
}

func TestCovAcct_AccountFeaturesSettings_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "account settings - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/settings": cmdtest.NewMockResponse(`{
				"restrict_student_past_view":false,
				"restrict_student_future_view":false,
				"hide_distribution_graphs":false,
				"lock_all_announcements":false,
				"usage_rights_required":false,
				"default_due_date_restricted":false
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesSettingsCmd(), tc)
}

func TestCovAcct_AccountFeaturesSettings_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "account settings - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/settings": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesSettingsCmd(), tc)
}

func TestCovAcct_AccountFeaturesPermissions_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "account permissions - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/permissions": cmdtest.NewMockResponse(`{
				"manage_courses":true,"manage_users":false
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesPermissionsCmd(), tc)
}

func TestCovAcct_AccountFeaturesPermissions_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "account permissions - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/permissions": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountFeaturesPermissionsCmd(), tc)
}

// ---------------------------------------------------------------------------
// grading_period_sets
// ---------------------------------------------------------------------------

func TestCovAcct_GradingPeriodSetsList_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list grading period sets - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_period_sets": cmdtest.NewMockResponse(`[
				{"id":5,"title":"2024-2025","weighted_grading_periods":false}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsListCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list grading period sets - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_period_sets": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsListCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list grading period sets - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_period_sets": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsListCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsCreate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create grading period set - ok",
		Args: []string{"1", "--title", "2024-2025"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_period_sets": cmdtest.NewMockResponse(`{
				"id":5,"title":"2024-2025","weighted_grading_periods":false
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsCreateCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsCreate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create grading period set - error",
		Args: []string{"1", "--title", "2024-2025"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_period_sets": cmdtest.NewErrorResponse(422, "invalid"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsCreateCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsUpdate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update grading period set - ok",
		Args: []string{"1", "5", "--title", "2024-2025 Updated"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_period_sets/5": cmdtest.NewMockResponse(`{
				"id":5,"title":"2024-2025 Updated","weighted_grading_periods":false
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsUpdateCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsUpdate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update grading period set - error",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_period_sets/5": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsUpdateCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsDelete_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete grading period set - ok",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_period_sets/5": cmdtest.NewMockResponse(`{}`),
		},
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsDeleteCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsDelete_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete grading period set - error",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_period_sets/5": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsDeleteCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsListPeriods_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list grading periods - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_periods": cmdtest.NewMockResponse(`[
				{"id":1,"title":"Q1","start_date":"2024-08-01","end_date":"2024-10-31"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsListPeriodsCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsListPeriods_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list grading periods - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_periods": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsListPeriodsCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsListPeriods_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list grading periods - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_periods": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsListPeriodsCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsDeletePeriod_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete grading period - ok",
		Args: []string{"1", "3"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_periods/3": cmdtest.NewMockResponse(`{}`),
		},
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsDeletePeriodCmd(), tc)
}

func TestCovAcct_GradingPeriodSetsDeletePeriod_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete grading period - error",
		Args: []string{"1", "3"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/grading_periods/3": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newGradingPeriodSetsDeletePeriodCmd(), tc)
}

// ---------------------------------------------------------------------------
// account_reports
// ---------------------------------------------------------------------------

func TestCovAcct_AccountReportsList_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account reports - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports": cmdtest.NewMockResponse(`[
				{"report":"course_storage_csv","title":"Course Storage"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountReportsListCmd(), tc)
}

func TestCovAcct_AccountReportsList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account reports - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountReportsListCmd(), tc)
}

func TestCovAcct_AccountReportsList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account reports - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountReportsListCmd(), tc)
}

func TestCovAcct_AccountReportsRuns_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list report runs - ok",
		Args: []string{"1", "course_storage_csv"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv": cmdtest.NewMockResponse(`[
				{"id":10,"report":"course_storage_csv","status":"complete","progress":100}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountReportsRunsCmd(), tc)
}

func TestCovAcct_AccountReportsRuns_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list report runs - empty",
		Args: []string{"1", "course_storage_csv"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountReportsRunsCmd(), tc)
}

func TestCovAcct_AccountReportsRuns_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list report runs - error",
		Args: []string{"1", "course_storage_csv"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountReportsRunsCmd(), tc)
}

func TestCovAcct_AccountReportsStart_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "start report - ok",
		Args: []string{"1", "course_storage_csv"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv": cmdtest.NewMockResponse(`{
				"id":11,"report":"course_storage_csv","status":"created","progress":0
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountReportsStartCmd(), tc)
}

func TestCovAcct_AccountReportsStart_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "start report - error",
		Args: []string{"1", "course_storage_csv"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountReportsStartCmd(), tc)
}

func TestCovAcct_AccountReportsGet_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get report run - ok",
		Args: []string{"1", "course_storage_csv", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv/10": cmdtest.NewMockResponse(`{
				"id":10,"report":"course_storage_csv","status":"complete","progress":100,
				"created_at":"2026-01-01T00:00:00Z",
				"started_at":"2026-01-01T00:00:01Z",
				"ended_at":"2026-01-01T00:01:00Z",
				"attachment":{"url":"https://example.com/file.csv"},
				"parameters":{"param1":"value1"}
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountReportsGetCmd(), tc)
}

func TestCovAcct_AccountReportsGet_MinimalResponse(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get report run - minimal response",
		Args: []string{"1", "course_storage_csv", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv/10": cmdtest.NewMockResponse(`{
				"id":10,"report":"course_storage_csv","status":"running","progress":50,
				"created_at":"2026-01-01T00:00:00Z"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountReportsGetCmd(), tc)
}

func TestCovAcct_AccountReportsGet_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get report run - error",
		Args: []string{"1", "course_storage_csv", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv/10": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountReportsGetCmd(), tc)
}

func TestCovAcct_AccountReportsDelete_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete report run - ok",
		Args: []string{"1", "course_storage_csv", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv/10": cmdtest.NewMockResponse(`{}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountReportsDeleteCmd(), tc)
}

func TestCovAcct_AccountReportsDelete_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete report run - error",
		Args: []string{"1", "course_storage_csv", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv/10": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountReportsDeleteCmd(), tc)
}

func TestCovAcct_AccountReportsAbort_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "abort report run - ok",
		Args: []string{"1", "course_storage_csv", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv/10": cmdtest.NewMockResponse(`{
				"id":10,"report":"course_storage_csv","status":"aborted","progress":50
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountReportsAbortCmd(), tc)
}

func TestCovAcct_AccountReportsAbort_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "abort report run - error",
		Args: []string{"1", "course_storage_csv", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/reports/course_storage_csv/10": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountReportsAbortCmd(), tc)
}

// ---------------------------------------------------------------------------
// enrollment_terms
// ---------------------------------------------------------------------------

func TestCovAcct_EnrollmentTermsList_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list enrollment terms - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms": cmdtest.NewMockResponse(`{
				"enrollment_terms": [
					{"id":42,"name":"Fall 2025","start_at":"2025-08-01","end_at":"2025-12-15"}
				]
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsListCmd(), tc)
}

func TestCovAcct_EnrollmentTermsList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list enrollment terms - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms": cmdtest.NewMockResponse(`{"enrollment_terms":[]}`),
		},
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsListCmd(), tc)
}

func TestCovAcct_EnrollmentTermsList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list enrollment terms - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsListCmd(), tc)
}

func TestCovAcct_EnrollmentTermsGet_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get enrollment term - ok",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms/42": cmdtest.NewMockResponse(`{
				"id":42,"name":"Fall 2025","start_at":"2025-08-01","end_at":"2025-12-15"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsGetCmd(), tc)
}

func TestCovAcct_EnrollmentTermsGet_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get enrollment term - error",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms/42": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsGetCmd(), tc)
}

func TestCovAcct_EnrollmentTermsCreate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create enrollment term - ok",
		Args: []string{"1", "--name", "Fall 2025", "--start-at", "2025-08-01", "--end-at", "2025-12-15"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms": cmdtest.NewMockResponse(`{
				"id":42,"name":"Fall 2025","start_at":"2025-08-01","end_at":"2025-12-15"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsCreateCmd(), tc)
}

func TestCovAcct_EnrollmentTermsCreate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create enrollment term - error",
		Args: []string{"1", "--name", "Fall 2025"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms": cmdtest.NewErrorResponse(422, "invalid"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsCreateCmd(), tc)
}

func TestCovAcct_EnrollmentTermsUpdate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update enrollment term - ok",
		Args: []string{"1", "42", "--name", "Fall 2025 Updated"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms/42": cmdtest.NewMockResponse(`{
				"id":42,"name":"Fall 2025 Updated","start_at":"2025-08-01","end_at":"2025-12-15"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsUpdateCmd(), tc)
}

func TestCovAcct_EnrollmentTermsUpdate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update enrollment term - error",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms/42": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsUpdateCmd(), tc)
}

func TestCovAcct_EnrollmentTermsDelete_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete enrollment term - ok",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms/42": cmdtest.NewMockResponse(`{}`),
		},
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsDeleteCmd(), tc)
}

func TestCovAcct_EnrollmentTermsDelete_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete enrollment term - error",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/terms/42": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newEnrollmentTermsDeleteCmd(), tc)
}

// ---------------------------------------------------------------------------
// account_notifications
// ---------------------------------------------------------------------------

func TestCovAcct_AccountNotificationsList_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account notifications - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications": cmdtest.NewMockResponse(`[
				{"id":5,"subject":"Maintenance","message":"Down for maintenance","start_at":"2026-01-01","end_at":"2026-01-02"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsListCmd(), tc)
}

func TestCovAcct_AccountNotificationsList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account notifications - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsListCmd(), tc)
}

func TestCovAcct_AccountNotificationsList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account notifications - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsListCmd(), tc)
}

func TestCovAcct_AccountNotificationsGet_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get account notification - ok",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications/5": cmdtest.NewMockResponse(`{
				"id":5,"subject":"Maintenance","message":"Down for maintenance",
				"start_at":"2026-01-01","end_at":"2026-01-02",
				"icon":"warning","roles":["StudentEnrollment"]
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsGetCmd(), tc)
}

func TestCovAcct_AccountNotificationsGet_MinimalFields(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get account notification - minimal fields",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications/5": cmdtest.NewMockResponse(`{
				"id":5,"subject":"Info","message":"Hello","start_at":"2026-01-01","end_at":"2026-01-02"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsGetCmd(), tc)
}

func TestCovAcct_AccountNotificationsGet_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get account notification - error",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications/5": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsGetCmd(), tc)
}

func TestCovAcct_AccountNotificationsCreate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create account notification - ok",
		Args: []string{"1", "--subject", "Maint", "--message", "Down", "--start-at", "2026-01-01", "--end-at", "2026-01-02"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications": cmdtest.NewMockResponse(`{
				"id":6,"subject":"Maint","message":"Down","start_at":"2026-01-01","end_at":"2026-01-02"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsCreateCmd(), tc)
}

func TestCovAcct_AccountNotificationsCreate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create account notification - error",
		Args: []string{"1", "--subject", "Maint", "--message", "Down", "--start-at", "2026-01-01", "--end-at", "2026-01-02"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications": cmdtest.NewErrorResponse(422, "invalid"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsCreateCmd(), tc)
}

func TestCovAcct_AccountNotificationsUpdate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update account notification - ok",
		Args: []string{"1", "5", "--subject", "Updated"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications/5": cmdtest.NewMockResponse(`{
				"id":5,"subject":"Updated","message":"Down","start_at":"2026-01-01","end_at":"2026-01-02"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsUpdateCmd(), tc)
}

func TestCovAcct_AccountNotificationsUpdate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update account notification - error",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications/5": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsUpdateCmd(), tc)
}

func TestCovAcct_AccountNotificationsDelete_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete account notification - ok",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications/5": cmdtest.NewMockResponse(`{}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsDeleteCmd(), tc)
}

func TestCovAcct_AccountNotificationsDelete_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete account notification - error",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/account_notifications/5": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountNotificationsDeleteCmd(), tc)
}

// ---------------------------------------------------------------------------
// account_content_migrations
// ---------------------------------------------------------------------------

func TestCovAcct_AccountContentMigrationsList_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account content migrations - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations": cmdtest.NewMockResponse(`[
				{"id":42,"migration_type":"course_copy_importer","workflow_state":"completed","started_at":"2026-01-01T00:00:00Z"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsListCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsList_LongType(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account content migrations - long migration type truncated",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations": cmdtest.NewMockResponse(`[
				{"id":42,"migration_type":"a_very_long_migration_type_name_that_exceeds_thirty_characters","workflow_state":"running","started_at":""}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsListCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account content migrations - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsListCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account content migrations - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsListCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsGet_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get account content migration - ok",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations/42": cmdtest.NewMockResponse(`{
				"id":42,"migration_type":"course_copy_importer","workflow_state":"completed",
				"started_at":"2026-01-01T00:00:00Z","finished_at":"2026-01-01T00:05:00Z",
				"migration_issues_count":2
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsGetCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsGet_MinimalResponse(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get account content migration - minimal response",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations/42": cmdtest.NewMockResponse(`{
				"id":42,"migration_type":"course_copy_importer","workflow_state":"running"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsGetCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsGet_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get account content migration - error",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations/42": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsGetCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsCreate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create account content migration - ok",
		Args: []string{"1", "--type", "course_copy_importer", "--source-course-id", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations": cmdtest.NewMockResponse(`{
				"id":43,"migration_type":"course_copy_importer","workflow_state":"created"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsCreateCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsCreate_NoSourceCourse(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create account content migration - no source course",
		Args: []string{"1", "--type", "course_copy_importer"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations": cmdtest.NewMockResponse(`{
				"id":43,"migration_type":"course_copy_importer","workflow_state":"created"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsCreateCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsCreate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create account content migration - error",
		Args: []string{"1", "--type", "course_copy_importer"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations": cmdtest.NewErrorResponse(422, "invalid"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsCreateCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsMigrators_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list migrators - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations/migrators": cmdtest.NewMockResponse(`[
				{"type":"course_copy_importer","name":"Course Copy","requires_file_upload":false}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsMigratorsCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsMigrators_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list migrators - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations/migrators": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsMigratorsCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsMigrators_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list migrators - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations/migrators": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsMigratorsCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsIssues_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list migration issues - ok",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations/42/migration_issues": cmdtest.NewMockResponse(`[
				{"id":1,"issue_type":"warning","workflow_state":"active","description":"Some issue occurred"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsIssuesCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsIssues_LongDesc(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list migration issues - long description truncated",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations/42/migration_issues": cmdtest.NewMockResponse(`[
				{"id":1,"issue_type":"warning","workflow_state":"active","description":"A very long description that exceeds the display width limit for truncation"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsIssuesCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsIssues_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list migration issues - empty",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations/42/migration_issues": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsIssuesCmd(), tc)
}

func TestCovAcct_AccountContentMigrationsIssues_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list migration issues - error",
		Args: []string{"1", "42"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/content_migrations/42/migration_issues": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountContentMigrationsIssuesCmd(), tc)
}

// ---------------------------------------------------------------------------
// temporary_enrollment_pairings
// ---------------------------------------------------------------------------

func TestCovAcct_TemporaryEnrollmentPairingsList_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list temporary enrollment pairings - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/temporary_enrollment_pairings": cmdtest.NewMockResponse(`[
				{"id":5,"root_account_id":1,"workflow_state":"active","starting_enrollment_state":"active"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsListCmd(), tc)
}

func TestCovAcct_TemporaryEnrollmentPairingsList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list temporary enrollment pairings - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/temporary_enrollment_pairings": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsListCmd(), tc)
}

func TestCovAcct_TemporaryEnrollmentPairingsList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list temporary enrollment pairings - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/temporary_enrollment_pairings": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsListCmd(), tc)
}

func TestCovAcct_TemporaryEnrollmentPairingsGet_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get temporary enrollment pairing - ok",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/temporary_enrollment_pairings/5": cmdtest.NewMockResponse(`{
				"id":5,"root_account_id":1,"workflow_state":"active",
				"starting_enrollment_state":"active",
				"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T01:00:00Z"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsGetCmd(), tc)
}

func TestCovAcct_TemporaryEnrollmentPairingsGet_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get temporary enrollment pairing - error",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/temporary_enrollment_pairings/5": cmdtest.NewErrorResponse(404, "not found"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsGetCmd(), tc)
}

func TestCovAcct_TemporaryEnrollmentPairingsCreate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create temporary enrollment pairing - ok",
		Args: []string{"1", "--enrollment-state", "active", "--role-id", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/temporary_enrollment_pairings": cmdtest.NewMockResponse(`{
				"id":6,"root_account_id":1,"workflow_state":"active","starting_enrollment_state":"active"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsCreateCmd(), tc)
}

func TestCovAcct_TemporaryEnrollmentPairingsCreate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create temporary enrollment pairing - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/temporary_enrollment_pairings": cmdtest.NewErrorResponse(422, "invalid"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsCreateCmd(), tc)
}

func TestCovAcct_TemporaryEnrollmentPairingsDelete_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete temporary enrollment pairing - ok",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/temporary_enrollment_pairings/5": cmdtest.NewMockResponse(`{}`),
		},
	}
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsDeleteCmd(), tc)
}

func TestCovAcct_TemporaryEnrollmentPairingsDelete_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "delete temporary enrollment pairing - error",
		Args: []string{"1", "5"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/temporary_enrollment_pairings/5": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsDeleteCmd(), tc)
}

// ---------------------------------------------------------------------------
// account_logins
// ---------------------------------------------------------------------------

func TestCovAcct_AccountLoginsList_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account logins - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/logins": cmdtest.NewMockResponse(`[
				{"id":10,"unique_id":"alice@example.com","user_id":100}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountLoginsListCmd(), tc)
}

func TestCovAcct_AccountLoginsList_WithUserFilter(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account logins - with user filter",
		Args: []string{"1", "--user-id", "100"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/logins": cmdtest.NewMockResponse(`[
				{"id":10,"unique_id":"alice@example.com","user_id":100}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountLoginsListCmd(), tc)
}

func TestCovAcct_AccountLoginsList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account logins - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/logins": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountLoginsListCmd(), tc)
}

func TestCovAcct_AccountLoginsList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list account logins - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/logins": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountLoginsListCmd(), tc)
}

func TestCovAcct_AccountLoginsCreate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create account login - ok",
		Args: []string{"1", "--user-id", "100", "--unique-id", "alice@example.com", "--password", "secret"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/logins": cmdtest.NewMockResponse(`{
				"id":10,"unique_id":"alice@example.com","user_id":100
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountLoginsCreateCmd(), tc)
}

func TestCovAcct_AccountLoginsCreate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create account login - error",
		Args: []string{"1", "--user-id", "100", "--unique-id", "alice@example.com"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/logins": cmdtest.NewErrorResponse(422, "invalid"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountLoginsCreateCmd(), tc)
}

func TestCovAcct_AccountLoginsUpdate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update account login - ok",
		Args: []string{"1", "10", "--unique-id", "new@example.com"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/logins/10": cmdtest.NewMockResponse(`{
				"id":10,"unique_id":"new@example.com","user_id":100
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountLoginsUpdateCmd(), tc)
}

func TestCovAcct_AccountLoginsUpdate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update account login - error",
		Args: []string{"1", "10"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/logins/10": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountLoginsUpdateCmd(), tc)
}

// ---------------------------------------------------------------------------
// developer_keys
// ---------------------------------------------------------------------------

func TestCovAcct_DeveloperKeysList_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list developer keys - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/developer_keys": cmdtest.NewMockResponse(`[
				{"id":10,"name":"My Key","email":"dev@example.com","workflow_state":"active"}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newDeveloperKeysListCmd(), tc)
}

func TestCovAcct_DeveloperKeysList_Empty(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list developer keys - empty",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/developer_keys": cmdtest.NewMockResponse(`[]`),
		},
	}
	cmdtest.RunCommandTest(t, newDeveloperKeysListCmd(), tc)
}

func TestCovAcct_DeveloperKeysList_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list developer keys - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/developer_keys": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newDeveloperKeysListCmd(), tc)
}

func TestCovAcct_DeveloperKeysCreate_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create developer key - ok",
		Args: []string{"1", "--name", "My Key", "--email", "dev@example.com", "--redirect-uri", "https://app.example.com/callback"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/developer_keys": cmdtest.NewMockResponse(`{
				"id":10,"name":"My Key","email":"dev@example.com","workflow_state":"active"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newDeveloperKeysCreateCmd(), tc)
}

func TestCovAcct_DeveloperKeysCreate_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create developer key - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/developer_keys": cmdtest.NewErrorResponse(422, "invalid"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newDeveloperKeysCreateCmd(), tc)
}

func TestCovAcct_DeveloperKeysBind_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bind developer key - ok",
		Args: []string{"1", "10", "--workflow-state", "on"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/developer_keys/10/developer_key_account_bindings": cmdtest.NewMockResponse(`{
				"id":1,"developer_key_id":10,"account_id":1,"workflow_state":"on"
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newDeveloperKeysBindCmd(), tc)
}

func TestCovAcct_DeveloperKeysBind_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "bind developer key - error",
		Args: []string{"1", "10", "--workflow-state", "on"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/developer_keys/10/developer_key_account_bindings": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newDeveloperKeysBindCmd(), tc)
}

// ---------------------------------------------------------------------------
// csp_settings
// ---------------------------------------------------------------------------

func TestCovAcct_CSPSettingsGet_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get CSP settings - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/csp_settings": cmdtest.NewMockResponse(`{
				"status":"enabled","locked":false,"domains":["example.com","sub.example.com"]
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newCSPSettingsGetCmd(), tc)
}

func TestCovAcct_CSPSettingsGet_NoDomains(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get CSP settings - no domains",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/csp_settings": cmdtest.NewMockResponse(`{
				"status":"disabled","locked":true,"locked_by":"parent","domains":[]
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newCSPSettingsGetCmd(), tc)
}

func TestCovAcct_CSPSettingsGet_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get CSP settings - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/csp_settings": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCSPSettingsGetCmd(), tc)
}

func TestCovAcct_CSPSettingsAddDomain_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "add CSP domain - ok",
		Args: []string{"1", "--domain", "example.com"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/csp_settings/domains": cmdtest.NewMockResponse(`{
				"status":"enabled","locked":false,"domains":["example.com"]
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newCSPSettingsAddDomainCmd(), tc)
}

func TestCovAcct_CSPSettingsAddDomain_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "add CSP domain - error",
		Args: []string{"1", "--domain", "example.com"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/csp_settings/domains": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCSPSettingsAddDomainCmd(), tc)
}

func TestCovAcct_CSPSettingsRemoveDomain_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "remove CSP domain - ok",
		Args: []string{"1", "--domain", "example.com"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/csp_settings/domains": cmdtest.NewMockResponse(`{
				"status":"enabled","locked":false,"domains":[]
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newCSPSettingsRemoveDomainCmd(), tc)
}

func TestCovAcct_CSPSettingsRemoveDomain_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "remove CSP domain - error",
		Args: []string{"1", "--domain", "example.com"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/csp_settings/domains": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCSPSettingsRemoveDomainCmd(), tc)
}

func TestCovAcct_CSPSettingsLock_Lock(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "lock CSP settings - ok",
		Args: []string{"1", "--locked"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/csp_settings": cmdtest.NewMockResponse(`{
				"status":"enabled","locked":true,"domains":[]
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newCSPSettingsLockCmd(), tc)
}

func TestCovAcct_CSPSettingsLock_Unlock(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "unlock CSP settings - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/csp_settings": cmdtest.NewMockResponse(`{
				"status":"enabled","locked":false,"domains":[]
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newCSPSettingsLockCmd(), tc)
}

func TestCovAcct_CSPSettingsLock_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "lock CSP settings - error",
		Args: []string{"1", "--locked"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/csp_settings": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newCSPSettingsLockCmd(), tc)
}

// ---------------------------------------------------------------------------
// account_analytics
// ---------------------------------------------------------------------------

func TestCovAcct_AccountAnalyticsTermActivity_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "term activity - ok",
		Args: []string{"1", "2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/analytics/terms/2/activity": cmdtest.NewMockResponse(`[
				{"date":"2026-01-01","views":100,"participations":50}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermActivityCmd(), tc)
}

func TestCovAcct_AccountAnalyticsTermActivity_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "term activity - error",
		Args: []string{"1", "2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/analytics/terms/2/activity": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermActivityCmd(), tc)
}

func TestCovAcct_AccountAnalyticsTermGrades_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "term grades - ok",
		Args: []string{"1", "2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/analytics/terms/2/grades": cmdtest.NewMockResponse(`[
				{"bucket":"0-20","mean_current_score":15.0}
			]`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermGradesCmd(), tc)
}

func TestCovAcct_AccountAnalyticsTermGrades_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "term grades - error",
		Args: []string{"1", "2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/analytics/terms/2/grades": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermGradesCmd(), tc)
}

func TestCovAcct_AccountAnalyticsTermStatistics_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "term statistics - ok",
		Args: []string{"1", "2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/analytics/terms/2/statistics": cmdtest.NewMockResponse(`{
				"courses":10,"subaccounts":2,"teachers":5
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermStatisticsCmd(), tc)
}

func TestCovAcct_AccountAnalyticsTermStatistics_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "term statistics - error",
		Args: []string{"1", "2"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/analytics/terms/2/statistics": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermStatisticsCmd(), tc)
}

func TestCovAcct_AccountAnalyticsCompletedStatistics_OK(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "completed statistics - ok",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/analytics/completed/statistics": cmdtest.NewMockResponse(`{
				"courses":5,"teachers":3
			}`),
		},
	}
	cmdtest.RunCommandTest(t, newAccountAnalyticsCompletedStatisticsCmd(), tc)
}

func TestCovAcct_AccountAnalyticsCompletedStatistics_APIError(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "completed statistics - error",
		Args: []string{"1"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/analytics/completed/statistics": cmdtest.NewErrorResponse(500, "error"),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newAccountAnalyticsCompletedStatisticsCmd(), tc)
}

// ---------------------------------------------------------------------------
// Invalid argument tests — trigger strconv.ParseInt error branches in newXxx
// functions so the command file reaches ≥80% statement coverage.
// No mock server needed since argument validation happens before HTTP calls.
// ---------------------------------------------------------------------------

func TestCovAcct_AuthProviders_InvalidAccountID(t *testing.T) {
	// list: invalid account-id
	cmdtest.RunCommandTest(t, newAuthProvidersListCmd(), cmdtest.CommandTestCase{Name: "list - bad account-id", Args: []string{"notanumber"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	// get: invalid account-id
	cmdtest.RunCommandTest(t, newAuthProvidersGetCmd(), cmdtest.CommandTestCase{Name: "get - bad account-id", Args: []string{"notanumber", "10"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	// get: invalid provider-id
	cmdtest.RunCommandTest(t, newAuthProvidersGetCmd(), cmdtest.CommandTestCase{Name: "get - bad provider-id", Args: []string{"1", "notanumber"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	// create: invalid account-id
	cmdtest.RunCommandTest(t, newAuthProvidersCreateCmd(), cmdtest.CommandTestCase{Name: "create - bad account-id", Args: []string{"notanumber", "--auth-type", "saml"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	// delete: invalid account-id
	cmdtest.RunCommandTest(t, newAuthProvidersDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad account-id", Args: []string{"notanumber", "10"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	// delete: invalid provider-id
	cmdtest.RunCommandTest(t, newAuthProvidersDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad provider-id", Args: []string{"1", "notanumber"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	// restore: invalid account-id
	cmdtest.RunCommandTest(t, newAuthProvidersRestoreCmd(), cmdtest.CommandTestCase{Name: "restore - bad account-id", Args: []string{"notanumber", "10"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	// restore: invalid provider-id
	cmdtest.RunCommandTest(t, newAuthProvidersRestoreCmd(), cmdtest.CommandTestCase{Name: "restore - bad provider-id", Args: []string{"1", "notanumber"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	// force-password-reset: invalid account-id
	cmdtest.RunCommandTest(t, newAuthProvidersForcePasswordResetCmd(), cmdtest.CommandTestCase{Name: "force-pw-reset - bad account-id", Args: []string{"notanumber"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	// sso-settings: invalid account-id
	cmdtest.RunCommandTest(t, newAuthProvidersSSOSettingsCmd(), cmdtest.CommandTestCase{Name: "sso-settings - bad account-id", Args: []string{"notanumber"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_AccountFeatures_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newAccountFeaturesListCmd(), cmdtest.CommandTestCase{Name: "list - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountFeaturesListEnabledCmd(), cmdtest.CommandTestCase{Name: "list-enabled - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountFeaturesGetFlagCmd(), cmdtest.CommandTestCase{Name: "get-flag - bad id", Args: []string{"x", "feat"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountFeaturesSetFlagCmd(), cmdtest.CommandTestCase{Name: "set-flag - bad id", Args: []string{"x", "feat", "--state", "on"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountFeaturesDeleteFlagCmd(), cmdtest.CommandTestCase{Name: "del-flag - bad id", Args: []string{"x", "feat"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountFeaturesSettingsCmd(), cmdtest.CommandTestCase{Name: "settings - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountFeaturesPermissionsCmd(), cmdtest.CommandTestCase{Name: "perms - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_GradingPeriodSets_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newGradingPeriodSetsListCmd(), cmdtest.CommandTestCase{Name: "list - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newGradingPeriodSetsCreateCmd(), cmdtest.CommandTestCase{Name: "create - bad id", Args: []string{"x", "--title", "T"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newGradingPeriodSetsUpdateCmd(), cmdtest.CommandTestCase{Name: "update - bad account-id", Args: []string{"x", "5"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newGradingPeriodSetsUpdateCmd(), cmdtest.CommandTestCase{Name: "update - bad set-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newGradingPeriodSetsDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad account-id", Args: []string{"x", "5"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newGradingPeriodSetsDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad set-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newGradingPeriodSetsListPeriodsCmd(), cmdtest.CommandTestCase{Name: "list-periods - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newGradingPeriodSetsDeletePeriodCmd(), cmdtest.CommandTestCase{Name: "del-period - bad account-id", Args: []string{"x", "3"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newGradingPeriodSetsDeletePeriodCmd(), cmdtest.CommandTestCase{Name: "del-period - bad period-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_AccountReports_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newAccountReportsListCmd(), cmdtest.CommandTestCase{Name: "list - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountReportsRunsCmd(), cmdtest.CommandTestCase{Name: "runs - bad id", Args: []string{"x", "r"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountReportsStartCmd(), cmdtest.CommandTestCase{Name: "start - bad id", Args: []string{"x", "r"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountReportsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad account-id", Args: []string{"x", "r", "10"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountReportsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad run-id", Args: []string{"1", "r", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountReportsDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad account-id", Args: []string{"x", "r", "10"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountReportsDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad run-id", Args: []string{"1", "r", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountReportsAbortCmd(), cmdtest.CommandTestCase{Name: "abort - bad account-id", Args: []string{"x", "r", "10"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountReportsAbortCmd(), cmdtest.CommandTestCase{Name: "abort - bad run-id", Args: []string{"1", "r", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_EnrollmentTerms_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newEnrollmentTermsListCmd(), cmdtest.CommandTestCase{Name: "list - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newEnrollmentTermsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad account-id", Args: []string{"x", "42"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newEnrollmentTermsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad term-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newEnrollmentTermsCreateCmd(), cmdtest.CommandTestCase{Name: "create - bad id", Args: []string{"x", "--name", "T"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newEnrollmentTermsUpdateCmd(), cmdtest.CommandTestCase{Name: "update - bad account-id", Args: []string{"x", "42"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newEnrollmentTermsUpdateCmd(), cmdtest.CommandTestCase{Name: "update - bad term-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newEnrollmentTermsDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad account-id", Args: []string{"x", "42"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newEnrollmentTermsDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad term-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_AccountNotifications_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newAccountNotificationsListCmd(), cmdtest.CommandTestCase{Name: "list - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountNotificationsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad account-id", Args: []string{"x", "5"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountNotificationsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad notif-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountNotificationsCreateCmd(), cmdtest.CommandTestCase{Name: "create - bad id", Args: []string{"x", "--subject", "S", "--message", "M", "--start-at", "A", "--end-at", "B"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountNotificationsUpdateCmd(), cmdtest.CommandTestCase{Name: "update - bad account-id", Args: []string{"x", "5"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountNotificationsUpdateCmd(), cmdtest.CommandTestCase{Name: "update - bad notif-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountNotificationsDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad account-id", Args: []string{"x", "5"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountNotificationsDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad notif-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_AccountContentMigrations_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newAccountContentMigrationsListCmd(), cmdtest.CommandTestCase{Name: "list - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountContentMigrationsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad account-id", Args: []string{"x", "42"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountContentMigrationsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad migration-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountContentMigrationsCreateCmd(), cmdtest.CommandTestCase{Name: "create - bad id", Args: []string{"x", "--type", "t"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountContentMigrationsMigratorsCmd(), cmdtest.CommandTestCase{Name: "migrators - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountContentMigrationsIssuesCmd(), cmdtest.CommandTestCase{Name: "issues - bad account-id", Args: []string{"x", "42"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountContentMigrationsIssuesCmd(), cmdtest.CommandTestCase{Name: "issues - bad migration-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_TemporaryEnrollmentPairings_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsListCmd(), cmdtest.CommandTestCase{Name: "list - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad account-id", Args: []string{"x", "5"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsCreateCmd(), cmdtest.CommandTestCase{Name: "create - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad account-id", Args: []string{"x", "5"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newTemporaryEnrollmentPairingsDeleteCmd(), cmdtest.CommandTestCase{Name: "delete - bad id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_AccountLogins_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newAccountLoginsListCmd(), cmdtest.CommandTestCase{Name: "list - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountLoginsCreateCmd(), cmdtest.CommandTestCase{Name: "create - bad id", Args: []string{"x", "--user-id", "100", "--unique-id", "a"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountLoginsUpdateCmd(), cmdtest.CommandTestCase{Name: "update - bad account-id", Args: []string{"x", "10"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountLoginsUpdateCmd(), cmdtest.CommandTestCase{Name: "update - bad login-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_DeveloperKeys_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newDeveloperKeysListCmd(), cmdtest.CommandTestCase{Name: "list - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newDeveloperKeysCreateCmd(), cmdtest.CommandTestCase{Name: "create - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newDeveloperKeysBindCmd(), cmdtest.CommandTestCase{Name: "bind - bad account-id", Args: []string{"x", "10", "--workflow-state", "on"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newDeveloperKeysBindCmd(), cmdtest.CommandTestCase{Name: "bind - bad key-id", Args: []string{"1", "x", "--workflow-state", "on"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_CSPSettings_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newCSPSettingsGetCmd(), cmdtest.CommandTestCase{Name: "get - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newCSPSettingsAddDomainCmd(), cmdtest.CommandTestCase{Name: "add-domain - bad id", Args: []string{"x", "--domain", "e.com"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newCSPSettingsRemoveDomainCmd(), cmdtest.CommandTestCase{Name: "rm-domain - bad id", Args: []string{"x", "--domain", "e.com"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newCSPSettingsLockCmd(), cmdtest.CommandTestCase{Name: "lock - bad id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}

func TestCovAcct_AccountAnalytics_InvalidArgs(t *testing.T) {
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermActivityCmd(), cmdtest.CommandTestCase{Name: "term-activity - bad account-id", Args: []string{"x", "2"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermActivityCmd(), cmdtest.CommandTestCase{Name: "term-activity - bad term-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermGradesCmd(), cmdtest.CommandTestCase{Name: "term-grades - bad account-id", Args: []string{"x", "2"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermGradesCmd(), cmdtest.CommandTestCase{Name: "term-grades - bad term-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermStatisticsCmd(), cmdtest.CommandTestCase{Name: "term-stats - bad account-id", Args: []string{"x", "2"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountAnalyticsTermStatisticsCmd(), cmdtest.CommandTestCase{Name: "term-stats - bad term-id", Args: []string{"1", "x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
	cmdtest.RunCommandTest(t, newAccountAnalyticsCompletedStatisticsCmd(), cmdtest.CommandTestCase{Name: "completed-stats - bad account-id", Args: []string{"x"}, MockResponses: map[string]cmdtest.MockResponse{}, ExpectError: true})
}
