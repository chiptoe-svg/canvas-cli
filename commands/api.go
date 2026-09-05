package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

func init() {
	rootCmd.AddCommand(newAPICmd())
}

func newAPICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Read any Canvas endpoint this CLI has no command for yet",
		Long: `Read-only escape hatch: GET any path under /api/v1/ and print the
JSON. Use it to look at something no command covers; if you need it more
than once, ask for a command. This CLI cannot write through this path.

Examples:
  canvas api get /api/v1/courses/123/settings
  canvas api get "/api/v1/courses/123/assignments?bucket=upcoming"`,
	}
	cmd.AddCommand(newAPIGetCmd())
	return cmd
}

// newAPIGetCmd is the GET-only "canvas api" subcommand. Because it can never
// mutate state, it is safe to advertise to read-only MCP clients: the shared
// classifier (classifyCanvasCommand) buckets "api get" as a read, so it carries
// readOnlyHint=true. It gives broad Canvas read coverage from a single tool
// schema instead of allowlisting every typed read tool. See issue #60.
func newAPIGetCmd() *cobra.Command {
	opts := &options.APIOptions{}

	cmd := &cobra.Command{
		Use:   "get <PATH>",
		Short: "Make a read-only GET request to Canvas",
		Long: `Make a raw GET request to any Canvas API endpoint.

Examples:
  # List all courses
  canvas api get /api/v1/courses

  # Search users with query parameters
  canvas api get /api/v1/users -q "search_term=john" -q "per_page=50"

  # Follow pagination
  canvas api get /api/v1/courses --paginate`,
		Args: ExactArgsWithUsage(1, "path"),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runAPICommand(cmd, args[0], client, opts)
		},
	}

	cmd.Flags().StringArrayVarP(&opts.Query, "query", "q", nil, "Query parameters (key=value, repeatable)")
	cmd.Flags().StringArrayVarP(&opts.Headers, "header", "H", nil, "Custom headers (key:value, repeatable)")
	cmd.Flags().BoolVar(&opts.Paginate, "paginate", false, "Follow pagination links")
	cmd.Flags().BoolVar(&opts.RawOutput, "raw", false, "Output raw response without formatting")
	cmd.Flags().BoolVar(&opts.ShowHeaders, "show-headers", false, "Include response headers in output")

	return cmd
}

// runAPICommand performs a read-only GET request against path. It is the
// shared runner for "canvas api get"; the CLI has no write-capable path.
func runAPICommand(cmd *cobra.Command, path string, client *api.Client, opts *options.APIOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctx := cmd.Context()

	logger.LogCommandStart(ctx, "api.request", map[string]interface{}{
		"method": http.MethodGet,
		"path":   path,
	})

	service := api.NewRawService(client)

	// Build request options
	reqOpts := &api.RawRequestOptions{
		Paginate: opts.Paginate,
	}

	// Parse query parameters.
	// Guard with Changed() to avoid accumulating values from previous Execute() calls
	// when the cobra command is reused (e.g., in tests).
	if cmd.Flags().Changed("query") {
		query := make(map[string][]string)
		for _, q := range opts.Query {
			parts := strings.SplitN(q, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid query parameter format: %s (use key=value)", q)
			}
			key := parts[0]
			value := parts[1]
			query[key] = append(query[key], value)
		}
		reqOpts.Query = query
	}

	// Parse custom headers.
	// Same Changed() guard as for query params.
	if cmd.Flags().Changed("header") {
		headers := make(map[string]string)
		for _, h := range opts.Headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid header format: %s (use key:value)", h)
			}
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
		reqOpts.Headers = headers
	}

	// Make the request
	resp, err := service.Request(ctx, http.MethodGet, path, reqOpts)
	if err != nil {
		logger.LogCommandError(ctx, "api.request", err, map[string]interface{}{
			"method": http.MethodGet,
			"path":   path,
		})
		return err
	}

	logger.LogCommandComplete(ctx, "api.request", 1)
	return outputAPIResponse(cmd, resp, opts)
}

func outputAPIResponse(cmd *cobra.Command, resp *api.RawResponse, opts *options.APIOptions) error {
	// If raw output, just print the body
	if opts.RawOutput {
		cmd.Println(string(resp.Body))
		return nil
	}

	// Build output structure
	out := make(map[string]interface{})
	out["status_code"] = resp.StatusCode

	if opts.ShowHeaders {
		headers := make(map[string]string)
		for key, values := range resp.Headers {
			if len(values) > 0 {
				headers[key] = values[0]
			}
		}
		out["headers"] = headers
	}

	// Parse body as JSON if possible
	if len(resp.Body) > 0 {
		var body interface{}
		if err := json.Unmarshal(resp.Body, &body); err == nil {
			out["body"] = body
		} else {
			out["body"] = string(resp.Body)
		}
	}

	// Add pagination info if available
	if resp.Pagination != nil && resp.Pagination.HasNextPage() {
		out["pagination"] = map[string]interface{}{
			"has_next": resp.Pagination.HasNextPage(),
			"next":     resp.Pagination.Next,
		}
	}

	// Format output based on output format flag
	return formatOutput(out, nil)
}
