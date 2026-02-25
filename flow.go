package illygen

import (
	"fmt"

	"github.com/leraniode/illygen/internal/graph"
)

// Flow is a net of connected nodes — the reasoning pipeline.
// Like a neural network, it is a directed weighted graph.
// Nodes are added with Add(), connected with Link().
//
// Example:
//
//	flow := illygen.NewFlow().
//	    Add(inputNode).
//	    Add(outputNode).
//	    Link("input", "output", 1.0)
type Flow struct {
	nodes map[string]*Node
	graph *graph.Graph
	entry string
}

// NewFlow creates a new empty Flow.
func NewFlow() *Flow {
	return &Flow{
		nodes: make(map[string]*Node),
		graph: graph.New(),
	}
}

// Add registers a node into the flow.
// The first node added becomes the entry point automatically.
// Returns the Flow for chaining.
func (f *Flow) Add(node *Node) *Flow {
	f.nodes[node.ID()] = node
	if f.entry == "" {
		f.entry = node.ID()
	}
	return f
}

// Link connects two nodes with a weight.
// Weight represents the strength of this connection (0.0 to 1.0).
// Higher weight connections are preferred by the engine.
// Returns the Flow for chaining.
func (f *Flow) Link(from, to string, weight float64) *Flow {
	if err := f.graph.Add(from, to, weight); err != nil {
		// edge already exists — silently skip for fluent API usability
		_ = err
	}
	return f
}

// Entry explicitly sets which node the flow starts from.
// Useful when the first node added is not the intended entry point.
func (f *Flow) Entry(nodeID string) *Flow {
	f.entry = nodeID
	return f
}

// node retrieves a node by ID. Returns an error if not found.
func (f *Flow) node(id string) (*Node, error) {
	n, ok := f.nodes[id]
	if !ok {
		return nil, fmt.Errorf("illygen: node %q not found in flow", id)
	}
	return n, nil
}

// entryNode returns the entry node. Returns an error if no entry is defined.
func (f *Flow) entryNode() (*Node, error) {
	if f.entry == "" {
		return nil, fmt.Errorf("illygen: flow has no entry node — call Add() first")
	}
	return f.node(f.entry)
}
