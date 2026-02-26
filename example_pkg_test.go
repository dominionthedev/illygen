package illygen_test

import (
	"fmt"

	illygen "github.com/leraniode/illygen"
)

// Basic example showing how to build a flow and run it with an Engine.
func ExampleNewFlow_basic() {
	node := illygen.NewNode("hello", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "hello world", Confidence: 1.0}
	})

	flow := illygen.NewFlow().Add(node)
	engine := illygen.NewEngine()

	res, _ := engine.Run(flow, illygen.Context{})
	fmt.Println(res.Value)
	// Output: hello world
}

// Example demonstrating the KnowledgeStore API and Domain lookup.
func ExampleKnowledgeStore_basic() {
	s := illygen.NewKnowledgeStore()
	_ = s.Add("g1", "greetings", map[string]any{"response": "hi"})

	units := s.Domain("greetings")
	fmt.Println(len(units))
	// Output: 1
}

// Example showing how to attach a KnowledgeStore to an Engine and query it from a Node.
func ExampleEngine() {
	s := illygen.NewKnowledgeStore()
	_ = s.Add("g1", "greetings", map[string]any{"keywords": []string{"hello"}, "response": "hi"})

	node := illygen.NewNode("n", func(ctx illygen.Context) illygen.Result {
		ks := illygen.Knowledge(ctx)
		if ks == nil {
			return illygen.Result{Value: "no"}
		}
		units := ks.Domain("greetings")
		if len(units) > 0 {
			return illygen.Result{Value: units[0].Fact("response")}
		}
		return illygen.Result{Value: "none"}
	})

	flow := illygen.NewFlow().Add(node)
	e := illygen.NewEngine(s)
	res, _ := e.Run(flow, illygen.Context{})
	fmt.Println(res.Value)
	// Output: hi
}
