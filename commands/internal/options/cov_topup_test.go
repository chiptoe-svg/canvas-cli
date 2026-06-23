package options

import "testing"

// TestTopup_APIOptions_Validate covers the APIOptions.Validate no-op.
func TestTopup_APIOptions_Validate(t *testing.T) {
	opts := &APIOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("APIOptions.Validate() = %v, want nil", err)
	}
}

// TestTopup_CacheClearOptions_Validate covers CacheClearOptions.Validate no-op.
func TestTopup_CacheClearOptions_Validate(t *testing.T) {
	opts := &CacheClearOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("CacheClearOptions.Validate() = %v, want nil", err)
	}
}

// TestTopup_CompletionOptions_Validate covers CompletionOptions.Validate no-op.
func TestTopup_CompletionOptions_Validate(t *testing.T) {
	opts := &CompletionOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("CompletionOptions.Validate() = %v, want nil", err)
	}
}

// TestTopup_ReplOptions_Validate covers ReplOptions.Validate no-op.
func TestTopup_ReplOptions_Validate(t *testing.T) {
	opts := &ReplOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("ReplOptions.Validate() = %v, want nil", err)
	}
}

// TestTopup_SyncAssignmentsOptions_Validate covers SyncAssignmentsOptions.Validate no-op.
func TestTopup_SyncAssignmentsOptions_Validate(t *testing.T) {
	opts := &SyncAssignmentsOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("SyncAssignmentsOptions.Validate() = %v, want nil", err)
	}
}

// TestTopup_SyncCourseOptions_Validate covers SyncCourseOptions.Validate no-op.
func TestTopup_SyncCourseOptions_Validate(t *testing.T) {
	opts := &SyncCourseOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("SyncCourseOptions.Validate() = %v, want nil", err)
	}
}

// TestTopup_TelemetryEnableOptions_Validate covers TelemetryEnableOptions.Validate no-op.
func TestTopup_TelemetryEnableOptions_Validate(t *testing.T) {
	opts := &TelemetryEnableOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("TelemetryEnableOptions.Validate() = %v, want nil", err)
	}
}

// TestTopup_TelemetryDisableOptions_Validate covers TelemetryDisableOptions.Validate no-op.
func TestTopup_TelemetryDisableOptions_Validate(t *testing.T) {
	opts := &TelemetryDisableOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("TelemetryDisableOptions.Validate() = %v, want nil", err)
	}
}

// TestTopup_TelemetryStatusOptions_Validate covers TelemetryStatusOptions.Validate no-op.
func TestTopup_TelemetryStatusOptions_Validate(t *testing.T) {
	opts := &TelemetryStatusOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("TelemetryStatusOptions.Validate() = %v, want nil", err)
	}
}

// TestTopup_TelemetryShowOptions_Validate covers TelemetryShowOptions.Validate no-op.
func TestTopup_TelemetryShowOptions_Validate(t *testing.T) {
	opts := &TelemetryShowOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("TelemetryShowOptions.Validate() = %v, want nil", err)
	}
}

// TestTopup_TelemetryClearOptions_Validate covers TelemetryClearOptions.Validate no-op.
func TestTopup_TelemetryClearOptions_Validate(t *testing.T) {
	opts := &TelemetryClearOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("TelemetryClearOptions.Validate() = %v, want nil", err)
	}
}
