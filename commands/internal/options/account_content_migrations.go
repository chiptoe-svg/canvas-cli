package options

// AccountContentMigrationsListOptions contains options for listing account content migrations
type AccountContentMigrationsListOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AccountContentMigrationsListOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountContentMigrationsGetOptions contains options for getting an account content migration
type AccountContentMigrationsGetOptions struct {
	AccountID   int64
	MigrationID int64
}

// Validate validates the options
func (o *AccountContentMigrationsGetOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("migration-id", o.MigrationID)
}

// AccountContentMigrationsCreateOptions contains options for creating an account content migration
type AccountContentMigrationsCreateOptions struct {
	AccountID      int64
	MigrationType  string
	SourceCourseID int64
}

// Validate validates the options
func (o *AccountContentMigrationsCreateOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("type", o.MigrationType)
}

// AccountContentMigrationsMigratorsOptions contains options for listing account migrators
type AccountContentMigrationsMigratorsOptions struct {
	AccountID int64
}

// Validate validates the options
func (o *AccountContentMigrationsMigratorsOptions) Validate() error {
	return ValidateRequired("account-id", o.AccountID)
}

// AccountContentMigrationsIssuesOptions contains options for listing migration issues
type AccountContentMigrationsIssuesOptions struct {
	AccountID   int64
	MigrationID int64
}

// Validate validates the options
func (o *AccountContentMigrationsIssuesOptions) Validate() error {
	if err := ValidateRequired("account-id", o.AccountID); err != nil {
		return err
	}
	return ValidateRequired("migration-id", o.MigrationID)
}
