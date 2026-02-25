// Package illygen is a lightweight intelligence engine for building AI-like systems in Go.
//
// Illygen lets you build systems that reason, make decisions, and learn —
// without needing expensive AI models, GPUs, or cloud services.
//
// # Core Concepts
//
// A Node is a single unit of reasoning. You define its logic as a plain Go function.
// A Flow is a net of connected nodes — the reasoning pipeline.
// An Engine runs flows and returns results.
// A KnowledgeStore holds structured facts that nodes can query.
//
// # Minimal Example
//
//	input := illygen.NewNode("input", func(ctx illygen.Context) illygen.Result {
//	    return illygen.Result{Next: "output", Confidence: 1.0}
//	})
//
//	output := illygen.NewNode("output", func(ctx illygen.Context) illygen.Result {
//	    return illygen.Result{Value: "Hi! I'm Illygen.", Confidence: 1.0}
//	})
//
//	flow := illygen.NewFlow().
//	    Add(input).
//	    Add(output).
//	    Link("input", "output", 1.0)
//
//	engine := illygen.NewEngine()
//
//	result, err := engine.Run(flow, illygen.Context{"input": "hello"})
//	fmt.Println(result.Value) // Hi! I'm Illygen.
//
// # Public API (v0.1)
//
//	illygen.NewNode(id, fn)    → *Node
//	illygen.NewFlow()          → *Flow
//	illygen.NewEngine()        → *Engine
//	illygen.NewKnowledgeStore() → *KnowledgeStore
//	illygen.Knowledge(ctx)     → *KnowledgeStore  (inside a NodeFunc)
//
//	flow.Add(node)             → *Flow
//	flow.Link(from, to, w)    → *Flow
//	flow.Entry(nodeID)         → *Flow
//
//	engine.Run(flow, ctx)      → (Result, error)
//
//	ctx.Get(key)               → any
//	ctx.Set(key, value)
//	ctx.String(key)            → string
//	ctx.Has(key)               → bool
//
//	result.Value               → any
//	result.Confidence          → float64
//
// Everything else is internal.
package illygen
