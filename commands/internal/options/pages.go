package options

import "fmt"

// pagesValidateContext checks that exactly one of courseID or groupID is set.
func pagesValidateContext(courseID, groupID int64) error {
	if courseID > 0 && groupID > 0 {
		return fmt.Errorf("specify either --course-id or --group-id, not both")
	}
	if courseID <= 0 && groupID <= 0 {
		return fmt.Errorf("one of --course-id or --group-id is required")
	}
	return nil
}

// PagesListOptions contains options for listing pages
type PagesListOptions struct {
	CourseID   int64
	GroupID    int64
	Sort       string
	Order      string
	SearchTerm string
	Published  string
	Include    []string
}

// Validate validates the options
func (o *PagesListOptions) Validate() error {
	return pagesValidateContext(o.CourseID, o.GroupID)
}

// PagesGetOptions contains options for getting a page
type PagesGetOptions struct {
	CourseID int64
	GroupID  int64
	URLOrID  string
}

// Validate validates the options
func (o *PagesGetOptions) Validate() error {
	if err := pagesValidateContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.URLOrID == "" {
		return fmt.Errorf("page url or id is required")
	}
	return nil
}

// PagesFrontOptions contains options for getting the front page
type PagesFrontOptions struct {
	CourseID int64
	GroupID  int64
}

// Validate validates the options
func (o *PagesFrontOptions) Validate() error {
	return pagesValidateContext(o.CourseID, o.GroupID)
}

// PagesCreateOptions contains options for creating a page
type PagesCreateOptions struct {
	CourseID     int64
	GroupID      int64
	Title        string
	Body         string
	EditingRoles string
	NotifyUpdate bool
	Published    bool
	FrontPage    bool
	PublishAt    string
}

// Validate validates the options
func (o *PagesCreateOptions) Validate() error {
	if err := pagesValidateContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

// PagesUpdateOptions contains options for updating a page
type PagesUpdateOptions struct {
	CourseID     int64
	GroupID      int64
	URLOrID      string
	Title        string
	Body         string
	EditingRoles string
	NotifyUpdate bool
	Published    bool
	FrontPage    bool
	PublishAt    string
	// Track which fields were actually set
	TitleSet        bool
	BodySet         bool
	EditingRolesSet bool
	NotifyUpdateSet bool
	PublishedSet    bool
	FrontPageSet    bool
	PublishAtSet    bool
}

// Validate validates the options
func (o *PagesUpdateOptions) Validate() error {
	if err := pagesValidateContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.URLOrID == "" {
		return fmt.Errorf("page url or id is required")
	}
	return nil
}

// PagesDeleteOptions contains options for deleting a page
type PagesDeleteOptions struct {
	CourseID int64
	GroupID  int64
	URLOrID  string
	Force    bool
}

// Validate validates the options
func (o *PagesDeleteOptions) Validate() error {
	if err := pagesValidateContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.URLOrID == "" {
		return fmt.Errorf("page url or id is required")
	}
	return nil
}

// PagesDuplicateOptions contains options for duplicating a page (course-only)
type PagesDuplicateOptions struct {
	CourseID int64
	URLOrID  string
}

// Validate validates the options
func (o *PagesDuplicateOptions) Validate() error {
	if o.CourseID <= 0 {
		return fmt.Errorf("course-id is required and must be greater than 0")
	}
	if o.URLOrID == "" {
		return fmt.Errorf("page url or id is required")
	}
	return nil
}

// PagesRevisionsOptions contains options for listing page revisions
type PagesRevisionsOptions struct {
	CourseID int64
	GroupID  int64
	URLOrID  string
}

// Validate validates the options
func (o *PagesRevisionsOptions) Validate() error {
	if err := pagesValidateContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.URLOrID == "" {
		return fmt.Errorf("page url or id is required")
	}
	return nil
}

// PagesRevertOptions contains options for reverting to a revision
type PagesRevertOptions struct {
	CourseID   int64
	GroupID    int64
	URLOrID    string
	RevisionID int64
}

// Validate validates the options
func (o *PagesRevertOptions) Validate() error {
	if err := pagesValidateContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.URLOrID == "" {
		return fmt.Errorf("page url or id is required")
	}
	if o.RevisionID <= 0 {
		return fmt.Errorf("revision-id is required and must be greater than 0")
	}
	return nil
}
