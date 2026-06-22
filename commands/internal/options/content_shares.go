package options

import "fmt"

// ContentSharesListSentOptions contains options for listing sent content shares
type ContentSharesListSentOptions struct {
	UserID int64
}

// Validate validates the options
func (o *ContentSharesListSentOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// ContentSharesListReceivedOptions contains options for listing received content shares
type ContentSharesListReceivedOptions struct {
	UserID int64
}

// Validate validates the options
func (o *ContentSharesListReceivedOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// ContentSharesGetOptions contains options for getting a content share
type ContentSharesGetOptions struct {
	UserID int64
	ID     int64
}

// Validate validates the options
func (o *ContentSharesGetOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.ID <= 0 {
		return fmt.Errorf("id is required and must be greater than 0")
	}
	return nil
}

// ContentSharesDeleteOptions contains options for deleting a content share
type ContentSharesDeleteOptions struct {
	UserID int64
	ID     int64
}

// Validate validates the options
func (o *ContentSharesDeleteOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.ID <= 0 {
		return fmt.Errorf("id is required and must be greater than 0")
	}
	return nil
}
