package options

// EportfoliosListOptions holds options for listing ePortfolios.
type EportfoliosListOptions struct {
	UserID int64
}

// Validate validates the options.
func (o *EportfoliosListOptions) Validate() error {
	return ValidateRequired("user-id", o.UserID)
}

// EportfolioGetOptions holds options for getting an ePortfolio.
type EportfolioGetOptions struct {
	ID int64
}

// Validate validates the options.
func (o *EportfolioGetOptions) Validate() error {
	return ValidateRequired("id", o.ID)
}

// EportfolioDeleteOptions holds options for deleting an ePortfolio.
type EportfolioDeleteOptions struct {
	ID    int64
	Force bool
}

// Validate validates the options.
func (o *EportfolioDeleteOptions) Validate() error {
	return ValidateRequired("id", o.ID)
}

// EportfolioPagesListOptions holds options for listing ePortfolio pages.
type EportfolioPagesListOptions struct {
	EportfolioID int64
}

// Validate validates the options.
func (o *EportfolioPagesListOptions) Validate() error {
	return ValidateRequired("eportfolio-id", o.EportfolioID)
}
