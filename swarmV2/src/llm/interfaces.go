package llm

// LanguageModelProvider defines the interface for language model providers.
type LanguageModelProvider interface {
	GenerateResponse(prompt string) (string, error)
	SetModel(model string)
}

// MultiModelProvider defines the interface for providers that can handle multiple models.
type MultiModelProvider interface {
	LanguageModelProvider
	GetAvailableModels() []string
}

// LLMRouter defines the interface for routing requests to the appropriate language model provider.
type LLMRouter interface {
	RouteRequest(providerName string, prompt string) (string, error)
}

// ModelInfo contains information about a language model
type ModelInfo struct {
	Name        string
	Version     string
	Provider    string
	MaxTokens   int
	Description string
}

// Request represents a request to a language model provider
type Request struct {
	ProviderName string
	Prompt       string
	Model        string
	MaxTokens    int
	Temperature  float64
	Parameters   map[string]interface{}
}

// Response represents a response from a language model provider
type Response struct {
	Content    string
	TokensUsed int
	Model      string
	Provider   string
	Error      error
}

// ModelProvider defines the interface for handling model requests
type ModelProvider interface {
	HandleRequest(agentName string, request Request) (Response, error)
	GetCapabilities() []string
}
