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

// Bool returns a context value as a bool.
// Returns false if the key doesn't exist or is not a bool.
func (c Context) Bool(key string) bool {
	v, _ := c[key].(bool)
	return v
}

// Int returns a context value as an int.
// Returns 0 if the key doesn't exist or is not an int.
func (c Context) Int(key string) int {
	v, _ := c[key].(int)
	return v
}

// Float returns a context value as a float64.
// Returns 0 if the key doesn't exist or is not a float64.
func (c Context) Float(key string) float64 {
	v, _ := c[key].(float64)
	return v
}
