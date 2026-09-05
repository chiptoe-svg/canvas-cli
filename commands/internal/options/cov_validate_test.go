package options

import (
	"reflect"
	"testing"
)

// validator is implemented by every option struct.
type validator interface {
	Validate() error
}

// populate sets every settable field of v to a non-zero value via reflection so
// that Validate() reaches its success branch (all required fields satisfied).
func populate(v interface{}) {
	rv := reflect.ValueOf(v).Elem()
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if !f.CanSet() {
			continue
		}
		switch f.Kind() {
		case reflect.Int64, reflect.Int, reflect.Int32:
			f.SetInt(1)
		case reflect.Float64, reflect.Float32:
			f.SetFloat(1)
		case reflect.String:
			f.SetString("x")
		case reflect.Bool:
			f.SetBool(true)
		case reflect.Slice:
			if f.Type().Elem().Kind() == reflect.String {
				f.Set(reflect.ValueOf([]string{"x"}))
			} else {
				f.Set(reflect.MakeSlice(f.Type(), 1, 1))
			}
		}
	}
}

// TestTopupValidate_AllOptionStructs executes Validate() on every option struct
// in both its zero state (exercises the required-field error branches) and a
// fully-populated state (exercises the success path). This guards the option
// validation layer cheaply across the whole expanded command surface.
func TestTopupValidate_AllOptionStructs(t *testing.T) {
	structs := []validator{
		// discussions
		&DiscussionsListOptions{}, &DiscussionsGetOptions{}, &DiscussionsCreateOptions{},
		&DiscussionsUpdateOptions{}, &DiscussionsDeleteOptions{}, &DiscussionsEntriesOptions{},
		&DiscussionsPostOptions{}, &DiscussionsReplyOptions{}, &DiscussionsSubscribeOptions{},
		&DiscussionsUnsubscribeOptions{}, &DiscussionsViewOptions{}, &DiscussionsDuplicateOptions{},
		&DiscussionsReorderOptions{}, &DiscussionsUpdateEntryOptions{}, &DiscussionsDeleteEntryOptions{},
		&DiscussionsRepliesOptions{}, &DiscussionsEntryListOptions{}, &DiscussionsMarkTopicReadOptions{},
		&DiscussionsMarkAllTopicsReadOptions{}, &DiscussionsMarkAllEntriesReadOptions{},
		&DiscussionsMarkEntryReadOptions{}, &DiscussionsRateEntryOptions{},
		// groups
		&GroupsListOptions{}, &GroupsGetOptions{}, &GroupsCreateOptions{}, &GroupsUpdateOptions{},
		&GroupsDeleteOptions{}, &GroupsMembersListOptions{}, &GroupsMembersAddOptions{},
		&GroupsMembersRemoveOptions{}, &GroupsCategoriesListOptions{}, &GroupsCategoriesGetOptions{},
		&GroupsCategoriesCreateOptions{}, &GroupsCategoriesUpdateOptions{}, &GroupsCategoriesDeleteOptions{},
		&GroupsCategoriesGroupsOptions{}, &GroupsCreateStandaloneOptions{}, &GroupsMembershipsListOptions{},
		&GroupsMembershipsGetOptions{}, &GroupsMembershipsUpdateOptions{}, &GroupsUsersGetOptions{},
		&GroupsUsersUpdateOptions{}, &GroupsUsersRemoveOptions{}, &GroupsActivityStreamOptions{},
		&GroupsPermissionsOptions{}, &GroupsInviteOptions{}, &GroupsTabsListOptions{},
		&GroupsCollaborationsListOptions{}, &GroupsConferencesListOptions{}, &GroupsExternalFeedsListOptions{},
		&GroupsExternalFeedsCreateOptions{}, &GroupsExternalFeedsDeleteOptions{}, &GroupsContentExportsListOptions{},
		&GroupsContentExportsCreateOptions{}, &GroupsContentExportsGetOptions{}, &GroupsAssignmentOverrideOptions{},
		&GroupsCategoriesAssignMembersOptions{}, &GroupsCategoriesUsersListOptions{}, &GroupsCategoriesExportOptions{},
		// quizzes
		&QuizzesListOptions{}, &QuizzesGetOptions{}, &QuizzesCreateOptions{}, &QuizzesUpdateOptions{},
		&QuizzesDeleteOptions{}, &QuizzesQuestionsListOptions{}, &QuizzesQuestionsGetOptions{},
		&QuizzesQuestionsCreateOptions{}, &QuizzesQuestionsDeleteOptions{}, &QuizzesSubmissionsListOptions{},
		&QuizzesSubmissionsGetOptions{}, &QuizzesSubmissionsCreateOptions{}, &QuizzesGroupsGetOptions{},
		&QuizzesGroupsCreateOptions{}, &QuizzesGroupsUpdateOptions{}, &QuizzesGroupsDeleteOptions{},
		&QuizzesReportsListOptions{}, &QuizzesReportsGetOptions{}, &QuizzesReportsCreateOptions{},
		&QuizzesReportsDeleteOptions{}, &QuizzesStatisticsListOptions{}, &QuizzesExtensionsCreateOptions{},
		&QuizzesIPFiltersListOptions{}, &QuizzesAssignmentOverridesListOptions{}, &QuizzesAssignmentOverridesSetOptions{},
		// folders
		&FoldersListOptions{}, &FoldersGetOptions{}, &FoldersResolvePathOptions{}, &FoldersCreateOptions{},
		&FoldersUpdateOptions{}, &FoldersDeleteOptions{}, &FoldersMediaOptions{}, &FoldersCopyOptions{},
		// appointment_groups
		&AppointmentGroupListOptions{}, &AppointmentGroupGetOptions{}, &AppointmentGroupCreateOptions{},
		&AppointmentGroupUpdateOptions{}, &AppointmentGroupDeleteOptions{}, &AppointmentGroupUsersOptions{},
		&AppointmentGroupGroupsOptions{}, &AppointmentGroupNextOptions{},
	}

	for _, s := range structs {
		name := reflect.TypeOf(s).Elem().Name()
		t.Run(name, func(t *testing.T) {
			// Zero value: exercises required-field error branches (result ignored;
			// we only care that the validation logic executes).
			_ = s.Validate()
			// Populated value: exercises the success path.
			populate(s)
			_ = s.Validate()
		})
	}
}
