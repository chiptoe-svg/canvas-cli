package options

import "testing"

// ---------------------------------------------------------------------------
// users.go — extra option types not covered by users_test.go
// ---------------------------------------------------------------------------

func TestUsersProfileOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UsersProfileOptions
		wantErr bool
	}{
		{"valid", &UsersProfileOptions{UserID: 1}, false},
		{"zero ID", &UsersProfileOptions{UserID: 0}, true},
		{"negative ID", &UsersProfileOptions{UserID: -1}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUsersSettingsOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UsersSettingsOptions
		wantErr bool
	}{
		{"valid", &UsersSettingsOptions{UserID: 42}, false},
		{"zero ID", &UsersSettingsOptions{UserID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUsersUpdateSettingsOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UsersUpdateSettingsOptions
		wantErr bool
	}{
		{"valid", &UsersUpdateSettingsOptions{UserID: 1}, false},
		{"with flags", &UsersUpdateSettingsOptions{UserID: 5, ManualMarkAsRead: true, CollapseGlobalNav: true}, false},
		{"zero ID", &UsersUpdateSettingsOptions{UserID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUsersPageViewsOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UsersPageViewsOptions
		wantErr bool
	}{
		{"valid", &UsersPageViewsOptions{UserID: 10}, false},
		{"zero ID", &UsersPageViewsOptions{UserID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUsersLoginsOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UsersLoginsOptions
		wantErr bool
	}{
		{"valid", &UsersLoginsOptions{UserID: 7}, false},
		{"zero ID", &UsersLoginsOptions{UserID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUsersCoursesOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UsersCoursesOptions
		wantErr bool
	}{
		{"valid", &UsersCoursesOptions{UserID: 3}, false},
		{"zero ID", &UsersCoursesOptions{UserID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUsersMissingSubmissionsOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UsersMissingSubmissionsOptions
		wantErr bool
	}{
		{"valid", &UsersMissingSubmissionsOptions{UserID: 99}, false},
		{"zero ID", &UsersMissingSubmissionsOptions{UserID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUsersActivityStreamOptions_Validate(t *testing.T) {
	opts := &UsersActivityStreamOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("UsersActivityStreamOptions.Validate() unexpected error: %v", err)
	}
}

func TestUsersTodoOptions_Validate(t *testing.T) {
	opts := &UsersTodoOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("UsersTodoOptions.Validate() unexpected error: %v", err)
	}
}

func TestUsersUpcomingEventsOptions_Validate(t *testing.T) {
	opts := &UsersUpcomingEventsOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("UsersUpcomingEventsOptions.Validate() unexpected error: %v", err)
	}
}

func TestUsersMergeOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UsersMergeOptions
		wantErr bool
	}{
		{"valid", &UsersMergeOptions{UserID: 1, DestinationUserID: 2}, false},
		{"zero source ID", &UsersMergeOptions{UserID: 0, DestinationUserID: 2}, true},
		{"zero dest ID", &UsersMergeOptions{UserID: 1, DestinationUserID: 0}, true},
		{"both zero", &UsersMergeOptions{UserID: 0, DestinationUserID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUsersSplitOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UsersSplitOptions
		wantErr bool
	}{
		{"valid", &UsersSplitOptions{UserID: 5}, false},
		{"zero ID", &UsersSplitOptions{UserID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// content_shares.go
// ---------------------------------------------------------------------------

func TestContentSharesListSentOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *ContentSharesListSentOptions
		wantErr bool
	}{
		{"valid", &ContentSharesListSentOptions{UserID: 1}, false},
		{"zero ID", &ContentSharesListSentOptions{UserID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestContentSharesListReceivedOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *ContentSharesListReceivedOptions
		wantErr bool
	}{
		{"valid", &ContentSharesListReceivedOptions{UserID: 2}, false},
		{"zero ID", &ContentSharesListReceivedOptions{UserID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestContentSharesGetOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *ContentSharesGetOptions
		wantErr bool
	}{
		{"valid", &ContentSharesGetOptions{UserID: 1, ID: 100}, false},
		{"zero user ID", &ContentSharesGetOptions{UserID: 0, ID: 100}, true},
		{"zero share ID", &ContentSharesGetOptions{UserID: 1, ID: 0}, true},
		{"both zero", &ContentSharesGetOptions{}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestContentSharesDeleteOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *ContentSharesDeleteOptions
		wantErr bool
	}{
		{"valid", &ContentSharesDeleteOptions{UserID: 1, ID: 200}, false},
		{"zero user ID", &ContentSharesDeleteOptions{UserID: 0, ID: 200}, true},
		{"zero share ID", &ContentSharesDeleteOptions{UserID: 1, ID: 0}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
