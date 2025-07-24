package main

import (
	"fmt"

	"github.com/benozo/neuron/src/agents/base"
)

// Custom Workflow Example
// This example shows how to create custom agent interactions and workflows.
func main() {
	fmt.Println("=== Custom Workflow Demo ===")

	// Create different types of agents
	coordinator := base.NewCoordinator()
	specialist1 := base.NewSpecialist("DataProcessor", "Process and transform data")
	specialist2 := base.NewSpecialist("Validator", "Validate processed data")

	// Setup the coordinator with specialists
	fmt.Println("Setting up custom workflow...")
	coordinator.RegisterAgent(specialist1)
	coordinator.RegisterAgent(specialist2)

	// Display coordinator information
	fmt.Printf("\nCoordinator Status: %s\n", coordinator.GetStatus())
	fmt.Printf("Registered Agents: %d\n", len(coordinator.GetAgents()))

	// List all agents managed by coordinator
	fmt.Println("\nManaged Agents:")
	for name, agent := range coordinator.GetAgents() {
		fmt.Printf("- %s: %s\n", name, agent.GetName())
	}

	// Create and execute a custom workflow
	workflow := base.NewWorkflow("CustomDataWorkflow", []base.Agent{specialist1, specialist2})

	fmt.Println("\nExecuting custom workflow...")
	err := workflow.Execute()
	if err != nil {
		fmt.Printf("Workflow execution error: %v\n", err)
		return
	}

	fmt.Println("Custom workflow completed successfully!")

	// Show final metrics
	metrics := coordinator.GetMetrics()
	fmt.Printf("\nFinal Metrics:\n")
	fmt.Printf("- Workflows Executed: %d\n", metrics.WorkflowsExecuted)
	fmt.Printf("- Agents Managed: %d\n", metrics.AgentsRegistered)
}
