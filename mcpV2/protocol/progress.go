package protocol

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ProgressTracker manages progress notifications for MCP operations
type ProgressTracker struct {
	activeTokens map[string]*ProgressToken
	mu           sync.RWMutex
	notifier     ProgressNotifier
}

// ProgressToken represents an active progress tracking token
type ProgressToken struct {
	Token     string    `json:"token"`
	Progress  float64   `json:"progress"` // 0.0 to 1.0
	Total     int64     `json:"total,omitempty"`
	StartTime time.Time `json:"startTime"`
	ctx       context.Context
	cancel    context.CancelFunc
}

// ProgressNotifier sends progress notifications
type ProgressNotifier interface {
	SendProgress(ctx context.Context, token string, progress float64, total int64) error
}

// ProgressNotification represents a progress notification message
type ProgressNotification struct {
	ProgressToken string  `json:"progressToken"`
	Progress      float64 `json:"progress"`
	Total         int64   `json:"total,omitempty"`
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(notifier ProgressNotifier) *ProgressTracker {
	return &ProgressTracker{
		activeTokens: make(map[string]*ProgressToken),
		notifier:     notifier,
	}
}

// StartProgress begins tracking a new progress token
func (pt *ProgressTracker) StartProgress(ctx context.Context, token string) *ProgressToken {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	// Cancel existing token with same name
	if existing, exists := pt.activeTokens[token]; exists {
		existing.cancel()
	}

	ctx, cancel := context.WithCancel(ctx)
	progressToken := &ProgressToken{
		Token:     token,
		Progress:  0.0,
		StartTime: time.Now(),
		ctx:       ctx,
		cancel:    cancel,
	}

	pt.activeTokens[token] = progressToken
	return progressToken
}

// UpdateProgress updates the progress for a token
func (pt *ProgressTracker) UpdateProgress(token string, progress float64, total int64) error {
	pt.mu.RLock()
	progressToken, exists := pt.activeTokens[token]
	pt.mu.RUnlock()

	if !exists {
		return NewRPCError(InvalidRequest, "progress token not found")
	}

	progressToken.Progress = progress
	progressToken.Total = total

	if pt.notifier != nil {
		return pt.notifier.SendProgress(progressToken.ctx, token, progress, total)
	}

	return nil
}

// CompleteProgress marks a progress token as complete and removes it
func (pt *ProgressTracker) CompleteProgress(token string) error {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	progressToken, exists := pt.activeTokens[token]
	if !exists {
		return NewRPCError(InvalidRequest, "progress token not found")
	}

	// Send final progress update
	if pt.notifier != nil {
		pt.notifier.SendProgress(progressToken.ctx, token, 1.0, progressToken.Total)
	}

	progressToken.cancel()
	delete(pt.activeTokens, token)
	return nil
}

// CancelProgress cancels a progress token
func (pt *ProgressTracker) CancelProgress(token string) error {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	progressToken, exists := pt.activeTokens[token]
	if !exists {
		return NewRPCError(InvalidRequest, "progress token not found")
	}

	progressToken.cancel()
	delete(pt.activeTokens, token)
	return nil
}

// GetProgress returns the current progress for a token
func (pt *ProgressTracker) GetProgress(token string) (*ProgressToken, error) {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	progressToken, exists := pt.activeTokens[token]
	if !exists {
		return nil, NewRPCError(InvalidRequest, "progress token not found")
	}

	// Return a copy to avoid race conditions
	return &ProgressToken{
		Token:     progressToken.Token,
		Progress:  progressToken.Progress,
		Total:     progressToken.Total,
		StartTime: progressToken.StartTime,
	}, nil
}

// ListActiveTokens returns all active progress tokens
func (pt *ProgressTracker) ListActiveTokens() []string {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	tokens := make([]string, 0, len(pt.activeTokens))
	for token := range pt.activeTokens {
		tokens = append(tokens, token)
	}
	return tokens
}

// GenerateProgressToken generates a unique progress token
func GenerateProgressToken() string {
	return fmt.Sprintf("progress_%d_%d", time.Now().UnixNano(), rand.Int63())
}

// WithProgress adds progress tracking to a Meta object
func WithProgress(meta *Meta, token string) *Meta {
	if meta == nil {
		meta = &Meta{}
	}
	meta.ProgressToken = token
	return meta
}
