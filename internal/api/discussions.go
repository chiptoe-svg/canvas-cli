package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// DiscussionTopic represents a Canvas discussion topic or announcement
type DiscussionTopic struct {
	ID                      int64            `json:"id"`
	Title                   string           `json:"title"`
	Message                 string           `json:"message"`
	HTMLURL                 string           `json:"html_url"`
	PostedAt                *time.Time       `json:"posted_at,omitempty"`
	LastReplyAt             *time.Time       `json:"last_reply_at,omitempty"`
	RequireInitialPost      bool             `json:"require_initial_post"`
	UserCanSeePosts         bool             `json:"user_can_see_posts"`
	DiscussionSubentryCount int              `json:"discussion_subentry_count"`
	ReadState               string           `json:"read_state"`
	UnreadCount             int              `json:"unread_count"`
	Subscribed              bool             `json:"subscribed"`
	SubscriptionHold        string           `json:"subscription_hold,omitempty"`
	AssignmentID            *int64           `json:"assignment_id,omitempty"`
	DelayedPostAt           *time.Time       `json:"delayed_post_at,omitempty"`
	Published               bool             `json:"published"`
	LockAt                  *time.Time       `json:"lock_at,omitempty"`
	Locked                  bool             `json:"locked"`
	Pinned                  bool             `json:"pinned"`
	LockedForUser           bool             `json:"locked_for_user"`
	LockInfo                *LockInfo        `json:"lock_info,omitempty"`
	LockExplanation         string           `json:"lock_explanation,omitempty"`
	UserName                string           `json:"user_name,omitempty"`
	RootTopicID             *int64           `json:"root_topic_id,omitempty"`
	PodcastURL              string           `json:"podcast_url,omitempty"`
	DiscussionType          string           `json:"discussion_type"`
	GroupCategoryID         *int64           `json:"group_category_id,omitempty"`
	Attachments             []FileAttachment `json:"attachments,omitempty"`
	Permissions             map[string]bool  `json:"permissions,omitempty"`
	AllowRating             bool             `json:"allow_rating"`
	OnlyGradersCanRate      bool             `json:"only_graders_can_rate"`
	SortByRating            bool             `json:"sort_by_rating"`
	ContextCode             string           `json:"context_code,omitempty"`
	Author                  *User            `json:"author,omitempty"`
	IsAnnouncement          bool             `json:"is_announcement,omitempty"`
}

// FileAttachment represents a file attachment
type FileAttachment struct {
	ContentType string `json:"content-type"`
	URL         string `json:"url"`
	Filename    string `json:"filename"`
	DisplayName string `json:"display_name"`
}

// DiscussionEntry represents an entry in a discussion
type DiscussionEntry struct {
	ID              int64             `json:"id"`
	UserID          int64             `json:"user_id"`
	ParentID        *int64            `json:"parent_id,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Message         string            `json:"message"`
	Rating          int               `json:"rating"`
	RatingSum       int               `json:"rating_sum"`
	ReadState       string            `json:"read_state"`
	ForcedReadState bool              `json:"forced_read_state"`
	User            *User             `json:"user,omitempty"`
	Replies         []DiscussionEntry `json:"replies,omitempty"`
}

// DiscussionView is the full cached view of a discussion topic returned by /view
type DiscussionView struct {
	UnreadEntries []int64           `json:"unread_entries"`
	EntryRatings  map[string]int    `json:"entry_ratings"`
	ForcedEntries []int64           `json:"forced_entries"`
	Participants  []DiscussionUser  `json:"participants"`
	View          []DiscussionEntry `json:"view"`
	NewEntries    []DiscussionEntry `json:"new_entries,omitempty"`
}

// DiscussionUser is the lightweight participant representation inside DiscussionView
type DiscussionUser struct {
	ID             int64  `json:"id"`
	DisplayName    string `json:"display_name"`
	AvatarImageURL string `json:"avatar_image_url"`
	HTMLURL        string `json:"html_url"`
}

// DiscussionSummary is an AI-generated summary for a discussion topic
type DiscussionSummary struct {
	ID        int64       `json:"id"`
	UserInput string      `json:"userInput,omitempty"`
	Text      string      `json:"text"`
	Usage     interface{} `json:"usage,omitempty"`
}

// SummaryFeedbackResult is returned after posting summary feedback
type SummaryFeedbackResult struct {
	Liked    bool `json:"liked"`
	Disliked bool `json:"disliked"`
}

// DiscussionsService handles discussion-related API calls
type DiscussionsService struct {
	client *Client
}

// NewDiscussionsService creates a new discussions service
func NewDiscussionsService(client *Client) *DiscussionsService {
	return &DiscussionsService{client: client}
}

// discussionContextPath returns the API prefix for course or group context.
// contextType must be "courses" or "groups".
func discussionContextPath(contextType string, contextID int64) string {
	return fmt.Sprintf("/api/v1/%s/%d/discussion_topics", contextType, contextID)
}

// ListDiscussionsOptions holds options for listing discussions
type ListDiscussionsOptions struct {
	Include           []string // all_dates, sections, sections_user_count, overrides
	OrderBy           string   // position, recent_activity, title
	Scope             string   // locked, unlocked, pinned, unpinned
	OnlyAnnouncements bool
	FilterBy          string // all, unread
	SearchTerm        string
	Page              int
	PerPage           int
}

// List retrieves all discussion topics for a course (course-context only, for backward compat).
// For group discussions use ListContext.
func (s *DiscussionsService) List(ctx context.Context, courseID int64, opts *ListDiscussionsOptions) ([]DiscussionTopic, error) {
	return s.ListContext(ctx, "courses", courseID, opts)
}

// ListContext retrieves all discussion topics for a course or group.
// contextType must be "courses" or "groups".
func (s *DiscussionsService) ListContext(ctx context.Context, contextType string, contextID int64, opts *ListDiscussionsOptions) ([]DiscussionTopic, error) {
	path := discussionContextPath(contextType, contextID)

	if opts != nil {
		query := url.Values{}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if opts.OrderBy != "" {
			query.Add("order_by", opts.OrderBy)
		}

		if opts.Scope != "" {
			query.Add("scope", opts.Scope)
		}

		if opts.OnlyAnnouncements {
			query.Add("only_announcements", "true")
		}

		if opts.FilterBy != "" {
			query.Add("filter_by", opts.FilterBy)
		}

		if opts.SearchTerm != "" {
			query.Add("search_term", opts.SearchTerm)
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var topics []DiscussionTopic
	if err := s.client.GetAllPages(ctx, path, &topics); err != nil {
		return nil, err
	}

	return topics, nil
}

// Get retrieves a single discussion topic under a course (for backward compat).
// For group discussions use GetContext.
func (s *DiscussionsService) Get(ctx context.Context, courseID, topicID int64, include []string) (*DiscussionTopic, error) {
	return s.GetContext(ctx, "courses", courseID, topicID, include)
}

// GetContext retrieves a single discussion topic under a course or group.
func (s *DiscussionsService) GetContext(ctx context.Context, contextType string, contextID, topicID int64, include []string) (*DiscussionTopic, error) {
	path := fmt.Sprintf("%s/%d", discussionContextPath(contextType, contextID), topicID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var topic DiscussionTopic
	if err := s.client.GetJSON(ctx, path, &topic); err != nil {
		return nil, err
	}

	return &topic, nil
}

// GetView returns the full cached view of a discussion topic (threaded entries tree).
func (s *DiscussionsService) GetView(ctx context.Context, contextType string, contextID, topicID int64) (*DiscussionView, error) {
	path := fmt.Sprintf("%s/%d/view", discussionContextPath(contextType, contextID), topicID)

	var view DiscussionView
	if err := s.client.GetJSON(ctx, path, &view); err != nil {
		return nil, err
	}

	return &view, nil
}

// CreateDiscussionParams holds parameters for creating a discussion
type CreateDiscussionParams struct {
	Title                  string
	Message                string
	DiscussionType         string // side_comment, threaded, not_threaded
	Published              bool
	DelayedPostAt          string
	AllowRating            bool
	LockAt                 string
	PodcastEnabled         bool
	PodcastHasStudentPosts bool
	RequireInitialPost     bool
	IsAnnouncement         bool
	Pinned                 bool
	PositionAfter          string
	GroupCategoryID        int64
	OnlyGradersCanRate     bool
	SpecificSections       string
}

// Create creates a new discussion topic under a course (for backward compat).
// For group discussions use CreateContext.
func (s *DiscussionsService) Create(ctx context.Context, courseID int64, params *CreateDiscussionParams) (*DiscussionTopic, error) {
	return s.CreateContext(ctx, "courses", courseID, params)
}

// CreateContext creates a new discussion topic under a course or group.
func (s *DiscussionsService) CreateContext(ctx context.Context, contextType string, contextID int64, params *CreateDiscussionParams) (*DiscussionTopic, error) {
	path := discussionContextPath(contextType, contextID)

	body := make(map[string]interface{})

	if params.Title != "" {
		body["title"] = params.Title
	}

	if params.Message != "" {
		body["message"] = params.Message
	}

	if params.DiscussionType != "" {
		body["discussion_type"] = params.DiscussionType
	}

	body["published"] = params.Published

	if params.DelayedPostAt != "" {
		body["delayed_post_at"] = params.DelayedPostAt
	}

	if params.AllowRating {
		body["allow_rating"] = true
	}

	if params.LockAt != "" {
		body["lock_at"] = params.LockAt
	}

	if params.PodcastEnabled {
		body["podcast_enabled"] = true
	}

	if params.PodcastHasStudentPosts {
		body["podcast_has_student_posts"] = true
	}

	if params.RequireInitialPost {
		body["require_initial_post"] = true
	}

	if params.IsAnnouncement {
		body["is_announcement"] = true
	}

	if params.Pinned {
		body["pinned"] = true
	}

	if params.PositionAfter != "" {
		body["position_after"] = params.PositionAfter
	}

	if params.GroupCategoryID > 0 {
		body["group_category_id"] = params.GroupCategoryID
	}

	if params.OnlyGradersCanRate {
		body["only_graders_can_rate"] = true
	}

	if params.SpecificSections != "" {
		body["specific_sections"] = params.SpecificSections
	}

	var topic DiscussionTopic
	if err := s.client.PostJSON(ctx, path, body, &topic); err != nil {
		return nil, err
	}

	return &topic, nil
}

// UpdateDiscussionParams holds parameters for updating a discussion
type UpdateDiscussionParams struct {
	Title              *string
	Message            *string
	DiscussionType     *string
	Published          *bool
	DelayedPostAt      *string
	AllowRating        *bool
	LockAt             *string
	PodcastEnabled     *bool
	RequireInitialPost *bool
	Pinned             *bool
	Locked             *bool
}

// Update updates an existing discussion topic under a course (for backward compat).
// For group discussions use UpdateContext.
func (s *DiscussionsService) Update(ctx context.Context, courseID, topicID int64, params *UpdateDiscussionParams) (*DiscussionTopic, error) {
	return s.UpdateContext(ctx, "courses", courseID, topicID, params)
}

// UpdateContext updates an existing discussion topic under a course or group.
func (s *DiscussionsService) UpdateContext(ctx context.Context, contextType string, contextID, topicID int64, params *UpdateDiscussionParams) (*DiscussionTopic, error) {
	path := fmt.Sprintf("%s/%d", discussionContextPath(contextType, contextID), topicID)

	body := make(map[string]interface{})

	if params.Title != nil {
		body["title"] = *params.Title
	}

	if params.Message != nil {
		body["message"] = *params.Message
	}

	if params.DiscussionType != nil {
		body["discussion_type"] = *params.DiscussionType
	}

	if params.Published != nil {
		body["published"] = *params.Published
	}

	if params.DelayedPostAt != nil {
		body["delayed_post_at"] = *params.DelayedPostAt
	}

	if params.AllowRating != nil {
		body["allow_rating"] = *params.AllowRating
	}

	if params.LockAt != nil {
		body["lock_at"] = *params.LockAt
	}

	if params.PodcastEnabled != nil {
		body["podcast_enabled"] = *params.PodcastEnabled
	}

	if params.RequireInitialPost != nil {
		body["require_initial_post"] = *params.RequireInitialPost
	}

	if params.Pinned != nil {
		body["pinned"] = *params.Pinned
	}

	if params.Locked != nil {
		body["locked"] = *params.Locked
	}

	var topic DiscussionTopic
	if err := s.client.PutJSON(ctx, path, body, &topic); err != nil {
		return nil, err
	}

	return &topic, nil
}

// Delete deletes a discussion topic under a course (for backward compat).
// For group discussions use DeleteContext.
func (s *DiscussionsService) Delete(ctx context.Context, courseID, topicID int64) error {
	return s.DeleteContext(ctx, "courses", courseID, topicID)
}

// DeleteContext deletes a discussion topic under a course or group.
func (s *DiscussionsService) DeleteContext(ctx context.Context, contextType string, contextID, topicID int64) error {
	path := fmt.Sprintf("%s/%d", discussionContextPath(contextType, contextID), topicID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// Duplicate duplicates a discussion topic under a course or group.
func (s *DiscussionsService) Duplicate(ctx context.Context, contextType string, contextID, topicID int64) (*DiscussionTopic, error) {
	path := fmt.Sprintf("%s/%d/duplicate", discussionContextPath(contextType, contextID), topicID)

	var topic DiscussionTopic
	if err := s.client.PostJSON(ctx, path, nil, &topic); err != nil {
		return nil, err
	}

	return &topic, nil
}

// Reorder reorders pinned discussion topics.
// order is the slice of topic IDs in the desired pinned order.
func (s *DiscussionsService) Reorder(ctx context.Context, contextType string, contextID int64, order []int64) error {
	path := fmt.Sprintf("%s/reorder", discussionContextPath(contextType, contextID))

	body := map[string]interface{}{
		"order": order,
	}

	return s.client.PostJSON(ctx, path, body, nil)
}

// ListEntries retrieves all top-level entries for a discussion topic under a course (backward compat).
// For group discussions use ListEntriesContext.
func (s *DiscussionsService) ListEntries(ctx context.Context, courseID, topicID int64) ([]DiscussionEntry, error) {
	return s.ListEntriesContext(ctx, "courses", courseID, topicID)
}

// ListEntriesContext retrieves all top-level entries for a discussion topic.
func (s *DiscussionsService) ListEntriesContext(ctx context.Context, contextType string, contextID, topicID int64) ([]DiscussionEntry, error) {
	path := fmt.Sprintf("%s/%d/entries", discussionContextPath(contextType, contextID), topicID)

	var entries []DiscussionEntry
	if err := s.client.GetAllPages(ctx, path, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetEntryList retrieves specific discussion entries by ID.
func (s *DiscussionsService) GetEntryList(ctx context.Context, contextType string, contextID, topicID int64, ids []int64) ([]DiscussionEntry, error) {
	path := fmt.Sprintf("%s/%d/entry_list", discussionContextPath(contextType, contextID), topicID)

	if len(ids) > 0 {
		query := url.Values{}
		for _, id := range ids {
			query.Add("ids[]", strconv.FormatInt(id, 10))
		}
		path += "?" + query.Encode()
	}

	var entries []DiscussionEntry
	if err := s.client.GetAllPages(ctx, path, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// PostEntry posts a new top-level entry to a discussion topic under a course (backward compat).
// For group discussions use PostEntryContext.
func (s *DiscussionsService) PostEntry(ctx context.Context, courseID, topicID int64, message string) (*DiscussionEntry, error) {
	return s.PostEntryContext(ctx, "courses", courseID, topicID, message)
}

// PostEntryContext posts a new top-level entry to a discussion topic.
func (s *DiscussionsService) PostEntryContext(ctx context.Context, contextType string, contextID, topicID int64, message string) (*DiscussionEntry, error) {
	path := fmt.Sprintf("%s/%d/entries", discussionContextPath(contextType, contextID), topicID)

	body := map[string]interface{}{
		"message": message,
	}

	var entry DiscussionEntry
	if err := s.client.PostJSON(ctx, path, body, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// UpdateEntry updates the message of an existing discussion entry.
func (s *DiscussionsService) UpdateEntry(ctx context.Context, contextType string, contextID, topicID, entryID int64, message string) (*DiscussionEntry, error) {
	path := fmt.Sprintf("%s/%d/entries/%d", discussionContextPath(contextType, contextID), topicID, entryID)

	body := map[string]interface{}{
		"message": message,
	}

	var entry DiscussionEntry
	if err := s.client.PutJSON(ctx, path, body, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// DeleteEntry deletes a discussion entry.
func (s *DiscussionsService) DeleteEntry(ctx context.Context, contextType string, contextID, topicID, entryID int64) error {
	path := fmt.Sprintf("%s/%d/entries/%d", discussionContextPath(contextType, contextID), topicID, entryID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// PostReply posts a reply to an entry under a course (backward compat).
// For group discussions use PostReplyContext.
func (s *DiscussionsService) PostReply(ctx context.Context, courseID, topicID, entryID int64, message string) (*DiscussionEntry, error) {
	return s.PostReplyContext(ctx, "courses", courseID, topicID, entryID, message)
}

// PostReplyContext posts a reply to an entry in a discussion topic.
func (s *DiscussionsService) PostReplyContext(ctx context.Context, contextType string, contextID, topicID, entryID int64, message string) (*DiscussionEntry, error) {
	path := fmt.Sprintf("%s/%d/entries/%d/replies", discussionContextPath(contextType, contextID), topicID, entryID)

	body := map[string]interface{}{
		"message": message,
	}

	var entry DiscussionEntry
	if err := s.client.PostJSON(ctx, path, body, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// ListReplies retrieves the replies to a top-level entry.
func (s *DiscussionsService) ListReplies(ctx context.Context, contextType string, contextID, topicID, entryID int64) ([]DiscussionEntry, error) {
	path := fmt.Sprintf("%s/%d/entries/%d/replies", discussionContextPath(contextType, contextID), topicID, entryID)

	var entries []DiscussionEntry
	if err := s.client.GetAllPages(ctx, path, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// MarkTopicRead marks a topic's initial message as read under a course (backward compat).
// For group discussions use MarkTopicReadContext.
func (s *DiscussionsService) MarkTopicRead(ctx context.Context, courseID, topicID int64) error {
	return s.MarkTopicReadContext(ctx, "courses", courseID, topicID)
}

// MarkTopicReadContext marks a topic's initial message as read.
func (s *DiscussionsService) MarkTopicReadContext(ctx context.Context, contextType string, contextID, topicID int64) error {
	path := fmt.Sprintf("%s/%d/read", discussionContextPath(contextType, contextID), topicID)
	return s.client.PutJSON(ctx, path, nil, nil)
}

// MarkTopicUnread marks a topic's initial message as unread under a course (backward compat).
// For group discussions use MarkTopicUnreadContext.
func (s *DiscussionsService) MarkTopicUnread(ctx context.Context, courseID, topicID int64) error {
	return s.MarkTopicUnreadContext(ctx, "courses", courseID, topicID)
}

// MarkTopicUnreadContext marks a topic's initial message as unread.
func (s *DiscussionsService) MarkTopicUnreadContext(ctx context.Context, contextType string, contextID, topicID int64) error {
	path := fmt.Sprintf("%s/%d/read", discussionContextPath(contextType, contextID), topicID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// MarkAllTopicsRead marks all topics in the context as read.
func (s *DiscussionsService) MarkAllTopicsRead(ctx context.Context, contextType string, contextID int64) error {
	path := fmt.Sprintf("%s/read_all", discussionContextPath(contextType, contextID))
	return s.client.PutJSON(ctx, path, nil, nil)
}

// MarkAllEntriesRead marks all entries in a topic as read.
func (s *DiscussionsService) MarkAllEntriesRead(ctx context.Context, contextType string, contextID, topicID int64) error {
	path := fmt.Sprintf("%s/%d/read_all", discussionContextPath(contextType, contextID), topicID)
	return s.client.PutJSON(ctx, path, nil, nil)
}

// MarkAllEntriesUnread marks all entries in a topic as unread.
func (s *DiscussionsService) MarkAllEntriesUnread(ctx context.Context, contextType string, contextID, topicID int64) error {
	path := fmt.Sprintf("%s/%d/read_all", discussionContextPath(contextType, contextID), topicID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// MarkEntryRead marks a single discussion entry as read.
func (s *DiscussionsService) MarkEntryRead(ctx context.Context, contextType string, contextID, topicID, entryID int64) error {
	path := fmt.Sprintf("%s/%d/entries/%d/read", discussionContextPath(contextType, contextID), topicID, entryID)
	return s.client.PutJSON(ctx, path, nil, nil)
}

// MarkEntryUnread marks a single discussion entry as unread.
func (s *DiscussionsService) MarkEntryUnread(ctx context.Context, contextType string, contextID, topicID, entryID int64) error {
	path := fmt.Sprintf("%s/%d/entries/%d/read", discussionContextPath(contextType, contextID), topicID, entryID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// RateEntry rates a discussion entry. rating must be 0 (un-rate) or 1 (like).
func (s *DiscussionsService) RateEntry(ctx context.Context, contextType string, contextID, topicID, entryID int64, rating int) error {
	path := fmt.Sprintf("%s/%d/entries/%d/rating", discussionContextPath(contextType, contextID), topicID, entryID)

	body := map[string]interface{}{
		"rating": rating,
	}

	return s.client.PostJSON(ctx, path, body, nil)
}

// Subscribe subscribes the current user to a topic under a course (backward compat).
// For group discussions use SubscribeContext.
func (s *DiscussionsService) Subscribe(ctx context.Context, courseID, topicID int64) error {
	return s.SubscribeContext(ctx, "courses", courseID, topicID)
}

// SubscribeContext subscribes the current user to a topic.
func (s *DiscussionsService) SubscribeContext(ctx context.Context, contextType string, contextID, topicID int64) error {
	path := fmt.Sprintf("%s/%d/subscribed", discussionContextPath(contextType, contextID), topicID)
	return s.client.PutJSON(ctx, path, nil, nil)
}

// Unsubscribe unsubscribes the current user from a topic under a course (backward compat).
// For group discussions use UnsubscribeContext.
func (s *DiscussionsService) Unsubscribe(ctx context.Context, courseID, topicID int64) error {
	return s.UnsubscribeContext(ctx, "courses", courseID, topicID)
}

// UnsubscribeContext unsubscribes the current user from a topic.
func (s *DiscussionsService) UnsubscribeContext(ctx context.Context, contextType string, contextID, topicID int64) error {
	path := fmt.Sprintf("%s/%d/subscribed", discussionContextPath(contextType, contextID), topicID)
	_, err := s.client.Delete(ctx, path)
	return err
}

// GetSummary returns the last AI-generated summary for a discussion topic.
func (s *DiscussionsService) GetSummary(ctx context.Context, contextType string, contextID, topicID int64) (*DiscussionSummary, error) {
	path := fmt.Sprintf("%s/%d/summaries", discussionContextPath(contextType, contextID), topicID)

	var summary DiscussionSummary
	if err := s.client.GetJSON(ctx, path, &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

// CreateSummary generates (or returns a cached) AI summary for a discussion topic.
func (s *DiscussionsService) CreateSummary(ctx context.Context, contextType string, contextID, topicID int64, userInput string) (*DiscussionSummary, error) {
	path := fmt.Sprintf("%s/%d/summaries", discussionContextPath(contextType, contextID), topicID)

	var body map[string]interface{}
	if userInput != "" {
		body = map[string]interface{}{"userInput": userInput}
	}

	var summary DiscussionSummary
	if err := s.client.PostJSON(ctx, path, body, &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

// DisableSummary disables AI summaries for a discussion topic.
func (s *DiscussionsService) DisableSummary(ctx context.Context, contextType string, contextID, topicID int64) error {
	path := fmt.Sprintf("%s/%d/summaries/disable", discussionContextPath(contextType, contextID), topicID)
	return s.client.PutJSON(ctx, path, nil, nil)
}

// SummaryFeedback posts feedback on a specific summary.
// action must be one of: "seen", "like", "dislike", "reset_like", "regenerate", "disable_summary".
func (s *DiscussionsService) SummaryFeedback(ctx context.Context, contextType string, contextID, topicID, summaryID int64, action string) (*SummaryFeedbackResult, error) {
	path := fmt.Sprintf("%s/%d/summaries/%d/feedback", discussionContextPath(contextType, contextID), topicID, summaryID)

	body := map[string]interface{}{
		"_action": action,
	}

	var result SummaryFeedbackResult
	if err := s.client.PostJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AnnouncementsService handles announcement-specific API calls
type AnnouncementsService struct {
	client *Client
}

// NewAnnouncementsService creates a new announcements service
func NewAnnouncementsService(client *Client) *AnnouncementsService {
	return &AnnouncementsService{client: client}
}

// ListAnnouncementsOptions holds options for listing announcements
type ListAnnouncementsOptions struct {
	ContextCodes []string // course_123, course_456
	StartDate    string   // yyyy-mm-dd or ISO 8601
	EndDate      string
	ActiveOnly   bool
	LatestOnly   bool
	Include      []string // sections, sections_user_count
}

// List retrieves announcements for the given contexts
func (s *AnnouncementsService) List(ctx context.Context, opts *ListAnnouncementsOptions) ([]DiscussionTopic, error) {
	path := "/api/v1/announcements"

	if opts != nil {
		query := url.Values{}

		for _, code := range opts.ContextCodes {
			query.Add("context_codes[]", code)
		}

		if opts.StartDate != "" {
			query.Add("start_date", opts.StartDate)
		}

		if opts.EndDate != "" {
			query.Add("end_date", opts.EndDate)
		}

		if opts.ActiveOnly {
			query.Add("active_only", "true")
		}

		if opts.LatestOnly {
			query.Add("latest_only", "true")
		}

		for _, inc := range opts.Include {
			query.Add("include[]", inc)
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var announcements []DiscussionTopic
	if err := s.client.GetAllPages(ctx, path, &announcements); err != nil {
		return nil, err
	}

	return announcements, nil
}

// DiscussionDateDetails holds date availability details for a discussion topic
type DiscussionDateDetails struct {
	DueAt     *string              `json:"due_at,omitempty"`
	UnlockAt  *string              `json:"unlock_at,omitempty"`
	LockAt    *string              `json:"lock_at,omitempty"`
	Overrides []AssignmentOverride `json:"overrides,omitempty"`
}

// GetDateDetails retrieves date details for a discussion topic in a course
// Canvas path: GET /api/v1/courses/:course_id/discussion_topics/:discussion_topic_id/date_details
func (s *DiscussionsService) GetDateDetails(ctx context.Context, courseID, topicID int64) (*DiscussionDateDetails, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d/date_details", courseID, topicID)

	var result DiscussionDateDetails
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateDateDetails updates date details for a discussion topic in a course
// Canvas path: PUT /api/v1/courses/:course_id/discussion_topics/:discussion_topic_id/date_details
func (s *DiscussionsService) UpdateDateDetails(ctx context.Context, courseID, topicID int64, params *DiscussionDateDetails) (*DiscussionDateDetails, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/discussion_topics/%d/date_details", courseID, topicID)

	var result DiscussionDateDetails
	if err := s.client.PutJSON(ctx, path, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
