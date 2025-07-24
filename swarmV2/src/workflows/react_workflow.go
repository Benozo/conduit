package workflows

import (
	"github.com/benozo/neuron/src/agents/react"
	"github.com/benozo/neuron/src/core"
)

// ReactWorkflow orchestrates the interaction between React agents in a loop.
type ReactWorkflow struct {
	Coordinator core.Agent
	Reasoner    *react.Reasoner
	Actor       *react.Actor
	Observer    *react.Observer
}

// NewReactWorkflow initializes a new ReactWorkflow with the provided agents.
func NewReactWorkflow(coordinator core.Agent, reasoner *react.Reasoner, actor *react.Actor, observer *react.Observer) *ReactWorkflow {
	return &ReactWorkflow{
		Coordinator: coordinator,
		Reasoner:    reasoner,
		Actor:       actor,
		Observer:    observer,
	}
}

// Execute runs the React loop, coordinating the actions of the agents.
func (rw *ReactWorkflow) Execute() {
	// Reasoning phase
	situation := "Analyze current situation and determine next action"
	decision, err := rw.Reasoner.Analyze(situation)
	if err != nil {
		// Handle error gracefully
		decision = "No action determined due to reasoning error"
	}

	// Acting phase
	rw.Actor.Act(decision)

	// Observing phase - monitor the actor's actions
	observation := rw.Observer.Monitor(rw.Actor.Name, decision)
	if observation != nil {
		// Optionally provide feedback based on observation
		rw.Observer.ProvideFeedback(rw.Actor.Name, "Action completed")
	}
}
