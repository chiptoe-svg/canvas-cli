package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// Conference represents a Canvas web conference (BigBlueButton, Zoom, etc.).
type Conference struct {
	ID                  int64         `json:"id"`
	ConferenceType      string        `json:"conference_type,omitempty"`
	ConferenceKey       string        `json:"conference_key,omitempty"`
	Description         string        `json:"description,omitempty"`
	Duration            *int          `json:"duration,omitempty"`
	EndedAt             string        `json:"ended_at,omitempty"`
	StartedAt           string        `json:"started_at,omitempty"`
	Title               string        `json:"title"`
	URL                 string        `json:"url,omitempty"`
	ContextType         string        `json:"context_type,omitempty"`
	ContextID           int64         `json:"context_id,omitempty"`
	Participants        []interface{} `json:"participants,omitempty"`
	HasAdvancedSettings bool          `json:"has_advanced_settings,omitempty"`
	JoinURL             string        `json:"join_url,omitempty"`
	LongRunning         bool          `json:"long_running,omitempty"`
	UserSettings        interface{}   `json:"user_settings,omitempty"`
	Recordings          []interface{} `json:"recordings,omitempty"`
	Agenda              string        `json:"agenda,omitempty"`
}

// ConferencesService handles conference-related API calls.
type ConferencesService struct {
	client *Client
}

// NewConferencesService creates a new ConferencesService.
func NewConferencesService(client *Client) *ConferencesService {
	return &ConferencesService{client: client}
}

// ListConferencesOptions holds query parameters for listing conferences.
type ListConferencesOptions struct {
	State   string // live, ended
	PerPage int
}

func buildConferenceQuery(opts *ListConferencesOptions) string {
	if opts == nil {
		return ""
	}
	q := url.Values{}
	if opts.State != "" {
		q.Set("state", opts.State)
	}
	if opts.PerPage > 0 {
		q.Set("per_page", strconv.Itoa(opts.PerPage))
	}
	if len(q) > 0 {
		return "?" + q.Encode()
	}
	return ""
}

// List retrieves all conferences the current user participates in.
func (s *ConferencesService) List(ctx context.Context, opts *ListConferencesOptions) ([]Conference, error) {
	path := "/api/v1/conferences" + buildConferenceQuery(opts)

	var confs []Conference
	if err := s.client.GetAllPages(ctx, path, &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

// ListForCourse retrieves conferences in a course.
func (s *ConferencesService) ListForCourse(ctx context.Context, courseID int64, opts *ListConferencesOptions) ([]Conference, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/conferences", courseID) + buildConferenceQuery(opts)

	var confs []Conference
	if err := s.client.GetAllPages(ctx, path, &confs); err != nil {
		return nil, err
	}

	return confs, nil
}

// ListForGroup retrieves conferences in a group.
func (s *ConferencesService) ListForGroup(ctx context.Context, groupID int64, opts *ListConferencesOptions) ([]Conference, error) {
	path := fmt.Sprintf("/api/v1/groups/%d/conferences", groupID) + buildConferenceQuery(opts)

	var confs []Conference
	if err := s.client.GetAllPages(ctx, path, &confs); err != nil {
		return nil, err
	}

	return confs, nil
}
