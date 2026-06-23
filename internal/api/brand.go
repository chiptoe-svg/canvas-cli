package api

import (
	"context"
	"fmt"
)

// BrandVariables represents the Canvas branding/theme variables for an account.
type BrandVariables struct {
	IcBrandPrimary           string `json:"ic-brand-primary,omitempty"`
	IcBrandButtonPrimaryBgd  string `json:"ic-brand-button--primary-bgd,omitempty"`
	IcBrandButtonPrimaryText string `json:"ic-brand-button--primary-text,omitempty"`
	IcBrandFontColorDark     string `json:"ic-brand-font-color-dark,omitempty"`
	IcLinkColor              string `json:"ic-link-color,omitempty"`
	IcBrandNavBgd            string `json:"ic-brand-nav-bgd,omitempty"`
	IcBrandNavTextColor      string `json:"ic-brand-nav-text-color,omitempty"`
	IcBrandNavIconFill       string `json:"ic-brand-nav-icon-fill,omitempty"`
	IcBrandLogoImage         string `json:"ic-brand-logo-image,omitempty"`
	IcBrandFaviconImage      string `json:"ic-brand-favicon,omitempty"`
}

// BrandService handles brand-variable API calls.
type BrandService struct {
	client *Client
}

// NewBrandService creates a new BrandService.
func NewBrandService(client *Client) *BrandService {
	return &BrandService{client: client}
}

// GetVariables retrieves the current Canvas theme variables (global).
func (s *BrandService) GetVariables(ctx context.Context) (*BrandVariables, error) {
	path := "/api/v1/brand_variables"

	var bv BrandVariables
	if err := s.client.GetJSON(ctx, path, &bv); err != nil {
		return nil, err
	}

	return &bv, nil
}

// GetVariablesForAccount retrieves the brand variables for a specific account.
func (s *BrandService) GetVariablesForAccount(ctx context.Context, accountID int64) (*BrandVariables, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/brand_variables", accountID)

	var bv BrandVariables
	if err := s.client.GetJSON(ctx, path, &bv); err != nil {
		return nil, err
	}

	return &bv, nil
}

// GetVariablesForCourse retrieves the brand variables for a specific course.
func (s *BrandService) GetVariablesForCourse(ctx context.Context, courseID int64) (*BrandVariables, error) {
	path := fmt.Sprintf("/api/v1/courses/%d/brand_variables", courseID)

	var bv BrandVariables
	if err := s.client.GetJSON(ctx, path, &bv); err != nil {
		return nil, err
	}

	return &bv, nil
}
