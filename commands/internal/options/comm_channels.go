package options

import "fmt"

// CommChannelsListOptions contains options for listing communication channels
type CommChannelsListOptions struct {
	UserID int64
}

// Validate validates the options
func (o *CommChannelsListOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	return nil
}

// CommChannelsCreateOptions contains options for creating a communication channel
type CommChannelsCreateOptions struct {
	UserID           int64
	Address          string
	Type             string
	SkipConfirmation bool
}

// Validate validates the options
func (o *CommChannelsCreateOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.Address == "" {
		return fmt.Errorf("address is required")
	}
	if o.Type == "" {
		return fmt.Errorf("type is required")
	}
	return nil
}

// CommChannelsDeleteOptions contains options for deleting a communication channel
type CommChannelsDeleteOptions struct {
	UserID    int64
	ChannelID int64
}

// Validate validates the options
func (o *CommChannelsDeleteOptions) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("user-id is required and must be greater than 0")
	}
	if o.ChannelID <= 0 {
		return fmt.Errorf("channel-id is required and must be greater than 0")
	}
	return nil
}
