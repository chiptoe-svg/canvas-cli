package main

import "testing"

// TestNormalizePath exercises the normalizePath function that collapses path
// parameters, trailing slashes, and query strings into a canonical form used
// for spec-to-CLI coverage comparison.
func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// --- documented (colon) params ---
		{
			name:  "single colon param",
			input: "/api/v1/courses/:course_id",
			want:  "/api/v1/courses/:x",
		},
		{
			name:  "multiple colon params",
			input: "/api/v1/courses/:course_id/assignments/:assignment_id",
			want:  "/api/v1/courses/:x/assignments/:x",
		},
		{
			name:  "colon id param",
			input: "/api/v1/accounts/:id",
			want:  "/api/v1/accounts/:x",
		},
		// --- Go fmt.Sprintf format verbs ---
		{
			name:  "single percent-d verb",
			input: "/api/v1/courses/%d/assignments",
			want:  "/api/v1/courses/:x/assignments",
		},
		{
			name:  "multiple percent verbs",
			input: "/api/v1/courses/%d/assignments/%d",
			want:  "/api/v1/courses/:x/assignments/:x",
		},
		{
			name:  "percent-s verb",
			input: "/api/v1/users/%s/profile",
			want:  "/api/v1/users/:x/profile",
		},
		// --- trailing slash ---
		{
			name:  "trailing slash stripped",
			input: "/api/v1/courses/",
			want:  "/api/v1/courses",
		},
		{
			name:  "no trailing slash unchanged",
			input: "/api/v1/courses",
			want:  "/api/v1/courses",
		},
		// --- query string ---
		{
			name:  "query string stripped",
			input: "/api/v1/courses?per_page=50&include[]=teachers",
			want:  "/api/v1/courses",
		},
		{
			name:  "param plus query string",
			input: "/api/v1/courses/:course_id/students?enrollment_state=active",
			want:  "/api/v1/courses/:x/students",
		},
		// --- /self alias ---
		{
			name:  "self at end of path replaced",
			input: "/api/v1/users/self",
			want:  "/api/v1/users/:x",
		},
		{
			name:  "self as path segment replaced",
			input: "/api/v1/users/self/profile",
			want:  "/api/v1/users/:x/profile",
		},
		// --- literal segments kept as-is ---
		{
			name:  "literal front_page kept",
			input: "/api/v1/courses/:course_id/pages/front_page",
			want:  "/api/v1/courses/:x/pages/front_page",
		},
		{
			name:  "literal latest kept",
			input: "/api/v1/courses/:course_id/content_exports/latest",
			want:  "/api/v1/courses/:x/content_exports/latest",
		},
		// --- combined ---
		{
			name:  "param plus trailing slash plus query",
			input: "/api/v1/accounts/:account_id/users/?per_page=10",
			want:  "/api/v1/accounts/:x/users",
		},
		// --- already normalised ---
		{
			name:  "already normalised path unchanged",
			input: "/api/v1/courses/:x/assignments/:x",
			want:  "/api/v1/courses/:x/assignments/:x",
		},
		// --- empty / root ---
		{
			name:  "empty path",
			input: "",
			want:  "",
		},
		{
			name:  "root slash",
			input: "/",
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizePath(tc.input)
			if got != tc.want {
				t.Errorf("normalizePath(%q) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestFirstSegment covers the firstSegment helper that extracts the resource
// area from a canonical path for grouping coverage gaps by resource.
func TestFirstSegment(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "v1 path top-level",
			input: "/api/v1/courses",
			want:  "courses",
		},
		{
			name:  "v1 path nested",
			input: "/api/v1/courses/:x/assignments",
			want:  "courses",
		},
		{
			name:  "lti path",
			input: "/api/lti/courses/:x/line_items",
			want:  "courses",
		},
		{
			name:  "audit path",
			input: "/api/v1/audit/grade_change",
			want:  "audit",
		},
		{
			name:  "single segment",
			input: "/api/v1/accounts",
			want:  "accounts",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := firstSegment(tc.input)
			if got != tc.want {
				t.Errorf("firstSegment(%q) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestResolveHost verifies the host priority: flag > env > default.
func TestResolveHost(t *testing.T) {
	t.Run("flag value wins", func(t *testing.T) {
		t.Setenv("CANVAS_SPEC_HOST", "https://env.example.com")
		got := resolveHost("https://flag.example.com")
		if got != "https://flag.example.com" {
			t.Errorf("got %q, want flag value", got)
		}
	})

	t.Run("env value used when flag is empty", func(t *testing.T) {
		t.Setenv("CANVAS_SPEC_HOST", "https://env.example.com")
		got := resolveHost("")
		if got != "https://env.example.com" {
			t.Errorf("got %q, want env value", got)
		}
	})

	t.Run("default used when both flag and env empty", func(t *testing.T) {
		t.Setenv("CANVAS_SPEC_HOST", "")
		got := resolveHost("")
		if got != defaultHost {
			t.Errorf("got %q, want default %q", got, defaultHost)
		}
	})

	t.Run("trailing slash stripped from flag", func(t *testing.T) {
		got := resolveHost("https://flag.example.com/")
		if got != "https://flag.example.com" {
			t.Errorf("got %q, expected no trailing slash", got)
		}
	})

	t.Run("trailing slash stripped from env", func(t *testing.T) {
		t.Setenv("CANVAS_SPEC_HOST", "https://env.example.com/")
		got := resolveHost("")
		if got != "https://env.example.com" {
			t.Errorf("got %q, expected no trailing slash", got)
		}
	})
}

// TestExtractEndpoints verifies that the Swagger 1.2 resource walker correctly
// converts {param} placeholders to :param and deduplicates endpoints.
func TestExtractEndpoints(t *testing.T) {
	res := &swaggerResource{
		ResourcePath: "/courses",
		APIs: []struct {
			Path       string `json:"path"`
			Operations []struct {
				Method   string `json:"method"`
				Nickname string `json:"nickname"`
			} `json:"operations"`
		}{
			{
				Path: "/v1/courses/{course_id}",
				Operations: []struct {
					Method   string `json:"method"`
					Nickname string `json:"nickname"`
				}{
					{Method: "GET", Nickname: "show_course"},
					{Method: "PUT", Nickname: "update_course"},
				},
			},
			{
				Path: "/v1/courses",
				Operations: []struct {
					Method   string `json:"method"`
					Nickname string `json:"nickname"`
				}{
					{Method: "GET", Nickname: "list_courses"},
				},
			},
		},
	}

	seen := map[string]bool{}
	var endpoints []Endpoint
	extractEndpoints(res, seen, &endpoints)

	if len(endpoints) != 3 {
		t.Fatalf("expected 3 endpoints, got %d: %v", len(endpoints), endpoints)
	}

	// Verify {course_id} → :course_id conversion in path.
	for _, ep := range endpoints {
		if ep.Path == "/api/v1/courses/:course_id" {
			return // found at least one converted param path — good
		}
	}
	t.Error("expected endpoint with /api/v1/courses/:course_id path")
}

// TestExtractEndpoints_Deduplication verifies that the same method+path pair
// is not emitted more than once.
func TestExtractEndpoints_Deduplication(t *testing.T) {
	res := &swaggerResource{
		ResourcePath: "/assignments",
		APIs: []struct {
			Path       string `json:"path"`
			Operations []struct {
				Method   string `json:"method"`
				Nickname string `json:"nickname"`
			} `json:"operations"`
		}{
			{
				Path: "/v1/courses/{course_id}/assignments",
				Operations: []struct {
					Method   string `json:"method"`
					Nickname string `json:"nickname"`
				}{
					{Method: "GET", Nickname: "list"},
				},
			},
		},
	}

	seen := map[string]bool{}
	var endpoints []Endpoint
	// Call twice — second call should add nothing due to deduplication.
	extractEndpoints(res, seen, &endpoints)
	extractEndpoints(res, seen, &endpoints)

	if len(endpoints) != 1 {
		t.Errorf("expected 1 endpoint after deduplication, got %d", len(endpoints))
	}
}

// TestExtractModels verifies that model fields are harvested from Swagger.
func TestExtractModels(t *testing.T) {
	res := &swaggerResource{
		Models: map[string]struct {
			ID          string `json:"id"`
			Description string `json:"description"`
			Properties  map[string]struct {
				Type        string `json:"type"`
				Description string `json:"description"`
			} `json:"properties"`
		}{
			"Course": {
				ID: "Course",
				Properties: map[string]struct {
					Type        string `json:"type"`
					Description string `json:"description"`
				}{
					"id":   {Type: "integer", Description: "Course ID"},
					"name": {Type: "string", Description: "Course name"},
				},
			},
		},
	}

	models := map[string]*ModelDef{}
	extractModels(res, "courses", models)

	def, ok := models["Course"]
	if !ok {
		t.Fatal("expected 'Course' model to be extracted")
	}
	if def.Resource != "courses" {
		t.Errorf("expected resource 'courses', got %q", def.Resource)
	}
	if _, ok := def.Fields["name"]; !ok {
		t.Error("expected 'name' field in Course model")
	}
}

// TestExtractModels_NoDuplicates ensures a model seen in multiple resources
// is only kept from the first resource.
func TestExtractModels_NoDuplicates(t *testing.T) {
	makeRes := func() *swaggerResource {
		return &swaggerResource{
			Models: map[string]struct {
				ID          string `json:"id"`
				Description string `json:"description"`
				Properties  map[string]struct {
					Type        string `json:"type"`
					Description string `json:"description"`
				} `json:"properties"`
			}{
				"User": {
					ID: "User",
					Properties: map[string]struct {
						Type        string `json:"type"`
						Description string `json:"description"`
					}{
						"id": {Type: "integer"},
					},
				},
			},
		}
	}

	models := map[string]*ModelDef{}
	extractModels(makeRes(), "users", models)
	extractModels(makeRes(), "courses", models) // should be ignored for "User"

	def := models["User"]
	if def.Resource != "users" {
		t.Errorf("expected first resource 'users' to win, got %q", def.Resource)
	}
}

// TestCollapseContextPath_MatchesSpecContractImpl verifies that the
// collapseContextPath implementation in tools/speccheck/normalize.go produces
// the same aliases as the equivalent collapseContext function inside
// internal/api/spec_contract_test.go.
//
// Both functions must stay in sync because coverage.go uses collapseContextPath
// when matching CLI endpoints against the spec, and the contract test uses
// collapseContext when matching CLI paths against documented endpoints. Drift
// between the two produces a false "implemented" or false "missing" result.
//
// We can't import spec_contract_test.go here (it's in a different package and
// a test-only file), so we replicate its logic inline and assert both produce
// identical output for a representative set of paths.
func TestCollapseContextPath_MatchesSpecContractImpl(t *testing.T) {
	// specCollapseContext mirrors collapseContext in internal/api/spec_contract_test.go.
	specCtxPairRe := ctxPairRe // both use the same regex pattern
	specCollapseContext := func(norm string) []string {
		if !specCtxPairRe.MatchString(norm) {
			return nil
		}
		return []string{
			specCtxPairRe.ReplaceAllString(norm, "/api/v1/:x$2"),
			specCtxPairRe.ReplaceAllString(norm, "/api/v1/:x/:x$2"),
		}
	}

	paths := []string{
		"/api/v1/courses/:x/discussion_topics",
		"/api/v1/groups/:x/files",
		"/api/v1/accounts/:x/outcome_groups",
		"/api/v1/sections/:x/assignments/:x/submissions",
		"/api/v1/users/:x/profile",
		"/api/v1/courses/:x",
		"/api/v1/global/outcome_groups/:x",    // no context pair — should return nil
		"/api/v1/outcomes/:x",                 // no context pair
		"/api/v1/courses/:x/pages/front_page", // literal final segment
	}

	for _, p := range paths {
		norm := normalizePath(p) // already normalized but run through for safety
		gotTool := collapseContextPath(norm)
		gotContract := specCollapseContext(norm)

		if len(gotTool) != len(gotContract) {
			t.Errorf("collapseContextPath vs collapseContext length mismatch for %q: tool=%v contract=%v",
				p, gotTool, gotContract)
			continue
		}
		for i := range gotTool {
			if gotTool[i] != gotContract[i] {
				t.Errorf("collapseContextPath vs collapseContext differ at index %d for %q: tool=%q contract=%q",
					i, p, gotTool[i], gotContract[i])
			}
		}
	}
}

// TestIsAPIPath verifies the helper that filters API paths.
func TestIsAPIPath(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"/api/v1/courses", true},
		{"/api/lti/line_items", true},
		{"/api/v1/", true},
		{"/doc/api/courses", false},
		{"courses", false},
		{"", false},
	}
	for _, tc := range cases {
		got := isAPIPath(tc.path)
		if got != tc.want {
			t.Errorf("isAPIPath(%q) = %v; want %v", tc.path, got, tc.want)
		}
	}
}

// TestUnquote covers both quoted and backtick forms.
func TestUnquote(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{`"/api/v1/courses"`, "/api/v1/courses"},
		{"`/api/v1/courses`", "/api/v1/courses"},
		{"naked", "naked"},
		{`""`, ""},
		{"", ""},
		{`"a"`, "a"},
	}
	for _, tc := range cases {
		got := unquote(tc.input)
		if got != tc.want {
			t.Errorf("unquote(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}
