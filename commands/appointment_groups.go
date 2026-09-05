package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/chiptoe-svg/canvas-cli/commands/internal/logging"
	"github.com/chiptoe-svg/canvas-cli/commands/internal/options"
	"github.com/chiptoe-svg/canvas-cli/internal/api"
)

// appointmentGroupsCmd is the root command group for appointment groups
var appointmentGroupsCmd = &cobra.Command{
	Use:   "appointment-groups",
	Short: "Manage Canvas appointment groups",
	Long: `Manage Canvas appointment groups (office hours, sign-up slots, etc.).

Appointment groups provide a way of creating a bundle of time slots that
users can sign up for (e.g. "Office Hours" or "Meet with professor about
Final Project"). Both time slots and reservations of time slots are stored
as Calendar Events.

Examples:
  canvas appointment-groups list
  canvas appointment-groups list --scope manageable
  canvas appointment-groups get 543
  canvas appointment-groups create --context course_123 --title "Office Hours"
  canvas appointment-groups users 543
  canvas appointment-groups next`,
}

func init() {
	rootCmd.AddCommand(appointmentGroupsCmd)

	appointmentGroupsCmd.AddCommand(newAppointmentGroupListCmd())
	appointmentGroupsCmd.AddCommand(newAppointmentGroupGetCmd())
	appointmentGroupsCmd.AddCommand(newAppointmentGroupCreateCmd())
	appointmentGroupsCmd.AddCommand(newAppointmentGroupUpdateCmd())
	appointmentGroupsCmd.AddCommand(newAppointmentGroupDeleteCmd())
	appointmentGroupsCmd.AddCommand(newAppointmentGroupUsersCmd())
	appointmentGroupsCmd.AddCommand(newAppointmentGroupGroupsCmd())
	appointmentGroupsCmd.AddCommand(newAppointmentGroupNextCmd())
}

func newAppointmentGroupListCmd() *cobra.Command {
	opts := &options.AppointmentGroupListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List appointment groups",
		Long: `List appointment groups that can be reserved or managed by the current user.

Examples:
  canvas appointment-groups list
  canvas appointment-groups list --scope manageable
  canvas appointment-groups list --scope reservable --include appointments,participant_count
  canvas appointment-groups list --past`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAppointmentGroupList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Scope, "scope", "", "Scope: reservable or manageable (default: reservable)")
	cmd.Flags().StringSliceVar(&opts.ContextCodes, "context", []string{}, "Context codes to filter by (e.g. course_123)")
	cmd.Flags().BoolVar(&opts.IncludePastAppointments, "past", false, "Include past appointment groups")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional fields: appointments, child_events, participant_count, reserved_times, all_context_codes")

	return cmd
}

func newAppointmentGroupGetCmd() *cobra.Command {
	opts := &options.AppointmentGroupGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get an appointment group",
		Long: `Get details of a specific appointment group.

Examples:
  canvas appointment-groups get 543
  canvas appointment-groups get 543 --include appointments,child_events`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.GroupID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAppointmentGroupGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional fields: appointments, child_events, all_context_codes")

	return cmd
}

func newAppointmentGroupCreateCmd() *cobra.Command {
	opts := &options.AppointmentGroupCreateOptions{}
	var contextCodesRaw []string
	var subContextCodesRaw []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an appointment group",
		Long: `Create a new appointment group with time slots.

Examples:
  canvas appointment-groups create --context course_123 --title "Office Hours"
  canvas appointment-groups create --context course_123 --title "Project Reviews" --publish --participants-per-slot 1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ContextCodes = contextCodesRaw
			opts.SubContextCodes = subContextCodesRaw

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAppointmentGroupCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&contextCodesRaw, "context", []string{}, "Context codes (e.g. course_123) — required")
	cmd.Flags().StringSliceVar(&subContextCodesRaw, "sub-context", []string{}, "Sub-context codes (sections or group category)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Title (required)")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Description")
	cmd.Flags().StringVar(&opts.LocationName, "location", "", "Location name")
	cmd.Flags().StringVar(&opts.LocationAddress, "address", "", "Location address")
	cmd.Flags().BoolVar(&opts.Publish, "publish", false, "Publish immediately")
	cmd.Flags().IntVar(&opts.ParticipantsPerAppointment, "participants-per-slot", 0, "Max participants per slot")
	cmd.Flags().IntVar(&opts.MinAppointmentsPerParticipant, "min-slots", 0, "Min slots per participant")
	cmd.Flags().IntVar(&opts.MaxAppointmentsPerParticipant, "max-slots", 0, "Max slots per participant")
	cmd.Flags().StringVar(&opts.ParticipantVisibility, "visibility", "", "Participant visibility: private or protected")
	cmd.Flags().BoolVar(&opts.AllowObserverSignup, "allow-observer", false, "Allow observer sign-up")
	mustMarkRequired(cmd, "title")

	return cmd
}

func newAppointmentGroupUpdateCmd() *cobra.Command {
	opts := &options.AppointmentGroupUpdateOptions{}
	var contextCodesRaw []string
	var allowObserver bool
	var allowObserverSet bool

	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update an appointment group",
		Long: `Update an existing appointment group.

Examples:
  canvas appointment-groups update 543 --title "Updated Title"
  canvas appointment-groups update 543 --publish`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.GroupID = id
			opts.ContextCodes = contextCodesRaw

			if allowObserverSet {
				opts.AllowObserverSignup = &allowObserver
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAppointmentGroupUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&contextCodesRaw, "context", []string{}, "Context codes (required per spec)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Updated title")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Updated description")
	cmd.Flags().StringVar(&opts.LocationName, "location", "", "Updated location name")
	cmd.Flags().StringVar(&opts.LocationAddress, "address", "", "Updated location address")
	cmd.Flags().BoolVar(&opts.Publish, "publish", false, "Publish the appointment group")
	cmd.Flags().IntVar(&opts.ParticipantsPerAppointment, "participants-per-slot", 0, "Max participants per slot")
	cmd.Flags().IntVar(&opts.MinAppointmentsPerParticipant, "min-slots", 0, "Min slots per participant")
	cmd.Flags().IntVar(&opts.MaxAppointmentsPerParticipant, "max-slots", 0, "Max slots per participant")
	cmd.Flags().StringVar(&opts.ParticipantVisibility, "visibility", "", "Participant visibility: private or protected")
	cmd.Flags().BoolVar(&allowObserver, "allow-observer", false, "Allow observer sign-up")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		allowObserverSet = cmd.Flags().Changed("allow-observer")
		return nil
	}

	return cmd
}

func newAppointmentGroupDeleteCmd() *cobra.Command {
	opts := &options.AppointmentGroupDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <group-id>",
		Short: "Delete an appointment group",
		Long: `Delete an appointment group and all associated time slots and reservations.

Examples:
  canvas appointment-groups delete 543
  canvas appointment-groups delete 543 --reason "Group cancelled" --force`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.GroupID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAppointmentGroupDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.CancelReason, "reason", "", "Cancellation reason")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func newAppointmentGroupUsersCmd() *cobra.Command {
	opts := &options.AppointmentGroupUsersOptions{}

	cmd := &cobra.Command{
		Use:   "users <group-id>",
		Short: "List appointment group participants",
		Long: `List users participating (or eligible to participate) in an appointment group.

Only returns results for groups with participant_type "User".

Examples:
  canvas appointment-groups users 543
  canvas appointment-groups users 543 --status registered`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.GroupID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAppointmentGroupUsers(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.RegistrationStatus, "status", "", "Registration status: all, registered (default: all)")

	return cmd
}

func newAppointmentGroupGroupsCmd() *cobra.Command {
	opts := &options.AppointmentGroupGroupsOptions{}

	cmd := &cobra.Command{
		Use:   "groups <group-id>",
		Short: "List appointment group student groups",
		Long: `List student groups participating in an appointment group.

Only returns results for groups with participant_type "Group".

Examples:
  canvas appointment-groups groups 543
  canvas appointment-groups groups 543 --status registered`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.GroupID = id

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAppointmentGroupGroups(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.RegistrationStatus, "status", "", "Registration status: all, registered (default: all)")

	return cmd
}

func newAppointmentGroupNextCmd() *cobra.Command {
	opts := &options.AppointmentGroupNextOptions{}
	var groupIDsRaw []string

	cmd := &cobra.Command{
		Use:   "next",
		Short: "Get next available appointment",
		Long: `Return the next appointment available to sign up for.

If no future appointments are available, an empty result is returned.

Examples:
  canvas appointment-groups next
  canvas appointment-groups next --group-ids 1,2,3`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse comma-separated or repeated IDs
			for _, raw := range groupIDsRaw {
				for _, part := range strings.Split(raw, ",") {
					part = strings.TrimSpace(part)
					if part == "" {
						continue
					}
					id, err := strconv.ParseInt(part, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid group ID: %s", part)
					}
					opts.AppointmentGroupIDs = append(opts.AppointmentGroupIDs, id)
				}
			}

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAppointmentGroupNext(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&groupIDsRaw, "group-ids", []string{}, "Appointment group IDs to search")

	return cmd
}

// ---- Run functions ----

func runAppointmentGroupList(ctx context.Context, client *api.Client, opts *options.AppointmentGroupListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "appointment-groups.list", map[string]interface{}{
		"scope": opts.Scope,
		"past":  opts.IncludePastAppointments,
	})

	svc := api.NewAppointmentGroupsService(client)

	groups, err := svc.List(ctx, &api.ListAppointmentGroupsOptions{
		Scope:                   opts.Scope,
		ContextCodes:            opts.ContextCodes,
		IncludePastAppointments: opts.IncludePastAppointments,
		Include:                 opts.Include,
	})
	if err != nil {
		logger.LogCommandError(ctx, "appointment-groups.list", err, nil)
		return fmt.Errorf("failed to list appointment groups: %w", err)
	}

	printVerbose("Found %d appointment groups:\n\n", len(groups))
	logger.LogCommandComplete(ctx, "appointment-groups.list", len(groups))
	return formatEmptyOrOutput(groups, "No appointment groups found")
}

func runAppointmentGroupGet(ctx context.Context, client *api.Client, opts *options.AppointmentGroupGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "appointment-groups.get", map[string]interface{}{"group_id": opts.GroupID})

	svc := api.NewAppointmentGroupsService(client)

	group, err := svc.Get(ctx, opts.GroupID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "appointment-groups.get", err, map[string]interface{}{"group_id": opts.GroupID})
		return fmt.Errorf("failed to get appointment group: %w", err)
	}

	logger.LogCommandComplete(ctx, "appointment-groups.get", 1)
	return formatOutput(group, nil)
}

func runAppointmentGroupCreate(ctx context.Context, client *api.Client, opts *options.AppointmentGroupCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "appointment-groups.create", map[string]interface{}{
		"title":         opts.Title,
		"context_codes": opts.ContextCodes,
	})

	svc := api.NewAppointmentGroupsService(client)

	group, err := svc.Create(ctx, &api.CreateAppointmentGroupParams{
		ContextCodes:                  opts.ContextCodes,
		SubContextCodes:               opts.SubContextCodes,
		Title:                         opts.Title,
		Description:                   opts.Description,
		LocationName:                  opts.LocationName,
		LocationAddress:               opts.LocationAddress,
		Publish:                       opts.Publish,
		ParticipantsPerAppointment:    opts.ParticipantsPerAppointment,
		MinAppointmentsPerParticipant: opts.MinAppointmentsPerParticipant,
		MaxAppointmentsPerParticipant: opts.MaxAppointmentsPerParticipant,
		ParticipantVisibility:         opts.ParticipantVisibility,
		AllowObserverSignup:           opts.AllowObserverSignup,
	})
	if err != nil {
		logger.LogCommandError(ctx, "appointment-groups.create", err, nil)
		return fmt.Errorf("failed to create appointment group: %w", err)
	}

	logger.LogCommandComplete(ctx, "appointment-groups.create", 1)
	return formatSuccessOutput(group, "Appointment group created successfully!")
}

func runAppointmentGroupUpdate(ctx context.Context, client *api.Client, opts *options.AppointmentGroupUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "appointment-groups.update", map[string]interface{}{"group_id": opts.GroupID})

	svc := api.NewAppointmentGroupsService(client)

	group, err := svc.Update(ctx, opts.GroupID, &api.UpdateAppointmentGroupParams{
		ContextCodes:                  opts.ContextCodes,
		SubContextCodes:               opts.SubContextCodes,
		Title:                         opts.Title,
		Description:                   opts.Description,
		LocationName:                  opts.LocationName,
		LocationAddress:               opts.LocationAddress,
		Publish:                       opts.Publish,
		ParticipantsPerAppointment:    opts.ParticipantsPerAppointment,
		MinAppointmentsPerParticipant: opts.MinAppointmentsPerParticipant,
		MaxAppointmentsPerParticipant: opts.MaxAppointmentsPerParticipant,
		ParticipantVisibility:         opts.ParticipantVisibility,
		AllowObserverSignup:           opts.AllowObserverSignup,
	})
	if err != nil {
		logger.LogCommandError(ctx, "appointment-groups.update", err, map[string]interface{}{"group_id": opts.GroupID})
		return fmt.Errorf("failed to update appointment group: %w", err)
	}

	logger.LogCommandComplete(ctx, "appointment-groups.update", 1)
	return formatSuccessOutput(group, "Appointment group updated successfully!")
}

func runAppointmentGroupDelete(ctx context.Context, client *api.Client, opts *options.AppointmentGroupDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "appointment-groups.delete", map[string]interface{}{
		"group_id": opts.GroupID,
		"force":    opts.Force,
	})

	confirmed, err := confirmDelete("appointment group", opts.GroupID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "appointment-groups.delete", err, map[string]interface{}{"group_id": opts.GroupID})
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		logger.LogCommandComplete(ctx, "appointment-groups.delete", 0)
		return nil
	}

	svc := api.NewAppointmentGroupsService(client)

	group, err := svc.Delete(ctx, opts.GroupID, opts.CancelReason)
	if err != nil {
		logger.LogCommandError(ctx, "appointment-groups.delete", err, map[string]interface{}{"group_id": opts.GroupID})
		return fmt.Errorf("failed to delete appointment group: %w", err)
	}

	printInfo("Appointment group %d deleted successfully\n", opts.GroupID)
	logger.LogCommandComplete(ctx, "appointment-groups.delete", 1)

	if group != nil && group.ID != 0 {
		return formatOutput(group, nil)
	}
	return nil
}

func runAppointmentGroupUsers(ctx context.Context, client *api.Client, opts *options.AppointmentGroupUsersOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "appointment-groups.users", map[string]interface{}{
		"group_id": opts.GroupID,
		"status":   opts.RegistrationStatus,
	})

	svc := api.NewAppointmentGroupsService(client)

	users, err := svc.ListUsers(ctx, opts.GroupID, opts.RegistrationStatus)
	if err != nil {
		logger.LogCommandError(ctx, "appointment-groups.users", err, map[string]interface{}{"group_id": opts.GroupID})
		return fmt.Errorf("failed to list appointment group users: %w", err)
	}

	printVerbose("Found %d users:\n\n", len(users))
	logger.LogCommandComplete(ctx, "appointment-groups.users", len(users))
	return formatEmptyOrOutput(users, "No users found")
}

func runAppointmentGroupGroups(ctx context.Context, client *api.Client, opts *options.AppointmentGroupGroupsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "appointment-groups.groups", map[string]interface{}{
		"group_id": opts.GroupID,
		"status":   opts.RegistrationStatus,
	})

	svc := api.NewAppointmentGroupsService(client)

	groups, err := svc.ListGroups(ctx, opts.GroupID, opts.RegistrationStatus)
	if err != nil {
		logger.LogCommandError(ctx, "appointment-groups.groups", err, map[string]interface{}{"group_id": opts.GroupID})
		return fmt.Errorf("failed to list appointment group student groups: %w", err)
	}

	printVerbose("Found %d student groups:\n\n", len(groups))
	logger.LogCommandComplete(ctx, "appointment-groups.groups", len(groups))
	return formatEmptyOrOutput(groups, "No student groups found")
}

func runAppointmentGroupNext(ctx context.Context, client *api.Client, opts *options.AppointmentGroupNextOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "appointment-groups.next", map[string]interface{}{
		"group_ids": opts.AppointmentGroupIDs,
	})

	svc := api.NewAppointmentGroupsService(client)

	events, err := svc.NextAppointment(ctx, opts.AppointmentGroupIDs)
	if err != nil {
		logger.LogCommandError(ctx, "appointment-groups.next", err, nil)
		return fmt.Errorf("failed to get next appointment: %w", err)
	}

	logger.LogCommandComplete(ctx, "appointment-groups.next", len(events))
	return formatEmptyOrOutput(events, "No upcoming appointments available")
}
