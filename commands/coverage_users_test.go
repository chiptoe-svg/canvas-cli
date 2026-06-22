package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cmdtest "github.com/jjuanrivvera/canvas-cli/commands/internal/testing"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// ---------------------------------------------------------------------------
// readUserJSON helper
// ---------------------------------------------------------------------------

func TestReadUserJSON_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "user.json")
	content := []byte(`{"name":"From File"}`)
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatal(err)
	}

	got, err := readUserJSON(path, false)
	if err != nil {
		t.Fatalf("readUserJSON: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("got %q, want %q", got, content)
	}
}

func TestReadUserJSON_NeitherFileNorStdin(t *testing.T) {
	got, err := readUserJSON("", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %q", got)
	}
}

func TestReadUserJSON_FileNotFound(t *testing.T) {
	_, err := readUserJSON("/nonexistent/path/user.json", false)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

// ---------------------------------------------------------------------------
// parseUserCreateJSON
// ---------------------------------------------------------------------------

func TestParseUserCreateJSON_AllFields(t *testing.T) {
	data := []byte(`{
		"name": "Full Name",
		"short_name": "Short",
		"sortable_name": "Name, Full",
		"email": "full@example.com",
		"login_id": "full.name",
		"password": "secret",
		"sis_user_id": "SIS123",
		"time_zone": "America/Chicago",
		"locale": "es",
		"skip_registration": true,
		"skip_confirmation": true
	}`)

	params := &api.CreateUserParams{}
	if err := parseUserCreateJSON(data, params); err != nil {
		t.Fatalf("parseUserCreateJSON: %v", err)
	}

	if params.Name != "Full Name" {
		t.Errorf("Name = %q, want %q", params.Name, "Full Name")
	}
	if params.ShortName != "Short" {
		t.Errorf("ShortName = %q, want %q", params.ShortName, "Short")
	}
	if params.SortableName != "Name, Full" {
		t.Errorf("SortableName = %q, want %q", params.SortableName, "Name, Full")
	}
	if params.Email != "full@example.com" {
		t.Errorf("Email = %q, want %q", params.Email, "full@example.com")
	}
	if params.UniqueID != "full.name" {
		t.Errorf("UniqueID = %q, want %q", params.UniqueID, "full.name")
	}
	if params.Password != "secret" {
		t.Errorf("Password = %q, want %q", params.Password, "secret")
	}
	if params.SISUserID != "SIS123" {
		t.Errorf("SISUserID = %q, want %q", params.SISUserID, "SIS123")
	}
	if params.TimeZone != "America/Chicago" {
		t.Errorf("TimeZone = %q, want %q", params.TimeZone, "America/Chicago")
	}
	if params.Locale != "es" {
		t.Errorf("Locale = %q, want %q", params.Locale, "es")
	}
	if !params.SkipRegistration {
		t.Error("Expected SkipRegistration = true")
	}
	if !params.SkipConfirmation {
		t.Error("Expected SkipConfirmation = true")
	}
}

func TestParseUserCreateJSON_InvalidJSON(t *testing.T) {
	params := &api.CreateUserParams{}
	err := parseUserCreateJSON([]byte(`not json`), params)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestParseUserCreateJSON_EmptyFields_NotOverwritten(t *testing.T) {
	// When JSON has empty strings the existing params should not be overwritten
	params := &api.CreateUserParams{Name: "Existing Name"}
	data := []byte(`{}`)
	if err := parseUserCreateJSON(data, params); err != nil {
		t.Fatalf("parseUserCreateJSON: %v", err)
	}
	// Name should remain unchanged because JSON had no "name" field
	if params.Name != "Existing Name" {
		t.Errorf("Name = %q, expected %q", params.Name, "Existing Name")
	}
}

// ---------------------------------------------------------------------------
// parseUserUpdateJSON
// ---------------------------------------------------------------------------

func TestParseUserUpdateJSON_AllFields(t *testing.T) {
	data := []byte(`{
		"name": "Updated Name",
		"short_name": "Upd",
		"sortable_name": "Name, Updated",
		"email": "updated@example.com",
		"time_zone": "America/Denver",
		"locale": "fr"
	}`)

	params := &api.UpdateUserParams{}
	if err := parseUserUpdateJSON(data, params); err != nil {
		t.Fatalf("parseUserUpdateJSON: %v", err)
	}

	if params.Name != "Updated Name" {
		t.Errorf("Name = %q, want %q", params.Name, "Updated Name")
	}
	if params.ShortName != "Upd" {
		t.Errorf("ShortName = %q, want %q", params.ShortName, "Upd")
	}
	if params.SortableName != "Name, Updated" {
		t.Errorf("SortableName = %q, want %q", params.SortableName, "Name, Updated")
	}
	if params.Email != "updated@example.com" {
		t.Errorf("Email = %q, want %q", params.Email, "updated@example.com")
	}
	if params.TimeZone != "America/Denver" {
		t.Errorf("TimeZone = %q, want %q", params.TimeZone, "America/Denver")
	}
	if params.Locale != "fr" {
		t.Errorf("Locale = %q, want %q", params.Locale, "fr")
	}
}

func TestParseUserUpdateJSON_InvalidJSON(t *testing.T) {
	params := &api.UpdateUserParams{}
	err := parseUserUpdateJSON([]byte(`{bad json`), params)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestParseUserUpdateJSON_EmptyFields_NotOverwritten(t *testing.T) {
	params := &api.UpdateUserParams{Name: "Keep Me"}
	data := []byte(`{"short_name":"Only Short"}`)
	if err := parseUserUpdateJSON(data, params); err != nil {
		t.Fatalf("parseUserUpdateJSON: %v", err)
	}
	if params.Name != "Keep Me" {
		t.Errorf("Name = %q, expected %q", params.Name, "Keep Me")
	}
	if params.ShortName != "Only Short" {
		t.Errorf("ShortName = %q, expected %q", params.ShortName, "Only Short")
	}
}

// ---------------------------------------------------------------------------
// Command: users create — JSON stdin path via file (stdin is harder to inject)
// ---------------------------------------------------------------------------

func TestUsersCreateCmd_JSONFileWithSkipFlags(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "user.json")
	jsonContent := `{
		"name":"Stdin User",
		"email":"stdin@example.com",
		"skip_registration": true,
		"skip_confirmation": true
	}`
	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0600); err != nil {
		t.Fatal(err)
	}

	tc := cmdtest.CommandTestCase{
		Name: "create user from JSON file with skip flags",
		Args: []string{"--account-id", "1", "--json-file", jsonPath},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewMockResponse(`{
				"id": 300,
				"name": "Stdin User",
				"email": "stdin@example.com"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Stdin User") {
				t.Error("Expected 'Stdin User' in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersCreateCmd(), tc)
}

func TestUsersCreateCmd_InvalidJSONContents(t *testing.T) {
	dir := t.TempDir()
	badPath := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(badPath, []byte(`not valid json at all`), 0600); err != nil {
		t.Fatal(err)
	}

	tc := cmdtest.CommandTestCase{
		Name: "create user from bad JSON file",
		Args: []string{"--account-id", "1", "--json-file", badPath},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersCreateCmd(), tc)
}

func TestUsersCreateCmd_AllOptionalFlags(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "create user with all optional flags",
		Args: []string{
			"--account-id", "1",
			"--name", "Full User",
			"--short-name", "Full",
			"--sortable-name", "User, Full",
			"--email", "full@example.com",
			"--login-id", "full.user",
			"--password", "S3cr3t!",
			"--sis-user-id", "SIS001",
			"--timezone", "UTC",
			"--locale", "en",
			"--skip-registration",
			"--skip-confirmation",
		},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewMockResponse(`{
				"id": 42,
				"name": "Full User",
				"login_id": "full.user",
				"email": "full@example.com"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "User created successfully") {
				t.Error("Expected success message in output")
			}
			if !strings.Contains(output, "full@example.com") {
				t.Error("Expected email in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersCreateCmd(), tc)
}

// ---------------------------------------------------------------------------
// Command: users update — JSON file invalid contents path
// ---------------------------------------------------------------------------

func TestUsersUpdateCmd_InvalidJSONContents(t *testing.T) {
	dir := t.TempDir()
	badPath := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(badPath, []byte(`{invalid json`), 0600); err != nil {
		t.Fatal(err)
	}

	tc := cmdtest.CommandTestCase{
		Name: "update user from bad JSON file",
		Args: []string{"100", "--json", badPath},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/100": cmdtest.NewMockResponse(`{}`),
		},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersUpdateCmd(), tc)
}

func TestUsersUpdateCmd_AllOptionalFlags(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "update user with all optional flags",
		Args: []string{
			"50",
			"--name", "Updated Full",
			"--short-name", "UpFull",
			"--sortable-name", "Full, Updated",
			"--email", "upfull@example.com",
			"--timezone", "America/Los_Angeles",
			"--locale", "de",
		},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/50": cmdtest.NewMockResponse(`{
				"id": 50,
				"name": "Updated Full",
				"email": "upfull@example.com"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "User updated successfully") {
				t.Error("Expected success message in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersUpdateCmd(), tc)
}

// ---------------------------------------------------------------------------
// Command: users get — invalid arg
// ---------------------------------------------------------------------------

func TestUsersGetCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "get user - non-numeric ID",
		Args:        []string{"not-a-number"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersGetCmd(), tc)
}

// ---------------------------------------------------------------------------
// Command: users update — invalid arg
// ---------------------------------------------------------------------------

func TestUsersUpdateCmd_InvalidID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name:        "update user - non-numeric ID",
		Args:        []string{"not-a-number", "--name", "Test"},
		ExpectError: true,
	}
	cmdtest.RunCommandTest(t, newUsersUpdateCmd(), tc)
}

// ---------------------------------------------------------------------------
// Command: users list — with search flag (account context)
// ---------------------------------------------------------------------------

func TestUsersListCmd_WithSearchAndEnrollmentFilters(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "list users with search and enrollment filters",
		Args: []string{
			"--account-id", "1",
			"--search", "alice",
			"--enrollment-type", "student",
			"--enrollment-state", "active",
			"--include", "email,avatar_url",
		},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/accounts/1/users": cmdtest.NewMockResponse(`[
				{"id": 7, "name": "Alice Smith", "email": "alice@example.com"}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Alice Smith") {
				t.Error("Expected 'Alice Smith' in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersListCmd(), tc)
}

// ---------------------------------------------------------------------------
// Command: users me — verifies login_id present output
// ---------------------------------------------------------------------------

func TestUsersMeCmd_WithLoginID(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "get current user with login_id",
		Args: []string{},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/users/self": cmdtest.NewMockResponse(`{
				"id": 1,
				"name": "Me User",
				"login_id": "meuser",
				"email": "me@example.com"
			}`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "Me User") {
				t.Error("Expected 'Me User' in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersMeCmd(), tc)
}

// ---------------------------------------------------------------------------
// Command: users search — verifies non-empty result path fully executes
// ---------------------------------------------------------------------------

func TestUsersSearchCmd_MultipleResults(t *testing.T) {
	tc := cmdtest.CommandTestCase{
		Name: "search users - multiple results",
		Args: []string{"john"},
		MockResponses: map[string]cmdtest.MockResponse{
			"/api/v1/search/recipients": cmdtest.NewMockResponse(`[
				{"id": 1, "name": "John Alpha"},
				{"id": 2, "name": "John Beta"}
			]`),
		},
		ExpectError: false,
		ValidateOutput: func(t *testing.T, output string) {
			if !strings.Contains(output, "John Alpha") {
				t.Error("Expected 'John Alpha' in output")
			}
		},
	}
	cmdtest.RunCommandTest(t, newUsersSearchCmd(), tc)
}
