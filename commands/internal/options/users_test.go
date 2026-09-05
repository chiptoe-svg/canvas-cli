package options

import "testing"

func TestUsersListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *UsersListOptions
		wantErr bool
	}{
		{
			name:    "valid - course ID",
			opts:    &UsersListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &UsersListOptions{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsersGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *UsersGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &UsersGetOptions{UserID: 1},
			wantErr: false,
		},
		{
			name:    "zero user ID",
			opts:    &UsersGetOptions{UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsersMeOptions_Validate(t *testing.T) {
	opts := &UsersMeOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("UsersMeOptions.Validate() error = %v, want nil", err)
	}
}

func TestUsersSearchOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *UsersSearchOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &UsersSearchOptions{SearchTerm: "john"},
			wantErr: false,
		},
		{
			name:    "missing search term",
			opts:    &UsersSearchOptions{SearchTerm: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UsersSearchOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
