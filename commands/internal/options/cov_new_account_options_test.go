package options

// cov_new_account_options_test.go — table-driven Validate() tests for the
// Wave-4 options structs added in feature/spec-compliance:
//   auth_providers, account_features, grading_period_sets, account_reports,
//   enrollment_terms, account_notifications, account_content_migrations,
//   temporary_enrollment_pairings, account_logins, developer_keys,
//   csp_settings, account_analytics
//
// Prefix TestCovAcctOpts_ prevents collisions with existing test functions.

import (
	"testing"
)

// ---------------------------------------------------------------------------
// auth_providers options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_AuthProvidersListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AuthProvidersListOptions
		wantErr bool
	}{
		{"valid", AuthProvidersListOptions{AccountID: 1}, false},
		{"missing account-id", AuthProvidersListOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AuthProvidersGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AuthProvidersGetOptions
		wantErr bool
	}{
		{"valid", AuthProvidersGetOptions{AccountID: 1, ProviderID: 10}, false},
		{"missing account-id", AuthProvidersGetOptions{ProviderID: 10}, true},
		{"missing provider-id", AuthProvidersGetOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AuthProvidersCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AuthProvidersCreateOptions
		wantErr bool
	}{
		{"valid", AuthProvidersCreateOptions{AccountID: 1, AuthType: "saml"}, false},
		{"missing account-id", AuthProvidersCreateOptions{AuthType: "saml"}, true},
		{"missing auth-type", AuthProvidersCreateOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AuthProvidersDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AuthProvidersDeleteOptions
		wantErr bool
	}{
		{"valid", AuthProvidersDeleteOptions{AccountID: 1, ProviderID: 5}, false},
		{"missing account-id", AuthProvidersDeleteOptions{ProviderID: 5}, true},
		{"missing provider-id", AuthProvidersDeleteOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AuthProvidersRestoreOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AuthProvidersRestoreOptions
		wantErr bool
	}{
		{"valid", AuthProvidersRestoreOptions{AccountID: 1, ProviderID: 5}, false},
		{"missing account-id", AuthProvidersRestoreOptions{ProviderID: 5}, true},
		{"missing provider-id", AuthProvidersRestoreOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AuthProvidersForcePasswordResetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AuthProvidersForcePasswordResetOptions
		wantErr bool
	}{
		{"valid", AuthProvidersForcePasswordResetOptions{AccountID: 1}, false},
		{"missing account-id", AuthProvidersForcePasswordResetOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AuthProvidersSSOSettingsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AuthProvidersSSOSettingsOptions
		wantErr bool
	}{
		{"valid", AuthProvidersSSOSettingsOptions{AccountID: 1}, false},
		{"missing account-id", AuthProvidersSSOSettingsOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// account_features options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_AccountFeaturesListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountFeaturesListOptions
		wantErr bool
	}{
		{"valid", AccountFeaturesListOptions{AccountID: 1}, false},
		{"missing account-id", AccountFeaturesListOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountFeaturesListEnabledOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountFeaturesListEnabledOptions
		wantErr bool
	}{
		{"valid", AccountFeaturesListEnabledOptions{AccountID: 1}, false},
		{"missing account-id", AccountFeaturesListEnabledOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountFeaturesGetFlagOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountFeaturesGetFlagOptions
		wantErr bool
	}{
		{"valid", AccountFeaturesGetFlagOptions{AccountID: 1, Feature: "analytics_2"}, false},
		{"missing account-id", AccountFeaturesGetFlagOptions{Feature: "analytics_2"}, true},
		{"missing feature", AccountFeaturesGetFlagOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountFeaturesSetFlagOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountFeaturesSetFlagOptions
		wantErr bool
	}{
		{"valid", AccountFeaturesSetFlagOptions{AccountID: 1, Feature: "analytics_2", State: "on"}, false},
		{"missing account-id", AccountFeaturesSetFlagOptions{Feature: "analytics_2", State: "on"}, true},
		{"missing feature", AccountFeaturesSetFlagOptions{AccountID: 1, State: "on"}, true},
		{"missing state", AccountFeaturesSetFlagOptions{AccountID: 1, Feature: "analytics_2"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountFeaturesDeleteFlagOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountFeaturesDeleteFlagOptions
		wantErr bool
	}{
		{"valid", AccountFeaturesDeleteFlagOptions{AccountID: 1, Feature: "analytics_2"}, false},
		{"missing account-id", AccountFeaturesDeleteFlagOptions{Feature: "analytics_2"}, true},
		{"missing feature", AccountFeaturesDeleteFlagOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountFeaturesSettingsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountFeaturesSettingsOptions
		wantErr bool
	}{
		{"valid", AccountFeaturesSettingsOptions{AccountID: 1}, false},
		{"missing account-id", AccountFeaturesSettingsOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountFeaturesPermissionsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountFeaturesPermissionsOptions
		wantErr bool
	}{
		{"valid", AccountFeaturesPermissionsOptions{AccountID: 1}, false},
		{"valid with perms", AccountFeaturesPermissionsOptions{AccountID: 1, Permissions: []string{"manage_courses"}}, false},
		{"missing account-id", AccountFeaturesPermissionsOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// grading_period_sets options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_GradingPeriodSetsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    GradingPeriodSetsListOptions
		wantErr bool
	}{
		{"valid", GradingPeriodSetsListOptions{AccountID: 1}, false},
		{"missing account-id", GradingPeriodSetsListOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_GradingPeriodSetsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    GradingPeriodSetsCreateOptions
		wantErr bool
	}{
		{"valid", GradingPeriodSetsCreateOptions{AccountID: 1, Title: "2024-2025"}, false},
		{"missing account-id", GradingPeriodSetsCreateOptions{Title: "2024-2025"}, true},
		{"missing title", GradingPeriodSetsCreateOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_GradingPeriodSetsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    GradingPeriodSetsUpdateOptions
		wantErr bool
	}{
		{"valid", GradingPeriodSetsUpdateOptions{AccountID: 1, SetID: 5}, false},
		{"missing account-id", GradingPeriodSetsUpdateOptions{SetID: 5}, true},
		{"missing set-id", GradingPeriodSetsUpdateOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_GradingPeriodSetsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    GradingPeriodSetsDeleteOptions
		wantErr bool
	}{
		{"valid", GradingPeriodSetsDeleteOptions{AccountID: 1, SetID: 5}, false},
		{"missing account-id", GradingPeriodSetsDeleteOptions{SetID: 5}, true},
		{"missing set-id", GradingPeriodSetsDeleteOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_GradingPeriodSetsListPeriodsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    GradingPeriodSetsListPeriodsOptions
		wantErr bool
	}{
		{"valid", GradingPeriodSetsListPeriodsOptions{AccountID: 1}, false},
		{"missing account-id", GradingPeriodSetsListPeriodsOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_GradingPeriodSetsDeletePeriodOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    GradingPeriodSetsDeletePeriodOptions
		wantErr bool
	}{
		{"valid", GradingPeriodSetsDeletePeriodOptions{AccountID: 1, PeriodID: 3}, false},
		{"missing account-id", GradingPeriodSetsDeletePeriodOptions{PeriodID: 3}, true},
		{"missing period-id", GradingPeriodSetsDeletePeriodOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// account_reports options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_AccountReportsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountReportsListOptions
		wantErr bool
	}{
		{"valid", AccountReportsListOptions{AccountID: 1}, false},
		{"missing account-id", AccountReportsListOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountReportsRunsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountReportsRunsOptions
		wantErr bool
	}{
		{"valid", AccountReportsRunsOptions{AccountID: 1, ReportName: "course_storage_csv"}, false},
		{"missing account-id", AccountReportsRunsOptions{ReportName: "course_storage_csv"}, true},
		{"missing report-name", AccountReportsRunsOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountReportsStartOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountReportsStartOptions
		wantErr bool
	}{
		{"valid", AccountReportsStartOptions{AccountID: 1, ReportName: "course_storage_csv"}, false},
		{"missing account-id", AccountReportsStartOptions{ReportName: "course_storage_csv"}, true},
		{"missing report-name", AccountReportsStartOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountReportsGetRunOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountReportsGetRunOptions
		wantErr bool
	}{
		{"valid", AccountReportsGetRunOptions{AccountID: 1, ReportName: "r", RunID: 10}, false},
		{"missing account-id", AccountReportsGetRunOptions{ReportName: "r", RunID: 10}, true},
		{"missing report-name", AccountReportsGetRunOptions{AccountID: 1, RunID: 10}, true},
		{"missing run-id", AccountReportsGetRunOptions{AccountID: 1, ReportName: "r"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountReportsDeleteRunOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountReportsDeleteRunOptions
		wantErr bool
	}{
		{"valid", AccountReportsDeleteRunOptions{AccountID: 1, ReportName: "r", RunID: 10}, false},
		{"missing account-id", AccountReportsDeleteRunOptions{ReportName: "r", RunID: 10}, true},
		{"missing report-name", AccountReportsDeleteRunOptions{AccountID: 1, RunID: 10}, true},
		{"missing run-id", AccountReportsDeleteRunOptions{AccountID: 1, ReportName: "r"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountReportsAbortRunOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountReportsAbortRunOptions
		wantErr bool
	}{
		{"valid", AccountReportsAbortRunOptions{AccountID: 1, ReportName: "r", RunID: 10}, false},
		{"missing account-id", AccountReportsAbortRunOptions{ReportName: "r", RunID: 10}, true},
		{"missing report-name", AccountReportsAbortRunOptions{AccountID: 1, RunID: 10}, true},
		{"missing run-id", AccountReportsAbortRunOptions{AccountID: 1, ReportName: "r"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// enrollment_terms options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_EnrollmentTermsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    EnrollmentTermsListOptions
		wantErr bool
	}{
		{"valid", EnrollmentTermsListOptions{AccountID: 1}, false},
		{"missing account-id", EnrollmentTermsListOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_EnrollmentTermsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    EnrollmentTermsGetOptions
		wantErr bool
	}{
		{"valid", EnrollmentTermsGetOptions{AccountID: 1, TermID: 42}, false},
		{"missing account-id", EnrollmentTermsGetOptions{TermID: 42}, true},
		{"missing term-id", EnrollmentTermsGetOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_EnrollmentTermsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    EnrollmentTermsCreateOptions
		wantErr bool
	}{
		{"valid", EnrollmentTermsCreateOptions{AccountID: 1, Name: "Fall 2025"}, false},
		{"missing account-id", EnrollmentTermsCreateOptions{Name: "Fall 2025"}, true},
		{"missing name", EnrollmentTermsCreateOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_EnrollmentTermsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    EnrollmentTermsUpdateOptions
		wantErr bool
	}{
		{"valid", EnrollmentTermsUpdateOptions{AccountID: 1, TermID: 42}, false},
		{"missing account-id", EnrollmentTermsUpdateOptions{TermID: 42}, true},
		{"missing term-id", EnrollmentTermsUpdateOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_EnrollmentTermsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    EnrollmentTermsDeleteOptions
		wantErr bool
	}{
		{"valid", EnrollmentTermsDeleteOptions{AccountID: 1, TermID: 42}, false},
		{"missing account-id", EnrollmentTermsDeleteOptions{TermID: 42}, true},
		{"missing term-id", EnrollmentTermsDeleteOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// account_notifications options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_AccountNotificationsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountNotificationsListOptions
		wantErr bool
	}{
		{"valid", AccountNotificationsListOptions{AccountID: 1}, false},
		{"missing account-id", AccountNotificationsListOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountNotificationsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountNotificationsGetOptions
		wantErr bool
	}{
		{"valid", AccountNotificationsGetOptions{AccountID: 1, NotificationID: 5}, false},
		{"missing account-id", AccountNotificationsGetOptions{NotificationID: 5}, true},
		{"missing notification-id", AccountNotificationsGetOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountNotificationsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountNotificationsCreateOptions
		wantErr bool
	}{
		{
			"valid",
			AccountNotificationsCreateOptions{
				AccountID: 1, Subject: "s", Message: "m", StartAt: "2026-01-01", EndAt: "2026-01-02",
			},
			false,
		},
		{"missing account-id", AccountNotificationsCreateOptions{Subject: "s", Message: "m", StartAt: "a", EndAt: "b"}, true},
		{"missing subject", AccountNotificationsCreateOptions{AccountID: 1, Message: "m", StartAt: "a", EndAt: "b"}, true},
		{"missing message", AccountNotificationsCreateOptions{AccountID: 1, Subject: "s", StartAt: "a", EndAt: "b"}, true},
		{"missing start-at", AccountNotificationsCreateOptions{AccountID: 1, Subject: "s", Message: "m", EndAt: "b"}, true},
		{"missing end-at", AccountNotificationsCreateOptions{AccountID: 1, Subject: "s", Message: "m", StartAt: "a"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountNotificationsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountNotificationsUpdateOptions
		wantErr bool
	}{
		{"valid", AccountNotificationsUpdateOptions{AccountID: 1, NotificationID: 5}, false},
		{"missing account-id", AccountNotificationsUpdateOptions{NotificationID: 5}, true},
		{"missing notification-id", AccountNotificationsUpdateOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountNotificationsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountNotificationsDeleteOptions
		wantErr bool
	}{
		{"valid", AccountNotificationsDeleteOptions{AccountID: 1, NotificationID: 5}, false},
		{"missing account-id", AccountNotificationsDeleteOptions{NotificationID: 5}, true},
		{"missing notification-id", AccountNotificationsDeleteOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// account_content_migrations options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_AccountContentMigrationsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountContentMigrationsListOptions
		wantErr bool
	}{
		{"valid", AccountContentMigrationsListOptions{AccountID: 1}, false},
		{"missing account-id", AccountContentMigrationsListOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountContentMigrationsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountContentMigrationsGetOptions
		wantErr bool
	}{
		{"valid", AccountContentMigrationsGetOptions{AccountID: 1, MigrationID: 42}, false},
		{"missing account-id", AccountContentMigrationsGetOptions{MigrationID: 42}, true},
		{"missing migration-id", AccountContentMigrationsGetOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountContentMigrationsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountContentMigrationsCreateOptions
		wantErr bool
	}{
		{"valid", AccountContentMigrationsCreateOptions{AccountID: 1, MigrationType: "course_copy_importer"}, false},
		{"missing account-id", AccountContentMigrationsCreateOptions{MigrationType: "course_copy_importer"}, true},
		{"missing type", AccountContentMigrationsCreateOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountContentMigrationsMigratorsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountContentMigrationsMigratorsOptions
		wantErr bool
	}{
		{"valid", AccountContentMigrationsMigratorsOptions{AccountID: 1}, false},
		{"missing account-id", AccountContentMigrationsMigratorsOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountContentMigrationsIssuesOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountContentMigrationsIssuesOptions
		wantErr bool
	}{
		{"valid", AccountContentMigrationsIssuesOptions{AccountID: 1, MigrationID: 42}, false},
		{"missing account-id", AccountContentMigrationsIssuesOptions{MigrationID: 42}, true},
		{"missing migration-id", AccountContentMigrationsIssuesOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// temporary_enrollment_pairings options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_TemporaryEnrollmentPairingsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    TemporaryEnrollmentPairingsListOptions
		wantErr bool
	}{
		{"valid", TemporaryEnrollmentPairingsListOptions{AccountID: 1}, false},
		{"missing account-id", TemporaryEnrollmentPairingsListOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_TemporaryEnrollmentPairingsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    TemporaryEnrollmentPairingsGetOptions
		wantErr bool
	}{
		{"valid", TemporaryEnrollmentPairingsGetOptions{AccountID: 1, ID: 5}, false},
		{"missing account-id", TemporaryEnrollmentPairingsGetOptions{ID: 5}, true},
		{"missing id", TemporaryEnrollmentPairingsGetOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_TemporaryEnrollmentPairingsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    TemporaryEnrollmentPairingsCreateOptions
		wantErr bool
	}{
		{"valid", TemporaryEnrollmentPairingsCreateOptions{AccountID: 1}, false},
		{"missing account-id", TemporaryEnrollmentPairingsCreateOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_TemporaryEnrollmentPairingsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    TemporaryEnrollmentPairingsDeleteOptions
		wantErr bool
	}{
		{"valid", TemporaryEnrollmentPairingsDeleteOptions{AccountID: 1, ID: 5}, false},
		{"missing account-id", TemporaryEnrollmentPairingsDeleteOptions{ID: 5}, true},
		{"missing id", TemporaryEnrollmentPairingsDeleteOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// account_logins options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_AccountLoginsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountLoginsListOptions
		wantErr bool
	}{
		{"valid", AccountLoginsListOptions{AccountID: 1}, false},
		{"valid with user-id", AccountLoginsListOptions{AccountID: 1, UserID: 100}, false},
		{"missing account-id", AccountLoginsListOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountLoginsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountLoginsCreateOptions
		wantErr bool
	}{
		{"valid", AccountLoginsCreateOptions{AccountID: 1, UserID: 100, UniqueID: "alice@example.com"}, false},
		{"missing account-id", AccountLoginsCreateOptions{UserID: 100, UniqueID: "alice@example.com"}, true},
		{"missing user-id", AccountLoginsCreateOptions{AccountID: 1, UniqueID: "alice@example.com"}, true},
		{"missing unique-id", AccountLoginsCreateOptions{AccountID: 1, UserID: 100}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountLoginsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountLoginsUpdateOptions
		wantErr bool
	}{
		{"valid", AccountLoginsUpdateOptions{AccountID: 1, LoginID: 10}, false},
		{"missing account-id", AccountLoginsUpdateOptions{LoginID: 10}, true},
		{"missing login-id", AccountLoginsUpdateOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// developer_keys options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_DeveloperKeysListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    DeveloperKeysListOptions
		wantErr bool
	}{
		{"valid", DeveloperKeysListOptions{AccountID: 1}, false},
		{"missing account-id", DeveloperKeysListOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_DeveloperKeysCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    DeveloperKeysCreateOptions
		wantErr bool
	}{
		{"valid", DeveloperKeysCreateOptions{AccountID: 1}, false},
		{"valid with name", DeveloperKeysCreateOptions{AccountID: 1, Name: "My Key"}, false},
		{"missing account-id", DeveloperKeysCreateOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_DeveloperKeysBindOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    DeveloperKeysBindOptions
		wantErr bool
	}{
		{"valid", DeveloperKeysBindOptions{AccountID: 1, DeveloperKeyID: 10, WorkflowState: "on"}, false},
		{"missing account-id", DeveloperKeysBindOptions{DeveloperKeyID: 10, WorkflowState: "on"}, true},
		{"missing developer-key-id", DeveloperKeysBindOptions{AccountID: 1, WorkflowState: "on"}, true},
		{"missing workflow-state", DeveloperKeysBindOptions{AccountID: 1, DeveloperKeyID: 10}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// csp_settings options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_CSPSettingsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    CSPSettingsGetOptions
		wantErr bool
	}{
		{"valid", CSPSettingsGetOptions{AccountID: 1}, false},
		{"missing account-id", CSPSettingsGetOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_CSPSettingsAddDomainOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    CSPSettingsAddDomainOptions
		wantErr bool
	}{
		{"valid", CSPSettingsAddDomainOptions{AccountID: 1, Domain: "example.com"}, false},
		{"missing account-id", CSPSettingsAddDomainOptions{Domain: "example.com"}, true},
		{"missing domain", CSPSettingsAddDomainOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_CSPSettingsRemoveDomainOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    CSPSettingsRemoveDomainOptions
		wantErr bool
	}{
		{"valid", CSPSettingsRemoveDomainOptions{AccountID: 1, Domain: "example.com"}, false},
		{"missing account-id", CSPSettingsRemoveDomainOptions{Domain: "example.com"}, true},
		{"missing domain", CSPSettingsRemoveDomainOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_CSPSettingsLockOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    CSPSettingsLockOptions
		wantErr bool
	}{
		{"valid locked", CSPSettingsLockOptions{AccountID: 1, Locked: true}, false},
		{"valid unlocked", CSPSettingsLockOptions{AccountID: 1, Locked: false}, false},
		{"missing account-id", CSPSettingsLockOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// account_analytics options
// ---------------------------------------------------------------------------

func TestCovAcctOpts_AccountAnalyticsTermOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountAnalyticsTermOptions
		wantErr bool
	}{
		{"valid", AccountAnalyticsTermOptions{AccountID: 1, TermID: 2}, false},
		{"missing account-id", AccountAnalyticsTermOptions{TermID: 2}, true},
		{"missing term-id", AccountAnalyticsTermOptions{AccountID: 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovAcctOpts_AccountAnalyticsCompletedOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    AccountAnalyticsCompletedOptions
		wantErr bool
	}{
		{"valid", AccountAnalyticsCompletedOptions{AccountID: 1}, false},
		{"missing account-id", AccountAnalyticsCompletedOptions{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
