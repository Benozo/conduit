package memory

import (
	"sync"
)

// SharedMemory provides a structure for agents to share data.
type SharedMemory struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewSharedMemory initializes a new SharedMemory instance.
func NewSharedMemory() *SharedMemory {
	return &SharedMemory{
		data: make(map[string]interface{}),
	}
}

// Set stores a value in shared memory.
func (sm *SharedMemory) Set(key string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

// Get retrieves a value from shared memory.
func (sm *SharedMemory) Get(key string) (interface{}, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, exists := sm.data[key]
	return value, exists
}

// Delete removes a value from shared memory.
func (sm *SharedMemory) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

// Clear removes all values from shared memory.
func (sm *SharedMemory) Clear() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data = make(map[string]interface{})
}