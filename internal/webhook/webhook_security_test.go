package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestListener_NoVerification_AcceptsRequests verifies that a listener
// without any secret or JWK set accepts incoming events (the verification
// warning is emitted by the command layer, not the library).
func TestListener_NoVerification_AcceptsRequests(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	listener := New(&Config{
		Addr:   "127.0.0.1:0",
		Logger: logger,
	})

	eventReceived := false
	listener.On("test_event", func(ctx context.Context, event *Event) error {
		eventReceived = true
		return nil
	})

	payload, _ := json.Marshal(Event{
		ID:        "e1",
		EventType: "test_event",
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	listener.handleWebhook(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if !eventReceived {
		t.Error("expected event handler to be called")
	}
}

// TestListener_BodyLogging_NoFullBody verifies that the default logger does not
// dump the full event body, protecting potential PII at default verbosity.
func TestListener_BodyLogging_NoFullBody(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	listener := New(&Config{
		Addr:   "127.0.0.1:0",
		Secret: "secret",
		Logger: logger,
	})
	listener.On("sensitive_event", func(ctx context.Context, event *Event) error {
		return nil
	})

	eventData := map[string]interface{}{
		"id":         "evt-123",
		"event_type": "sensitive_event",
		"body": map[string]interface{}{
			"user_ssn":   "123-45-6789",
			"user_email": "secret@example.com",
		},
	}
	payload, _ := json.Marshal(eventData)

	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write(payload)
	sig := hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(payload))
	req.Header.Set("X-Canvas-Signature", sig)
	rr := httptest.NewRecorder()

	listener.handleWebhook(rr, req)

	logged := logBuf.String()
	// The full body values must not appear in log output at default verbosity.
	if strings.Contains(logged, "123-45-6789") {
		t.Errorf("SSN should not appear in log output: %s", logged)
	}
	if strings.Contains(logged, "secret@example.com") {
		t.Errorf("email should not appear in log output: %s", logged)
	}
	// The event type and ID should be logged (metadata is fine).
	if !strings.Contains(logged, "sensitive_event") {
		t.Errorf("event type should appear in log output: %s", logged)
	}
}

// TestListener_DefaultAddr documents the expected default bind address.
// The actual binding is handled by the command layer; this test verifies
// that a Config created with an explicit loopback addr is accepted.
func TestListener_DefaultAddr_Loopback(t *testing.T) {
	listener := New(&Config{
		Addr: "127.0.0.1:8080",
	})
	if listener.addr != "127.0.0.1:8080" {
		t.Errorf("expected addr 127.0.0.1:8080, got %s", listener.addr)
	}
}
