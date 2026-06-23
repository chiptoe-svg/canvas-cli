package api

import (
	"context"
	"fmt"
	"io"
)

// AccountReportsService handles account report API calls
type AccountReportsService struct {
	client *Client
}

// NewAccountReportsService creates a new account reports service
func NewAccountReportsService(client *Client) *AccountReportsService {
	return &AccountReportsService{client: client}
}

// AccountReport represents a Canvas account report type
type AccountReport struct {
	ID         string                 `json:"report"`
	Title      string                 `json:"title"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// AccountReportRun represents a single run of an account report
type AccountReportRun struct {
	ID         int64                  `json:"id"`
	Report     string                 `json:"report"`
	Status     string                 `json:"status"`
	CreatedAt  string                 `json:"created_at"`
	StartedAt  string                 `json:"started_at,omitempty"`
	EndedAt    string                 `json:"ended_at,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Progress   int                    `json:"progress"`
	Attachment *struct {
		ID  int64  `json:"id"`
		URL string `json:"url"`
	} `json:"attachment,omitempty"`
}

// List retrieves available report types for an account
func (s *AccountReportsService) List(ctx context.Context, accountID int64) ([]AccountReport, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/reports", accountID)

	var reports []AccountReport
	if err := s.client.GetAllPages(ctx, path, &reports); err != nil {
		return nil, fmt.Errorf("list account reports: %w", err)
	}

	return reports, nil
}

// ListRuns retrieves all runs for a specific report type
func (s *AccountReportsService) ListRuns(ctx context.Context, accountID int64, report string) ([]AccountReportRun, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/reports/%s", accountID, report)

	var runs []AccountReportRun
	if err := s.client.GetAllPages(ctx, path, &runs); err != nil {
		return nil, fmt.Errorf("list account report runs: %w", err)
	}

	return runs, nil
}

// Start starts a new run of a specific report type
func (s *AccountReportsService) Start(ctx context.Context, accountID int64, report string, body interface{}) (*AccountReportRun, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/reports/%s", accountID, report)

	var run AccountReportRun
	if err := s.client.PostJSON(ctx, path, body, &run); err != nil {
		return nil, fmt.Errorf("start account report: %w", err)
	}

	return &run, nil
}

// GetRun retrieves a specific report run
func (s *AccountReportsService) GetRun(ctx context.Context, accountID int64, report string, id int64) (*AccountReportRun, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/reports/%s/%d", accountID, report, id)

	var run AccountReportRun
	if err := s.client.GetJSON(ctx, path, &run); err != nil {
		return nil, fmt.Errorf("get account report run: %w", err)
	}

	return &run, nil
}

// DeleteRun deletes a specific report run
func (s *AccountReportsService) DeleteRun(ctx context.Context, accountID int64, report string, id int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/reports/%s/%d", accountID, report, id)

	resp, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("delete account report run: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	return nil
}

// AbortRun aborts a running report
func (s *AccountReportsService) AbortRun(ctx context.Context, accountID int64, report string, id int64) (*AccountReportRun, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/reports/%s/%d/abort", accountID, report, id)

	var run AccountReportRun
	if err := s.client.PutJSON(ctx, path, nil, &run); err != nil {
		return nil, fmt.Errorf("abort account report run: %w", err)
	}

	return &run, nil
}
