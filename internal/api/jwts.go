package api

import "context"

// JWT represents a Canvas JWT token response.
type JWT struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	CanvasUserID int64  `json:"canvas_user_id,omitempty"`
	RealUserID   int64  `json:"real_user_id,omitempty"`
	PseudonymID  int64  `json:"pseudonym_id,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	CompactJWT   string `json:"compact_jwt,omitempty"`
}

// JWTsService handles JWT creation and refresh.
type JWTsService struct {
	client *Client
}

// NewJWTsService creates a new JWTsService.
func NewJWTsService(client *Client) *JWTsService {
	return &JWTsService{client: client}
}

// Create creates a new Canvas JWT for the current user.
// workflowState is optional; leave empty for the default.
func (s *JWTsService) Create(ctx context.Context, workflowState string) (*JWT, error) {
	path := "/api/v1/jwts"

	body := map[string]interface{}{}
	if workflowState != "" {
		body["workflows"] = []string{workflowState}
	}

	var jwt JWT
	if err := s.client.PostJSON(ctx, path, body, &jwt); err != nil {
		return nil, err
	}

	return &jwt, nil
}

// Refresh refreshes an existing Canvas JWT.
func (s *JWTsService) Refresh(ctx context.Context, refreshToken string) (*JWT, error) {
	path := "/api/v1/jwts/refresh"

	body := map[string]interface{}{
		"jwt": refreshToken,
	}

	var jwt JWT
	if err := s.client.PostJSON(ctx, path, body, &jwt); err != nil {
		return nil, err
	}

	return &jwt, nil
}
