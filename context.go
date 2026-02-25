package illygen

// Context is a simple map that carries data through a flow execution.
// It is the single source of truth passed to every node during a run.
//
// Example:
//
//	ctx := illygen.Context{
//	    "input": "hello",
//	    "user":  "ada",
//	}
type Context map[string]any

// Get retrieves a value by key. Returns nil if the key doesn't exist.
func (c Context) Get(key string) any {
	return c[key]
}

// Set stores a value under the given key.
func (c Context) Set(key string, value any) {
	c[key] = value
}

// Has reports whether a key exists in the context.
func (c Context) Has(key string) bool {
	_, ok := c[key]
	return ok
}

// String is a convenience method that returns a context value as a string.
// Returns an empty string if the key doesn't exist or is not a string.
func (c Context) String(key string) string {
	v, _ := c[key].(string)
	return v
}
