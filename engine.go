package illygen

import (
	"fmt"

	"github.com/leraniode/illygen/internal/runtime"
)

// Engine is the execution core of Illygen.
// It runs flows, walks nodes, and returns results.
// Create one engine and reuse it — it is stateless.
//
// Example:
//
//	engine := illygen.NewEngine()
//	result := engine.Run(flow, illygen.Context{"input": "hello"})
//	fmt.Println(result.Value)
type Engine struct {
	knowledge *KnowledgeStore
}

// NewEngine creates a new Engine.
// Optionally attach a KnowledgeStore so nodes can query knowledge during execution.
//
//	engine := illygen.NewEngine()
//	engine := illygen.NewEngine(store) // with knowledge
func NewEngine(store ...*KnowledgeStore) *Engine {
	e := &Engine{}
	if len(store) > 0 {
		e.knowledge = store[0]
	}
	return e
}

// Run executes a flow with the given context and returns the final Result.
// Execution starts at the flow's entry node and walks the graph until
// a node returns an empty Next or no linked nodes remain.
//
// Run is synchronous and safe to call concurrently from multiple goroutines.
func (e *Engine) Run(flow *Flow, ctx Context) (Result, error) {
	// Inject knowledge store into context so nodes can access it
	if e.knowledge != nil {
		ctx.Set("__knowledge__", e.knowledge)
	}

	// Get entry node
	entry, err := flow.entryNode()
	if err != nil {
		return Result{}, err
	}

	// Build the executor function — bridges internal runtime with public types
	executor := func(nodeID string) (any, float64, string, error) {
		node, err := flow.node(nodeID)
		if err != nil {
			return nil, 0, "", err
		}

		result := node.execute(ctx)

		// If result specifies Next, use it.
		// Otherwise, follow the highest-weight edge from the graph.
		next := result.Next
		if next == "" {
			edges := flow.graph.From(nodeID)
			if len(edges) > 0 {
				// Highest-weight edge is first (graph returns sorted)
				next = edges[0].To
			}
		}

		// Validate that next node actually exists
		if next != "" {
			if _, err := flow.node(next); err != nil {
				return nil, 0, "", fmt.Errorf("illygen: node %q routed to unknown node %q", nodeID, next)
			}
		}

		return result.Value, result.Confidence, next, nil
	}

	trace, err := runtime.Execute(entry.ID(), executor)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Value:      trace.Final.Value,
		Confidence: trace.Final.Confidence,
	}, nil
}

// Knowledge returns a helper for querying the engine's KnowledgeStore from within a NodeFunc.
// Usage inside a node:
//
//	units := illygen.Knowledge(ctx).Domain("greetings")
func Knowledge(ctx Context) *KnowledgeStore {
	store, _ := ctx.Get("__knowledge__").(*KnowledgeStore)
	return store
}
