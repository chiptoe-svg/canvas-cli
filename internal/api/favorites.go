package api

import (
	"context"
	"fmt"
)

// FavoritesService handles favorites-related API calls
type FavoritesService struct {
	client *Client
}

// NewFavoritesService creates a new favorites service
func NewFavoritesService(client *Client) *FavoritesService {
	return &FavoritesService{client: client}
}

// ListCourses retrieves the current user's favorite courses
func (s *FavoritesService) ListCourses(ctx context.Context) ([]Course, error) {
	path := "/api/v1/users/self/favorites/courses"
	var result []Course
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing favorite courses: %w", err)
	}
	return result, nil
}

// AddCourse adds a course to the current user's favorites
func (s *FavoritesService) AddCourse(ctx context.Context, id int64) (*Course, error) {
	path := fmt.Sprintf("/api/v1/users/self/favorites/courses/%d", id)
	var result Course
	if err := s.client.PostJSON(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("adding course %d to favorites: %w", id, err)
	}
	return &result, nil
}

// RemoveCourse removes a course from the current user's favorites
func (s *FavoritesService) RemoveCourse(ctx context.Context, id int64) (*Course, error) {
	path := fmt.Sprintf("/api/v1/users/self/favorites/courses/%d", id)
	var result Course
	if err := s.client.DeleteJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("removing course %d from favorites: %w", id, err)
	}
	return &result, nil
}

// ResetCourses resets the current user's favorite courses to the default list
func (s *FavoritesService) ResetCourses(ctx context.Context) ([]Course, error) {
	path := "/api/v1/users/self/favorites/courses"
	var result []Course
	if err := s.client.DeleteJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("resetting favorite courses: %w", err)
	}
	return result, nil
}

// ListGroups retrieves the current user's favorite groups
func (s *FavoritesService) ListGroups(ctx context.Context) ([]Group, error) {
	path := "/api/v1/users/self/favorites/groups"
	var result []Group
	if err := s.client.GetAllPages(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("listing favorite groups: %w", err)
	}
	return result, nil
}

// AddGroup adds a group to the current user's favorites
func (s *FavoritesService) AddGroup(ctx context.Context, id int64) (*Group, error) {
	path := fmt.Sprintf("/api/v1/users/self/favorites/groups/%d", id)
	var result Group
	if err := s.client.PostJSON(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("adding group %d to favorites: %w", id, err)
	}
	return &result, nil
}

// RemoveGroup removes a group from the current user's favorites
func (s *FavoritesService) RemoveGroup(ctx context.Context, id int64) (*Group, error) {
	path := fmt.Sprintf("/api/v1/users/self/favorites/groups/%d", id)
	var result Group
	if err := s.client.DeleteJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("removing group %d from favorites: %w", id, err)
	}
	return &result, nil
}

// ResetGroups resets the current user's favorite groups to the default list
func (s *FavoritesService) ResetGroups(ctx context.Context) ([]Group, error) {
	path := "/api/v1/users/self/favorites/groups"
	var result []Group
	if err := s.client.DeleteJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("resetting favorite groups: %w", err)
	}
	return result, nil
}
