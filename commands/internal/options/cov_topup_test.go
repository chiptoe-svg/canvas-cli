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
