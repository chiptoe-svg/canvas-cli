package options

// TelemetryOptions is the base type for telemetry sub-commands that have no flags.
// Separate structs are used to keep per-command intent clear.

// TelemetryEnableOptions contains options for the telemetry enable command
type TelemetryEnableOptions struct{}

// Validate validates the options.
func (o *TelemetryEnableOptions) Validate() error { return nil }

// TelemetryDisableOptions contains options for the telemetry disable command
type TelemetryDisableOptions struct{}

// Validate validates the options.
func (o *TelemetryDisableOptions) Validate() error { return nil }

// TelemetryStatusOptions contains options for the telemetry status command
type TelemetryStatusOptions struct{}

// Validate validates the options.
func (o *TelemetryStatusOptions) Validate() error { return nil }

// TelemetryShowOptions contains options for the telemetry show command
type TelemetryShowOptions struct{}

// Validate validates the options.
func (o *TelemetryShowOptions) Validate() error { return nil }

// TelemetryClearOptions contains options for the telemetry clear command
type TelemetryClearOptions struct{}

// Validate validates the options.
func (o *TelemetryClearOptions) Validate() error { return nil }
