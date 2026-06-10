package auth

import (
	"testing"
)

// TestStaticToken_SaveLoadDelete exercises the full lifecycle of a static API
// token stored via the file-backed token store.
func TestStaticToken_SaveLoadDelete(t *testing.T) {
	t.Setenv("CANVAS_CLI_MACHINE_ID", "test-machine-id")

	store := NewFileTokenStore(t.TempDir())

	const instance = "myschool"
	const apiToken = "7~some-test-api-token-that-is-long-enough"

	// Nothing stored yet
	if StaticTokenExists(store, instance) {
		t.Fatal("expected no static token before Save")
	}

	loaded, err := LoadStaticToken(store, instance)
	if err != nil {
		t.Fatalf("LoadStaticToken (not found): unexpected error: %v", err)
	}
	if loaded != "" {
		t.Fatalf("LoadStaticToken (not found): expected empty string, got %q", loaded)
	}

	// Save
	if err := SaveStaticToken(store, instance, apiToken); err != nil {
		t.Fatalf("SaveStaticToken: %v", err)
	}

	// Exists
	if !StaticTokenExists(store, instance) {
		t.Fatal("expected static token to exist after Save")
	}

	// Load
	got, err := LoadStaticToken(store, instance)
	if err != nil {
		t.Fatalf("LoadStaticToken: %v", err)
	}
	if got != apiToken {
		t.Errorf("LoadStaticToken: got %q, want %q", got, apiToken)
	}

	// Delete
	if err := DeleteStaticToken(store, instance); err != nil {
		t.Fatalf("DeleteStaticToken: %v", err)
	}

	// Gone
	if StaticTokenExists(store, instance) {
		t.Fatal("expected static token to be gone after Delete")
	}

	loaded2, err := LoadStaticToken(store, instance)
	if err != nil {
		t.Fatalf("LoadStaticToken after delete: unexpected error: %v", err)
	}
	if loaded2 != "" {
		t.Errorf("LoadStaticToken after delete: expected empty string, got %q", loaded2)
	}
}

// TestStaticToken_DoesNotCollideWithOAuthKey confirms that storing a static
// token does not affect the OAuth token for the same instance name, and vice
// versa.
func TestStaticToken_DoesNotCollideWithOAuthKey(t *testing.T) {
	t.Setenv("CANVAS_CLI_MACHINE_ID", "test-machine-id")

	store := NewFileTokenStore(t.TempDir())

	const instance = "myschool"
	const apiToken = "7~static-api-token-that-is-long-enough-here"

	// Save a static token
	if err := SaveStaticToken(store, instance, apiToken); err != nil {
		t.Fatalf("SaveStaticToken: %v", err)
	}

	// The OAuth store key (plain instance name) should be absent
	if store.Exists(instance) {
		t.Error("SaveStaticToken should not write to the plain instance key")
	}

	// The static key (instance + suffix) should be present
	if !StaticTokenExists(store, instance) {
		t.Error("SaveStaticToken should write to the static instance key")
	}
}

// TestStaticTokenKey verifies the key name used for storage.
func TestStaticTokenKey(t *testing.T) {
	got := staticInstanceKey("myschool")
	want := "myschool" + staticTokenSuffix
	if got != want {
		t.Errorf("staticInstanceKey = %q, want %q", got, want)
	}
}
