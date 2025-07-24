package memory

// Context manages context-specific information that agents can use during their operations.
type Context struct {
    Variables map[string]interface{}
}

// NewContext creates a new context with initialized variables.
func NewContext() *Context {
    return &Context{
        Variables: make(map[string]interface{}),
    }
}

// SetVariable sets a variable in the context.
func (c *Context) SetVariable(key string, value interface{}) {
    c.Variables[key] = value
}

// GetVariable retrieves a variable from the context.
func (c *Context) GetVariable(key string) (interface{}, bool) {
    value, exists := c.Variables[key]
    return value, exists
}

// Clear clears all variables in the context.
func (c *Context) Clear() {
    c.Variables = make(map[string]interface{})
}