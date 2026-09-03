package options

import (
	"fmt"
	"strings"
)

// SubmissionsExcuseOptions contains options for `submissions excuse`.
type SubmissionsExcuseOptions struct {
	CourseID   int64
	Student    string // id, exact name, sortable name, login, or a unique substring
	Assignment string // assignment id, quiz id, exact name, or a unique substring
	Unexcuse   bool
	Force      bool
	DryRun     bool
}

// Validate validates the options.
func (o *SubmissionsExcuseOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if strings.TrimSpace(o.Student) == "" {
		return fmt.Errorf("--student is required (an id, a name, or a unique part of a name)")
	}
	if strings.TrimSpace(o.Assignment) == "" {
		return fmt.Errorf("--assignment is required (an assignment or quiz id, a name, or a unique part of a name)")
	}
	return nil
}

// Excused is the state the command sets.
func (o *SubmissionsExcuseOptions) Excused() bool {
	return !o.Unexcuse
}
