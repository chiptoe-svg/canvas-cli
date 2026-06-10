package config

import (
	"strings"
	"testing"
)

// TestValidateInstance_HTTPSchemeRejectedForNonLoopback verifies that http://
// is rejected for non-loopback hosts to prevent Bearer tokens being sent
// in cleartext over the network.
func TestValidateInstance_HTTPSchemeRejectedForNonLoopback(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
		errFrag string
	}{
		// Non-loopback http:// must be rejected
		{
			name:    "http remote host rejected",
			url:     "http://canvas.example.com",
			wantErr: true,
			errFrag: "http://",
		},
		{
			name:    "http remote IP rejected",
			url:     "http://192.168.1.100",
			wantErr: true,
			errFrag: "http://",
		},
		// Loopback http:// must be accepted
		{
			name:    "http localhost allowed",
			url:     "http://localhost",
			wantErr: false,
		},
		{
			name:    "http localhost with port allowed",
			url:     "http://localhost:3000",
			wantErr: false,
		},
		{
			name:    "http 127.0.0.1 allowed",
			url:     "http://127.0.0.1",
			wantErr: false,
		},
		{
			name:    "http 127.0.0.1 with port allowed",
			url:     "http://127.0.0.1:8080",
			wantErr: false,
		},
		{
			name:    "http ::1 allowed",
			url:     "http://[::1]:3000",
			wantErr: false,
		},
		// https:// is always accepted
		{
			name:    "https remote host allowed",
			url:     "https://canvas.example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &Instance{
				Name: "test",
				URL:  tt.url,
			}
			err := ValidateInstance(inst)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateInstance(%q): expected error, got nil", tt.url)
					return
				}
				if tt.errFrag != "" && !strings.Contains(err.Error(), tt.errFrag) {
					t.Errorf("error %q should mention %q", err.Error(), tt.errFrag)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateInstance(%q): unexpected error: %v", tt.url, err)
				}
			}
		})
	}
}

// TestIsLoopbackHost covers the loopback detection helper directly.
func TestIsLoopbackHost(t *testing.T) {
	tests := []struct {
		host string
		want bool
	}{
		{"localhost", true},
		{"127.0.0.1", true},
		{"127.1.2.3", true},
		{"::1", true},
		{"canvas.instructure.com", false},
		{"192.168.1.1", false},
		{"10.0.0.1", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			got := isLoopbackHost(tt.host)
			if got != tt.want {
				t.Errorf("isLoopbackHost(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}
