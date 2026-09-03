package options

import (
	"strings"
	"testing"
)

func validScheduleOptions() *ScheduleOptions {
	return &ScheduleOptions{CourseID: 1, ID: 456, Due: "4:50pm"}
}

func TestScheduleOptions_Validate(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(o *ScheduleOptions)
		wantErr string
	}{
		{"valid by id", func(o *ScheduleOptions) {}, ""},
		{"valid by match", func(o *ScheduleOptions) { o.ID = 0; o.Match = "attendance" }, ""},
		{"valid regex match", func(o *ScheduleOptions) { o.ID = 0; o.Match = "/^Attendance/" }, ""},
		{"valid clear only", func(o *ScheduleOptions) { o.Due = ""; o.Clear = []string{"closed"} }, ""},
		{"valid clear list", func(o *ScheduleOptions) { o.Clear = []string{"available,closed"} }, ""},
		{"valid type", func(o *ScheduleOptions) { o.Type = "quiz" }, ""},
		{"no course", func(o *ScheduleOptions) { o.CourseID = 0 }, "course-id is required"},
		{"neither id nor match", func(o *ScheduleOptions) { o.ID = 0 }, "one of --id or --match"},
		{"both id and match", func(o *ScheduleOptions) { o.Match = "x" }, "mutually exclusive"},
		{"negative id", func(o *ScheduleOptions) { o.ID = -3 }, "greater than 0"},
		{"bad type", func(o *ScheduleOptions) { o.Type = "page" }, "unknown --type"},
		{"bad regex", func(o *ScheduleOptions) { o.ID = 0; o.Match = "/(/" }, "invalid --match regex"},
		{"bad clear", func(o *ScheduleOptions) { o.Clear = []string{"unlock_at"} }, "unknown --clear"},
		{"set and clear same field", func(o *ScheduleOptions) { o.Clear = []string{"due"} }, "--due and --clear due are mutually exclusive"},
		{"nothing to change", func(o *ScheduleOptions) { o.Due = "" }, "nothing to change"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			o := validScheduleOptions()
			c.mutate(o)
			err := o.Validate()
			if c.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("error = %v, want %q", err, c.wantErr)
			}
		})
	}
}

func TestScheduleOptions_Helpers(t *testing.T) {
	o := &ScheduleOptions{}
	if o.ItemType() != ScheduleTypeAll {
		t.Errorf("default type = %q", o.ItemType())
	}
	o.Type = "assignment"
	if o.ItemType() != "assignment" {
		t.Errorf("type = %q", o.ItemType())
	}
	if o.HasChanges() {
		t.Error("HasChanges on empty options")
	}
	o.Clear = []string{" Due ", "closed,available"}
	set, err := o.ClearSet()
	if err != nil || !set["due"] || !set["closed"] || !set["available"] || len(set) != 3 {
		t.Errorf("ClearSet = %v, %v", set, err)
	}
	if !o.HasChanges() {
		t.Error("HasChanges with --clear")
	}
}

func TestCompileNameMatch(t *testing.T) {
	all, err := CompileNameMatch("--match", "")
	if err != nil || !all("anything") {
		t.Errorf("empty match: %v", err)
	}
	sub, _ := CompileNameMatch("--match", "ATTEND")
	if !sub("Attendance 9/9") || sub("Essay") {
		t.Error("substring match is case-insensitive")
	}
	re, _ := CompileNameMatch("--match", "/^Attendance \\d+/")
	if !re("Attendance 12") || re("attendance 12") || re("Weekly Attendance 12") {
		t.Error("regex match is anchored and case-sensitive")
	}
	if _, err := CompileNameMatch("--match", "/[/"); err == nil || !strings.Contains(err.Error(), "invalid --match regex") {
		t.Errorf("bad regex: %v", err)
	}
}
