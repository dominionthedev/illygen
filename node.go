package illygen

// NodeFunc is the function signature every node must implement.
// It receives the current Context and returns a Result.
// This is the only thing a user needs to implement to create a node.
type NodeFunc func(ctx Context) Result

// Node is a single unit of reasoning in Illygen.
// It holds knowledge relevant to its domain and uses it to produce a Result.
// Think of it as a neuron — it gets consulted, fires, and passes a signal forward.
type Node struct {
	id       string
	fn       NodeFunc
}

// NewNode creates a new Node with the given ID and logic function.
//
// Example:
//
//	greeter := illygen.NewNode("greeter", func(ctx illygen.Context) illygen.Result {
//	    return illygen.Result{Value: "Hi! I'm Illygen.", Confidence: 1.0}
//	})
func NewNode(id string, fn NodeFunc) *Node {
	return &Node{
		id: id,
		fn: fn,
	}
}

// ID returns the node's unique identifier.
func (n *Node) ID() string {
	return n.id
}

// execute runs the node's logic against the given context.
// This is called internally by the engine — not by users.
func (n *Node) execute(ctx Context) Result {
	return n.fn(ctx)
}
