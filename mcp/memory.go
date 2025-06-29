package mcp

import "sync"

type Memory struct {
	mu    sync.RWMutex
	store map[string]interface{}
}

func NewMemory() *Memory {
	return &Memory{store: make(map[string]interface{})}
}

func (m *Memory) Set(key string, val interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = val
}

func (m *Memory) Get(key string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.store[key]
}

func (m *Memory) All() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.store
}
