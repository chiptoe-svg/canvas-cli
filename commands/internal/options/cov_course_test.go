package options

import "testing"

// ---------------------------------------------------------------------------
// grading.go — GradingPeriodsListOptions
// ---------------------------------------------------------------------------

func TestCovCourse_GradingPeriodsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingPeriodsListOptions
		wantErr bool
	}{
		{"valid", &GradingPeriodsListOptions{CourseID: 1}, false},
		{"missing course-id", &GradingPeriodsListOptions{CourseID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingPeriodsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingPeriodsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingPeriodsGetOptions
		wantErr bool
	}{
		{"valid", &GradingPeriodsGetOptions{CourseID: 1, ID: 5}, false},
		{"missing course-id", &GradingPeriodsGetOptions{CourseID: 0, ID: 5}, true},
		{"missing id", &GradingPeriodsGetOptions{CourseID: 1, ID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingPeriodsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingPeriodsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingPeriodsUpdateOptions
		wantErr bool
	}{
		{"valid", &GradingPeriodsUpdateOptions{CourseID: 1, ID: 5}, false},
		{"missing course-id", &GradingPeriodsUpdateOptions{CourseID: 0, ID: 5}, true},
		{"missing id", &GradingPeriodsUpdateOptions{CourseID: 1, ID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingPeriodsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingPeriodsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingPeriodsDeleteOptions
		wantErr bool
	}{
		{"valid", &GradingPeriodsDeleteOptions{CourseID: 1, ID: 5}, false},
		{"missing course-id", &GradingPeriodsDeleteOptions{CourseID: 0, ID: 5}, true},
		{"missing id", &GradingPeriodsDeleteOptions{CourseID: 1, ID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingPeriodsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingStandardsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingStandardsListOptions
		wantErr bool
	}{
		{"valid", &GradingStandardsListOptions{CourseID: 1}, false},
		{"missing course-id", &GradingStandardsListOptions{CourseID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingStandardsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingStandardsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingStandardsGetOptions
		wantErr bool
	}{
		{"valid", &GradingStandardsGetOptions{CourseID: 1, StandardID: 5}, false},
		{"missing course-id", &GradingStandardsGetOptions{CourseID: 0, StandardID: 5}, true},
		{"missing standard-id", &GradingStandardsGetOptions{CourseID: 1, StandardID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingStandardsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingStandardsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingStandardsCreateOptions
		wantErr bool
	}{
		{"valid", &GradingStandardsCreateOptions{CourseID: 1, Title: "Letter Grades"}, false},
		{"missing course-id", &GradingStandardsCreateOptions{CourseID: 0, Title: "Letter Grades"}, true},
		{"missing title", &GradingStandardsCreateOptions{CourseID: 1, Title: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingStandardsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingStandardsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingStandardsDeleteOptions
		wantErr bool
	}{
		{"valid", &GradingStandardsDeleteOptions{CourseID: 1, StandardID: 5}, false},
		{"missing course-id", &GradingStandardsDeleteOptions{CourseID: 0, StandardID: 5}, true},
		{"missing standard-id", &GradingStandardsDeleteOptions{CourseID: 1, StandardID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingStandardsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// files.go — additional validators not yet in files_test.go
// ---------------------------------------------------------------------------

func TestCovCourse_FilesResetVerifierOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesResetVerifierOptions
		wantErr bool
	}{
		{"valid", &FilesResetVerifierOptions{FileID: 1}, false},
		{"missing file-id", &FilesResetVerifierOptions{FileID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesResetVerifierOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_FilesCopyOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesCopyOptions
		wantErr bool
	}{
		{"valid", &FilesCopyOptions{DestFolderID: 10, SourceFileID: 5}, false},
		{"missing dest-folder-id", &FilesCopyOptions{DestFolderID: 0, SourceFileID: 5}, true},
		{"missing source-file-id", &FilesCopyOptions{DestFolderID: 10, SourceFileID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesCopyOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_FilesUsageRightsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesUsageRightsOptions
		wantErr bool
	}{
		{
			"valid with file IDs",
			&FilesUsageRightsOptions{
				CourseID:         1,
				FileIDs:          []int64{1},
				UseJustification: "own_copyright",
			},
			false,
		},
		{
			"valid with folder IDs",
			&FilesUsageRightsOptions{
				CourseID:         1,
				FolderIDs:        []int64{1},
				UseJustification: "own_copyright",
			},
			false,
		},
		{
			"no context",
			&FilesUsageRightsOptions{
				FileIDs:          []int64{1},
				UseJustification: "own_copyright",
			},
			true,
		},
		{
			"multiple contexts",
			&FilesUsageRightsOptions{
				CourseID:         1,
				GroupID:          2,
				FileIDs:          []int64{1},
				UseJustification: "own_copyright",
			},
			true,
		},
		{
			"no file or folder IDs",
			&FilesUsageRightsOptions{
				CourseID:         1,
				UseJustification: "own_copyright",
			},
			true,
		},
		{
			"missing use-justification",
			&FilesUsageRightsOptions{
				CourseID: 1,
				FileIDs:  []int64{1},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesUsageRightsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_FilesRemoveUsageRightsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesRemoveUsageRightsOptions
		wantErr bool
	}{
		{
			"valid with file IDs",
			&FilesRemoveUsageRightsOptions{CourseID: 1, FileIDs: []int64{1}},
			false,
		},
		{
			"valid with folder IDs",
			&FilesRemoveUsageRightsOptions{CourseID: 1, FolderIDs: []int64{1}},
			false,
		},
		{
			"no context",
			&FilesRemoveUsageRightsOptions{FileIDs: []int64{1}},
			true,
		},
		{
			"multiple contexts",
			&FilesRemoveUsageRightsOptions{CourseID: 1, GroupID: 2, FileIDs: []int64{1}},
			true,
		},
		{
			"no file or folder IDs",
			&FilesRemoveUsageRightsOptions{CourseID: 1},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesRemoveUsageRightsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_FilesLicensesOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *FilesLicensesOptions
		wantErr bool
	}{
		{"valid course-id", &FilesLicensesOptions{CourseID: 1}, false},
		{"valid group-id", &FilesLicensesOptions{GroupID: 1}, false},
		{"valid user-id", &FilesLicensesOptions{UserID: 1}, false},
		{"no context", &FilesLicensesOptions{}, true},
		{"multiple contexts", &FilesLicensesOptions{CourseID: 1, GroupID: 2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FilesLicensesOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// media_objects.go
// ---------------------------------------------------------------------------

func TestCovCourse_MediaObjectsListOptions_Validate(t *testing.T) {
	// Validate always returns nil (no required fields)
	opts := &MediaObjectsListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("MediaObjectsListOptions.Validate() expected nil, got %v", err)
	}
	opts.CourseID = 1
	if err := opts.Validate(); err != nil {
		t.Errorf("MediaObjectsListOptions.Validate() with CourseID expected nil, got %v", err)
	}
}

func TestCovCourse_MediaObjectUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *MediaObjectUpdateOptions
		wantErr bool
	}{
		{"valid", &MediaObjectUpdateOptions{MediaID: "m-abc", Title: "Title"}, false},
		{"missing media-id", &MediaObjectUpdateOptions{MediaID: "", Title: "Title"}, true},
		{"missing title", &MediaObjectUpdateOptions{MediaID: "m-abc", Title: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MediaObjectUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_MediaTracksListOptions_Validate(t *testing.T) {
	// Validate always returns nil
	opts := &MediaTracksListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("MediaTracksListOptions.Validate() expected nil, got %v", err)
	}
}

func TestCovCourse_MediaAttachmentsListOptions_Validate(t *testing.T) {
	// Validate always returns nil
	opts := &MediaAttachmentsListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("MediaAttachmentsListOptions.Validate() expected nil, got %v", err)
	}
}

func TestCovCourse_MediaAttachmentUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *MediaAttachmentUpdateOptions
		wantErr bool
	}{
		{"valid", &MediaAttachmentUpdateOptions{AttachmentID: 5, Title: "Title"}, false},
		{"missing attachment-id", &MediaAttachmentUpdateOptions{AttachmentID: 0, Title: "Title"}, true},
		{"missing title", &MediaAttachmentUpdateOptions{AttachmentID: 5, Title: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MediaAttachmentUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// enrollment_terms.go
// ---------------------------------------------------------------------------

func TestCovCourse_EnrollmentTermsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentTermsListOptions
		wantErr bool
	}{
		{"valid", &EnrollmentTermsListOptions{AccountID: 1}, false},
		{"missing account-id", &EnrollmentTermsListOptions{AccountID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentTermsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_EnrollmentTermsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentTermsGetOptions
		wantErr bool
	}{
		{"valid", &EnrollmentTermsGetOptions{AccountID: 1, TermID: 5}, false},
		{"missing account-id", &EnrollmentTermsGetOptions{AccountID: 0, TermID: 5}, true},
		{"missing term-id", &EnrollmentTermsGetOptions{AccountID: 1, TermID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentTermsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_EnrollmentTermsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentTermsCreateOptions
		wantErr bool
	}{
		{"valid", &EnrollmentTermsCreateOptions{AccountID: 1, Name: "Fall 2024"}, false},
		{"missing account-id", &EnrollmentTermsCreateOptions{AccountID: 0, Name: "Fall 2024"}, true},
		{"missing name", &EnrollmentTermsCreateOptions{AccountID: 1, Name: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentTermsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_EnrollmentTermsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentTermsUpdateOptions
		wantErr bool
	}{
		{"valid", &EnrollmentTermsUpdateOptions{AccountID: 1, TermID: 5}, false},
		{"missing account-id", &EnrollmentTermsUpdateOptions{AccountID: 0, TermID: 5}, true},
		{"missing term-id", &EnrollmentTermsUpdateOptions{AccountID: 1, TermID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentTermsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_EnrollmentTermsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EnrollmentTermsDeleteOptions
		wantErr bool
	}{
		{"valid", &EnrollmentTermsDeleteOptions{AccountID: 1, TermID: 5}, false},
		{"missing account-id", &EnrollmentTermsDeleteOptions{AccountID: 0, TermID: 5}, true},
		{"missing term-id", &EnrollmentTermsDeleteOptions{AccountID: 1, TermID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentTermsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// blackout_dates.go
// ---------------------------------------------------------------------------

func TestCovCourse_BlackoutDatesListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlackoutDatesListOptions
		wantErr bool
	}{
		{"valid", &BlackoutDatesListOptions{CourseID: 1}, false},
		{"missing course-id", &BlackoutDatesListOptions{CourseID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlackoutDatesListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_BlackoutDatesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlackoutDatesGetOptions
		wantErr bool
	}{
		{"valid", &BlackoutDatesGetOptions{CourseID: 1, ID: 5}, false},
		{"missing course-id", &BlackoutDatesGetOptions{CourseID: 0, ID: 5}, true},
		{"missing id", &BlackoutDatesGetOptions{CourseID: 1, ID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlackoutDatesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_BlackoutDatesCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlackoutDatesCreateOptions
		wantErr bool
	}{
		{
			"valid",
			&BlackoutDatesCreateOptions{CourseID: 1, StartDate: "2024-12-24", EndDate: "2024-12-26", EventTitle: "Break"},
			false,
		},
		{
			"missing course-id",
			&BlackoutDatesCreateOptions{CourseID: 0, StartDate: "2024-12-24", EndDate: "2024-12-26", EventTitle: "Break"},
			true,
		},
		{
			"missing start-date",
			&BlackoutDatesCreateOptions{CourseID: 1, StartDate: "", EndDate: "2024-12-26", EventTitle: "Break"},
			true,
		},
		{
			"missing end-date",
			&BlackoutDatesCreateOptions{CourseID: 1, StartDate: "2024-12-24", EndDate: "", EventTitle: "Break"},
			true,
		},
		{
			"missing title",
			&BlackoutDatesCreateOptions{CourseID: 1, StartDate: "2024-12-24", EndDate: "2024-12-26", EventTitle: ""},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlackoutDatesCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_BlackoutDatesUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlackoutDatesUpdateOptions
		wantErr bool
	}{
		{"valid", &BlackoutDatesUpdateOptions{CourseID: 1, ID: 5}, false},
		{"missing course-id", &BlackoutDatesUpdateOptions{CourseID: 0, ID: 5}, true},
		{"missing id", &BlackoutDatesUpdateOptions{CourseID: 1, ID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlackoutDatesUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_BlackoutDatesDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BlackoutDatesDeleteOptions
		wantErr bool
	}{
		{"valid", &BlackoutDatesDeleteOptions{CourseID: 1, ID: 5}, false},
		{"missing course-id", &BlackoutDatesDeleteOptions{CourseID: 0, ID: 5}, true},
		{"missing id", &BlackoutDatesDeleteOptions{CourseID: 1, ID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BlackoutDatesDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// account_reports.go
// ---------------------------------------------------------------------------

func TestCovCourse_AccountReportsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AccountReportsListOptions
		wantErr bool
	}{
		{"valid", &AccountReportsListOptions{AccountID: 1}, false},
		{"missing account-id", &AccountReportsListOptions{AccountID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountReportsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_AccountReportsRunsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AccountReportsRunsOptions
		wantErr bool
	}{
		{"valid", &AccountReportsRunsOptions{AccountID: 1, ReportName: "grade_export_csv"}, false},
		{"missing account-id", &AccountReportsRunsOptions{AccountID: 0, ReportName: "grade_export_csv"}, true},
		{"missing report-name", &AccountReportsRunsOptions{AccountID: 1, ReportName: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountReportsRunsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_AccountReportsStartOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AccountReportsStartOptions
		wantErr bool
	}{
		{"valid", &AccountReportsStartOptions{AccountID: 1, ReportName: "grade_export_csv"}, false},
		{"missing account-id", &AccountReportsStartOptions{AccountID: 0, ReportName: "grade_export_csv"}, true},
		{"missing report-name", &AccountReportsStartOptions{AccountID: 1, ReportName: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountReportsStartOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_AccountReportsGetRunOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AccountReportsGetRunOptions
		wantErr bool
	}{
		{"valid", &AccountReportsGetRunOptions{AccountID: 1, ReportName: "grade_export_csv", RunID: 100}, false},
		{"missing account-id", &AccountReportsGetRunOptions{AccountID: 0, ReportName: "grade_export_csv", RunID: 100}, true},
		{"missing report-name", &AccountReportsGetRunOptions{AccountID: 1, ReportName: "", RunID: 100}, true},
		{"missing run-id", &AccountReportsGetRunOptions{AccountID: 1, ReportName: "grade_export_csv", RunID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountReportsGetRunOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_AccountReportsDeleteRunOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AccountReportsDeleteRunOptions
		wantErr bool
	}{
		{"valid", &AccountReportsDeleteRunOptions{AccountID: 1, ReportName: "grade_export_csv", RunID: 100}, false},
		{"missing account-id", &AccountReportsDeleteRunOptions{AccountID: 0, ReportName: "grade_export_csv", RunID: 100}, true},
		{"missing report-name", &AccountReportsDeleteRunOptions{AccountID: 1, ReportName: "", RunID: 100}, true},
		{"missing run-id", &AccountReportsDeleteRunOptions{AccountID: 1, ReportName: "grade_export_csv", RunID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountReportsDeleteRunOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_AccountReportsAbortRunOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AccountReportsAbortRunOptions
		wantErr bool
	}{
		{"valid", &AccountReportsAbortRunOptions{AccountID: 1, ReportName: "grade_export_csv", RunID: 100}, false},
		{"missing account-id", &AccountReportsAbortRunOptions{AccountID: 0, ReportName: "grade_export_csv", RunID: 100}, true},
		{"missing report-name", &AccountReportsAbortRunOptions{AccountID: 1, ReportName: "", RunID: 100}, true},
		{"missing run-id", &AccountReportsAbortRunOptions{AccountID: 1, ReportName: "grade_export_csv", RunID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountReportsAbortRunOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// grading_period_sets.go
// ---------------------------------------------------------------------------

func TestCovCourse_GradingPeriodSetsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingPeriodSetsListOptions
		wantErr bool
	}{
		{"valid", &GradingPeriodSetsListOptions{AccountID: 1}, false},
		{"missing account-id", &GradingPeriodSetsListOptions{AccountID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingPeriodSetsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingPeriodSetsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingPeriodSetsCreateOptions
		wantErr bool
	}{
		{"valid", &GradingPeriodSetsCreateOptions{AccountID: 1, Title: "2024 Periods"}, false},
		{"missing account-id", &GradingPeriodSetsCreateOptions{AccountID: 0, Title: "2024 Periods"}, true},
		{"missing title", &GradingPeriodSetsCreateOptions{AccountID: 1, Title: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingPeriodSetsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingPeriodSetsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingPeriodSetsUpdateOptions
		wantErr bool
	}{
		{"valid", &GradingPeriodSetsUpdateOptions{AccountID: 1, SetID: 5}, false},
		{"missing account-id", &GradingPeriodSetsUpdateOptions{AccountID: 0, SetID: 5}, true},
		{"missing set-id", &GradingPeriodSetsUpdateOptions{AccountID: 1, SetID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingPeriodSetsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingPeriodSetsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingPeriodSetsDeleteOptions
		wantErr bool
	}{
		{"valid", &GradingPeriodSetsDeleteOptions{AccountID: 1, SetID: 5}, false},
		{"missing account-id", &GradingPeriodSetsDeleteOptions{AccountID: 0, SetID: 5}, true},
		{"missing set-id", &GradingPeriodSetsDeleteOptions{AccountID: 1, SetID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingPeriodSetsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingPeriodSetsListPeriodsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingPeriodSetsListPeriodsOptions
		wantErr bool
	}{
		{"valid", &GradingPeriodSetsListPeriodsOptions{AccountID: 1}, false},
		{"missing account-id", &GradingPeriodSetsListPeriodsOptions{AccountID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingPeriodSetsListPeriodsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_GradingPeriodSetsDeletePeriodOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *GradingPeriodSetsDeletePeriodOptions
		wantErr bool
	}{
		{"valid", &GradingPeriodSetsDeletePeriodOptions{AccountID: 1, PeriodID: 5}, false},
		{"missing account-id", &GradingPeriodSetsDeletePeriodOptions{AccountID: 0, PeriodID: 5}, true},
		{"missing period-id", &GradingPeriodSetsDeletePeriodOptions{AccountID: 1, PeriodID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GradingPeriodSetsDeletePeriodOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// content_exports.go
// ---------------------------------------------------------------------------

func TestCovCourse_ContentExportsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContentExportsListOptions
		wantErr bool
	}{
		{"valid", &ContentExportsListOptions{CourseID: 1}, false},
		{"missing course-id", &ContentExportsListOptions{CourseID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentExportsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_ContentExportsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContentExportsGetOptions
		wantErr bool
	}{
		{"valid", &ContentExportsGetOptions{CourseID: 1, ID: 5}, false},
		{"missing course-id", &ContentExportsGetOptions{CourseID: 0, ID: 5}, true},
		{"missing id", &ContentExportsGetOptions{CourseID: 1, ID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentExportsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_ContentExportsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ContentExportsCreateOptions
		wantErr bool
	}{
		{"valid", &ContentExportsCreateOptions{CourseID: 1, ExportType: "common_cartridge"}, false},
		{"missing course-id", &ContentExportsCreateOptions{CourseID: 0, ExportType: "common_cartridge"}, true},
		{"missing export-type", &ContentExportsCreateOptions{CourseID: 1, ExportType: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentExportsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_EpubExportsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EpubExportsGetOptions
		wantErr bool
	}{
		{"valid", &EpubExportsGetOptions{CourseID: 1, ID: 5}, false},
		{"missing course-id", &EpubExportsGetOptions{CourseID: 0, ID: 5}, true},
		{"missing id", &EpubExportsGetOptions{CourseID: 1, ID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EpubExportsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// course_pacing.go
// ---------------------------------------------------------------------------

func TestCovCourse_CoursePacingGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CoursePacingGetOptions
		wantErr bool
	}{
		{"valid", &CoursePacingGetOptions{CourseID: 1, PaceID: 5}, false},
		{"missing course-id", &CoursePacingGetOptions{CourseID: 0, PaceID: 5}, true},
		{"missing pace-id", &CoursePacingGetOptions{CourseID: 1, PaceID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CoursePacingGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_CoursePacingCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CoursePacingCreateOptions
		wantErr bool
	}{
		{"valid", &CoursePacingCreateOptions{CourseID: 1}, false},
		{"missing course-id", &CoursePacingCreateOptions{CourseID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CoursePacingCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_CoursePacingUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CoursePacingUpdateOptions
		wantErr bool
	}{
		{"valid", &CoursePacingUpdateOptions{CourseID: 1, PaceID: 5}, false},
		{"missing course-id", &CoursePacingUpdateOptions{CourseID: 0, PaceID: 5}, true},
		{"missing pace-id", &CoursePacingUpdateOptions{CourseID: 1, PaceID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CoursePacingUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCovCourse_CoursePacingDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CoursePacingDeleteOptions
		wantErr bool
	}{
		{"valid", &CoursePacingDeleteOptions{CourseID: 1, PaceID: 5}, false},
		{"missing course-id", &CoursePacingDeleteOptions{CourseID: 0, PaceID: 5}, true},
		{"missing pace-id", &CoursePacingDeleteOptions{CourseID: 1, PaceID: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CoursePacingDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
