package options

import "fmt"

// resolveDiscussionContext validates that exactly one of courseID or groupID is set
// and returns the Canvas API context type ("courses" or "groups") and the ID.
// This enforces that --course-id and --group-id are mutually exclusive.
func resolveDiscussionContext(courseID, groupID int64) (contextType string, contextID int64, err error) {
	switch {
	case courseID > 0 && groupID > 0:
		return "", 0, fmt.Errorf("--course-id and --group-id are mutually exclusive; specify only one")
	case courseID > 0:
		return "courses", courseID, nil
	case groupID > 0:
		return "groups", groupID, nil
	default:
		return "", 0, fmt.Errorf("one of --course-id or --group-id is required")
	}
}

// DiscussionsListOptions contains options for listing discussions
type DiscussionsListOptions struct {
	CourseID          int64
	GroupID           int64
	OrderBy           string
	Scope             string
	OnlyAnnouncements bool
	FilterBy          string
	SearchTerm        string
	Include           []string
}

// Validate validates the options
func (o *DiscussionsListOptions) Validate() error {
	_, _, err := resolveDiscussionContext(o.CourseID, o.GroupID)
	return err
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsListOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsGetOptions contains options for getting a discussion
type DiscussionsGetOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
	Include  []string
}

// Validate validates the options
func (o *DiscussionsGetOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsGetOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsCreateOptions contains options for creating a discussion
type DiscussionsCreateOptions struct {
	CourseID           int64
	GroupID            int64
	Title              string
	Message            string
	DiscussionType     string
	Published          bool
	DelayedPostAt      string
	AllowRating        bool
	LockAt             string
	RequireInitialPost bool
	Pinned             bool
}

// Validate validates the options
func (o *DiscussionsCreateOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsCreateOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsUpdateOptions contains options for updating a discussion
type DiscussionsUpdateOptions struct {
	CourseID           int64
	GroupID            int64
	TopicID            int64
	Title              string
	Message            string
	DiscussionType     string
	Published          bool
	DelayedPostAt      string
	AllowRating        bool
	LockAt             string
	RequireInitialPost bool
	Pinned             bool
	Locked             bool
	// Track which fields were actually set
	TitleSet              bool
	MessageSet            bool
	DiscussionTypeSet     bool
	PublishedSet          bool
	DelayedPostAtSet      bool
	AllowRatingSet        bool
	LockAtSet             bool
	RequireInitialPostSet bool
	PinnedSet             bool
	LockedSet             bool
}

// Validate validates the options
func (o *DiscussionsUpdateOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsUpdateOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsDeleteOptions contains options for deleting a discussion
type DiscussionsDeleteOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
	Force    bool
}

// Validate validates the options
func (o *DiscussionsDeleteOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsDeleteOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsEntriesOptions contains options for listing discussion entries
type DiscussionsEntriesOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
}

// Validate validates the options
func (o *DiscussionsEntriesOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsEntriesOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsPostOptions contains options for posting an entry
type DiscussionsPostOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
	Message  string
}

// Validate validates the options
func (o *DiscussionsPostOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	if o.Message == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsPostOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsReplyOptions contains options for replying to an entry
type DiscussionsReplyOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
	EntryID  int64
	Message  string
}

// Validate validates the options
func (o *DiscussionsReplyOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	if o.EntryID <= 0 {
		return fmt.Errorf("entry-id is required and must be greater than 0")
	}
	if o.Message == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsReplyOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsSubscribeOptions contains options for subscribing to a discussion
type DiscussionsSubscribeOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
}

// Validate validates the options
func (o *DiscussionsSubscribeOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsSubscribeOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsUnsubscribeOptions contains options for unsubscribing from a discussion
type DiscussionsUnsubscribeOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
}

// Validate validates the options
func (o *DiscussionsUnsubscribeOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsUnsubscribeOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsViewOptions contains options for getting the full topic view
type DiscussionsViewOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
}

// Validate validates the options
func (o *DiscussionsViewOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsViewOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsDuplicateOptions contains options for duplicating a topic
type DiscussionsDuplicateOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
}

// Validate validates the options
func (o *DiscussionsDuplicateOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsDuplicateOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsReorderOptions contains options for reordering pinned topics
type DiscussionsReorderOptions struct {
	CourseID int64
	GroupID  int64
	Order    []int64
}

// Validate validates the options
func (o *DiscussionsReorderOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if len(o.Order) == 0 {
		return fmt.Errorf("order is required (list of topic IDs)")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsReorderOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsUpdateEntryOptions contains options for updating a discussion entry
type DiscussionsUpdateEntryOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
	EntryID  int64
	Message  string
}

// Validate validates the options
func (o *DiscussionsUpdateEntryOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	if o.EntryID <= 0 {
		return fmt.Errorf("entry-id is required and must be greater than 0")
	}
	if o.Message == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsUpdateEntryOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsDeleteEntryOptions contains options for deleting a discussion entry
type DiscussionsDeleteEntryOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
	EntryID  int64
	Force    bool
}

// Validate validates the options
func (o *DiscussionsDeleteEntryOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	if o.EntryID <= 0 {
		return fmt.Errorf("entry-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsDeleteEntryOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsRepliesOptions contains options for listing entry replies
type DiscussionsRepliesOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
	EntryID  int64
}

// Validate validates the options
func (o *DiscussionsRepliesOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	if o.EntryID <= 0 {
		return fmt.Errorf("entry-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsRepliesOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsEntryListOptions contains options for fetching specific entries by ID
type DiscussionsEntryListOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
	IDs      []int64
}

// Validate validates the options
func (o *DiscussionsEntryListOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsEntryListOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsMarkTopicReadOptions contains options for marking a topic read/unread
type DiscussionsMarkTopicReadOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
}

// Validate validates the options
func (o *DiscussionsMarkTopicReadOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsMarkTopicReadOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsMarkAllTopicsReadOptions contains options for marking all topics read
type DiscussionsMarkAllTopicsReadOptions struct {
	CourseID int64
	GroupID  int64
}

// Validate validates the options
func (o *DiscussionsMarkAllTopicsReadOptions) Validate() error {
	_, _, err := resolveDiscussionContext(o.CourseID, o.GroupID)
	return err
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsMarkAllTopicsReadOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsMarkAllEntriesReadOptions contains options for marking all entries read/unread
type DiscussionsMarkAllEntriesReadOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
}

// Validate validates the options
func (o *DiscussionsMarkAllEntriesReadOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsMarkAllEntriesReadOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsMarkEntryReadOptions contains options for marking a single entry read/unread
type DiscussionsMarkEntryReadOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
	EntryID  int64
}

// Validate validates the options
func (o *DiscussionsMarkEntryReadOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	if o.EntryID <= 0 {
		return fmt.Errorf("entry-id is required and must be greater than 0")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsMarkEntryReadOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}

// DiscussionsRateEntryOptions contains options for rating a discussion entry
type DiscussionsRateEntryOptions struct {
	CourseID int64
	GroupID  int64
	TopicID  int64
	EntryID  int64
	Rating   int
}

// Validate validates the options
func (o *DiscussionsRateEntryOptions) Validate() error {
	if _, _, err := resolveDiscussionContext(o.CourseID, o.GroupID); err != nil {
		return err
	}
	if o.TopicID <= 0 {
		return fmt.Errorf("topic-id is required and must be greater than 0")
	}
	if o.EntryID <= 0 {
		return fmt.Errorf("entry-id is required and must be greater than 0")
	}
	if o.Rating != 0 && o.Rating != 1 {
		return fmt.Errorf("rating must be 0 (un-rate) or 1 (like)")
	}
	return nil
}

// ContextType returns "courses" or "groups" based on which ID is set.
func (o *DiscussionsRateEntryOptions) ContextType() (string, int64, error) {
	return resolveDiscussionContext(o.CourseID, o.GroupID)
}
