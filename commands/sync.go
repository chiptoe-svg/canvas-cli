package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"github.com/jjuanrivvera/canvas-cli/internal/auth"
	"github.com/jjuanrivvera/canvas-cli/internal/batch"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize resources between Canvas instances",
	Long: `Synchronize courses, assignments, and other resources between different Canvas instances.

This is useful for:
- Migrating content between instances
- Backing up course content
- Copying course structures
- Synchronizing development and production environments`,
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.AddCommand(newSyncAssignmentsCmd())
	syncCmd.AddCommand(newSyncCourseCmd())
}

func newSyncAssignmentsCmd() *cobra.Command {
	opts := &options.SyncAssignmentsOptions{}

	cmd := &cobra.Command{
		Use:   "assignments <source-instance> <source-course-id> <target-instance> <target-course-id>",
		Short: "Sync assignments between instances",
		Long: `Synchronize all assignments from a source course to a target course.

The source and target can be on different Canvas instances.

Examples:
  # Sync assignments from production to staging
  canvas sync assignments prod 12345 staging 67890

  # Sync with interactive conflict resolution
  canvas sync assignments prod 12345 staging 67890 --interactive`,
		Args: ExactArgsWithUsage(4, "source-instance", "source-course-id", "target-instance", "target-course-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSyncAssignments(cmd, args, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Enable interactive conflict resolution")
	return cmd
}

func newSyncCourseCmd() *cobra.Command {
	opts := &options.SyncCourseOptions{}

	cmd := &cobra.Command{
		Use:   "course <source-instance> <source-course-id> <target-instance> <target-course-id>",
		Short: "Sync entire course between instances",
		Long: `Synchronize an entire course structure including assignments, files, and settings.

The source and target can be on different Canvas instances.

Examples:
  # Sync course from production to staging
  canvas sync course prod 12345 staging 67890

  # Sync with interactive conflict resolution
  canvas sync course prod 12345 staging 67890 --interactive`,
		Args: ExactArgsWithUsage(4, "source-instance", "source-course-id", "target-instance", "target-course-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSyncCourse(cmd, args, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Enable interactive conflict resolution")
	return cmd
}

func runSyncAssignments(cmd *cobra.Command, args []string, opts *options.SyncAssignmentsOptions) error {
	ctx := cmd.Context()
	logger := logging.NewCommandLogger(verbose)

	sourceInstance := args[0]
	sourceCourseID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid source course ID: %w", err)
	}

	targetInstance := args[2]
	targetCourseID, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid target course ID: %w", err)
	}

	logger.LogCommandStart(ctx, "sync.assignments", map[string]interface{}{
		"source_instance":  sourceInstance,
		"source_course_id": sourceCourseID,
		"target_instance":  targetInstance,
		"target_course_id": targetCourseID,
		"interactive":      opts.Interactive,
	})

	// Create API clients for both instances
	sourceClient, err := getAPIClientForInstance(sourceInstance)
	if err != nil {
		return fmt.Errorf("failed to create source client: %w", err)
	}

	targetClient, err := getAPIClientForInstance(targetInstance)
	if err != nil {
		return fmt.Errorf("failed to create target client: %w", err)
	}

	// Create sync operation
	syncOp := batch.NewSyncOperation(sourceClient, targetClient, opts.Interactive)

	fmt.Printf("Syncing assignments from %s (course %d) to %s (course %d)\n\n",
		sourceInstance, sourceCourseID, targetInstance, targetCourseID)

	// Perform sync
	result, err := syncOp.SyncAssignments(ctx, sourceCourseID, targetCourseID)
	if err != nil {
		logger.LogCommandError(ctx, "sync.assignments", err, nil)
		fmt.Printf("\nSync failed: %v\n", err)
		return err
	}

	// Display results
	fmt.Printf("\nSync complete!\n")
	fmt.Printf("Total assignments: %d\n", result.TotalItems)
	fmt.Printf("Synced: %d\n", result.SyncedItems)
	fmt.Printf("Skipped: %d\n", result.SkippedItems)
	fmt.Printf("Failed: %d\n", result.FailedItems)

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, syncErr := range result.Errors {
			fmt.Printf("  - %v\n", syncErr)
		}
	}

	logger.LogCommandComplete(ctx, "sync.assignments", result.SyncedItems)
	return nil
}

func runSyncCourse(cmd *cobra.Command, args []string, opts *options.SyncCourseOptions) error {
	ctx := cmd.Context()
	logger := logging.NewCommandLogger(verbose)

	sourceInstance := args[0]
	sourceCourseID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid source course ID: %w", err)
	}

	targetInstance := args[2]
	targetCourseID, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid target course ID: %w", err)
	}

	logger.LogCommandStart(ctx, "sync.course", map[string]interface{}{
		"source_instance":  sourceInstance,
		"source_course_id": sourceCourseID,
		"target_instance":  targetInstance,
		"target_course_id": targetCourseID,
		"interactive":      opts.Interactive,
	})

	// Create API clients for both instances
	sourceClient, err := getAPIClientForInstance(sourceInstance)
	if err != nil {
		return fmt.Errorf("failed to create source client: %w", err)
	}

	targetClient, err := getAPIClientForInstance(targetInstance)
	if err != nil {
		return fmt.Errorf("failed to create target client: %w", err)
	}

	// Create sync operation
	syncOp := batch.NewSyncOperation(sourceClient, targetClient, opts.Interactive)

	fmt.Printf("Syncing course from %s (course %d) to %s (course %d)\n\n",
		sourceInstance, sourceCourseID, targetInstance, targetCourseID)

	// Perform sync
	err = syncOp.CopyCourse(ctx, sourceCourseID, targetCourseID)
	if err != nil {
		logger.LogCommandError(ctx, "sync.course", err, nil)
		fmt.Printf("\nSync failed: %v\n", err)
		return err
	}

	fmt.Printf("\nCourse sync complete!\n")
	logger.LogCommandComplete(ctx, "sync.course", 1)
	return nil
}

// getAPIClientForInstance creates an API client for a specific instance name
func getAPIClientForInstance(instanceName string) (*api.Client, error) {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Get instance by name
	instance, err := cfg.GetInstance(instanceName)
	if err != nil {
		return nil, fmt.Errorf("instance not found: %w", err)
	}

	// Get config directory
	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	// Load token
	tokenStore := auth.NewFallbackTokenStore(configDir)
	token, err := tokenStore.Load(instance.Name)
	if err != nil {
		return nil, fmt.Errorf("not authenticated with %s. Run 'canvas auth login' first", instance.Name)
	}

	// Create auto-refreshing token source if we have OAuth credentials
	var clientConfig api.ClientConfig
	if instance.ClientID != "" && instance.ClientSecret != "" {
		// Create oauth2 config for token refresh
		oauth2Config := auth.CreateOAuth2ConfigForInstance(instance.URL, instance.ClientID, instance.ClientSecret)
		tokenSource := auth.NewAutoRefreshTokenSource(oauth2Config, tokenStore, instance.Name, token)

		clientConfig = api.ClientConfig{
			BaseURL:        instance.URL,
			TokenSource:    tokenSource,
			RequestsPerSec: cfg.Settings.RequestsPerSecond,
		}
	} else {
		// Fall back to static token (no auto-refresh)
		clientConfig = api.ClientConfig{
			BaseURL:        instance.URL,
			Token:          token.AccessToken,
			RequestsPerSec: cfg.Settings.RequestsPerSecond,
		}
	}

	// Create API client
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	return client, nil
}
