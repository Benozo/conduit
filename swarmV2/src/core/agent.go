package core

// AgentFunction represents a callable function available to an agent
type AgentFunction struct {
	Name        string                                 // Function name
	Description string                                 // Function description
	Parameters  map[string]interface{}                 // Function parameters schema
	Handler     func(interface{}) (interface{}, error) // Function implementation
}

// BaseAgent is a concrete implementation of the Agent interface
type BaseAgent struct {
	Name         string          // Agent identifier
	AgentType    AgentType       // Agent type
	Instructions string          // Behavior instructions
	Functions    []AgentFunction // Available functions
	Model        string          // LLM model to use
	ToolChoice   string          // Tool selection strategy
}

// NewAgent initializes a new BaseAgent with the given parameters.
func NewAgent(name, instructions string, functions []AgentFunction, model, toolChoice string) *BaseAgent {
	return &BaseAgent{
		Name:         name,
		AgentType:    Specialist, // Default to Specialist type
		Instructions: instructions,
		Functions:    functions,
		Model:        model,
		ToolChoice:   toolChoice,
	}
}

// GetName returns the agent's name (implements Agent interface)
func (a *BaseAgent) GetName() string {
	return a.Name
}

// GetType returns the agent's type (implements Agent interface)
func (a *BaseAgent) GetType() AgentType {
	return a.AgentType
}

// Execute processes input and returns output (implements Agent interface)
func (a *BaseAgent) Execute(input interface{}) (interface{}, error) {
	// Implementation of execution logic goes here
	return input, nil
}

// Interact allows the agent to interact with other agents or tools.
func (a *BaseAgent) Interact(input string) (string, error) {
	// Implementation of interaction logic goes here
	return "", nil
}

// RegisterFunction adds a new function to the agent's capabilities.
func (a *BaseAgent) RegisterFunction(fn AgentFunction) {
	a.Functions = append(a.Functions, fn)
}
