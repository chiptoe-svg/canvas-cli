package options

import "fmt"

// BookmarksListOptions contains options for listing bookmarks
type BookmarksListOptions struct{}

// Validate validates the options
func (o *BookmarksListOptions) Validate() error { return nil }

// BookmarksGetOptions contains options for getting a bookmark
type BookmarksGetOptions struct {
	ID int64
}

// Validate validates the options
func (o *BookmarksGetOptions) Validate() error {
	if o.ID <= 0 {
		return fmt.Errorf("id is required and must be greater than 0")
	}
	return nil
}

// BookmarksCreateOptions contains options for creating a bookmark
type BookmarksCreateOptions struct {
	Name     string
	URL      string
	Position int
}

// Validate validates the options
func (o *BookmarksCreateOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("name is required")
	}
	if o.URL == "" {
		return fmt.Errorf("url is required")
	}
	return nil
}

// BookmarksUpdateOptions contains options for updating a bookmark
type BookmarksUpdateOptions struct {
	ID       int64
	Name     string
	URL      string
	Position int
}

// Validate validates the options
func (o *BookmarksUpdateOptions) Validate() error {
	if o.ID <= 0 {
		return fmt.Errorf("id is required and must be greater than 0")
	}
	return nil
}

// BookmarksDeleteOptions contains options for deleting a bookmark
type BookmarksDeleteOptions struct {
	ID int64
}

// Validate validates the options
func (o *BookmarksDeleteOptions) Validate() error {
	if o.ID <= 0 {
		return fmt.Errorf("id is required and must be greater than 0")
	}
	return nil
}
