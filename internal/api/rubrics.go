package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// RubricsService handles rubric-related API calls
type RubricsService struct {
	client *Client
}

// NewRubricsService creates a new rubrics service
func NewRubricsService(client *Client) *RubricsService {
	return &RubricsService{client: client}
}

// Rubric represents a Canvas rubric
type Rubric struct {
	ID                        int64               `json:"id"`
	Title                     string              `json:"title"`
	ContextType               string              `json:"context_type"`
	ContextID                 int64               `json:"context_id"`
	PointsPossible            float64             `json:"points_possible"`
	Reusable                  bool                `json:"reusable"`
	ReadOnly                  bool                `json:"read_only"`
	FreeFormCriterionComments bool                `json:"free_form_criterion_comments"`
	HideScoreTotal            bool                `json:"hide_score_total"`
	Data                      []RubricCriterion   `json:"data,omitempty"`
	Assessments               []RubricAssessment  `json:"assessments,omitempty"`
	Associations              []RubricAssociation `json:"associations,omitempty"`
}

// RubricAssociation represents a rubric association with an assignment
type RubricAssociation struct {
	ID                 int64  `json:"id"`
	RubricID           int64  `json:"rubric_id"`
	AssociationID      int64  `json:"association_id"`
	AssociationType    string `json:"association_type"`
	UseForGrading      bool   `json:"use_for_grading"`
	SummaryData        string `json:"summary_data,omitempty"`
	Purpose            string `json:"purpose"`
	HideScoreTotal     bool   `json:"hide_score_total"`
	HidePoints         bool   `json:"hide_points"`
	HideOutcomeResults bool   `json:"hide_outcome_results"`
}

// rubricResponse is a wrapper for API responses that include a rubric
type rubricResponse struct {
	Rubric *Rubric `json:"rubric"`
}

// rubricAssociationResponse is a wrapper for API responses that include a rubric association
type rubricAssociationResponse struct {
	RubricAssociation *RubricAssociation `json:"rubric_association"`
}

// rubricOrDryRun unwraps a rubric response. Canvas wraps rubric write
// responses in {"rubric": ...}; in dry-run mode the client never sends the
// request and answers with an empty body, so the wrapper carries no rubric.
// That is not a failure: return an empty rubric, as other services do.
func (s *RubricsService) rubricOrDryRun(response *rubricResponse) (*Rubric, error) {
	if response.Rubric != nil {
		return response.Rubric, nil
	}
	if s.client.IsDryRun() {
		return &Rubric{}, nil
	}
	return nil, fmt.Errorf("rubric not returned in response")
}

// ListRubricsOptions holds options for listing rubrics
type ListRubricsOptions struct {
	Include []string // assessments, associations, assignment_associations
	Page    int
	PerPage int
}

// ListCourse retrieves all rubrics for a course
func (s *RubricsService) ListCourse(ctx context.Context, courseID int64, opts *ListRubricsOptions) ([]Rubric, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubrics", courseID)

	if opts != nil {
		query := url.Values{}

		for _, include := range opts.Include {
			query.Add("include[]", include)
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

	var rubrics []Rubric
	if err := s.client.GetAllPages(ctx, path, &rubrics); err != nil {
		return nil, err
	}

	return rubrics, nil
}

// ListAccount retrieves all rubrics for an account
func (s *RubricsService) ListAccount(ctx context.Context, accountID int64, opts *ListRubricsOptions) ([]Rubric, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/rubrics", accountID)

	if opts != nil {
		query := url.Values{}

		for _, include := range opts.Include {
			query.Add("include[]", include)
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

	var rubrics []Rubric
	if err := s.client.GetAllPages(ctx, path, &rubrics); err != nil {
		return nil, err
	}

	return rubrics, nil
}

// GetCourse retrieves a single rubric from a course
func (s *RubricsService) GetCourse(ctx context.Context, courseID, rubricID int64, include []string) (*Rubric, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubrics/%d", courseID, rubricID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var rubric Rubric
	if err := s.client.GetJSON(ctx, path, &rubric); err != nil {
		return nil, err
	}

	return &rubric, nil
}

// GetAccount retrieves a single rubric from an account
func (s *RubricsService) GetAccount(ctx context.Context, accountID, rubricID int64, include []string) (*Rubric, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/rubrics/%d", accountID, rubricID)

	if len(include) > 0 {
		query := url.Values{}
		for _, inc := range include {
			query.Add("include[]", inc)
		}
		path += "?" + query.Encode()
	}

	var rubric Rubric
	if err := s.client.GetJSON(ctx, path, &rubric); err != nil {
		return nil, err
	}

	return &rubric, nil
}

// CreateRubricParams holds parameters for creating a rubric
type CreateRubricParams struct {
	Title                     string
	PointsPossible            float64
	FreeFormCriterionComments bool
	HideScoreTotal            bool
	Criteria                  []RubricCriterion
}

// Create creates a new rubric in a course
func (s *RubricsService) Create(ctx context.Context, courseID int64, params *CreateRubricParams) (*Rubric, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubrics", courseID)

	body := map[string]interface{}{
		"rubric": make(map[string]interface{}),
	}

	rubricData, ok := body["rubric"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid rubric data structure")
	}

	if params.Title != "" {
		rubricData["title"] = params.Title
	}

	if params.PointsPossible > 0 {
		rubricData["points_possible"] = params.PointsPossible
	}

	if params.FreeFormCriterionComments {
		rubricData["free_form_criterion_comments"] = params.FreeFormCriterionComments
	}

	if params.HideScoreTotal {
		rubricData["hide_score_total"] = params.HideScoreTotal
	}

	if len(params.Criteria) > 0 {
		rubricData["criteria"] = encodeRubricCriteria(params.Criteria)
	}

	var response rubricResponse
	if err := s.client.PostJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	return s.rubricOrDryRun(&response)
}

// encodeRubricCriteria converts criteria into the indexed-hash form Canvas
// expects: rubric[criteria][0][description], rubric[criteria][0][ratings][1][points], ...
func encodeRubricCriteria(criteria []RubricCriterion) map[string]interface{} {
	encoded := make(map[string]interface{}, len(criteria))
	for i, c := range criteria {
		criterionData := map[string]interface{}{
			"description":      c.Description,
			"long_description": c.LongDescription,
			"points":           c.Points,
		}
		// Ids come from a `rubrics get` round-trip. Canvas keys existing
		// rubric assessments by criterion id and mints new ids for criteria
		// sent without one, so forwarding them keeps graded work attached.
		if c.ID != "" {
			criterionData["id"] = c.ID
		}
		if c.CriterionUseRange {
			criterionData["criterion_use_range"] = true
		}

		if len(c.Ratings) > 0 {
			ratings := make(map[string]interface{}, len(c.Ratings))
			for j, r := range c.Ratings {
				ratingData := map[string]interface{}{
					"description":      r.Description,
					"long_description": r.LongDescription,
					"points":           r.Points,
				}
				if r.ID != "" {
					ratingData["id"] = r.ID
				}
				ratings[strconv.Itoa(j)] = ratingData
			}
			criterionData["ratings"] = ratings
		}

		encoded[strconv.Itoa(i)] = criterionData
	}
	return encoded
}

// UpdateRubricParams holds parameters for updating a rubric
type UpdateRubricParams struct {
	Title                     *string
	PointsPossible            *float64
	FreeFormCriterionComments *bool
	HideScoreTotal            *bool
	// Criteria, when non-empty, replaces the rubric's criteria wholesale.
	// Left empty, the existing criteria are untouched.
	Criteria []RubricCriterion
}

// Update updates an existing rubric
func (s *RubricsService) Update(ctx context.Context, courseID, rubricID int64, params *UpdateRubricParams) (*Rubric, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubrics/%d", courseID, rubricID)

	body := map[string]interface{}{
		"rubric": make(map[string]interface{}),
	}

	rubricData, ok := body["rubric"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("internal error: invalid rubric data structure")
	}

	if params.Title != nil {
		rubricData["title"] = *params.Title
	}

	if params.PointsPossible != nil {
		rubricData["points_possible"] = *params.PointsPossible
	}

	if params.FreeFormCriterionComments != nil {
		rubricData["free_form_criterion_comments"] = *params.FreeFormCriterionComments
	}

	if params.HideScoreTotal != nil {
		rubricData["hide_score_total"] = *params.HideScoreTotal
	}

	if len(params.Criteria) > 0 {
		rubricData["criteria"] = encodeRubricCriteria(params.Criteria)
	}

	var response rubricResponse
	if err := s.client.PutJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	return s.rubricOrDryRun(&response)
}

// Delete deletes a rubric
func (s *RubricsService) Delete(ctx context.Context, courseID, rubricID int64) (*Rubric, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubrics/%d", courseID, rubricID)

	var response rubricResponse
	if err := s.client.DeleteJSON(ctx, path, &response); err != nil {
		return nil, err
	}

	return s.rubricOrDryRun(&response)
}

// AssociateParams holds parameters for associating a rubric
type AssociateParams struct {
	AssociationType string // "Assignment"
	AssociationID   int64
	UseForGrading   bool
	HideScoreTotal  bool
	HidePoints      bool
	Purpose         string // "grading", "bookmark"
}

// Associate associates a rubric with an assignment
func (s *RubricsService) Associate(ctx context.Context, courseID, rubricID int64, params *AssociateParams) (*RubricAssociation, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/rubric_associations", courseID)

	body := map[string]interface{}{
		"rubric_association": map[string]interface{}{
			"rubric_id":        rubricID,
			"association_type": params.AssociationType,
			"association_id":   params.AssociationID,
			"use_for_grading":  params.UseForGrading,
			"hide_score_total": params.HideScoreTotal,
			"hide_points":      params.HidePoints,
			"purpose":          params.Purpose,
		},
	}

	var response rubricAssociationResponse
	if err := s.client.PostJSON(ctx, path, body, &response); err != nil {
		return nil, err
	}

	if response.RubricAssociation == nil {
		if s.client.IsDryRun() {
			return &RubricAssociation{}, nil
		}
		return nil, fmt.Errorf("rubric association not returned in response")
	}

	return response.RubricAssociation, nil
}
