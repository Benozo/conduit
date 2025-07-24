package react

import (
	"fmt"
	"strings"
	"time"
)

// ReasoningStep represents a step in the reasoning process
type ReasoningStep struct {
	Step       int                    `json:"step"`
	Thought    string                 `json:"thought"`
	Reasoning  string                 `json:"reasoning"`
	Confidence float64                `json:"confidence"`
	Context    map[string]interface{} `json:"context"`
	Timestamp  time.Time              `json:"timestamp"`
}

// ReasoningResult contains the complete reasoning process and decision
type ReasoningResult struct {
	Decision   string                 `json:"decision"`
	Confidence float64                `json:"confidence"`
	Steps      []*ReasoningStep       `json:"steps"`
	TotalTime  time.Duration          `json:"total_time"`
	Context    map[string]interface{} `json:"context"`
}

// Reasoner represents an agent that analyzes situations and makes decisions based on predefined logic.
type Reasoner struct {
	name                string
	reasoningSteps      []*ReasoningStep
	maxSteps            int
	confidenceThreshold float64
}

// NewReasoner creates a new Reasoner agent with the specified name.
func NewReasoner(name string) *Reasoner {
	return &Reasoner{
		name:                name,
		reasoningSteps:      make([]*ReasoningStep, 0),
		maxSteps:            10,
		confidenceThreshold: 0.7,
	}
}

// NewReasonerWithConfig creates a reasoner with custom configuration
func NewReasonerWithConfig(name string, maxSteps int, confidenceThreshold float64) *Reasoner {
	return &Reasoner{
		name:                name,
		reasoningSteps:      make([]*ReasoningStep, 0),
		maxSteps:            maxSteps,
		confidenceThreshold: confidenceThreshold,
	}
}

// Analyze evaluates the given situation and returns a decision based on predefined logic.
func (r *Reasoner) Analyze(situation string) (string, error) {
	startTime := time.Now()

	// Clear previous reasoning steps
	r.reasoningSteps = make([]*ReasoningStep, 0)

	// Implement analysis logic here
	if situation == "" {
		return "", fmt.Errorf("situation cannot be empty")
	}

	// Perform multi-step reasoning
	result, err := r.performReasoningProcess(situation)
	if err != nil {
		return "", err
	}

	result.TotalTime = time.Since(startTime)

	fmt.Printf("Reasoner %s completed analysis in %v with confidence %.2f\n",
		r.name, result.TotalTime, result.Confidence)

	return result.Decision, nil
}

// performReasoningProcess conducts the step-by-step reasoning
func (r *Reasoner) performReasoningProcess(situation string) (*ReasoningResult, error) {
	context := map[string]interface{}{
		"situation": situation,
		"reasoner":  r.name,
	}

	// Step 1: Initial assessment
	step1 := &ReasoningStep{
		Step:       1,
		Thought:    "Analyzing the given situation",
		Reasoning:  fmt.Sprintf("Initial assessment of situation: %s", situation),
		Confidence: 0.5,
		Context:    context,
		Timestamp:  time.Now(),
	}
	r.reasoningSteps = append(r.reasoningSteps, step1)

	// Step 2: Categorize situation
	category := r.categorizeSituation(situation)
	step2 := &ReasoningStep{
		Step:       2,
		Thought:    fmt.Sprintf("Situation categorized as: %s", category),
		Reasoning:  fmt.Sprintf("Based on keywords and context, this is a %s situation", category),
		Confidence: 0.7,
		Context:    context,
		Timestamp:  time.Now(),
	}
	r.reasoningSteps = append(r.reasoningSteps, step2)

	// Step 3: Generate decision based on category
	decision, confidence := r.generateDecision(category, situation)
	step3 := &ReasoningStep{
		Step:       3,
		Thought:    fmt.Sprintf("Decision reached: %s", decision),
		Reasoning:  fmt.Sprintf("Based on %s categorization, the appropriate action is: %s", category, decision),
		Confidence: confidence,
		Context:    context,
		Timestamp:  time.Now(),
	}
	r.reasoningSteps = append(r.reasoningSteps, step3)

	return &ReasoningResult{
		Decision:   decision,
		Confidence: confidence,
		Steps:      r.reasoningSteps,
		Context:    context,
	}, nil
}

// categorizeSituation categorizes the situation based on keywords
func (r *Reasoner) categorizeSituation(situation string) string {
	lower := strings.ToLower(situation)

	urgentKeywords := []string{"urgent", "emergency", "critical", "immediate", "now", "asap"}
	normalKeywords := []string{"normal", "regular", "standard", "routine", "typical"}
	analysisKeywords := []string{"analyze", "study", "research", "investigate", "examine"}
	actionKeywords := []string{"execute", "perform", "do", "act", "implement"}

	for _, keyword := range urgentKeywords {
		if strings.Contains(lower, keyword) {
			return "urgent"
		}
	}

	for _, keyword := range analysisKeywords {
		if strings.Contains(lower, keyword) {
			return "analytical"
		}
	}

	for _, keyword := range actionKeywords {
		if strings.Contains(lower, keyword) {
			return "action-required"
		}
	}

	for _, keyword := range normalKeywords {
		if strings.Contains(lower, keyword) {
			return "normal"
		}
	}

	return "unknown"
}

// generateDecision creates a decision based on the situation category
func (r *Reasoner) generateDecision(category, situation string) (string, float64) {
	switch category {
	case "urgent":
		return "Take immediate action", 0.9
	case "normal":
		return "Proceed with caution", 0.8
	case "analytical":
		return "Conduct thorough analysis", 0.85
	case "action-required":
		return "Execute planned action", 0.8
	case "unknown":
		return "Gather more information", 0.6
	default:
		return "No action required", 0.7
	}
}

// GetReasoningSteps returns the reasoning steps from the last analysis
func (r *Reasoner) GetReasoningSteps() []*ReasoningStep {
	return r.reasoningSteps
}

// SetMaxSteps configures the maximum reasoning steps
func (r *Reasoner) SetMaxSteps(maxSteps int) {
	r.maxSteps = maxSteps
}

// SetConfidenceThreshold sets the minimum confidence required for decisions
func (r *Reasoner) SetConfidenceThreshold(threshold float64) {
	r.confidenceThreshold = threshold
}

// GetName returns the name of the Reasoner agent.
func (r *Reasoner) GetName() string {
	return r.name
}

// CanReason checks if the reasoner is ready to analyze
func (r *Reasoner) CanReason() bool {
	return r.name != ""
}

// String returns a string representation of the reasoner
func (r *Reasoner) String() string {
	return fmt.Sprintf("Reasoner{Name: %s, Steps: %d, MaxSteps: %d, Confidence: %.2f}",
		r.name, len(r.reasoningSteps), r.maxSteps, r.confidenceThreshold)
}
