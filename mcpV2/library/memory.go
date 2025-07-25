package library

import (
	"errors"
	"sync"
	"time"
)

// Common memory errors
var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExpired  = errors.New("key expired")
)

// Memory provides key-value storage with multiple backend options
type Memory interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	List() ([]string, error)
	Clear() error
	Stats() (*MemoryStats, error)
	Close() error
}

// MemoryStats provides memory usage statistics
type MemoryStats struct {
	TotalKeys   int64                  `json:"total_keys"`
	ActiveKeys  int64                  `json:"active_keys"`
	MemoryUsage int64                  `json:"memory_usage_bytes"`
	Backend     string                 `json:"backend"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// MemoryOptions configures memory behavior
type MemoryOptions struct {
	Backend      string                 `json:"backend"`    // "inmemory", "redis", "badger", "bbolt", "sqlite"
	Persistent   bool                   `json:"persistent"` // Enable persistence for supported backends
	TTL          time.Duration          `json:"ttl"`        // Time-to-live for entries
	MaxSize      int64                  `json:"max_size"`   // Maximum memory usage
	CompressData bool                   `json:"compress"`   // Compress stored data
	Config       map[string]interface{} `json:"config"`     // Backend-specific configuration
}

// MemoryEntry represents a stored entry
type MemoryEntry struct {
	Value     interface{}            `json:"value"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// InMemoryBackend provides in-memory storage
type InMemoryBackend struct {
	data map[string]*MemoryEntry
	mu   sync.RWMutex
	opts *MemoryOptions
}

// NewInMemoryBackend creates in-memory storage
func NewInMemoryBackend(opts *MemoryOptions) *InMemoryBackend {
	if opts == nil {
		opts = &MemoryOptions{
			Backend: "inmemory",
		}
	}

	return &InMemoryBackend{
		data: make(map[string]*MemoryEntry),
		opts: opts,
	}
}

// Set stores a value in memory
func (m *InMemoryBackend) Set(key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry := &MemoryEntry{
		Value:     value,
		CreatedAt: time.Now(),
	}

	if m.opts.TTL > 0 {
		expiry := time.Now().Add(m.opts.TTL)
		entry.ExpiresAt = &expiry
	}

	m.data[key] = entry
	return nil
}

// Get retrieves a value from memory
func (m *InMemoryBackend) Get(key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.data[key]
	if !exists {
		return nil, ErrKeyNotFound
	}

	// Check expiration
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		delete(m.data, key)
		return nil, ErrKeyExpired
	}

	return entry.Value, nil
}

// Delete removes a key from memory
func (m *InMemoryBackend) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[key]; !exists {
		return ErrKeyNotFound
	}

	delete(m.data, key)
	return nil
}

// List returns all keys in memory
func (m *InMemoryBackend) List() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0, len(m.data))
	now := time.Now()

	for key, entry := range m.data {
		// Skip expired entries
		if entry.ExpiresAt != nil && now.After(*entry.ExpiresAt) {
			continue
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// Clear removes all keys from memory
func (m *InMemoryBackend) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]*MemoryEntry)
	return nil
}

// Stats returns memory statistics
func (m *InMemoryBackend) Stats() (*MemoryStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	activeKeys := int64(0)
	now := time.Now()

	for _, entry := range m.data {
		// Count only non-expired entries
		if entry.ExpiresAt == nil || now.Before(*entry.ExpiresAt) {
			activeKeys++
		}
	}

	return &MemoryStats{
		TotalKeys:   int64(len(m.data)),
		ActiveKeys:  activeKeys,
		MemoryUsage: 0, // TODO: Calculate actual memory usage
		Backend:     "inmemory",
		Metadata: map[string]interface{}{
			"ttl": m.opts.TTL.String(),
		},
	}, nil
}

// Close closes the memory backend
func (m *InMemoryBackend) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = nil
	return nil
}

// NewMemory creates a memory instance with specified backend
func NewMemory(opts *MemoryOptions) (Memory, error) {
	if opts == nil {
		opts = &MemoryOptions{Backend: "inmemory"}
	}

	switch opts.Backend {
	case "inmemory":
		return NewInMemoryBackend(opts), nil
	default:
		return NewInMemoryBackend(opts), nil
	}
}
