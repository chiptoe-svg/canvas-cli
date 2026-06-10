package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/cache"
	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage Canvas CLI cache",
	Long: `Manage the Canvas CLI response cache.

The cache stores API responses to reduce load on the Canvas server and improve
response times for repeated requests.

Examples:
  canvas cache stats                    # Show cache statistics
  canvas cache clear                    # Clear expired cache entries
  canvas cache clear --all              # Clear all cache entries`,
}

func init() {
	rootCmd.AddCommand(cacheCmd)
	cacheCmd.AddCommand(newCacheStatsCmd())
	cacheCmd.AddCommand(newCacheClearCmd())
}

func newCacheStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show cache statistics",
		Long:  `Display statistics about the cache including size, entry counts, and hit rates.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCacheStats(cmd)
		},
	}
}

func newCacheClearCmd() *cobra.Command {
	opts := &options.CacheClearOptions{}

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear cache entries",
		Long: `Clear cached responses to free up disk space or force fresh data.

By default, only expired entries are cleared. Use --all to clear everything.

Examples:
  canvas cache clear          # Clear expired entries only
  canvas cache clear --all    # Clear all entries`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCacheClear(cmd, opts)
		},
	}

	cmd.Flags().BoolVar(&opts.All, "all", false, "Clear all cache entries (not just expired)")
	return cmd
}

func getCacheDir() (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}
	return filepath.Join(configDir, "cache"), nil
}

func runCacheStats(cmd *cobra.Command) error {
	logger := logging.NewCommandLogger(verbose)
	ctx := cmd.Context()
	logger.LogCommandStart(ctx, "cache.stats", nil)

	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	// Check if cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		fmt.Println("Cache is empty (no cache directory).")
		logger.LogCommandComplete(ctx, "cache.stats", 0)
		return nil
	}

	// Get disk cache stats
	diskCache, err := cache.NewDiskCache(cacheDir, 0)
	if err != nil {
		return fmt.Errorf("failed to open cache: %w", err)
	}

	stats, err := diskCache.Stats()
	if err != nil {
		return fmt.Errorf("failed to get cache stats: %w", err)
	}

	// Calculate cache size
	size, err := getCacheDirSize(cacheDir)
	if err != nil {
		size = 0
	}

	fmt.Println("Cache Statistics")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("\nLocation: %s\n", cacheDir)
	fmt.Printf("Size: %s\n", formatBytes(size))
	fmt.Println()
	fmt.Printf("Total entries:   %d\n", stats.Total)
	fmt.Printf("Active entries:  %d\n", stats.Active)
	fmt.Printf("Expired entries: %d\n", stats.Expired)

	if stats.Total > 0 {
		activePercent := float64(stats.Active) / float64(stats.Total) * 100
		fmt.Printf("\nActive rate: %.1f%%\n", activePercent)
	}

	logger.LogCommandComplete(ctx, "cache.stats", stats.Total)
	return nil
}

func runCacheClear(cmd *cobra.Command, opts *options.CacheClearOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctx := cmd.Context()
	logger.LogCommandStart(ctx, "cache.clear", map[string]interface{}{
		"all": opts.All,
	})

	cacheDir, err := getCacheDir()
	if err != nil {
		return err
	}

	// Check if cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		fmt.Println("Cache is already empty.")
		logger.LogCommandComplete(ctx, "cache.clear", 0)
		return nil
	}

	diskCache, err := cache.NewDiskCache(cacheDir, 0)
	if err != nil {
		return fmt.Errorf("failed to open cache: %w", err)
	}

	if opts.All {
		// Confirm clearing all
		fmt.Print("Are you sure you want to clear ALL cache entries? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			logger.LogCommandComplete(ctx, "cache.clear", 0)
			return nil
		}

		// Get count before clearing
		stats, _ := diskCache.Stats()
		totalBefore := stats.Total

		if err := diskCache.Clear(); err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}

		fmt.Printf("Cleared %d cache entries.\n", totalBefore)
		logger.LogCommandComplete(ctx, "cache.clear", totalBefore)
	} else {
		// Clear only expired entries
		stats, _ := diskCache.Stats()
		expiredBefore := stats.Expired

		if expiredBefore == 0 {
			fmt.Println("No expired entries to clear.")
			logger.LogCommandComplete(ctx, "cache.clear", 0)
			return nil
		}

		// Clear expired entries by reading each file
		if err := clearExpiredEntries(cacheDir); err != nil {
			return fmt.Errorf("failed to clear expired entries: %w", err)
		}

		fmt.Printf("Cleared %d expired cache entries.\n", expiredBefore)
		logger.LogCommandComplete(ctx, "cache.clear", expiredBefore)
	}

	return nil
}

func getCacheDirSize(dir string) (int64, error) {
	var size int64

	err := filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

func clearExpiredEntries(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		// Read and check expiration
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		if isExpiredCacheFile(data) {
			os.Remove(path)
		}
	}

	return nil
}

func isExpiredCacheFile(data []byte) bool {
	type cacheItem struct {
		Expiration time.Time `json:"expiration"`
	}

	var item cacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		return false
	}

	return time.Now().After(item.Expiration)
}
