package options

import "fmt"

// FavoritesListCoursesOptions contains options for listing favorite courses
type FavoritesListCoursesOptions struct{}

// Validate validates the options
func (o *FavoritesListCoursesOptions) Validate() error { return nil }

// FavoritesAddCourseOptions contains options for adding a favorite course
type FavoritesAddCourseOptions struct {
	ID int64
}

// Validate validates the options
func (o *FavoritesAddCourseOptions) Validate() error {
	if o.ID <= 0 {
		return fmt.Errorf("id is required and must be greater than 0")
	}
	return nil
}

// FavoritesRemoveCourseOptions contains options for removing a favorite course
type FavoritesRemoveCourseOptions struct {
	ID int64
}

// Validate validates the options
func (o *FavoritesRemoveCourseOptions) Validate() error {
	if o.ID <= 0 {
		return fmt.Errorf("id is required and must be greater than 0")
	}
	return nil
}

// FavoritesListGroupsOptions contains options for listing favorite groups
type FavoritesListGroupsOptions struct{}

// Validate validates the options
func (o *FavoritesListGroupsOptions) Validate() error { return nil }

// FavoritesAddGroupOptions contains options for adding a favorite group
type FavoritesAddGroupOptions struct {
	ID int64
}

// Validate validates the options
func (o *FavoritesAddGroupOptions) Validate() error {
	if o.ID <= 0 {
		return fmt.Errorf("id is required and must be greater than 0")
	}
	return nil
}

// FavoritesRemoveGroupOptions contains options for removing a favorite group
type FavoritesRemoveGroupOptions struct {
	ID int64
}

// Validate validates the options
func (o *FavoritesRemoveGroupOptions) Validate() error {
	if o.ID <= 0 {
		return fmt.Errorf("id is required and must be greater than 0")
	}
	return nil
}
