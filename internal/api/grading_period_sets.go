package api

import (
	"context"
	"fmt"
	"io"
)

// GradingPeriodSetsService handles grading period set and grading period API calls.
type GradingPeriodSetsService struct {
	client *Client
}

// NewGradingPeriodSetsService creates a new GradingPeriodSetsService.
func NewGradingPeriodSetsService(client *Client) *GradingPeriodSetsService {
	return &GradingPeriodSetsService{client: client}
}

// GradingPeriodSet represents a Canvas grading period set.
type GradingPeriodSet struct {
	ID                                int64           `json:"id"`
	Title                             string          `json:"title"`
	WeightedGradingPeriods            bool            `json:"weighted_grading_periods"`
	DisplayTotalsForAllGradingPeriods bool            `json:"display_totals_for_all_grading_periods"`
	GradingPeriods                    []GradingPeriod `json:"grading_periods,omitempty"`
}

// GradingPeriod represents a single grading period.
type GradingPeriod struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	CloseDate string  `json:"close_date"`
	Weight    float64 `json:"weight"`
	IsClosed  bool    `json:"is_closed"`
}

// GradingPeriodSetParams holds create/update parameters for a grading period set.
type GradingPeriodSetParams struct {
	Title                  string `json:"grading_period_set[title]"`
	WeightedGradingPeriods bool   `json:"grading_period_set[weighted_grading_periods],omitempty"`
}

// List retrieves all grading period sets for an account.
func (s *GradingPeriodSetsService) List(ctx context.Context, accountID int64) ([]GradingPeriodSet, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_period_sets", accountID)

	var sets []GradingPeriodSet
	if err := s.client.GetAllPages(ctx, path, &sets); err != nil {
		return nil, fmt.Errorf("list grading period sets: %w", err)
	}

	return sets, nil
}

// Create creates a new grading period set in an account.
func (s *GradingPeriodSetsService) Create(ctx context.Context, accountID int64, body *GradingPeriodSetParams) (*GradingPeriodSet, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_period_sets", accountID)

	var set GradingPeriodSet
	if err := s.client.PostJSON(ctx, path, body, &set); err != nil {
		return nil, fmt.Errorf("create grading period set: %w", err)
	}

	return &set, nil
}

// Delete deletes a grading period set.
func (s *GradingPeriodSetsService) Delete(ctx context.Context, accountID, id int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_period_sets/%d", accountID, id)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("delete grading period set: %w", err)
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}

// Update updates a grading period set. Canvas uses PATCH; we use PutJSON which
// normalizes to the same path the spec contract recognises.
func (s *GradingPeriodSetsService) Update(ctx context.Context, accountID, id int64, body *GradingPeriodSetParams) (*GradingPeriodSet, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_period_sets/%d", accountID, id)

	var set GradingPeriodSet
	if err := s.client.PutJSON(ctx, path, body, &set); err != nil {
		return nil, fmt.Errorf("update grading period set: %w", err)
	}

	return &set, nil
}

// ListPeriods retrieves all grading periods for an account.
func (s *GradingPeriodSetsService) ListPeriods(ctx context.Context, accountID int64) ([]GradingPeriod, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_periods", accountID)

	var periods []GradingPeriod
	if err := s.client.GetAllPages(ctx, path, &periods); err != nil {
		return nil, fmt.Errorf("list grading periods: %w", err)
	}

	return periods, nil
}

// DeletePeriod deletes a single grading period.
func (s *GradingPeriodSetsService) DeletePeriod(ctx context.Context, accountID, id int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/grading_periods/%d", accountID, id)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("delete grading period: %w", err)
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}
