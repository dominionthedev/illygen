package illygen_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	illygen "github.com/leraniode/illygen"
)

// ─────────────────────────────────────────────
//  Node
// ─────────────────────────────────────────────

func TestNewNode_Basic(t *testing.T) {
	node := illygen.NewNode("greeter", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "hello", Confidence: 1.0}
	})

	if node.ID() != "greeter" {
		t.Errorf("expected ID %q, got %q", "greeter", node.ID())
	}
}

func TestNewNode_PanicsOnEmptyID(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty node ID, got none")
		}
	}()
	illygen.NewNode("", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{}
	})
}

func TestNewNode_PanicsOnNilFunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil NodeFunc, got none")
		}
	}()
	illygen.NewNode("node", nil)
}

// ─────────────────────────────────────────────
//  Context
// ─────────────────────────────────────────────

func TestContext_SetAndGet(t *testing.T) {
	ctx := illygen.Context{}
	ctx.Set("key", "value")

	if got := ctx.Get("key"); got != "value" {
		t.Errorf("expected %q, got %v", "value", got)
	}
}

func TestContext_GetMissing(t *testing.T) {
	ctx := illygen.Context{}
	if got := ctx.Get("missing"); got != nil {
		t.Errorf("expected nil for missing key, got %v", got)
	}
}

func TestContext_Has(t *testing.T) {
	ctx := illygen.Context{"x": 1}
	if !ctx.Has("x") {
		t.Error("expected Has to return true for existing key")
	}
	if ctx.Has("y") {
		t.Error("expected Has to return false for missing key")
	}
}

func TestContext_String(t *testing.T) {
	ctx := illygen.Context{"name": "ada"}
	if got := ctx.String("name"); got != "ada" {
		t.Errorf("expected %q, got %q", "ada", got)
	}
	if got := ctx.String("missing"); got != "" {
		t.Errorf("expected empty string for missing key, got %q", got)
	}
}

func TestContext_Bool(t *testing.T) {
	ctx := illygen.Context{"flag": true}
	if !ctx.Bool("flag") {
		t.Error("expected true")
	}
	if ctx.Bool("missing") {
		t.Error("expected false for missing key")
	}
}

func TestContext_Int(t *testing.T) {
	ctx := illygen.Context{"count": 42}
	if got := ctx.Int("count"); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
	if got := ctx.Int("missing"); got != 0 {
		t.Errorf("expected 0 for missing key, got %d", got)
	}
}

func TestContext_Float(t *testing.T) {
	ctx := illygen.Context{"score": 0.95}
	if got := ctx.Float("score"); got != 0.95 {
		t.Errorf("expected 0.95, got %f", got)
	}
}

// ─────────────────────────────────────────────
//  Flow
// ─────────────────────────────────────────────

func TestFlow_FirstAddBecomesEntry(t *testing.T) {
	a := illygen.NewNode("a", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "from a"}
	})
	b := illygen.NewNode("b", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "from b"}
	})

	flow := illygen.NewFlow().Add(a).Add(b)
	engine := illygen.NewEngine()

	result, err := engine.Run(flow, illygen.Context{})
	if err != nil {
		t.Fatal(err)
	}
	// No links — should stop at entry node "a"
	if result.Value != "from a" {
		t.Errorf("expected entry node result, got %v", result.Value)
	}
}

func TestFlow_ExplicitEntry(t *testing.T) {
	a := illygen.NewNode("a", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "from a"}
	})
	b := illygen.NewNode("b", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "from b"}
	})

	flow := illygen.NewFlow().Add(a).Add(b).Entry("b")
	engine := illygen.NewEngine()

	result, err := engine.Run(flow, illygen.Context{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Value != "from b" {
		t.Errorf("expected result from explicit entry node b, got %v", result.Value)
	}
}

func TestFlow_AddNilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when adding nil node, got none")
		}
	}()
	illygen.NewFlow().Add(nil)
}

func TestFlow_DuplicateLinkIgnored(t *testing.T) {
	a := illygen.NewNode("a", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "ok"}
	})
	// Should not panic or error on duplicate Link
	flow := illygen.NewFlow().
		Add(a).
		Link("a", "a", 1.0).
		Link("a", "a", 0.5)
	if flow == nil {
		t.Error("expected flow to be non-nil after duplicate Link")
	}
}

// ─────────────────────────────────────────────
//  Engine
// ─────────────────────────────────────────────

func TestEngine_Run_SingleNode(t *testing.T) {
	node := illygen.NewNode("only", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "done", Confidence: 0.9}
	})
	flow := illygen.NewFlow().Add(node)
	engine := illygen.NewEngine()

	result, err := engine.Run(flow, illygen.Context{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Value != "done" {
		t.Errorf("expected %q, got %v", "done", result.Value)
	}
	if result.Confidence != 0.9 {
		t.Errorf("expected confidence 0.9, got %f", result.Confidence)
	}
}

func TestEngine_Run_NextRouting(t *testing.T) {
	a := illygen.NewNode("a", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Next: "b", Confidence: 0.8}
	})
	b := illygen.NewNode("b", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "reached b", Confidence: 1.0}
	})

	flow := illygen.NewFlow().Add(a).Add(b).Link("a", "b", 1.0)
	engine := illygen.NewEngine()

	result, err := engine.Run(flow, illygen.Context{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Value != "reached b" {
		t.Errorf("expected %q, got %v", "reached b", result.Value)
	}
}

func TestEngine_Run_GraphLinkFallback(t *testing.T) {
	// Next is NOT set — engine should follow the graph Link automatically
	a := illygen.NewNode("a", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Confidence: 0.5} // no Next
	})
	b := illygen.NewNode("b", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "auto-routed", Confidence: 1.0}
	})

	flow := illygen.NewFlow().Add(a).Add(b).Link("a", "b", 1.0)
	engine := illygen.NewEngine()

	result, err := engine.Run(flow, illygen.Context{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Value != "auto-routed" {
		t.Errorf("expected auto-routed via graph link, got %v", result.Value)
	}
}

func TestEngine_Run_ContextPassedToNodes(t *testing.T) {
	node := illygen.NewNode("reader", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: ctx.String("greeting")}
	})
	flow := illygen.NewFlow().Add(node)
	engine := illygen.NewEngine()

	result, err := engine.Run(flow, illygen.Context{"greeting": "hi from context"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Value != "hi from context" {
		t.Errorf("expected context value, got %v", result.Value)
	}
}

func TestEngine_Run_NilContext(t *testing.T) {
	// Nil context must not panic — should be treated as empty
	node := illygen.NewNode("safe", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: "ok"}
	})
	flow := illygen.NewFlow().Add(node)
	engine := illygen.NewEngine()

	result, err := engine.Run(flow, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Value != "ok" {
		t.Errorf("expected ok, got %v", result.Value)
	}
}

func TestEngine_Run_EmptyFlow(t *testing.T) {
	flow := illygen.NewFlow()
	engine := illygen.NewEngine()

	_, err := engine.Run(flow, illygen.Context{})
	if err == nil {
		t.Error("expected error for empty flow, got nil")
	}
}

func TestEngine_Run_UnknownNextNode(t *testing.T) {
	node := illygen.NewNode("bad", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Next: "does-not-exist"}
	})
	flow := illygen.NewFlow().Add(node)
	engine := illygen.NewEngine()

	_, err := engine.Run(flow, illygen.Context{})
	if err == nil {
		t.Error("expected error for unknown Next node, got nil")
	}
	if !strings.Contains(err.Error(), "does-not-exist") {
		t.Errorf("expected error to mention missing node, got: %v", err)
	}
}

func TestEngine_Run_MultiStep(t *testing.T) {
	// a → b → c — each node appends to a slice in context
	a := illygen.NewNode("a", func(ctx illygen.Context) illygen.Result {
		ctx.Set("trace", "a")
		return illygen.Result{Next: "b"}
	})
	b := illygen.NewNode("b", func(ctx illygen.Context) illygen.Result {
		ctx.Set("trace", ctx.String("trace")+"-b")
		return illygen.Result{Next: "c"}
	})
	c := illygen.NewNode("c", func(ctx illygen.Context) illygen.Result {
		ctx.Set("trace", ctx.String("trace")+"-c")
		return illygen.Result{Value: ctx.String("trace"), Confidence: 1.0}
	})

	flow := illygen.NewFlow().
		Add(a).Add(b).Add(c).
		Link("a", "b", 1.0).
		Link("b", "c", 1.0)

	result, err := illygen.NewEngine().Run(flow, illygen.Context{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Value != "a-b-c" {
		t.Errorf("expected trace %q, got %v", "a-b-c", result.Value)
	}
}

func TestEngine_Run_Concurrent(t *testing.T) {
	// Same engine and flow, called from 50 goroutines simultaneously
	node := illygen.NewNode("concurrent", func(ctx illygen.Context) illygen.Result {
		return illygen.Result{Value: ctx.String("id"), Confidence: 1.0}
	})
	flow := illygen.NewFlow().Add(node)
	engine := illygen.NewEngine()

	const workers = 50
	results := make([]string, workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("worker-%d", n)
			res, err := engine.Run(flow, illygen.Context{"id": id})
			if err != nil {
				t.Errorf("worker %d error: %v", n, err)
				return
			}
			results[n] = res.Value.(string)
		}(i)
	}

	wg.Wait()

	for i, r := range results {
		expected := fmt.Sprintf("worker-%d", i)
		if r != expected {
			t.Errorf("worker %d: expected %q, got %q", i, expected, r)
		}
	}
}

// ─────────────────────────────────────────────
//  KnowledgeStore
// ─────────────────────────────────────────────

func TestKnowledgeStore_Add(t *testing.T) {
	store := illygen.NewKnowledgeStore()
	err := store.Add("k1", "test", map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.Size() != 1 {
		t.Errorf("expected size 1, got %d", store.Size())
	}
}

func TestKnowledgeStore_Add_EmptyID(t *testing.T) {
	store := illygen.NewKnowledgeStore()
	err := store.Add("", "domain", map[string]any{})
	if err == nil {
		t.Error("expected error for empty id, got nil")
	}
}

func TestKnowledgeStore_Add_EmptyDomain(t *testing.T) {
	store := illygen.NewKnowledgeStore()
	err := store.Add("k1", "", map[string]any{})
	if err == nil {
		t.Error("expected error for empty domain, got nil")
	}
}

func TestKnowledgeStore_Add_Duplicate(t *testing.T) {
	store := illygen.NewKnowledgeStore()
	_ = store.Add("k1", "test", map[string]any{})
	err := store.Add("k1", "test", map[string]any{})
	if err == nil {
		t.Error("expected error for duplicate id, got nil")
	}
}

func TestKnowledgeStore_Get(t *testing.T) {
	store := illygen.NewKnowledgeStore()
	_ = store.Add("k1", "facts", map[string]any{"answer": "42"})

	unit, ok := store.Get("k1")
	if !ok {
		t.Fatal("expected unit to exist")
	}
	if unit.Fact("answer") != "42" {
		t.Errorf("expected fact %q, got %v", "42", unit.Fact("answer"))
	}
}

func TestKnowledgeStore_Get_Missing(t *testing.T) {
	store := illygen.NewKnowledgeStore()
	_, ok := store.Get("nope")
	if ok {
		t.Error("expected ok=false for missing unit")
	}
}

func TestKnowledgeStore_Domain(t *testing.T) {
	store := illygen.NewKnowledgeStore()
	_ = store.Add("a", "greetings", map[string]any{"r": "hi"})
	_ = store.Add("b", "greetings", map[string]any{"r": "hello"})
	_ = store.Add("c", "farewells", map[string]any{"r": "bye"})

	units := store.Domain("greetings")
	if len(units) != 2 {
		t.Errorf("expected 2 units in greetings domain, got %d", len(units))
	}
	for _, u := range units {
		if u.Domain != "greetings" {
			t.Errorf("expected domain %q, got %q", "greetings", u.Domain)
		}
	}
}

func TestKnowledgeStore_Domain_Empty(t *testing.T) {
	store := illygen.NewKnowledgeStore()
	units := store.Domain("nobody")
	if len(units) != 0 {
		t.Errorf("expected empty slice for unknown domain, got %d units", len(units))
	}
}

func TestKnowledgeStore_Domain_SortedByWeight(t *testing.T) {
	store := illygen.NewKnowledgeStore()
	_ = store.Add("low", "test", map[string]any{})
	_ = store.Add("high", "test", map[string]any{})

	// Manually adjust weights via Get
	low, _ := store.Get("low")
	low.Weight = 0.3
	high, _ := store.Get("high")
	high.Weight = 0.9

	units := store.Domain("test")
	if len(units) < 2 {
		t.Fatal("expected 2 units")
	}
	if units[0].ID != "high" {
		t.Errorf("expected highest weight unit first, got %q", units[0].ID)
	}
}

// ─────────────────────────────────────────────
//  Knowledge injection
// ─────────────────────────────────────────────

func TestEngine_Knowledge_InjectedIntoContext(t *testing.T) {
	store := illygen.NewKnowledgeStore()
	_ = store.Add("fact1", "science", map[string]any{"answer": "42"})

	node := illygen.NewNode("lookup", func(ctx illygen.Context) illygen.Result {
		ks := illygen.Knowledge(ctx)
		if ks == nil {
			return illygen.Result{Value: "no store"}
		}
		units := ks.Domain("science")
		if len(units) == 0 {
			return illygen.Result{Value: "no units"}
		}
		return illygen.Result{Value: units[0].Fact("answer"), Confidence: 1.0}
	})

	flow := illygen.NewFlow().Add(node)
	engine := illygen.NewEngine(store)

	result, err := engine.Run(flow, illygen.Context{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Value != "42" {
		t.Errorf("expected knowledge fact %q, got %v", "42", result.Value)
	}
}

func TestKnowledge_NilWhenNoStore(t *testing.T) {
	node := illygen.NewNode("n", func(ctx illygen.Context) illygen.Result {
		ks := illygen.Knowledge(ctx)
		if ks != nil {
			return illygen.Result{Value: "unexpected store"}
		}
		return illygen.Result{Value: "no store as expected"}
	})

	flow := illygen.NewFlow().Add(node)
	engine := illygen.NewEngine() // no store

	result, err := engine.Run(flow, illygen.Context{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Value != "no store as expected" {
		t.Errorf("unexpected: %v", result.Value)
	}
}
