package react

import (
	"fmt"
)

// Action represents an action that can be performed by an Actor
type Action struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Result      interface{}            `json:"result,omitempty"`
	Error       error                  `json:"error,omitempty"`
}

// Actor represents an agent that performs actions based on decisions made by the Reasoner.
type Actor struct {
	Name          string
	Instructions  string
	Reasoner      *Reasoner
	Observer      *Observer
	ActionHistory []*Action
}

// NewActor initializes a new Actor with the given name and instructions.
func NewActor(name, instructions string, reasoner *Reasoner) *Actor {
	return &Actor{
		Name:          name,
		Instructions:  instructions,
		Reasoner:      reasoner,
		ActionHistory: make([]*Action, 0),
	}
}

// Act performs an action based on the Reasoner's analysis of the situation.
func (a *Actor) Act(situation string) (*Action, error) {
	if a.Reasoner == nil {
		return nil, fmt.Errorf("no reasoner available for actor %s", a.Name)
	}

	// Get decision from reasoner
	decision, err := a.Reasoner.Analyze(situation)
	if err != nil {
		return nil, fmt.Errorf("reasoner analysis failed: %w", err)
	}

	if decision == "" {
		return nil, fmt.Errorf("no decision made by reasoner for situation: %s", situation)
	}

	// Create and execute action
	action := &Action{
		Type:        "decision_based",
		Description: decision,
		Parameters: map[string]interface{}{
			"situation": situation,
			"actor":     a.Name,
		},
	}

	// Perform the action based on the decision
	fmt.Printf("Actor %s is performing action: %s (based on situation: %s)\n",
		a.Name, decision, situation)

	// Simulate action execution
	action.Result = fmt.Sprintf("Action '%s' completed successfully", decision)

	// Store action in history
	a.ActionHistory = append(a.ActionHistory, action)

	// Notify observer if available
	if a.Observer != nil {
		a.Observer.Monitor(a.Name, decision)
	}

	return action, nil
}

// ActWithParameters performs a specific action with given parameters
func (a *Actor) ActWithParameters(actionType string, parameters map[string]interface{}) (*Action, error) {
	action := &Action{
		Type:        actionType,
		Description: fmt.Sprintf("Executing %s", actionType),
		Parameters:  parameters,
	}

	fmt.Printf("Actor %s is performing action: %s with parameters: %v\n",
		a.Name, actionType, parameters)

	// Simulate action execution
	switch actionType {
	case "search":
		action.Result = fmt.Sprintf("Search completed for query: %v", parameters["query"])
	case "calculate":
		action.Result = fmt.Sprintf("Calculation performed: %v", parameters["expression"])
	case "analyze":
		action.Result = fmt.Sprintf("Analysis completed for data: %v", parameters["data"])
	default:
		action.Result = fmt.Sprintf("Generic action '%s' completed", actionType)
	}

	// Store action in history
	a.ActionHistory = append(a.ActionHistory, action)

	// Notify observer if available
	if a.Observer != nil {
		a.Observer.Monitor(a.Name, action.Description)
	}

	return action, nil
}

// SetReasoner updates the Reasoner for the Actor.
func (a *Actor) SetReasoner(reasoner *Reasoner) {
	a.Reasoner = reasoner
}

// SetObserver sets an observer for this actor
func (a *Actor) SetObserver(observer *Observer) {
	a.Observer = observer
}

// GetActionHistory returns the history of actions performed by this actor
func (a *Actor) GetActionHistory() []*Action {
	return a.ActionHistory
}

// ClearActionHistory clears the action history
func (a *Actor) ClearActionHistory() {
	a.ActionHistory = make([]*Action, 0)
}

// GetLastAction returns the most recent action performed
func (a *Actor) GetLastAction() *Action {
	if len(a.ActionHistory) == 0 {
		return nil
	}
	return a.ActionHistory[len(a.ActionHistory)-1]
}

// CanAct checks if the actor is ready to perform actions
func (a *Actor) CanAct() bool {
	return a.Reasoner != nil
}

// String returns a string representation of the actor
func (a *Actor) String() string {
	return fmt.Sprintf("Actor{Name: %s, Actions: %d, HasReasoner: %t, HasObserver: %t}",
		a.Name, len(a.ActionHistory), a.Reasoner != nil, a.Observer != nil)
}
