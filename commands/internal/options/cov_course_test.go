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
