package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// MediaObject represents a Canvas media object (video/audio).
type MediaObject struct {
	ID               string       `json:"id,omitempty"`
	MediaID          string       `json:"media_id,omitempty"`
	MediaType        string       `json:"media_type,omitempty"`
	Duration         int          `json:"duration,omitempty"`
	Title            string       `json:"title,omitempty"`
	UserEnteredTitle string       `json:"user_entered_title,omitempty"`
	EmbeddedURL      string       `json:"embedded_url,omitempty"`
	MediaTracks      []MediaTrack `json:"media_tracks,omitempty"`
}

// MediaTrack represents a media track (subtitle/caption file) for a media object.
type MediaTrack struct {
	ID        int64  `json:"id,omitempty"`
	MediaID   string `json:"media_id,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Locale    string `json:"locale,omitempty"`
	Content   string `json:"content,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	WEBVTTURL string `json:"webvtt_content,omitempty"`
}

// MediaAttachment represents a Canvas media attachment (file-backed media).
type MediaAttachment struct {
	ID           int64        `json:"id"`
	ContentType  string       `json:"content_type,omitempty"`
	Filename     string       `json:"filename,omitempty"`
	DisplayName  string       `json:"display_name,omitempty"`
	Size         int64        `json:"size,omitempty"`
	MediaEntryID string       `json:"media_entry_id,omitempty"`
	MediaObject  *MediaObject `json:"media_object,omitempty"`
}

// MediaObjectsService handles media object and attachment API calls.
type MediaObjectsService struct {
	client *Client
}

// NewMediaObjectsService creates a new MediaObjectsService.
func NewMediaObjectsService(client *Client) *MediaObjectsService {
	return &MediaObjectsService{client: client}
}

// ListMediaObjectsOptions holds query parameters for listing media objects.
type ListMediaObjectsOptions struct {
	Sort    string // title, created_at, updated_at, user_name
	Order   string // asc, desc
	Exclude []string
	PerPage int
}

func buildMediaObjectsQuery(opts *ListMediaObjectsOptions) string {
	if opts == nil {
		return ""
	}
	q := url.Values{}
	if opts.Sort != "" {
		q.Set("sort", opts.Sort)
	}
	if opts.Order != "" {
		q.Set("order", opts.Order)
	}
	for _, ex := range opts.Exclude {
		q.Add("exclude[]", ex)
	}
	if opts.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(opts.PerPage))
	}
	if len(q) > 0 {
		return "?" + q.Encode()
	}
	return ""
}

// List retrieves all media objects for the current user.
func (s *MediaObjectsService) List(ctx context.Context, opts *ListMediaObjectsOptions) ([]MediaObject, error) {
	path := "/api/v1/media_objects" + buildMediaObjectsQuery(opts)

	var objs []MediaObject
	if err := s.client.GetAllPages(ctx, path, &objs); err != nil {
		return nil, err
	}

	return objs, nil
}

// ListForCourse retrieves media objects for a course.
func (s *MediaObjectsService) ListForCourse(ctx context.Context, courseID int64, opts *ListMediaObjectsOptions) ([]MediaObject, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/media_objects", courseID) + buildMediaObjectsQuery(opts)

	var objs []MediaObject
	if err := s.client.GetAllPages(ctx, path, &objs); err != nil {
		return nil, err
	}

	return objs, nil
}

// ListForGroup retrieves media objects for a group.
func (s *MediaObjectsService) ListForGroup(ctx context.Context, groupID int64, opts *ListMediaObjectsOptions) ([]MediaObject, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/media_objects", groupID) + buildMediaObjectsQuery(opts)

	var objs []MediaObject
	if err := s.client.GetAllPages(ctx, path, &objs); err != nil {
		return nil, err
	}

	return objs, nil
}

// Update updates the title of a media object.
func (s *MediaObjectsService) Update(ctx context.Context, mediaObjectID string, title string) (*MediaObject, error) {
	path := fmt.Sprintf("/api/v1/media_objects/%s", mediaObjectID)

	body := map[string]interface{}{}
	if title != "" {
		body["user_entered_title"] = title
	}

	var obj MediaObject
	if err := s.client.PutJSON(ctx, path, body, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

// GetMediaTracks retrieves the media tracks for a media object.
func (s *MediaObjectsService) GetMediaTracks(ctx context.Context, mediaObjectID string) ([]MediaTrack, error) {
	path := fmt.Sprintf("/api/v1/media_objects/%s/media_tracks", mediaObjectID)

	var tracks []MediaTrack
	if err := s.client.GetAllPages(ctx, path, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

// UpdateMediaTracks updates the media tracks for a media object.
// content is the full WEBVTT content for a track with the given locale and kind.
func (s *MediaObjectsService) UpdateMediaTracks(ctx context.Context, mediaObjectID string, tracks []map[string]string) ([]MediaTrack, error) {
	path := fmt.Sprintf("/api/v1/media_objects/%s/media_tracks", mediaObjectID)

	body := map[string]interface{}{
		"media_tracks": tracks,
	}

	var result []MediaTrack
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// --- Media Attachments ---

// ListAttachments retrieves all media attachments for the current user.
func (s *MediaObjectsService) ListAttachments(ctx context.Context, opts *ListMediaObjectsOptions) ([]MediaAttachment, error) {
	path := "/api/v1/media_attachments" + buildMediaObjectsQuery(opts)

	var atts []MediaAttachment
	if err := s.client.GetAllPages(ctx, path, &atts); err != nil {
		return nil, err
	}

	return atts, nil
}

// ListAttachmentsForCourse retrieves media attachments for a course.
func (s *MediaObjectsService) ListAttachmentsForCourse(ctx context.Context, courseID int64, opts *ListMediaObjectsOptions) ([]MediaAttachment, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/media_attachments", courseID) + buildMediaObjectsQuery(opts)

	var atts []MediaAttachment
	if err := s.client.GetAllPages(ctx, path, &atts); err != nil {
		return nil, err
	}

	return atts, nil
}

// ListAttachmentsForGroup retrieves media attachments for a group.
func (s *MediaObjectsService) ListAttachmentsForGroup(ctx context.Context, groupID int64, opts *ListMediaObjectsOptions) ([]MediaAttachment, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/media_attachments", groupID) + buildMediaObjectsQuery(opts)

	var atts []MediaAttachment
	if err := s.client.GetAllPages(ctx, path, &atts); err != nil {
		return nil, err
	}

	return atts, nil
}

// UpdateAttachment updates the display_name of a media attachment.
func (s *MediaObjectsService) UpdateAttachment(ctx context.Context, attachmentID int64, title string) (*MediaAttachment, error) {
	path := fmt.Sprintf("/api/v1/media_attachments/%d", attachmentID)

	body := map[string]interface{}{}
	if title != "" {
		body["user_entered_title"] = title
	}

	var att MediaAttachment
	if err := s.client.PutJSON(ctx, path, body, &att); err != nil {
		return nil, err
	}

	return &att, nil
}

// GetAttachmentMediaTracks retrieves the media tracks for a media attachment.
func (s *MediaObjectsService) GetAttachmentMediaTracks(ctx context.Context, attachmentID int64) ([]MediaTrack, error) {
	path := fmt.Sprintf("/api/v1/media_attachments/%d/media_tracks", attachmentID)

	var tracks []MediaTrack
	if err := s.client.GetAllPages(ctx, path, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}

// UpdateAttachmentMediaTracks updates the media tracks for a media attachment.
func (s *MediaObjectsService) UpdateAttachmentMediaTracks(ctx context.Context, attachmentID int64, tracks []map[string]string) ([]MediaTrack, error) {
	path := fmt.Sprintf("/api/v1/media_attachments/%d/media_tracks", attachmentID)

	body := map[string]interface{}{
		"media_tracks": tracks,
	}

	var result []MediaTrack
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
