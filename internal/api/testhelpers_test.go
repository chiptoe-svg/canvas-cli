package api

import "testing"

// newTestClient constructs a Client pointed at serverURL with standard test
// credentials (token "test-token", 10 requests/sec). It calls t.Fatal if
// construction fails.
func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	client, err := NewClient(ClientConfig{
		BaseURL:        serverURL,
		Token:          "test-token",
		RequestsPerSec: 10,
	})
	if err != nil {
		t.Fatalf("newTestClient: %v", err)
	}
	return client
}
