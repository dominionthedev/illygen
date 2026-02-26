package illygen

import (
	"fmt"

	"github.com/leraniode/illygen/internal/runtime"
)

// Engine is the execution core of Illygen.
// It runs flows, walks the node graph, and returns the final Result.
//
// An Engine is stateless and safe to reuse across multiple flows
// and concurrent goroutines. Create one and share it freely.
//
// Example:
//
//	engine := illygen.NewEngine()
//	result, err := engine.Run(flow, illygen.Context{"input": "hello"})
//	fmt.Println(result.Value)
type Engine struct {
	knowledge *KnowledgeStore
}

// NewEngine creates a new Engine.
// Optionally attach a KnowledgeStore — nodes can then query it
// during execution via illygen.Knowledge(ctx).
//
//	engine := illygen.NewEngine()          // no knowledge
//	engine := illygen.NewEngine(store)     // with knowledge
func NewEngine(store ...*KnowledgeStore) *Engine {
	e := &Engine{}
	if len(store) > 0 {
		e.knowledge = store[0]
	}
	return e
}

// Run executes a flow with the given context and returns the final Result.
//
// Execution starts at the flow's entry node and walks the graph:
//   - If a node's Result specifies Next, that node is consulted next.
//   - If Next is empty, the engine follows the highest-weight Link.
//   - Execution stops when there is no next node.
//
// A nil Context is treated as an empty Context — no panic.
// Run is safe to call concurrently from multiple goroutines.
func (e *Engine) Run(flow *Flow, ctx Context) (Result, error) {
	// Guard against nil context — treat it as empty rather than panicking.
	if ctx == nil {
		ctx = Context{}
	}

	// Inject knowledge store into context so nodes can access it.
	if e.knowledge != nil {
		ctx.Set("__knowledge__", e.knowledge)
	}

	// Resolve entry node.
	entry, err := flow.entryNode()
	if err != nil {
		return Result{}, err
	}

	// executor bridges the internal runtime with the public illygen types.
	executor := func(nodeID string) (any, float64, string, error) {
		node, err := flow.node(nodeID)
		if err != nil {
			return nil, 0, "", err
		}

		result := node.execute(ctx)

		// Result.Next takes priority. If not set, follow the highest-weight link.
		next := result.Next
		if next == "" {
			if edges := flow.graph.From(nodeID); len(edges) > 0 {
				next = edges[0].To // graph returns edges sorted by weight desc
			}
		}

		// Validate that the next node was registered in this flow.
		if next != "" {
			if _, err := flow.node(next); err != nil {
				return nil, 0, "", fmt.Errorf(
					"illygen: node %q routed to %q which is not in the flow — did you call flow.Add()?",
					nodeID, next,
				)
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

// Knowledge returns the KnowledgeStore attached to this engine's context.
// Call this inside a NodeFunc to query knowledge by domain.
//
// Returns nil if no KnowledgeStore was attached to the engine.
//
// Example:
//
//	node := illygen.NewNode("lookup", func(ctx illygen.Context) illygen.Result {
//	    units := illygen.Knowledge(ctx).Domain("greetings")
//	    // ...
//	})
func Knowledge(ctx Context) *KnowledgeStore {
	store, _ := ctx.Get("__knowledge__").(*KnowledgeStore)
	return store
}
