package options

// JWTCreateOptions holds options for creating a Canvas JWT.
type JWTCreateOptions struct {
	Workflow string
}

// Validate validates the options.
func (o *JWTCreateOptions) Validate() error {
	return nil
}

// JWTRefreshOptions holds options for refreshing a Canvas JWT.
type JWTRefreshOptions struct {
	Token string
}

// Validate validates the options.
func (o *JWTRefreshOptions) Validate() error {
	return ValidateRequired("token", o.Token)
}
