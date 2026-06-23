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
// favorites.go
// ---------------------------------------------------------------------------

func TestFavoritesListCoursesOptions_Validate(t *testing.T) {
	opts := &FavoritesListCoursesOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("FavoritesListCoursesOptions.Validate() unexpected error: %v", err)
	}
}

func TestFavoritesAddCourseOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *FavoritesAddCourseOptions
		wantErr bool
	}{
		{"valid", &FavoritesAddCourseOptions{ID: 1}, false},
		{"zero ID", &FavoritesAddCourseOptions{ID: 0}, true},
		{"negative ID", &FavoritesAddCourseOptions{ID: -5}, true},
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

func TestFavoritesRemoveCourseOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *FavoritesRemoveCourseOptions
		wantErr bool
	}{
		{"valid", &FavoritesRemoveCourseOptions{ID: 10}, false},
		{"zero ID", &FavoritesRemoveCourseOptions{ID: 0}, true},
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

func TestFavoritesListGroupsOptions_Validate(t *testing.T) {
	opts := &FavoritesListGroupsOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("FavoritesListGroupsOptions.Validate() unexpected error: %v", err)
	}
}

func TestFavoritesAddGroupOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *FavoritesAddGroupOptions
		wantErr bool
	}{
		{"valid", &FavoritesAddGroupOptions{ID: 456}, false},
		{"zero ID", &FavoritesAddGroupOptions{ID: 0}, true},
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

func TestFavoritesRemoveGroupOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *FavoritesRemoveGroupOptions
		wantErr bool
	}{
		{"valid", &FavoritesRemoveGroupOptions{ID: 789}, false},
		{"zero ID", &FavoritesRemoveGroupOptions{ID: 0}, true},
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
// bookmarks.go
// ---------------------------------------------------------------------------

func TestBookmarksListOptions_Validate(t *testing.T) {
	opts := &BookmarksListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("BookmarksListOptions.Validate() unexpected error: %v", err)
	}
}

func TestBookmarksGetOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *BookmarksGetOptions
		wantErr bool
	}{
		{"valid", &BookmarksGetOptions{ID: 1}, false},
		{"zero ID", &BookmarksGetOptions{ID: 0}, true},
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

func TestBookmarksCreateOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *BookmarksCreateOptions
		wantErr bool
	}{
		{"valid", &BookmarksCreateOptions{Name: "My Page", URL: "https://example.com"}, false},
		{"missing name", &BookmarksCreateOptions{Name: "", URL: "https://example.com"}, true},
		{"missing url", &BookmarksCreateOptions{Name: "My Page", URL: ""}, true},
		{"both missing", &BookmarksCreateOptions{}, true},
		{"with position", &BookmarksCreateOptions{Name: "Page", URL: "https://x.com", Position: 3}, false},
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

func TestBookmarksUpdateOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *BookmarksUpdateOptions
		wantErr bool
	}{
		{"valid", &BookmarksUpdateOptions{ID: 5}, false},
		{"with fields", &BookmarksUpdateOptions{ID: 5, Name: "New", URL: "https://new.com"}, false},
		{"zero ID", &BookmarksUpdateOptions{ID: 0}, true},
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

func TestBookmarksDeleteOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *BookmarksDeleteOptions
		wantErr bool
	}{
		{"valid", &BookmarksDeleteOptions{ID: 123}, false},
		{"zero ID", &BookmarksDeleteOptions{ID: 0}, true},
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
// course_nicknames.go
// ---------------------------------------------------------------------------

func TestCourseNicknamesListOptions_Validate(t *testing.T) {
	opts := &CourseNicknamesListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("CourseNicknamesListOptions.Validate() unexpected error: %v", err)
	}
}

func TestCourseNicknamesGetOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *CourseNicknamesGetOptions
		wantErr bool
	}{
		{"valid", &CourseNicknamesGetOptions{CourseID: 10}, false},
		{"zero ID", &CourseNicknamesGetOptions{CourseID: 0}, true},
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

func TestCourseNicknamesSetOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *CourseNicknamesSetOptions
		wantErr bool
	}{
		{"valid", &CourseNicknamesSetOptions{CourseID: 1, Nickname: "My Course"}, false},
		{"zero course ID", &CourseNicknamesSetOptions{CourseID: 0, Nickname: "My Course"}, true},
		{"empty nickname", &CourseNicknamesSetOptions{CourseID: 1, Nickname: ""}, true},
		{"both missing", &CourseNicknamesSetOptions{}, true},
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

func TestCourseNicknamesDeleteOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *CourseNicknamesDeleteOptions
		wantErr bool
	}{
		{"valid", &CourseNicknamesDeleteOptions{CourseID: 7}, false},
		{"zero ID", &CourseNicknamesDeleteOptions{CourseID: 0}, true},
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
// observees.go
// ---------------------------------------------------------------------------

func TestObserveesListOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *ObserveesListOptions
		wantErr bool
	}{
		{"valid", &ObserveesListOptions{UserID: 123}, false},
		{"zero ID", &ObserveesListOptions{UserID: 0}, true},
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

func TestObserveesGetOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *ObserveesGetOptions
		wantErr bool
	}{
		{"valid", &ObserveesGetOptions{UserID: 1, ObserveeID: 2}, false},
		{"zero user ID", &ObserveesGetOptions{UserID: 0, ObserveeID: 2}, true},
		{"zero observee ID", &ObserveesGetOptions{UserID: 1, ObserveeID: 0}, true},
		{"both zero", &ObserveesGetOptions{}, true},
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

func TestObserveesDeleteOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *ObserveesDeleteOptions
		wantErr bool
	}{
		{"valid", &ObserveesDeleteOptions{UserID: 10, ObserveeID: 20}, false},
		{"zero user ID", &ObserveesDeleteOptions{UserID: 0, ObserveeID: 20}, true},
		{"zero observee ID", &ObserveesDeleteOptions{UserID: 10, ObserveeID: 0}, true},
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

func TestObserversListOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *ObserversListOptions
		wantErr bool
	}{
		{"valid", &ObserversListOptions{UserID: 50}, false},
		{"zero ID", &ObserversListOptions{UserID: 0}, true},
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
// user_features.go
// ---------------------------------------------------------------------------

func TestUserFeaturesListOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UserFeaturesListOptions
		wantErr bool
	}{
		{"valid", &UserFeaturesListOptions{UserID: 100}, false},
		{"zero ID", &UserFeaturesListOptions{UserID: 0}, true},
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

func TestUserFeaturesListEnabledOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UserFeaturesListEnabledOptions
		wantErr bool
	}{
		{"valid", &UserFeaturesListEnabledOptions{UserID: 5}, false},
		{"zero ID", &UserFeaturesListEnabledOptions{UserID: 0}, true},
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

func TestUserFeaturesGetFlagOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UserFeaturesGetFlagOptions
		wantErr bool
	}{
		{"valid", &UserFeaturesGetFlagOptions{UserID: 1, Feature: "new_gradebook"}, false},
		{"zero user ID", &UserFeaturesGetFlagOptions{UserID: 0, Feature: "feat"}, true},
		{"empty feature", &UserFeaturesGetFlagOptions{UserID: 1, Feature: ""}, true},
		{"both missing", &UserFeaturesGetFlagOptions{}, true},
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

func TestUserFeaturesSetFlagOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UserFeaturesSetFlagOptions
		wantErr bool
	}{
		{"valid", &UserFeaturesSetFlagOptions{UserID: 1, Feature: "new_gradebook", State: "on"}, false},
		{"zero user ID", &UserFeaturesSetFlagOptions{UserID: 0, Feature: "feat", State: "on"}, true},
		{"empty feature", &UserFeaturesSetFlagOptions{UserID: 1, Feature: "", State: "on"}, true},
		{"empty state", &UserFeaturesSetFlagOptions{UserID: 1, Feature: "feat", State: ""}, true},
		{"all missing", &UserFeaturesSetFlagOptions{}, true},
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

func TestUserFeaturesDeleteFlagOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *UserFeaturesDeleteFlagOptions
		wantErr bool
	}{
		{"valid", &UserFeaturesDeleteFlagOptions{UserID: 1, Feature: "new_gradebook"}, false},
		{"zero user ID", &UserFeaturesDeleteFlagOptions{UserID: 0, Feature: "feat"}, true},
		{"empty feature", &UserFeaturesDeleteFlagOptions{UserID: 1, Feature: ""}, true},
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
// comm_channels.go
// ---------------------------------------------------------------------------

func TestCommChannelsListOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *CommChannelsListOptions
		wantErr bool
	}{
		{"valid", &CommChannelsListOptions{UserID: 123}, false},
		{"zero ID", &CommChannelsListOptions{UserID: 0}, true},
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

func TestCommChannelsCreateOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *CommChannelsCreateOptions
		wantErr bool
	}{
		{"valid", &CommChannelsCreateOptions{UserID: 1, Address: "u@e.com", Type: "email"}, false},
		{"with skip", &CommChannelsCreateOptions{UserID: 1, Address: "+1555", Type: "sms", SkipConfirmation: true}, false},
		{"zero user ID", &CommChannelsCreateOptions{UserID: 0, Address: "u@e.com", Type: "email"}, true},
		{"empty address", &CommChannelsCreateOptions{UserID: 1, Address: "", Type: "email"}, true},
		{"empty type", &CommChannelsCreateOptions{UserID: 1, Address: "u@e.com", Type: ""}, true},
		{"all missing", &CommChannelsCreateOptions{}, true},
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

func TestCommChannelsDeleteOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		opts    *CommChannelsDeleteOptions
		wantErr bool
	}{
		{"valid", &CommChannelsDeleteOptions{UserID: 1, ChannelID: 456}, false},
		{"zero user ID", &CommChannelsDeleteOptions{UserID: 0, ChannelID: 456}, true},
		{"zero channel ID", &CommChannelsDeleteOptions{UserID: 1, ChannelID: 0}, true},
		{"both zero", &CommChannelsDeleteOptions{}, true},
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
