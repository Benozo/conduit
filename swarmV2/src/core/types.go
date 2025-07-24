package core

// Common types and interfaces used throughout the core components of the swarm framework.

type AgentType string

const (
	Coordinator AgentType = "Coordinator"
	Specialist  AgentType = "Specialist"
	Process     AgentType = "Process"
	Retriever   AgentType = "Retriever"
	Generator   AgentType = "Generator"
	Evaluator   AgentType = "Evaluator"
	Reasoner    AgentType = "Reasoner"
	Actor       AgentType = "Actor"
	Observer    AgentType = "Observer"
)

type Agent interface {
	GetName() string
	GetType() AgentType
	Execute(input interface{}) (interface{}, error)
}

type Tool interface {
	GetName() string
	GetDescription() string
	Execute(parameters map[string]interface{}) (interface{}, error)
	Validate(parameters map[string]interface{}) error
}

type WorkflowInterface interface {
	Start() error
	Stop() error
	AddAgent(agent Agent) error
	RemoveAgent(agentName string) error
}
