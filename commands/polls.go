package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// pollsCmd is the root command group for polls
var pollsCmd = &cobra.Command{
	Use:   "polls",
	Short: "Manage Canvas polls",
	Long: `Manage Canvas polls, poll choices, poll sessions, and poll submissions.

Polls allow instructors to gather feedback and assess understanding through
interactive question-and-answer activities.

Examples:
  canvas polls list
  canvas polls get 1023
  canvas polls create --question "What is your favourite language?"
  canvas polls choices list --poll-id 1023
  canvas polls sessions list --poll-id 1023
  canvas polls sessions open --poll-id 1023 --session-id 55`,
}

// pollChoicesCmd groups poll-choice subcommands
var pollChoicesCmd = &cobra.Command{
	Use:   "choices",
	Short: "Manage poll choices",
	Long:  `Manage poll choices (the individual answer options within a poll).`,
}

// pollSessionsCmd groups poll-session subcommands
var pollSessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Manage poll sessions",
	Long:  `Manage poll sessions (instances of a poll opened for a course).`,
}

// pollSubmissionsCmd groups poll-submission subcommands
var pollSubmissionsCmd = &cobra.Command{
	Use:   "submissions",
	Short: "Manage poll submissions",
	Long:  `Manage poll submissions (student votes within a poll session).`,
}

func init() {
	rootCmd.AddCommand(pollsCmd)

	// Top-level poll commands
	pollsCmd.AddCommand(newPollListCmd())
	pollsCmd.AddCommand(newPollGetCmd())
	pollsCmd.AddCommand(newPollCreateCmd())
	pollsCmd.AddCommand(newPollUpdateCmd())
	pollsCmd.AddCommand(newPollDeleteCmd())

	// choices subgroup
	pollsCmd.AddCommand(pollChoicesCmd)
	pollChoicesCmd.AddCommand(newPollChoiceListCmd())
	pollChoicesCmd.AddCommand(newPollChoiceGetCmd())
	pollChoicesCmd.AddCommand(newPollChoiceCreateCmd())
	pollChoicesCmd.AddCommand(newPollChoiceUpdateCmd())
	pollChoicesCmd.AddCommand(newPollChoiceDeleteCmd())

	// sessions subgroup
	pollsCmd.AddCommand(pollSessionsCmd)
	pollSessionsCmd.AddCommand(newPollSessionListCmd())
	pollSessionsCmd.AddCommand(newPollSessionGetCmd())
	pollSessionsCmd.AddCommand(newPollSessionCreateCmd())
	pollSessionsCmd.AddCommand(newPollSessionUpdateCmd())
	pollSessionsCmd.AddCommand(newPollSessionDeleteCmd())
	pollSessionsCmd.AddCommand(newPollSessionOpenCmd())
	pollSessionsCmd.AddCommand(newPollSessionCloseCmd())
	pollSessionsCmd.AddCommand(newPollSessionListOpenedCmd())
	pollSessionsCmd.AddCommand(newPollSessionListClosedCmd())

	// submissions subgroup
	pollsCmd.AddCommand(pollSubmissionsCmd)
	pollSubmissionsCmd.AddCommand(newPollSubmissionGetCmd())
	pollSubmissionsCmd.AddCommand(newPollSubmissionCreateCmd())
}

// ---- Poll commands ----

func newPollListCmd() *cobra.Command {
	opts := &options.PollListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List polls",
		Long: `List polls created by the current user.

Examples:
  canvas polls list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollList(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newPollGetCmd() *cobra.Command {
	opts := &options.PollGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <poll-id>",
		Short: "Get a poll",
		Long: `Get details of a specific poll.

Examples:
  canvas polls get 1023`,
		Args: ExactArgsWithUsage(1, "poll-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid poll ID: %s", args[0])
			}
			opts.PollID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newPollCreateCmd() *cobra.Command {
	opts := &options.PollCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a poll",
		Long: `Create a new poll.

Examples:
  canvas polls create --question "Which language do you prefer?"
  canvas polls create --question "Rate today's lecture" --description "Quick feedback"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Question, "question", "", "Poll question (required)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Poll description")
	mustMarkRequired(cmd, "question")

	return cmd
}

func newPollUpdateCmd() *cobra.Command {
	opts := &options.PollUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <poll-id>",
		Short: "Update a poll",
		Long: `Update an existing poll.

Examples:
  canvas polls update 1023 --question "Updated question"`,
		Args: ExactArgsWithUsage(1, "poll-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid poll ID: %s", args[0])
			}
			opts.PollID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Question, "question", "", "Updated question")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Updated description")

	return cmd
}

func newPollDeleteCmd() *cobra.Command {
	opts := &options.PollDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <poll-id>",
		Short: "Delete a poll",
		Long: `Delete a poll and all associated data.

Examples:
  canvas polls delete 1023
  canvas polls delete 1023 --force`,
		Args: ExactArgsWithUsage(1, "poll-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid poll ID: %s", args[0])
			}
			opts.PollID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

// ---- Poll choice commands ----

func newPollChoiceListCmd() *cobra.Command {
	opts := &options.PollChoiceListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List poll choices",
		Long: `List all choices for a poll.

Examples:
  canvas polls choices list --poll-id 1023`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollChoiceList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	mustMarkRequired(cmd, "poll-id")

	return cmd
}

func newPollChoiceGetCmd() *cobra.Command {
	opts := &options.PollChoiceGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <choice-id>",
		Short: "Get a poll choice",
		Long: `Get details of a specific poll choice.

Examples:
  canvas polls choices get 55 --poll-id 1023`,
		Args: ExactArgsWithUsage(1, "choice-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid choice ID: %s", args[0])
			}
			opts.ChoiceID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollChoiceGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	mustMarkRequired(cmd, "poll-id")

	return cmd
}

func newPollChoiceCreateCmd() *cobra.Command {
	opts := &options.PollChoiceCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a poll choice",
		Long: `Create a new choice for a poll.

Examples:
  canvas polls choices create --poll-id 1023 --text "Option A" --correct
  canvas polls choices create --poll-id 1023 --text "Option B" --position 2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollChoiceCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	cmd.Flags().StringVar(&opts.Text, "text", "", "Choice text (required)")
	cmd.Flags().BoolVar(&opts.IsCorrect, "correct", false, "Mark as correct answer")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Display position")
	mustMarkRequired(cmd, "poll-id", "text")

	return cmd
}

func newPollChoiceUpdateCmd() *cobra.Command {
	opts := &options.PollChoiceUpdateOptions{}
	var isCorrectFlag bool
	var isCorrectSet bool

	cmd := &cobra.Command{
		Use:   "update <choice-id>",
		Short: "Update a poll choice",
		Long: `Update an existing poll choice.

Examples:
  canvas polls choices update 55 --poll-id 1023 --text "Updated text"`,
		Args: ExactArgsWithUsage(1, "choice-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid choice ID: %s", args[0])
			}
			opts.ChoiceID = id

			if isCorrectSet {
				opts.IsCorrect = &isCorrectFlag
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollChoiceUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	cmd.Flags().StringVar(&opts.Text, "text", "", "Updated text")
	cmd.Flags().BoolVar(&isCorrectFlag, "correct", false, "Mark as correct answer")
	cmd.Flags().IntVar(&opts.Position, "position", 0, "Updated display position")
	mustMarkRequired(cmd, "poll-id")

	// Track whether --correct was explicitly set
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		isCorrectSet = cmd.Flags().Changed("correct")
		return nil
	}

	return cmd
}

func newPollChoiceDeleteCmd() *cobra.Command {
	opts := &options.PollChoiceDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <choice-id>",
		Short: "Delete a poll choice",
		Long: `Delete a poll choice.

Examples:
  canvas polls choices delete 55 --poll-id 1023 --force`,
		Args: ExactArgsWithUsage(1, "choice-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid choice ID: %s", args[0])
			}
			opts.ChoiceID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollChoiceDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")
	mustMarkRequired(cmd, "poll-id")

	return cmd
}

// ---- Poll session commands ----

func newPollSessionListCmd() *cobra.Command {
	opts := &options.PollSessionListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List poll sessions",
		Long: `List sessions for a poll.

Examples:
  canvas polls sessions list --poll-id 1023`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSessionList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	mustMarkRequired(cmd, "poll-id")

	return cmd
}

func newPollSessionGetCmd() *cobra.Command {
	opts := &options.PollSessionGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <session-id>",
		Short: "Get a poll session",
		Long: `Get details of a specific poll session.

Examples:
  canvas polls sessions get 77 --poll-id 1023`,
		Args: ExactArgsWithUsage(1, "session-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid session ID: %s", args[0])
			}
			opts.SessionID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSessionGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	mustMarkRequired(cmd, "poll-id")

	return cmd
}

func newPollSessionCreateCmd() *cobra.Command {
	opts := &options.PollSessionCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a poll session",
		Long: `Create a new poll session for a course.

Examples:
  canvas polls sessions create --poll-id 1023 --course-id 111
  canvas polls sessions create --poll-id 1023 --course-id 111 --public-results`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSessionCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64Var(&opts.CourseSectionID, "section-id", 0, "Course section ID")
	cmd.Flags().BoolVar(&opts.HasPublicResults, "public-results", false, "Make results publicly visible")
	mustMarkRequired(cmd, "poll-id", "course-id")

	return cmd
}

func newPollSessionUpdateCmd() *cobra.Command {
	opts := &options.PollSessionUpdateOptions{}
	var publicResults bool
	var publicResultsSet bool

	cmd := &cobra.Command{
		Use:   "update <session-id>",
		Short: "Update a poll session",
		Long: `Update an existing poll session.

Examples:
  canvas polls sessions update 77 --poll-id 1023 --public-results`,
		Args: ExactArgsWithUsage(1, "session-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid session ID: %s", args[0])
			}
			opts.SessionID = id

			if publicResultsSet {
				opts.HasPublicResults = &publicResults
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSessionUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.CourseSectionID, "section-id", 0, "Course section ID")
	cmd.Flags().BoolVar(&publicResults, "public-results", false, "Make results publicly visible")
	mustMarkRequired(cmd, "poll-id")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		publicResultsSet = cmd.Flags().Changed("public-results")
		return nil
	}

	return cmd
}

func newPollSessionDeleteCmd() *cobra.Command {
	opts := &options.PollSessionDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <session-id>",
		Short: "Delete a poll session",
		Long: `Delete a poll session.

Examples:
  canvas polls sessions delete 77 --poll-id 1023 --force`,
		Args: ExactArgsWithUsage(1, "session-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid session ID: %s", args[0])
			}
			opts.SessionID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSessionDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")
	mustMarkRequired(cmd, "poll-id")

	return cmd
}

func newPollSessionOpenCmd() *cobra.Command {
	opts := &options.PollSessionOpenOptions{}

	cmd := &cobra.Command{
		Use:   "open <session-id>",
		Short: "Open a poll session",
		Long: `Open a poll session for student participation.

Examples:
  canvas polls sessions open 77 --poll-id 1023`,
		Args: ExactArgsWithUsage(1, "session-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid session ID: %s", args[0])
			}
			opts.SessionID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSessionOpen(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	mustMarkRequired(cmd, "poll-id")

	return cmd
}

func newPollSessionCloseCmd() *cobra.Command {
	opts := &options.PollSessionCloseOptions{}

	cmd := &cobra.Command{
		Use:   "close <session-id>",
		Short: "Close a poll session",
		Long: `Close a poll session to stop accepting submissions.

Examples:
  canvas polls sessions close 77 --poll-id 1023`,
		Args: ExactArgsWithUsage(1, "session-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid session ID: %s", args[0])
			}
			opts.SessionID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSessionClose(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	mustMarkRequired(cmd, "poll-id")

	return cmd
}

func newPollSessionListOpenedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-opened",
		Short: "List all open poll sessions",
		Long: `List all currently open poll sessions across all polls.

Examples:
  canvas polls sessions list-opened`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSessionListOpened(cmd.Context(), client)
		},
	}

	return cmd
}

func newPollSessionListClosedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-closed",
		Short: "List all closed poll sessions",
		Long: `List all closed poll sessions across all polls.

Examples:
  canvas polls sessions list-closed`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSessionListClosed(cmd.Context(), client)
		},
	}

	return cmd
}

// ---- Poll submission commands ----

func newPollSubmissionGetCmd() *cobra.Command {
	opts := &options.PollSubmissionGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <submission-id>",
		Short: "Get a poll submission",
		Long: `Get details of a specific poll submission.

Examples:
  canvas polls submissions get 200 --poll-id 1023 --session-id 77`,
		Args: ExactArgsWithUsage(1, "submission-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid submission ID: %s", args[0])
			}
			opts.SubmissionID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSubmissionGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	cmd.Flags().Int64Var(&opts.SessionID, "session-id", 0, "Session ID (required)")
	mustMarkRequired(cmd, "poll-id", "session-id")

	return cmd
}

func newPollSubmissionCreateCmd() *cobra.Command {
	opts := &options.PollSubmissionCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Submit a vote",
		Long: `Submit a vote for an active poll session.

Examples:
  canvas polls submissions create --poll-id 1023 --session-id 77 --choice-id 55`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runPollSubmissionCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.PollID, "poll-id", 0, "Poll ID (required)")
	cmd.Flags().Int64Var(&opts.SessionID, "session-id", 0, "Session ID (required)")
	cmd.Flags().Int64Var(&opts.PollChoiceID, "choice-id", 0, "Poll choice ID (required)")
	mustMarkRequired(cmd, "poll-id", "session-id", "choice-id")

	return cmd
}

// ---- Run functions ----

func runPollList(ctx context.Context, client *api.Client, opts *options.PollListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.list", map[string]interface{}{})

	svc := api.NewPollsService(client)

	polls, err := svc.ListPolls(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "polls.list", err, nil)
		return fmt.Errorf("failed to list polls: %w", err)
	}

	printVerbose("Found %d polls:\n\n", len(polls))
	logger.LogCommandComplete(ctx, "polls.list", len(polls))
	return formatEmptyOrOutput(polls, "No polls found")
}

func runPollGet(ctx context.Context, client *api.Client, opts *options.PollGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.get", map[string]interface{}{"poll_id": opts.PollID})

	svc := api.NewPollsService(client)

	poll, err := svc.GetPoll(ctx, opts.PollID)
	if err != nil {
		logger.LogCommandError(ctx, "polls.get", err, map[string]interface{}{"poll_id": opts.PollID})
		return fmt.Errorf("failed to get poll: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.get", 1)
	return formatOutput(poll, nil)
}

func runPollCreate(ctx context.Context, client *api.Client, opts *options.PollCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.create", map[string]interface{}{"question": opts.Question})

	svc := api.NewPollsService(client)

	poll, err := svc.CreatePoll(ctx, &api.CreatePollParams{
		Question:    opts.Question,
		Description: opts.Description,
	})
	if err != nil {
		logger.LogCommandError(ctx, "polls.create", err, nil)
		return fmt.Errorf("failed to create poll: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.create", 1)
	return formatSuccessOutput(poll, "Poll created successfully!")
}

func runPollUpdate(ctx context.Context, client *api.Client, opts *options.PollUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.update", map[string]interface{}{"poll_id": opts.PollID})

	svc := api.NewPollsService(client)

	poll, err := svc.UpdatePoll(ctx, opts.PollID, &api.UpdatePollParams{
		Question:    opts.Question,
		Description: opts.Description,
	})
	if err != nil {
		logger.LogCommandError(ctx, "polls.update", err, map[string]interface{}{"poll_id": opts.PollID})
		return fmt.Errorf("failed to update poll: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.update", 1)
	return formatSuccessOutput(poll, "Poll updated successfully!")
}

func runPollDelete(ctx context.Context, client *api.Client, opts *options.PollDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.delete", map[string]interface{}{"poll_id": opts.PollID, "force": opts.Force})

	confirmed, err := confirmDelete("poll", opts.PollID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "polls.delete", err, map[string]interface{}{"poll_id": opts.PollID})
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		logger.LogCommandComplete(ctx, "polls.delete", 0)
		return nil
	}

	svc := api.NewPollsService(client)

	if err := svc.DeletePoll(ctx, opts.PollID); err != nil {
		logger.LogCommandError(ctx, "polls.delete", err, map[string]interface{}{"poll_id": opts.PollID})
		return fmt.Errorf("failed to delete poll: %w", err)
	}

	printInfo("Poll %d deleted successfully\n", opts.PollID)
	logger.LogCommandComplete(ctx, "polls.delete", 1)
	return nil
}

func runPollChoiceList(ctx context.Context, client *api.Client, opts *options.PollChoiceListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.choices.list", map[string]interface{}{"poll_id": opts.PollID})

	svc := api.NewPollsService(client)

	choices, err := svc.ListPollChoices(ctx, opts.PollID)
	if err != nil {
		logger.LogCommandError(ctx, "polls.choices.list", err, map[string]interface{}{"poll_id": opts.PollID})
		return fmt.Errorf("failed to list poll choices: %w", err)
	}

	printVerbose("Found %d poll choices:\n\n", len(choices))
	logger.LogCommandComplete(ctx, "polls.choices.list", len(choices))
	return formatEmptyOrOutput(choices, "No poll choices found")
}

func runPollChoiceGet(ctx context.Context, client *api.Client, opts *options.PollChoiceGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.choices.get", map[string]interface{}{
		"poll_id":   opts.PollID,
		"choice_id": opts.ChoiceID,
	})

	svc := api.NewPollsService(client)

	choice, err := svc.GetPollChoice(ctx, opts.PollID, opts.ChoiceID)
	if err != nil {
		logger.LogCommandError(ctx, "polls.choices.get", err, nil)
		return fmt.Errorf("failed to get poll choice: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.choices.get", 1)
	return formatOutput(choice, nil)
}

func runPollChoiceCreate(ctx context.Context, client *api.Client, opts *options.PollChoiceCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.choices.create", map[string]interface{}{
		"poll_id": opts.PollID,
		"text":    opts.Text,
	})

	svc := api.NewPollsService(client)

	choice, err := svc.CreatePollChoice(ctx, opts.PollID, &api.CreatePollChoiceParams{
		Text:      opts.Text,
		IsCorrect: opts.IsCorrect,
		Position:  opts.Position,
	})
	if err != nil {
		logger.LogCommandError(ctx, "polls.choices.create", err, nil)
		return fmt.Errorf("failed to create poll choice: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.choices.create", 1)
	return formatSuccessOutput(choice, "Poll choice created successfully!")
}

func runPollChoiceUpdate(ctx context.Context, client *api.Client, opts *options.PollChoiceUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.choices.update", map[string]interface{}{
		"poll_id":   opts.PollID,
		"choice_id": opts.ChoiceID,
	})

	svc := api.NewPollsService(client)

	choice, err := svc.UpdatePollChoice(ctx, opts.PollID, opts.ChoiceID, &api.UpdatePollChoiceParams{
		Text:      opts.Text,
		IsCorrect: opts.IsCorrect,
		Position:  opts.Position,
	})
	if err != nil {
		logger.LogCommandError(ctx, "polls.choices.update", err, nil)
		return fmt.Errorf("failed to update poll choice: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.choices.update", 1)
	return formatSuccessOutput(choice, "Poll choice updated successfully!")
}

func runPollChoiceDelete(ctx context.Context, client *api.Client, opts *options.PollChoiceDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.choices.delete", map[string]interface{}{
		"poll_id":   opts.PollID,
		"choice_id": opts.ChoiceID,
		"force":     opts.Force,
	})

	confirmed, err := confirmDelete("poll choice", opts.ChoiceID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "polls.choices.delete", err, nil)
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		logger.LogCommandComplete(ctx, "polls.choices.delete", 0)
		return nil
	}

	svc := api.NewPollsService(client)

	if err := svc.DeletePollChoice(ctx, opts.PollID, opts.ChoiceID); err != nil {
		logger.LogCommandError(ctx, "polls.choices.delete", err, nil)
		return fmt.Errorf("failed to delete poll choice: %w", err)
	}

	printInfo("Poll choice %d deleted successfully\n", opts.ChoiceID)
	logger.LogCommandComplete(ctx, "polls.choices.delete", 1)
	return nil
}

func runPollSessionList(ctx context.Context, client *api.Client, opts *options.PollSessionListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.sessions.list", map[string]interface{}{"poll_id": opts.PollID})

	svc := api.NewPollsService(client)

	sessions, err := svc.ListPollSessions(ctx, opts.PollID)
	if err != nil {
		logger.LogCommandError(ctx, "polls.sessions.list", err, map[string]interface{}{"poll_id": opts.PollID})
		return fmt.Errorf("failed to list poll sessions: %w", err)
	}

	printVerbose("Found %d poll sessions:\n\n", len(sessions))
	logger.LogCommandComplete(ctx, "polls.sessions.list", len(sessions))
	return formatEmptyOrOutput(sessions, "No poll sessions found")
}

func runPollSessionGet(ctx context.Context, client *api.Client, opts *options.PollSessionGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.sessions.get", map[string]interface{}{
		"poll_id":    opts.PollID,
		"session_id": opts.SessionID,
	})

	svc := api.NewPollsService(client)

	session, err := svc.GetPollSession(ctx, opts.PollID, opts.SessionID)
	if err != nil {
		logger.LogCommandError(ctx, "polls.sessions.get", err, nil)
		return fmt.Errorf("failed to get poll session: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.sessions.get", 1)
	return formatOutput(session, nil)
}

func runPollSessionCreate(ctx context.Context, client *api.Client, opts *options.PollSessionCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.sessions.create", map[string]interface{}{
		"poll_id":   opts.PollID,
		"course_id": opts.CourseID,
	})

	svc := api.NewPollsService(client)

	session, err := svc.CreatePollSession(ctx, opts.PollID, &api.CreatePollSessionParams{
		CourseID:         opts.CourseID,
		CourseSectionID:  opts.CourseSectionID,
		HasPublicResults: opts.HasPublicResults,
	})
	if err != nil {
		logger.LogCommandError(ctx, "polls.sessions.create", err, nil)
		return fmt.Errorf("failed to create poll session: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.sessions.create", 1)
	return formatSuccessOutput(session, "Poll session created successfully!")
}

func runPollSessionUpdate(ctx context.Context, client *api.Client, opts *options.PollSessionUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.sessions.update", map[string]interface{}{
		"poll_id":    opts.PollID,
		"session_id": opts.SessionID,
	})

	svc := api.NewPollsService(client)

	session, err := svc.UpdatePollSession(ctx, opts.PollID, opts.SessionID, &api.UpdatePollSessionParams{
		CourseID:         opts.CourseID,
		CourseSectionID:  opts.CourseSectionID,
		HasPublicResults: opts.HasPublicResults,
	})
	if err != nil {
		logger.LogCommandError(ctx, "polls.sessions.update", err, nil)
		return fmt.Errorf("failed to update poll session: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.sessions.update", 1)
	return formatSuccessOutput(session, "Poll session updated successfully!")
}

func runPollSessionDelete(ctx context.Context, client *api.Client, opts *options.PollSessionDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.sessions.delete", map[string]interface{}{
		"poll_id":    opts.PollID,
		"session_id": opts.SessionID,
		"force":      opts.Force,
	})

	confirmed, err := confirmDelete("poll session", opts.SessionID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "polls.sessions.delete", err, nil)
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		logger.LogCommandComplete(ctx, "polls.sessions.delete", 0)
		return nil
	}

	svc := api.NewPollsService(client)

	if err := svc.DeletePollSession(ctx, opts.PollID, opts.SessionID); err != nil {
		logger.LogCommandError(ctx, "polls.sessions.delete", err, nil)
		return fmt.Errorf("failed to delete poll session: %w", err)
	}

	printInfo("Poll session %d deleted successfully\n", opts.SessionID)
	logger.LogCommandComplete(ctx, "polls.sessions.delete", 1)
	return nil
}

func runPollSessionOpen(ctx context.Context, client *api.Client, opts *options.PollSessionOpenOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.sessions.open", map[string]interface{}{
		"poll_id":    opts.PollID,
		"session_id": opts.SessionID,
	})

	svc := api.NewPollsService(client)

	session, err := svc.OpenPollSession(ctx, opts.PollID, opts.SessionID)
	if err != nil {
		logger.LogCommandError(ctx, "polls.sessions.open", err, nil)
		return fmt.Errorf("failed to open poll session: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.sessions.open", 1)
	return formatSuccessOutput(session, "Poll session opened successfully!")
}

func runPollSessionClose(ctx context.Context, client *api.Client, opts *options.PollSessionCloseOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.sessions.close", map[string]interface{}{
		"poll_id":    opts.PollID,
		"session_id": opts.SessionID,
	})

	svc := api.NewPollsService(client)

	session, err := svc.ClosePollSession(ctx, opts.PollID, opts.SessionID)
	if err != nil {
		logger.LogCommandError(ctx, "polls.sessions.close", err, nil)
		return fmt.Errorf("failed to close poll session: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.sessions.close", 1)
	return formatSuccessOutput(session, "Poll session closed successfully!")
}

func runPollSessionListOpened(ctx context.Context, client *api.Client) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.sessions.list-opened", nil)

	svc := api.NewPollsService(client)

	sessions, err := svc.ListOpenedPollSessions(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "polls.sessions.list-opened", err, nil)
		return fmt.Errorf("failed to list opened poll sessions: %w", err)
	}

	printVerbose("Found %d open poll sessions:\n\n", len(sessions))
	logger.LogCommandComplete(ctx, "polls.sessions.list-opened", len(sessions))
	return formatEmptyOrOutput(sessions, "No open poll sessions found")
}

func runPollSessionListClosed(ctx context.Context, client *api.Client) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.sessions.list-closed", nil)

	svc := api.NewPollsService(client)

	sessions, err := svc.ListClosedPollSessions(ctx)
	if err != nil {
		logger.LogCommandError(ctx, "polls.sessions.list-closed", err, nil)
		return fmt.Errorf("failed to list closed poll sessions: %w", err)
	}

	printVerbose("Found %d closed poll sessions:\n\n", len(sessions))
	logger.LogCommandComplete(ctx, "polls.sessions.list-closed", len(sessions))
	return formatEmptyOrOutput(sessions, "No closed poll sessions found")
}

func runPollSubmissionGet(ctx context.Context, client *api.Client, opts *options.PollSubmissionGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.submissions.get", map[string]interface{}{
		"poll_id":       opts.PollID,
		"session_id":    opts.SessionID,
		"submission_id": opts.SubmissionID,
	})

	svc := api.NewPollsService(client)

	sub, err := svc.GetPollSubmission(ctx, opts.PollID, opts.SessionID, opts.SubmissionID)
	if err != nil {
		logger.LogCommandError(ctx, "polls.submissions.get", err, nil)
		return fmt.Errorf("failed to get poll submission: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.submissions.get", 1)
	return formatOutput(sub, nil)
}

func runPollSubmissionCreate(ctx context.Context, client *api.Client, opts *options.PollSubmissionCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "polls.submissions.create", map[string]interface{}{
		"poll_id":        opts.PollID,
		"session_id":     opts.SessionID,
		"poll_choice_id": opts.PollChoiceID,
	})

	svc := api.NewPollsService(client)

	sub, err := svc.CreatePollSubmission(ctx, opts.PollID, opts.SessionID, &api.CreatePollSubmissionParams{
		PollChoiceID: opts.PollChoiceID,
	})
	if err != nil {
		logger.LogCommandError(ctx, "polls.submissions.create", err, nil)
		return fmt.Errorf("failed to submit poll vote: %w", err)
	}

	logger.LogCommandComplete(ctx, "polls.submissions.create", 1)
	return formatSuccessOutput(sub, "Vote submitted successfully!")
}
