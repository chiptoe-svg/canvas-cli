package options

import "fmt"

// ObserveesListOptions contains options for listing observees
type ObserveesListOptions struct {
	UserID int64
}

// Validate validates the options
func (o *ObserveesListOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// ObserveesGetOptions contains options for getting an observee
type ObserveesGetOptions struct {
	UserID     int64
	ObserveeID int64
}

// Validate validates the options
func (o *ObserveesGetOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.ObserveeID <= 0 {
		return fmt.Errorf("observee-id is required and must be greater than 0")
	}
	return nil
}

// ObserveesDeleteOptions contains options for removing an observee
type ObserveesDeleteOptions struct {
	UserID     int64
	ObserveeID int64
}

// Validate validates the options
func (o *ObserveesDeleteOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.ObserveeID <= 0 {
		return fmt.Errorf("observee-id is required and must be greater than 0")
	}
	return nil
}

// ObserversListOptions contains options for listing observers
type ObserversListOptions struct {
	UserID int64
}

// Validate validates the options
func (o *ObserversListOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}
