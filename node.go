package illygen

import "fmt"

// NodeFunc is the function signature every node must implement.
// It receives the current Context and returns a Result.
// This is the only thing a user needs to write to create a node.
type NodeFunc func(ctx Context) Result

// Node is a single unit of reasoning in Illygen.
// Like a neuron in a neural network, a node gets consulted,
// produces a signal (Result), and optionally routes to the next node.
//
// Nodes are created with NewNode and registered into a Flow via Flow.Add.
// The engine calls each node in sequence during flow execution.
type Node struct {
	id string
	fn NodeFunc
}

// NewNode creates a new Node with the given ID and logic function.
// The ID must be non-empty and unique within a flow.
// The fn must not be nil.
//
// Example:
//
//	greeter := illygen.NewNode("greeter", func(ctx illygen.Context) illygen.Result {
//	    return illygen.Result{Value: "Hi! I'm Illygen.", Confidence: 1.0}
//	})
func NewNode(id string, fn NodeFunc) *Node {
	if id == "" {
		panic("illygen: NewNode called with empty id")
	}
	if fn == nil {
		panic(fmt.Sprintf("illygen: NewNode %q called with nil NodeFunc", id))
	}
	return &Node{id: id, fn: fn}
}

// ID returns the node's unique identifier.
func (n *Node) ID() string {
	return n.id
}

// execute runs the node's logic against the given context.
// Called internally by the engine — not by users.
func (n *Node) execute(ctx Context) Result {
	return n.fn(ctx)
}
