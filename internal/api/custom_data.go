package api

import (
	"context"
	"fmt"
	"net/url"
)

// CustomDataService handles user custom data API calls
type CustomDataService struct {
	client *Client
}

// NewCustomDataService creates a new custom data service
func NewCustomDataService(client *Client) *CustomDataService {
	return &CustomDataService{client: client}
}

// CustomDataResult wraps the Canvas custom_data response envelope
type CustomDataResult struct {
	Data map[string]interface{} `json:"data"`
}

// Get retrieves custom data stored for a user under the given namespace
func (s *CustomDataService) Get(ctx context.Context, userID int64, ns string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/users/%d/custom_data", userID)
	if ns != "" {
		query := url.Values{}
		query.Set("ns", ns)
		path += "?" + query.Encode()
	}
	var result CustomDataResult
	if err := s.client.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("getting custom data for user %d: %w", userID, err)
	}
	return result.Data, nil
}

// Set stores arbitrary custom data for a user under the given namespace
func (s *CustomDataService) Set(ctx context.Context, userID int64, ns string, data map[string]interface{}) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/users/%d/custom_data", userID)
	body := map[string]interface{}{
		"ns":   ns,
		"data": data,
	}
	var result CustomDataResult
	if err := s.client.PutJSON(ctx, path, body, &result); err != nil {
		return nil, fmt.Errorf("setting custom data for user %d: %w", userID, err)
	}
	return result.Data, nil
}

// Delete removes custom data stored for a user under the given namespace
func (s *CustomDataService) Delete(ctx context.Context, userID int64, ns string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/v1/users/%d/custom_data", userID)
	if ns != "" {
		query := url.Values{}
		query.Set("ns", ns)
		path += "?" + query.Encode()
	}
	var result CustomDataResult
	if err := s.client.DeleteJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("deleting custom data for user %d: %w", userID, err)
	}
	return result.Data, nil
}
