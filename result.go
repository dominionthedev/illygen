package illygen

// Result is what every node returns and what engine.Run ultimately gives back.
//
// The engine uses Result to decide what happens next:
//   - If Next is set, the engine consults that node next.
//   - If Next is empty, the engine follows the highest-weight Link from the flow graph.
//   - If neither exists, execution ends and this Result is returned to the caller.
type Result struct {
	// Value is the main output produced by this node.
	// It can be any type — a string, a struct, a number.
	// The Value from the last node executed is what engine.Run returns.
	Value any

	// Confidence is how certain this node is about its Result (0.0 to 1.0).
	// Used by the engine to report certainty and by future learning logic
	// to reinforce or weaken routes over time.
	Confidence float64

	// Next is the ID of the node to consult after this one.
	// Setting Next overrides any Links defined in the flow graph.
	// Leave empty to let the engine follow the highest-weight Link automatically,
	// or to end the flow if no Links exist.
	Next string
}
