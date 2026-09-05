package commands

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Canvas CLI configuration",
	Long: `Manage Canvas CLI configuration including Canvas instances and settings.

Examples:
  canvas config list                              # List all configured instances
  canvas config add prod --url https://canvas.example.com
  canvas config use prod                          # Switch to prod instance
  canvas config show                              # Show current configuration
  canvas config remove staging                    # Remove an instance`,
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(newConfigListCmd())
	configCmd.AddCommand(newConfigAddCmd())
	configCmd.AddCommand(newConfigUseCmd())
	configCmd.AddCommand(newConfigRemoveCmd())
	configCmd.AddCommand(newConfigShowCmd())
}

func newConfigListCmd() *cobra.Command {
	opts := &options.ConfigListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configured Canvas instances",
		Long:  `List all configured Canvas instances with their URLs and status.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runConfigList(cmd.Context(), opts)
		},
	}

	return cmd
}

func newConfigAddCmd() *cobra.Command {
	opts := &options.ConfigAddOptions{}

	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add a new Canvas instance",
		Long: `Add a new Canvas instance to the configuration.

Examples:
  canvas config add production --url https://canvas.example.com
  canvas config add staging --url https://canvas-staging.example.com --description "Staging environment"`,
		Args: ExactArgsWithUsage(1, "name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			return runConfigAdd(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.URL, "url", "", "Canvas instance URL (required)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Instance description")
	cmd.Flags().StringVar(&opts.ClientID, "client-id", "", "OAuth client ID")
	mustMarkRequired(cmd, "url")

	return cmd
}

func newConfigUseCmd() *cobra.Command {
	opts := &options.ConfigUseOptions{}

	cmd := &cobra.Command{
		Use:   "use [instance-name]",
		Short: "Switch to a different Canvas instance",
		Long: `Set the default Canvas instance to use for all commands.

Examples:
  canvas config use production
  canvas config use staging`,
		Args: ExactArgsWithUsage(1, "instance-name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.InstanceName = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			return runConfigUse(cmd.Context(), opts)
		},
	}

	return cmd
}

func newConfigRemoveCmd() *cobra.Command {
	opts := &options.ConfigRemoveOptions{}

	cmd := &cobra.Command{
		Use:   "remove [instance-name]",
		Short: "Remove a Canvas instance from configuration",
		Long: `Remove a configured Canvas instance. This will not delete any data from Canvas,
only remove the instance from your local configuration.

Examples:
  canvas config remove staging`,
		Args: ExactArgsWithUsage(1, "instance-name"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.InstanceName = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			return runConfigRemove(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func newConfigShowCmd() *cobra.Command {
	opts := &options.ConfigShowOptions{}

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration details",
		Long:  `Display the current configuration including default instance and settings.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return runConfigShow(cmd.Context(), opts)
		},
	}

	return cmd
}

func runConfigList(ctx context.Context, opts *options.ConfigListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "config.list", map[string]interface{}{})

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "config.list", err, map[string]interface{}{})
		return fmt.Errorf("failed to load config: %w", err)
	}

	instances := cfg.ListInstances()
	if len(instances) == 0 {
		logger.LogCommandComplete(ctx, "config.list", 0)
		fmt.Println("No instances configured.")
		fmt.Println("\nTo add an instance:")
		fmt.Println("  canvas config add <name> --url <canvas-url>")
		return nil
	}

	// Print header
	fmt.Printf("%-15s %-45s %-10s %s\n", "NAME", "URL", "DEFAULT", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 90))

	// Print instances
	for _, instance := range instances {
		isDefault := ""
		if instance.Name == cfg.DefaultInstance {
			isDefault = "*"
		}

		description := instance.Description
		if len(description) > 20 {
			description = description[:17] + "..."
		}

		fmt.Printf("%-15s %-45s %-10s %s\n",
			instance.Name,
			instance.URL,
			isDefault,
			description,
		)
	}

	fmt.Printf("\nTotal: %d instance(s)\n", len(instances))
	logger.LogCommandComplete(ctx, "config.list", len(instances))
	return nil
}

func runConfigAdd(ctx context.Context, opts *options.ConfigAddOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "config.add", map[string]interface{}{
		"instance_name": opts.Name,
		"url":           opts.URL,
	})

	// Validate URL
	parsedURL, err := url.Parse(opts.URL)
	if err != nil {
		logger.LogCommandError(ctx, "config.add", err, map[string]interface{}{
			"url": opts.URL,
		})
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme == "" {
		opts.URL = "https://" + opts.URL
		parsedURL, _ = url.Parse(opts.URL)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		err := fmt.Errorf("URL must use http or https scheme")
		logger.LogCommandError(ctx, "config.add", err, map[string]interface{}{
			"url":    opts.URL,
			"scheme": parsedURL.Scheme,
		})
		return err
	}

	// Remove trailing slash
	opts.URL = strings.TrimSuffix(opts.URL, "/")

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "config.add", err, map[string]interface{}{})
		return fmt.Errorf("failed to load config: %w", err)
	}

	instance := &config.Instance{
		Name:        opts.Name,
		URL:         opts.URL,
		Description: opts.Description,
		ClientID:    opts.ClientID,
	}

	if err := cfg.AddInstance(instance); err != nil {
		logger.LogCommandError(ctx, "config.add", err, map[string]interface{}{
			"instance_name": opts.Name,
		})
		return fmt.Errorf("failed to add instance: %w", err)
	}

	printInfo("Instance %q added successfully.\n", opts.Name)

	if cfg.DefaultInstance == opts.Name {
		printInfo("Set as default instance.\n")
	}

	printInfoln("\nNext step: Authenticate with this instance:")
	printInfo("  canvas auth login --instance %s\n", opts.Name)

	logger.LogCommandComplete(ctx, "config.add", 1)
	return nil
}

func runConfigUse(ctx context.Context, opts *options.ConfigUseOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "config.use", map[string]interface{}{
		"instance_name": opts.InstanceName,
	})

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "config.use", err, map[string]interface{}{})
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.SetDefaultInstance(opts.InstanceName); err != nil {
		logger.LogCommandError(ctx, "config.use", err, map[string]interface{}{
			"instance_name": opts.InstanceName,
		})
		return fmt.Errorf("failed to set default instance: %w", err)
	}

	fmt.Printf("Switched to instance %q.\n", opts.InstanceName)
	logger.LogCommandComplete(ctx, "config.use", 1)
	return nil
}

func runConfigRemove(ctx context.Context, opts *options.ConfigRemoveOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "config.remove", map[string]interface{}{
		"instance_name": opts.InstanceName,
		"force":         opts.Force,
	})

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "config.remove", err, map[string]interface{}{})
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if instance exists
	if _, err := cfg.GetInstance(opts.InstanceName); err != nil {
		logger.LogCommandError(ctx, "config.remove", err, map[string]interface{}{
			"instance_name": opts.InstanceName,
		})
		return fmt.Errorf("instance %q not found", opts.InstanceName)
	}

	// Confirm removal unless --force is used
	if !opts.Force {
		fmt.Printf("Are you sure you want to remove instance %q? [y/N]: ", opts.InstanceName)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			logger.LogCommandError(ctx, "config.remove", err, map[string]interface{}{})
			return fmt.Errorf("failed to read response: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			logger.LogCommandComplete(ctx, "config.remove", 0)
			fmt.Println("Cancelled.")
			return nil
		}
	}

	if err := cfg.RemoveInstance(opts.InstanceName); err != nil {
		logger.LogCommandError(ctx, "config.remove", err, map[string]interface{}{
			"instance_name": opts.InstanceName,
		})
		return fmt.Errorf("failed to remove instance: %w", err)
	}

	fmt.Printf("Instance %q removed.\n", opts.InstanceName)

	if cfg.DefaultInstance != "" {
		fmt.Printf("New default instance: %s\n", cfg.DefaultInstance)
	}

	logger.LogCommandComplete(ctx, "config.remove", 1)
	return nil
}

func runConfigShow(ctx context.Context, opts *options.ConfigShowOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "config.show", map[string]interface{}{})

	cfg, err := config.Load()
	if err != nil {
		logger.LogCommandError(ctx, "config.show", err, map[string]interface{}{})
		return fmt.Errorf("failed to load config: %w", err)
	}

	configPath, _ := config.GetConfigPath()

	fmt.Println("Canvas CLI Configuration")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Printf("\nConfig file: %s\n", configPath)
	fmt.Printf("Default instance: %s\n", valueOrNone(cfg.DefaultInstance))

	fmt.Println("\nInstances:")
	if len(cfg.Instances) == 0 {
		fmt.Println("  (none)")
	} else {
		for name, instance := range cfg.Instances {
			marker := "  "
			if name == cfg.DefaultInstance {
				marker = "* "
			}
			fmt.Printf("%s%s: %s\n", marker, name, instance.URL)
			if instance.Description != "" {
				fmt.Printf("    Description: %s\n", instance.Description)
			}
			if instance.ClientID != "" {
				fmt.Printf("    Client ID: %s\n", instance.ClientID)
			}
			if instance.HasDefaultAccountID() {
				fmt.Printf("    Default Account ID: %d\n", instance.DefaultAccountID)
			}
		}
	}

	fmt.Println("\nSettings:")
	if cfg.Settings != nil {
		fmt.Printf("  Output format: %s\n", cfg.Settings.DefaultOutputFormat)
		fmt.Printf("  Requests/sec: %.1f\n", cfg.Settings.RequestsPerSecond)
		fmt.Printf("  Cache enabled: %t\n", cfg.Settings.CacheEnabled)
		fmt.Printf("  Cache TTL: %d minutes\n", cfg.Settings.CacheTTL)
		fmt.Printf("  Telemetry: %t\n", cfg.Settings.TelemetryEnabled)
		fmt.Printf("  Log level: %s\n", cfg.Settings.LogLevel)
		fmt.Printf("  Auto-update: %t\n", cfg.Settings.AutoUpdateEnabled)
		fmt.Printf("  Auto-update interval: %d minutes\n", cfg.Settings.AutoUpdateIntervalMin)
		if cfg.Settings.Timezone != "" {
			fmt.Printf("  Timezone: %s\n", cfg.Settings.Timezone)
		}
	}

	logger.LogCommandComplete(ctx, "config.show", 1)
	return nil
}

func valueOrNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}
