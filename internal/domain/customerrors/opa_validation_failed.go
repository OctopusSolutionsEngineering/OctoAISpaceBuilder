package customerrors

import "fmt"

// OpaValidationFailed is a custom error type that provides details about OPA policy validation failures.
type OpaValidationFailed struct {
	ExitCode   int
	DecisionID string
	Path       string
	Message    string
}

// Error implements the error interface for OpaValidationFailed.
func (e OpaValidationFailed) Error() string {
	if e.DecisionID != "" && e.Path != "" {
		return fmt.Sprintf("OPA policy check failed: %s (decision: %s, path: %s, exit code: %d)",
			e.Message, e.DecisionID, e.Path, e.ExitCode)
	}
	return fmt.Sprintf("OPA policy check failed: %s (exit code: %d)", e.Message, e.ExitCode)
}
