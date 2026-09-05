package options

import "testing"

func TestCoursesUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *CoursesUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &CoursesUpdateOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &CoursesUpdateOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CoursesUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
