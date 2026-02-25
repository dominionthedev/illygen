package illygen

// Result is the output of a node's execution and ultimately the output of a flow.
// Every node returns a Result. The engine reads it, follows the next route,
// and the final Result is what the caller receives.
type Result struct {
	// Value is the main output of this node — can be any type.
	Value any

	// Confidence is how certain this node is about its Result (0.0 to 1.0).
	// The engine and future learning logic use this to evaluate and improve routes.
	Confidence float64

	// Next is the ID of the node to consult after this one.
	// If empty, the flow ends here and this Result is returned to the caller.
	Next string
}
