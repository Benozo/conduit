package react

import (
	"fmt"
	"strings"
	"time"
)

// Observation represents an observation made by the Observer
type Observation struct {
	AgentName string                 `json:"agent_name"`
	Action    string                 `json:"action"`
	Result    interface{}            `json:"result"`
	Feedback  string                 `json:"feedback"`
	Quality   ObservationQuality     `json:"quality"`
	Timestamp time.Time              `json:"timestamp"`
	Context   map[string]interface{} `json:"context"`
}

// ObservationQuality represents the quality rating of an observation
type ObservationQuality string

const (
	QualityExcellent ObservationQuality = "excellent"
	QualityGood      ObservationQuality = "good"
	QualityFair      ObservationQuality = "fair"
	QualityPoor      ObservationQuality = "poor"
)

// Observer represents an agent that monitors the actions of other agents.
type Observer struct {
	name          string
	observations  []*Observation
	feedbackRules map[string]string
}

// NewObserver creates a new Observer agent with the given name.
func NewObserver(name string) *Observer {
	return &Observer{
		name:          name,
		observations:  make([]*Observation, 0),
		feedbackRules: make(map[string]string),
	}
}

// Monitor observes the actions of the specified agent and provides feedback.
func (o *Observer) Monitor(agentName string, action string) *Observation {
	observation := &Observation{
		AgentName: agentName,
		Action:    action,
		Quality:   o.evaluateActionQuality(action),
		Timestamp: time.Now(),
		Context: map[string]interface{}{
			"observer": o.name,
		},
	}

	// Generate automatic feedback
	observation.Feedback = o.generateFeedback(action, observation.Quality)

	// Store observation
	o.observations = append(o.observations, observation)

	fmt.Printf("Observer %s is monitoring agent %s performing action: %s (Quality: %s)\n",
		o.name, agentName, action, observation.Quality)

	return observation
}

// MonitorWithResult observes an action and its result
func (o *Observer) MonitorWithResult(agentName string, action string, result interface{}) *Observation {
	observation := o.Monitor(agentName, action)
	observation.Result = result

	// Re-evaluate quality with result information
	observation.Quality = o.evaluateActionQualityWithResult(action, result)
	observation.Feedback = o.generateFeedback(action, observation.Quality)

	return observation
}

// ProvideFeedback gives feedback based on the observed actions of agents.
func (o *Observer) ProvideFeedback(agentName string, feedback string) {
	fmt.Printf("Observer %s provides feedback to agent %s: %s\n", o.name, agentName, feedback)
}

// evaluateActionQuality assesses the quality of an action
func (o *Observer) evaluateActionQuality(action string) ObservationQuality {
	// Simple heuristic-based quality assessment
	if len(action) < 5 {
		return QualityPoor
	}

	// Check for positive action indicators
	positiveIndicators := []string{"analyze", "execute", "implement", "optimize", "improve"}
	for _, indicator := range positiveIndicators {
		if contains(action, indicator) {
			return QualityGood
		}
	}

	// Check for excellent action indicators
	excellentIndicators := []string{"immediate action", "thorough analysis", "planned action"}
	for _, indicator := range excellentIndicators {
		if contains(action, indicator) {
			return QualityExcellent
		}
	}

	return QualityFair
}

// evaluateActionQualityWithResult assesses quality including the result
func (o *Observer) evaluateActionQualityWithResult(action string, result interface{}) ObservationQuality {
	baseQuality := o.evaluateActionQuality(action)

	// Upgrade quality if result indicates success
	if result != nil {
		resultStr := fmt.Sprintf("%v", result)
		if contains(resultStr, "success") || contains(resultStr, "completed") {
			switch baseQuality {
			case QualityFair:
				return QualityGood
			case QualityGood:
				return QualityExcellent
			}
		}
	}

	return baseQuality
}

// generateFeedback creates feedback based on action and quality
func (o *Observer) generateFeedback(action string, quality ObservationQuality) string {
	switch quality {
	case QualityExcellent:
		return "Excellent action! This demonstrates optimal decision-making."
	case QualityGood:
		return "Good action taken. The approach is sound and effective."
	case QualityFair:
		return "Fair action. Consider optimizing the approach for better results."
	case QualityPoor:
		return "Action needs improvement. Consider alternative approaches."
	default:
		return "Action observed and recorded."
	}
}

// AddFeedbackRule adds a custom feedback rule
func (o *Observer) AddFeedbackRule(actionPattern string, feedback string) {
	o.feedbackRules[actionPattern] = feedback
}

// GetObservations returns all observations made by this observer
func (o *Observer) GetObservations() []*Observation {
	return o.observations
}

// GetObservationsForAgent returns observations for a specific agent
func (o *Observer) GetObservationsForAgent(agentName string) []*Observation {
	var agentObservations []*Observation
	for _, obs := range o.observations {
		if obs.AgentName == agentName {
			agentObservations = append(agentObservations, obs)
		}
	}
	return agentObservations
}

// GetLastObservation returns the most recent observation
func (o *Observer) GetLastObservation() *Observation {
	if len(o.observations) == 0 {
		return nil
	}
	return o.observations[len(o.observations)-1]
}

// ClearObservations clears all stored observations
func (o *Observer) ClearObservations() {
	o.observations = make([]*Observation, 0)
}

// GetQualityStats returns statistics about observation quality
func (o *Observer) GetQualityStats() map[ObservationQuality]int {
	stats := make(map[ObservationQuality]int)
	for _, obs := range o.observations {
		stats[obs.Quality]++
	}
	return stats
}

// GetName returns the observer's name
func (o *Observer) GetName() string {
	return o.name
}

// String returns a string representation of the observer
func (o *Observer) String() string {
	return fmt.Sprintf("Observer{Name: %s, Observations: %d}", o.name, len(o.observations))
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}
