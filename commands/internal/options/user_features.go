package options

import "fmt"

// UserFeaturesListOptions contains options for listing user features
type UserFeaturesListOptions struct {
	UserID int64
}

// Validate validates the options
func (o *UserFeaturesListOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UserFeaturesListEnabledOptions contains options for listing enabled user features
type UserFeaturesListEnabledOptions struct {
	UserID int64
}

// Validate validates the options
func (o *UserFeaturesListEnabledOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// UserFeaturesGetFlagOptions contains options for getting a feature flag
type UserFeaturesGetFlagOptions struct {
	UserID  int64
	Feature string
}

// Validate validates the options
func (o *UserFeaturesGetFlagOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.Feature == "" {
		return fmt.Errorf("feature is required")
	}
	return nil
}

// UserFeaturesSetFlagOptions contains options for setting a feature flag
type UserFeaturesSetFlagOptions struct {
	UserID  int64
	Feature string
	State   string
}

// Validate validates the options
func (o *UserFeaturesSetFlagOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.Feature == "" {
		return fmt.Errorf("feature is required")
	}
	if o.State == "" {
		return fmt.Errorf("state is required")
	}
	return nil
}

// UserFeaturesDeleteFlagOptions contains options for deleting a feature flag
type UserFeaturesDeleteFlagOptions struct {
	UserID  int64
	Feature string
}

// Validate validates the options
func (o *UserFeaturesDeleteFlagOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.Feature == "" {
		return fmt.Errorf("feature is required")
	}
	return nil
}
