package base

import (
	"fmt"
)

// Specialist represents an agent that performs specialized tasks based on workflow requirements.
type Specialist struct {
	name           string
	specialization string
}

// NewSpecialist creates a new Specialist agent with the given name and specialization.
func NewSpecialist(name string, specialization string) *Specialist {
	return &Specialist{
		name:           name,
		specialization: specialization,
	}
}

// Execute implements the Agent interface
func (s *Specialist) Execute() error {
	fmt.Printf("Specialist '%s' executing task with specialization: %s\n", s.name, s.specialization)
	// Implement the logic for performing the specialized task here.
	// This is a placeholder for actual task execution logic.
	return nil
}

// GetName implements the Agent interface
func (s *Specialist) GetName() string {
	return s.name
}

// GetSpecialization returns the agent's specialization
func (s *Specialist) GetSpecialization() string {
	return s.specialization
}

// PerformTask executes a specialized task based on the provided parameters.
func (s *Specialist) PerformTask(params map[string]interface{}) (interface{}, error) {
	fmt.Printf("Specialist '%s' performing specialized task: %s\n", s.name, s.specialization)
	// Implement the logic for performing the specialized task here.
	// This is a placeholder for actual task execution logic.
	return nil, nil
}
