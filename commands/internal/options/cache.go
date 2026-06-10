package options

// CacheClearOptions contains options for the cache clear command
type CacheClearOptions struct {
	All bool
}

// Validate validates the options.
func (o *CacheClearOptions) Validate() error {
	return nil
}
