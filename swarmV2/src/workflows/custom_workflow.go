package workflows

import (
	"github.com/benozo/neuron/src/agents/base"
	"github.com/benozo/neuron/src/memory"
)

// CustomWorkflow defines a structure for custom workflows
type CustomWorkflow struct {
	Coordinator *base.Coordinator
	Specialists []*base.Specialist
	Memory      *memory.SharedMemory
}

// NewCustomWorkflow initializes a new custom workflow with the given coordinator and specialists
func NewCustomWorkflow(coordinator *base.Coordinator, specialists []*base.Specialist, memory *memory.SharedMemory) *CustomWorkflow {
	return &CustomWorkflow{
		Coordinator: coordinator,
		Specialists: specialists,
		Memory:      memory,
	}
}

// Execute runs the custom workflow, coordinating the specialists through the coordinator
func (cw *CustomWorkflow) Execute() error {
	// Implementation of the workflow execution logic
	// This could involve delegating tasks to specialists and managing their interactions
	return nil
}

// AddSpecialist adds a new specialist to the workflow
func (cw *CustomWorkflow) AddSpecialist(specialist *base.Specialist) {
	cw.Specialists = append(cw.Specialists, specialist)
}

// RemoveSpecialist removes a specialist from the workflow
func (cw *CustomWorkflow) RemoveSpecialist(specialist *base.Specialist) {
	for i, s := range cw.Specialists {
		if s == specialist {
			cw.Specialists = append(cw.Specialists[:i], cw.Specialists[i+1:]...)
			break
		}
	}
}
