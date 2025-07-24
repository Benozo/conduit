package llm

import (
	"fmt"
)

// Router is responsible for routing requests to the appropriate language model provider.
type Router struct {
	providers map[string]ModelProvider
}

// NewRouter creates a new Router instance with the given model providers.
func NewRouter() *Router {
	return &Router{
		providers: make(map[string]ModelProvider),
	}
}

// RegisterProvider registers a language model provider with the router.
func (r *Router) RegisterProvider(name string, provider ModelProvider) error {
	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %s already registered", name)
	}
	r.providers[name] = provider
	return nil
}

// RouteRequest routes a request to the appropriate language model provider based on the agent's requirements.
func (r *Router) RouteRequest(agentName string, request Request) (Response, error) {
	provider, exists := r.providers[request.ProviderName]
	if !exists {
		return Response{}, fmt.Errorf("no provider found for %s", request.ProviderName)
	}
	return provider.HandleRequest(agentName, request)
}

// GetRegisteredProviders returns a list of all registered provider names
func (r *Router) GetRegisteredProviders() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// UnregisterProvider removes a provider from the router
func (r *Router) UnregisterProvider(name string) error {
	if _, exists := r.providers[name]; !exists {
		return fmt.Errorf("provider %s not found", name)
	}
	delete(r.providers, name)
	return nil
}
