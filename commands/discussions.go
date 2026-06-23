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

// discussionsCmd represents the discussions command group
var discussionsCmd = &cobra.Command{
	Use:   "discussions",
	Short: "Manage Canvas discussion topics",
	Long: `Manage Canvas discussion topics including listing, viewing, creating, and updating discussions.

Discussion topics are threaded conversations associated with Courses or Groups in Canvas.
They can be used for class discussions, Q&A, and collaborative learning.

Examples:
  canvas discussions list --course-id 123
  canvas discussions list --group-id 456
  canvas discussions get --course-id 123 789
  canvas discussions create --course-id 123 --title "Week 1 Discussion"
  canvas discussions entries --course-id 123 789`,
}

// addContextFlags is a helper that adds both --course-id and --group-id flags.
func addContextFlags(cmd *cobra.Command, courseID, groupID *int64) {
	cmd.Flags().Int64Var(courseID, "course-id", 0, "Course ID (mutually exclusive with --group-id)")
	cmd.Flags().Int64Var(groupID, "group-id", 0, "Group ID (mutually exclusive with --course-id)")
}

func init() {
	rootCmd.AddCommand(discussionsCmd)
	discussionsCmd.AddCommand(newDiscussionsListCmd())
	discussionsCmd.AddCommand(newDiscussionsGetCmd())
	discussionsCmd.AddCommand(newDiscussionsViewCmd())
	discussionsCmd.AddCommand(newDiscussionsCreateCmd())
	discussionsCmd.AddCommand(newDiscussionsUpdateCmd())
	discussionsCmd.AddCommand(newDiscussionsDeleteCmd())
	discussionsCmd.AddCommand(newDiscussionsDuplicateCmd())
	discussionsCmd.AddCommand(newDiscussionsReorderCmd())
	discussionsCmd.AddCommand(newDiscussionsEntriesCmd())
	discussionsCmd.AddCommand(newDiscussionsPostCmd())
	discussionsCmd.AddCommand(newDiscussionsUpdateEntryCmd())
	discussionsCmd.AddCommand(newDiscussionsDeleteEntryCmd())
	discussionsCmd.AddCommand(newDiscussionsReplyCmd())
	discussionsCmd.AddCommand(newDiscussionsRepliesCmd())
	discussionsCmd.AddCommand(newDiscussionsEntryListCmd())
	discussionsCmd.AddCommand(newDiscussionsMarkTopicReadCmd())
	discussionsCmd.AddCommand(newDiscussionsMarkTopicUnreadCmd())
	discussionsCmd.AddCommand(newDiscussionsMarkAllTopicsReadCmd())
	discussionsCmd.AddCommand(newDiscussionsMarkAllEntriesReadCmd())
	discussionsCmd.AddCommand(newDiscussionsMarkAllEntriesUnreadCmd())
	discussionsCmd.AddCommand(newDiscussionsMarkEntryReadCmd())
	discussionsCmd.AddCommand(newDiscussionsMarkEntryUnreadCmd())
	discussionsCmd.AddCommand(newDiscussionsRateEntryCmd())
	discussionsCmd.AddCommand(newDiscussionsSubscribeCmd())
	discussionsCmd.AddCommand(newDiscussionsUnsubscribeCmd())
}

func newDiscussionsListCmd() *cobra.Command {
	opts := &options.DiscussionsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List discussion topics in a course or group",
		Long: `List all discussion topics in a Canvas course or group.

Examples:
  canvas discussions list --course-id 123
  canvas discussions list --group-id 456
  canvas discussions list --course-id 123 --order-by recent_activity
  canvas discussions list --course-id 123 --scope pinned
  canvas discussions list --course-id 123 --filter unread`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsList(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().StringVar(&opts.OrderBy, "order-by", "", "Order by: position, recent_activity, title")
	cmd.Flags().StringVar(&opts.Scope, "scope", "", "Scope: locked, unlocked, pinned, unpinned")
	cmd.Flags().BoolVar(&opts.OnlyAnnouncements, "announcements", false, "Only show announcements")
	cmd.Flags().StringVar(&opts.FilterBy, "filter", "", "Filter by: all, unread")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search term")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include: all_dates, sections, sections_user_count, overrides")

	return cmd
}

func newDiscussionsGetCmd() *cobra.Command {
	opts := &options.DiscussionsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <topic-id>",
		Short: "Get a specific discussion topic",
		Long: `Get details of a specific discussion topic.

Examples:
  canvas discussions get --course-id 123 456
  canvas discussions get --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsGet(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include additional data")

	return cmd
}

func newDiscussionsViewCmd() *cobra.Command {
	opts := &options.DiscussionsViewOptions{}

	cmd := &cobra.Command{
		Use:   "view <topic-id>",
		Short: "Get the full threaded view of a discussion topic",
		Long: `Get the full cached view of a discussion topic, including all entries and participants.

Examples:
  canvas discussions view --course-id 123 456
  canvas discussions view --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsView(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsCreateCmd() *cobra.Command {
	opts := &options.DiscussionsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new discussion topic",
		Long: `Create a new discussion topic in a course or group.

Examples:
  canvas discussions create --course-id 123 --title "Week 1 Discussion"
  canvas discussions create --group-id 456 --title "Group Q&A"
  canvas discussions create --course-id 123 --title "Q&A" --message "<p>Ask questions here</p>" --type threaded
  canvas discussions create --course-id 123 --title "Pinned" --pinned --published`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsCreate(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().StringVar(&opts.Title, "title", "", "Discussion title (required)")
	cmd.Flags().StringVar(&opts.Message, "message", "", "Discussion message (HTML)")
	cmd.Flags().StringVar(&opts.DiscussionType, "type", "", "Discussion type: side_comment, threaded, not_threaded")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish the discussion")
	cmd.Flags().StringVar(&opts.DelayedPostAt, "delayed-post-at", "", "Delay posting until (ISO 8601)")
	cmd.Flags().BoolVar(&opts.AllowRating, "allow-rating", false, "Allow rating of entries")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock at date (ISO 8601)")
	cmd.Flags().BoolVar(&opts.RequireInitialPost, "require-initial-post", false, "Require initial post before viewing")
	cmd.Flags().BoolVar(&opts.Pinned, "pinned", false, "Pin the discussion")
	mustMarkRequired(cmd, "title")

	return cmd
}

func newDiscussionsUpdateCmd() *cobra.Command {
	opts := &options.DiscussionsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <topic-id>",
		Short: "Update an existing discussion topic",
		Long: `Update an existing discussion topic in a course or group.

Examples:
  canvas discussions update --course-id 123 456 --title "New Title"
  canvas discussions update --group-id 789 456 --pinned
  canvas discussions update --course-id 123 456 --locked`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			// Track which fields were changed
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.MessageSet = cmd.Flags().Changed("message")
			opts.DiscussionTypeSet = cmd.Flags().Changed("type")
			opts.PublishedSet = cmd.Flags().Changed("published")
			opts.DelayedPostAtSet = cmd.Flags().Changed("delayed-post-at")
			opts.AllowRatingSet = cmd.Flags().Changed("allow-rating")
			opts.LockAtSet = cmd.Flags().Changed("lock-at")
			opts.RequireInitialPostSet = cmd.Flags().Changed("require-initial-post")
			opts.PinnedSet = cmd.Flags().Changed("pinned")
			opts.LockedSet = cmd.Flags().Changed("locked")
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsUpdate(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().StringVar(&opts.Title, "title", "", "New discussion title")
	cmd.Flags().StringVar(&opts.Message, "message", "", "New discussion message")
	cmd.Flags().StringVar(&opts.DiscussionType, "type", "", "Discussion type")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish the discussion")
	cmd.Flags().StringVar(&opts.DelayedPostAt, "delayed-post-at", "", "Delay posting until")
	cmd.Flags().BoolVar(&opts.AllowRating, "allow-rating", false, "Allow rating")
	cmd.Flags().StringVar(&opts.LockAt, "lock-at", "", "Lock at date")
	cmd.Flags().BoolVar(&opts.RequireInitialPost, "require-initial-post", false, "Require initial post")
	cmd.Flags().BoolVar(&opts.Pinned, "pinned", false, "Pin the discussion")
	cmd.Flags().BoolVar(&opts.Locked, "locked", false, "Lock the discussion")

	return cmd
}

func newDiscussionsDeleteCmd() *cobra.Command {
	opts := &options.DiscussionsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <topic-id>",
		Short: "Delete a discussion topic",
		Long: `Delete a discussion topic from a course or group.

Examples:
  canvas discussions delete --course-id 123 456
  canvas discussions delete --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsDelete(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func newDiscussionsDuplicateCmd() *cobra.Command {
	opts := &options.DiscussionsDuplicateOptions{}

	cmd := &cobra.Command{
		Use:   "duplicate <topic-id>",
		Short: "Duplicate a discussion topic",
		Long: `Duplicate a discussion topic in a course or group.

Examples:
  canvas discussions duplicate --course-id 123 456
  canvas discussions duplicate --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsDuplicate(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsReorderCmd() *cobra.Command {
	opts := &options.DiscussionsReorderOptions{}
	var orderStr []string

	cmd := &cobra.Command{
		Use:   "reorder",
		Short: "Reorder pinned discussion topics",
		Long: `Reorder pinned discussion topics in a course or group.
All pinned topics must be included in the order list.

Examples:
  canvas discussions reorder --course-id 123 --order 104,102,103
  canvas discussions reorder --group-id 456 --order 5,3,1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse order IDs from string slice
			for _, s := range orderStr {
				id, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid topic ID in order: %s", s)
				}
				opts.Order = append(opts.Order, id)
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsReorder(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().StringSliceVar(&orderStr, "order", []string{}, "Comma-separated list of topic IDs in desired order (required)")
	mustMarkRequired(cmd, "order")

	return cmd
}

func newDiscussionsEntriesCmd() *cobra.Command {
	opts := &options.DiscussionsEntriesOptions{}

	cmd := &cobra.Command{
		Use:   "entries <topic-id>",
		Short: "List top-level entries in a discussion",
		Long: `List all top-level entries (posts) in a discussion topic.

Examples:
  canvas discussions entries --course-id 123 456
  canvas discussions entries --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsEntries(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsEntryListCmd() *cobra.Command {
	opts := &options.DiscussionsEntryListOptions{}
	var idStrs []string

	cmd := &cobra.Command{
		Use:   "entry-list <topic-id>",
		Short: "Fetch specific discussion entries by ID",
		Long: `Retrieve a list of specific discussion entries by their IDs.

Examples:
  canvas discussions entry-list --course-id 123 456 --ids 1,2,3
  canvas discussions entry-list --group-id 789 456 --ids 10,11`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			for _, s := range idStrs {
				id, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid entry ID: %s", s)
				}
				opts.IDs = append(opts.IDs, id)
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsEntryList(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().StringSliceVar(&idStrs, "ids", []string{}, "Comma-separated list of entry IDs to fetch")

	return cmd
}

func newDiscussionsPostCmd() *cobra.Command {
	opts := &options.DiscussionsPostOptions{}

	cmd := &cobra.Command{
		Use:   "post <topic-id> [message]",
		Short: "Post a new entry to a discussion",
		Long: `Post a new top-level entry to a discussion topic.

The message can be provided as a positional argument or using the --message flag.

Examples:
  canvas discussions post --course-id 123 456 "My response to the discussion"
  canvas discussions post --course-id 123 456 --message "My response to the discussion"
  canvas discussions post --group-id 789 456 "Group discussion response"`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			// Get message from positional arg or --message flag
			if len(args) > 1 {
				opts.Message = args[1]
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsPost(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().StringVarP(&opts.Message, "message", "m", "", "Message content (alternative to positional argument)")

	return cmd
}

func newDiscussionsUpdateEntryCmd() *cobra.Command {
	opts := &options.DiscussionsUpdateEntryOptions{}

	cmd := &cobra.Command{
		Use:   "update-entry <topic-id> <entry-id>",
		Short: "Update a discussion entry",
		Long: `Update the message of an existing discussion entry.

Examples:
  canvas discussions update-entry --course-id 123 456 789 --message "Corrected response"
  canvas discussions update-entry --group-id 789 456 123 --message "Updated text"`,
		Args: ExactArgsWithUsage(2, "topic-id", "entry-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			entryID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid entry ID: %s", args[1])
			}
			opts.EntryID = entryID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsUpdateEntry(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().StringVarP(&opts.Message, "message", "m", "", "Updated message content (required)")
	mustMarkRequired(cmd, "message")

	return cmd
}

func newDiscussionsDeleteEntryCmd() *cobra.Command {
	opts := &options.DiscussionsDeleteEntryOptions{}

	cmd := &cobra.Command{
		Use:   "delete-entry <topic-id> <entry-id>",
		Short: "Delete a discussion entry",
		Long: `Delete an existing discussion entry.

Examples:
  canvas discussions delete-entry --course-id 123 456 789
  canvas discussions delete-entry --group-id 789 456 123 --force`,
		Args: ExactArgsWithUsage(2, "topic-id", "entry-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			entryID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid entry ID: %s", args[1])
			}
			opts.EntryID = entryID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsDeleteEntry(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func newDiscussionsReplyCmd() *cobra.Command {
	opts := &options.DiscussionsReplyOptions{}

	cmd := &cobra.Command{
		Use:   "reply <topic-id> <entry-id> [message]",
		Short: "Reply to an entry in a discussion",
		Long: `Reply to a specific entry in a discussion topic.

The message can be provided as a positional argument or using the --message flag.

Examples:
  canvas discussions reply --course-id 123 456 789 "My reply to this entry"
  canvas discussions reply --course-id 123 456 789 --message "My reply to this entry"
  canvas discussions reply --group-id 789 456 123 "Group reply"`,
		Args: cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			entryID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid entry ID: %s", args[1])
			}
			opts.EntryID = entryID
			// Get message from positional arg or --message flag
			if len(args) > 2 {
				opts.Message = args[2]
			}
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsReply(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().StringVarP(&opts.Message, "message", "m", "", "Message content (alternative to positional argument)")

	return cmd
}

func newDiscussionsRepliesCmd() *cobra.Command {
	opts := &options.DiscussionsRepliesOptions{}

	cmd := &cobra.Command{
		Use:   "replies <topic-id> <entry-id>",
		Short: "List replies to a discussion entry",
		Long: `List all replies to a specific top-level entry in a discussion topic.

Examples:
  canvas discussions replies --course-id 123 456 789
  canvas discussions replies --group-id 789 456 123`,
		Args: ExactArgsWithUsage(2, "topic-id", "entry-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			entryID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid entry ID: %s", args[1])
			}
			opts.EntryID = entryID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsReplies(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsMarkTopicReadCmd() *cobra.Command {
	opts := &options.DiscussionsMarkTopicReadOptions{}

	cmd := &cobra.Command{
		Use:   "mark-read <topic-id>",
		Short: "Mark a discussion topic as read",
		Long: `Mark the initial text of a discussion topic as read.

Examples:
  canvas discussions mark-read --course-id 123 456
  canvas discussions mark-read --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsMarkTopicRead(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsMarkTopicUnreadCmd() *cobra.Command {
	opts := &options.DiscussionsMarkTopicReadOptions{}

	cmd := &cobra.Command{
		Use:   "mark-unread <topic-id>",
		Short: "Mark a discussion topic as unread",
		Long: `Mark the initial text of a discussion topic as unread.

Examples:
  canvas discussions mark-unread --course-id 123 456
  canvas discussions mark-unread --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsMarkTopicUnread(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsMarkAllTopicsReadCmd() *cobra.Command {
	opts := &options.DiscussionsMarkAllTopicsReadOptions{}

	cmd := &cobra.Command{
		Use:   "mark-all-read",
		Short: "Mark all discussion topics as read",
		Long: `Mark all discussion topics in a course or group as read.

Examples:
  canvas discussions mark-all-read --course-id 123
  canvas discussions mark-all-read --group-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsMarkAllTopicsRead(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsMarkAllEntriesReadCmd() *cobra.Command {
	opts := &options.DiscussionsMarkAllEntriesReadOptions{}

	cmd := &cobra.Command{
		Use:   "mark-entries-read <topic-id>",
		Short: "Mark all entries in a topic as read",
		Long: `Mark all entries in a discussion topic as read.

Examples:
  canvas discussions mark-entries-read --course-id 123 456
  canvas discussions mark-entries-read --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsMarkAllEntriesRead(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsMarkAllEntriesUnreadCmd() *cobra.Command {
	opts := &options.DiscussionsMarkAllEntriesReadOptions{}

	cmd := &cobra.Command{
		Use:   "mark-entries-unread <topic-id>",
		Short: "Mark all entries in a topic as unread",
		Long: `Mark all entries in a discussion topic as unread.

Examples:
  canvas discussions mark-entries-unread --course-id 123 456
  canvas discussions mark-entries-unread --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsMarkAllEntriesUnread(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsMarkEntryReadCmd() *cobra.Command {
	opts := &options.DiscussionsMarkEntryReadOptions{}

	cmd := &cobra.Command{
		Use:   "mark-entry-read <topic-id> <entry-id>",
		Short: "Mark a discussion entry as read",
		Long: `Mark a specific discussion entry as read.

Examples:
  canvas discussions mark-entry-read --course-id 123 456 789
  canvas discussions mark-entry-read --group-id 789 456 123`,
		Args: ExactArgsWithUsage(2, "topic-id", "entry-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			entryID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid entry ID: %s", args[1])
			}
			opts.EntryID = entryID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsMarkEntryRead(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsMarkEntryUnreadCmd() *cobra.Command {
	opts := &options.DiscussionsMarkEntryReadOptions{}

	cmd := &cobra.Command{
		Use:   "mark-entry-unread <topic-id> <entry-id>",
		Short: "Mark a discussion entry as unread",
		Long: `Mark a specific discussion entry as unread.

Examples:
  canvas discussions mark-entry-unread --course-id 123 456 789
  canvas discussions mark-entry-unread --group-id 789 456 123`,
		Args: ExactArgsWithUsage(2, "topic-id", "entry-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			entryID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid entry ID: %s", args[1])
			}
			opts.EntryID = entryID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsMarkEntryUnread(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsRateEntryCmd() *cobra.Command {
	opts := &options.DiscussionsRateEntryOptions{}

	cmd := &cobra.Command{
		Use:   "rate-entry <topic-id> <entry-id>",
		Short: "Rate a discussion entry",
		Long: `Rate a discussion entry. Rating must be 0 (un-rate) or 1 (like).

Examples:
  canvas discussions rate-entry --course-id 123 456 789 --rating 1
  canvas discussions rate-entry --group-id 789 456 123 --rating 0`,
		Args: ExactArgsWithUsage(2, "topic-id", "entry-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			entryID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid entry ID: %s", args[1])
			}
			opts.EntryID = entryID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsRateEntry(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)
	cmd.Flags().IntVar(&opts.Rating, "rating", 1, "Rating value: 0 (un-rate) or 1 (like)")

	return cmd
}

func newDiscussionsSubscribeCmd() *cobra.Command {
	opts := &options.DiscussionsSubscribeOptions{}

	cmd := &cobra.Command{
		Use:   "subscribe <topic-id>",
		Short: "Subscribe to a discussion topic",
		Long: `Subscribe to receive notifications for a discussion topic.

Examples:
  canvas discussions subscribe --course-id 123 456
  canvas discussions subscribe --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsSubscribe(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

func newDiscussionsUnsubscribeCmd() *cobra.Command {
	opts := &options.DiscussionsUnsubscribeOptions{}

	cmd := &cobra.Command{
		Use:   "unsubscribe <topic-id>",
		Short: "Unsubscribe from a discussion topic",
		Long: `Unsubscribe from a discussion topic to stop receiving notifications.

Examples:
  canvas discussions unsubscribe --course-id 123 456
  canvas discussions unsubscribe --group-id 789 456`,
		Args: ExactArgsWithUsage(1, "topic-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			topicID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid topic ID: %s", args[0])
			}
			opts.TopicID = topicID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runDiscussionsUnsubscribe(cmd.Context(), client, opts)
		},
	}

	addContextFlags(cmd, &opts.CourseID, &opts.GroupID)

	return cmd
}

// --- run functions ---

func runDiscussionsList(ctx context.Context, client *api.Client, opts *options.DiscussionsListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.list", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"order_by":     opts.OrderBy,
		"scope":        opts.Scope,
		"filter_by":    opts.FilterBy,
		"search_term":  opts.SearchTerm,
	})

	discussionsService := api.NewDiscussionsService(client)

	listOpts := &api.ListDiscussionsOptions{
		Include:           opts.Include,
		OrderBy:           opts.OrderBy,
		Scope:             opts.Scope,
		OnlyAnnouncements: opts.OnlyAnnouncements,
		FilterBy:          opts.FilterBy,
		SearchTerm:        opts.SearchTerm,
	}

	topics, err := discussionsService.ListContext(ctx, ctxType, ctxID, listOpts)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.list", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
		})
		return fmt.Errorf("failed to list discussions: %w", err)
	}

	printVerbose("Found %d discussion topics:\n\n", len(topics))
	logger.LogCommandComplete(ctx, "discussions.list", len(topics))
	return formatEmptyOrOutput(topics, "No discussion topics found")
}

func runDiscussionsGet(ctx context.Context, client *api.Client, opts *options.DiscussionsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.get", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	topic, err := discussionsService.GetContext(ctx, ctxType, ctxID, opts.TopicID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.get", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to get discussion: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.get", 1)
	return formatOutput(topic, nil)
}

func runDiscussionsView(ctx context.Context, client *api.Client, opts *options.DiscussionsViewOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.view", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	view, err := discussionsService.GetView(ctx, ctxType, ctxID, opts.TopicID)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.view", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to get discussion view: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.view", 1)
	return formatOutput(view, nil)
}

func runDiscussionsCreate(ctx context.Context, client *api.Client, opts *options.DiscussionsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.create", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"title":        opts.Title,
	})

	discussionsService := api.NewDiscussionsService(client)

	params := &api.CreateDiscussionParams{
		Title:              opts.Title,
		Message:            opts.Message,
		DiscussionType:     opts.DiscussionType,
		Published:          opts.Published,
		DelayedPostAt:      opts.DelayedPostAt,
		AllowRating:        opts.AllowRating,
		LockAt:             opts.LockAt,
		RequireInitialPost: opts.RequireInitialPost,
		Pinned:             opts.Pinned,
	}

	topic, err := discussionsService.CreateContext(ctx, ctxType, ctxID, params)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.create", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"title":        opts.Title,
		})
		return fmt.Errorf("failed to create discussion: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.create", 1)
	return formatSuccessOutput(topic, "Discussion created successfully!")
}

func runDiscussionsUpdate(ctx context.Context, client *api.Client, opts *options.DiscussionsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.update", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	params := &api.UpdateDiscussionParams{}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.MessageSet {
		params.Message = &opts.Message
	}
	if opts.DiscussionTypeSet {
		params.DiscussionType = &opts.DiscussionType
	}
	if opts.PublishedSet {
		params.Published = &opts.Published
	}
	if opts.DelayedPostAtSet {
		params.DelayedPostAt = &opts.DelayedPostAt
	}
	if opts.AllowRatingSet {
		params.AllowRating = &opts.AllowRating
	}
	if opts.LockAtSet {
		params.LockAt = &opts.LockAt
	}
	if opts.RequireInitialPostSet {
		params.RequireInitialPost = &opts.RequireInitialPost
	}
	if opts.PinnedSet {
		params.Pinned = &opts.Pinned
	}
	if opts.LockedSet {
		params.Locked = &opts.Locked
	}

	topic, err := discussionsService.UpdateContext(ctx, ctxType, ctxID, opts.TopicID, params)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.update", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to update discussion: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.update", 1)
	return formatSuccessOutput(topic, "Discussion updated successfully!")
}

func runDiscussionsDelete(ctx context.Context, client *api.Client, opts *options.DiscussionsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.delete", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
		"force":        opts.Force,
	})

	// Confirm deletion
	confirmed, err := confirmDelete("discussion", opts.TopicID, opts.Force)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		return nil
	}

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.DeleteContext(ctx, ctxType, ctxID, opts.TopicID); err != nil {
		logger.LogCommandError(ctx, "discussions.delete", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to delete discussion: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.delete", 1)
	printInfo("Discussion %d deleted successfully\n", opts.TopicID)
	return nil
}

func runDiscussionsDuplicate(ctx context.Context, client *api.Client, opts *options.DiscussionsDuplicateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.duplicate", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	topic, err := discussionsService.Duplicate(ctx, ctxType, ctxID, opts.TopicID)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.duplicate", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to duplicate discussion: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.duplicate", 1)
	return formatSuccessOutput(topic, "Discussion duplicated successfully!")
}

func runDiscussionsReorder(ctx context.Context, client *api.Client, opts *options.DiscussionsReorderOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.reorder", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"order":        opts.Order,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.Reorder(ctx, ctxType, ctxID, opts.Order); err != nil {
		logger.LogCommandError(ctx, "discussions.reorder", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
		})
		return fmt.Errorf("failed to reorder discussions: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.reorder", len(opts.Order))
	printInfo("Pinned discussions reordered successfully\n")
	return nil
}

func runDiscussionsEntries(ctx context.Context, client *api.Client, opts *options.DiscussionsEntriesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.entries", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	entries, err := discussionsService.ListEntriesContext(ctx, ctxType, ctxID, opts.TopicID)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.entries", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to list entries: %w", err)
	}

	printVerbose("Found %d entries:\n\n", len(entries))
	logger.LogCommandComplete(ctx, "discussions.entries", len(entries))
	return formatEmptyOrOutput(entries, "No entries found")
}

func runDiscussionsEntryList(ctx context.Context, client *api.Client, opts *options.DiscussionsEntryListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.entry-list", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
		"ids":          opts.IDs,
	})

	discussionsService := api.NewDiscussionsService(client)

	entries, err := discussionsService.GetEntryList(ctx, ctxType, ctxID, opts.TopicID, opts.IDs)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.entry-list", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to get entry list: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.entry-list", len(entries))
	return formatEmptyOrOutput(entries, "No entries found")
}

func runDiscussionsPost(ctx context.Context, client *api.Client, opts *options.DiscussionsPostOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.post", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	entry, err := discussionsService.PostEntryContext(ctx, ctxType, ctxID, opts.TopicID, opts.Message)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.post", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to post entry: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.post", 1)
	return formatSuccessOutput(entry, "Entry posted successfully!")
}

func runDiscussionsUpdateEntry(ctx context.Context, client *api.Client, opts *options.DiscussionsUpdateEntryOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.update-entry", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
		"entry_id":     opts.EntryID,
	})

	discussionsService := api.NewDiscussionsService(client)

	entry, err := discussionsService.UpdateEntry(ctx, ctxType, ctxID, opts.TopicID, opts.EntryID, opts.Message)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.update-entry", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
			"entry_id":     opts.EntryID,
		})
		return fmt.Errorf("failed to update entry: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.update-entry", 1)
	return formatSuccessOutput(entry, "Entry updated successfully!")
}

func runDiscussionsDeleteEntry(ctx context.Context, client *api.Client, opts *options.DiscussionsDeleteEntryOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.delete-entry", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
		"entry_id":     opts.EntryID,
		"force":        opts.Force,
	})

	confirmed, err := confirmDelete("discussion entry", opts.EntryID, opts.Force)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		return nil
	}

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.DeleteEntry(ctx, ctxType, ctxID, opts.TopicID, opts.EntryID); err != nil {
		logger.LogCommandError(ctx, "discussions.delete-entry", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
			"entry_id":     opts.EntryID,
		})
		return fmt.Errorf("failed to delete entry: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.delete-entry", 1)
	printInfo("Entry %d deleted successfully\n", opts.EntryID)
	return nil
}

func runDiscussionsReply(ctx context.Context, client *api.Client, opts *options.DiscussionsReplyOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.reply", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
		"entry_id":     opts.EntryID,
	})

	discussionsService := api.NewDiscussionsService(client)

	entry, err := discussionsService.PostReplyContext(ctx, ctxType, ctxID, opts.TopicID, opts.EntryID, opts.Message)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.reply", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
			"entry_id":     opts.EntryID,
		})
		return fmt.Errorf("failed to post reply: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.reply", 1)
	return formatSuccessOutput(entry, "Reply posted successfully!")
}

func runDiscussionsReplies(ctx context.Context, client *api.Client, opts *options.DiscussionsRepliesOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.replies", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
		"entry_id":     opts.EntryID,
	})

	discussionsService := api.NewDiscussionsService(client)

	replies, err := discussionsService.ListReplies(ctx, ctxType, ctxID, opts.TopicID, opts.EntryID)
	if err != nil {
		logger.LogCommandError(ctx, "discussions.replies", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
			"entry_id":     opts.EntryID,
		})
		return fmt.Errorf("failed to list replies: %w", err)
	}

	printVerbose("Found %d replies:\n\n", len(replies))
	logger.LogCommandComplete(ctx, "discussions.replies", len(replies))
	return formatEmptyOrOutput(replies, "No replies found")
}

func runDiscussionsMarkTopicRead(ctx context.Context, client *api.Client, opts *options.DiscussionsMarkTopicReadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.mark-read", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.MarkTopicReadContext(ctx, ctxType, ctxID, opts.TopicID); err != nil {
		logger.LogCommandError(ctx, "discussions.mark-read", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to mark topic as read: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.mark-read", 1)
	printInfo("Discussion %d marked as read\n", opts.TopicID)
	return nil
}

func runDiscussionsMarkTopicUnread(ctx context.Context, client *api.Client, opts *options.DiscussionsMarkTopicReadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.mark-unread", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.MarkTopicUnreadContext(ctx, ctxType, ctxID, opts.TopicID); err != nil {
		logger.LogCommandError(ctx, "discussions.mark-unread", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to mark topic as unread: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.mark-unread", 1)
	printInfo("Discussion %d marked as unread\n", opts.TopicID)
	return nil
}

func runDiscussionsMarkAllTopicsRead(ctx context.Context, client *api.Client, opts *options.DiscussionsMarkAllTopicsReadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.mark-all-read", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.MarkAllTopicsRead(ctx, ctxType, ctxID); err != nil {
		logger.LogCommandError(ctx, "discussions.mark-all-read", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
		})
		return fmt.Errorf("failed to mark all topics as read: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.mark-all-read", 1)
	printInfo("All discussion topics marked as read\n")
	return nil
}

func runDiscussionsMarkAllEntriesRead(ctx context.Context, client *api.Client, opts *options.DiscussionsMarkAllEntriesReadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.mark-entries-read", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.MarkAllEntriesRead(ctx, ctxType, ctxID, opts.TopicID); err != nil {
		logger.LogCommandError(ctx, "discussions.mark-entries-read", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to mark all entries as read: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.mark-entries-read", 1)
	printInfo("All entries in discussion %d marked as read\n", opts.TopicID)
	return nil
}

func runDiscussionsMarkAllEntriesUnread(ctx context.Context, client *api.Client, opts *options.DiscussionsMarkAllEntriesReadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.mark-entries-unread", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.MarkAllEntriesUnread(ctx, ctxType, ctxID, opts.TopicID); err != nil {
		logger.LogCommandError(ctx, "discussions.mark-entries-unread", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to mark all entries as unread: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.mark-entries-unread", 1)
	printInfo("All entries in discussion %d marked as unread\n", opts.TopicID)
	return nil
}

func runDiscussionsMarkEntryRead(ctx context.Context, client *api.Client, opts *options.DiscussionsMarkEntryReadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.mark-entry-read", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
		"entry_id":     opts.EntryID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.MarkEntryRead(ctx, ctxType, ctxID, opts.TopicID, opts.EntryID); err != nil {
		logger.LogCommandError(ctx, "discussions.mark-entry-read", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
			"entry_id":     opts.EntryID,
		})
		return fmt.Errorf("failed to mark entry as read: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.mark-entry-read", 1)
	printInfo("Entry %d marked as read\n", opts.EntryID)
	return nil
}

func runDiscussionsMarkEntryUnread(ctx context.Context, client *api.Client, opts *options.DiscussionsMarkEntryReadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.mark-entry-unread", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
		"entry_id":     opts.EntryID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.MarkEntryUnread(ctx, ctxType, ctxID, opts.TopicID, opts.EntryID); err != nil {
		logger.LogCommandError(ctx, "discussions.mark-entry-unread", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
			"entry_id":     opts.EntryID,
		})
		return fmt.Errorf("failed to mark entry as unread: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.mark-entry-unread", 1)
	printInfo("Entry %d marked as unread\n", opts.EntryID)
	return nil
}

func runDiscussionsRateEntry(ctx context.Context, client *api.Client, opts *options.DiscussionsRateEntryOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.rate-entry", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
		"entry_id":     opts.EntryID,
		"rating":       opts.Rating,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.RateEntry(ctx, ctxType, ctxID, opts.TopicID, opts.EntryID, opts.Rating); err != nil {
		logger.LogCommandError(ctx, "discussions.rate-entry", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
			"entry_id":     opts.EntryID,
		})
		return fmt.Errorf("failed to rate entry: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.rate-entry", 1)
	printInfo("Entry %d rated %d\n", opts.EntryID, opts.Rating)
	return nil
}

func runDiscussionsSubscribe(ctx context.Context, client *api.Client, opts *options.DiscussionsSubscribeOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.subscribe", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.SubscribeContext(ctx, ctxType, ctxID, opts.TopicID); err != nil {
		logger.LogCommandError(ctx, "discussions.subscribe", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.subscribe", 1)
	fmt.Printf("Subscribed to discussion %d\n", opts.TopicID)
	return nil
}

func runDiscussionsUnsubscribe(ctx context.Context, client *api.Client, opts *options.DiscussionsUnsubscribeOptions) error {
	logger := logging.NewCommandLogger(verbose)
	ctxType, ctxID, _ := opts.ContextType()
	logger.LogCommandStart(ctx, "discussions.unsubscribe", map[string]interface{}{
		"context_type": ctxType,
		"context_id":   ctxID,
		"topic_id":     opts.TopicID,
	})

	discussionsService := api.NewDiscussionsService(client)

	if err := discussionsService.UnsubscribeContext(ctx, ctxType, ctxID, opts.TopicID); err != nil {
		logger.LogCommandError(ctx, "discussions.unsubscribe", err, map[string]interface{}{
			"context_type": ctxType,
			"context_id":   ctxID,
			"topic_id":     opts.TopicID,
		})
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	logger.LogCommandComplete(ctx, "discussions.unsubscribe", 1)
	fmt.Printf("Unsubscribed from discussion %d\n", opts.TopicID)
	return nil
}
