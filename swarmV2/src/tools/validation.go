package tools

import (
	"errors"
	"regexp"
)

// ValidateAgentName checks if the provided agent name is valid.
func ValidateAgentName(name string) error {
	if name == "" {
		return errors.New("agent name cannot be empty")
	}
	if len(name) > 50 {
		return errors.New("agent name cannot exceed 50 characters")
	}
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, name); !matched {
		return errors.New("agent name can only contain alphanumeric characters and underscores")
	}
	return nil
}

// ValidateWorkflow checks if the provided workflow is valid.
func ValidateWorkflow(workflow string) error {
	if workflow == "" {
		return errors.New("workflow cannot be empty")
	}
	// Additional validation logic can be added here
	return nil
}

// ValidateTool checks if the provided tool name is valid.
func ValidateTool(name string) error {
	if name == "" {
		return errors.New("tool name cannot be empty")
	}
	// Additional validation logic can be added here
	return nil
}