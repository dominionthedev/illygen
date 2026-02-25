// Intent Detection — Illygen v0.1 Official Example
//
// This demonstrates the complete Illygen API:
//
//	input → intent → action → output
//
// The flow detects what the user means, routes to the right action,
// and returns a response. Knowledge is stored in a KnowledgeStore
// and queried by nodes during execution.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	illygen "github.com/leraniode/illygen"
)

func main() {
	// ── Knowledge ──────────────────────────────────────────────────
	store := illygen.NewKnowledgeStore()

	_ = store.Add("greet-1", "greetings", map[string]any{
		"response": "Hi! I'm Illygen — a lightweight intelligence engine. How can I help?",
	})
	_ = store.Add("greet-2", "greetings", map[string]any{
		"response": "Hello! What can I do for you today?",
	})
	_ = store.Add("bye-1", "farewells", map[string]any{
		"response": "Goodbye! Come back anytime.",
	})
	_ = store.Add("illygen-1", "facts", map[string]any{
		"topic":    "illygen",
		"response": "Illygen is a lightweight intelligence engine built in Go. It uses flows, nodes, and knowledge to reason — no AI models needed.",
	})
	_ = store.Add("node-1", "facts", map[string]any{
		"topic":    "node",
		"response": "A node is the atomic unit of reasoning in Illygen. Define its logic as a Go function, connect it in a flow, and the engine does the rest.",
	})
	_ = store.Add("flow-1", "facts", map[string]any{
		"topic":    "flow",
		"response": "A flow is a connected net of nodes — Illygen's reasoning pipeline. It works like a neural network, reshaping as it learns.",
	})

	// ── Nodes ──────────────────────────────────────────────────────

	// InputNode: classifies the user's intent and routes accordingly
	inputNode := illygen.NewNode("input", func(ctx illygen.Context) illygen.Result {
		text := strings.ToLower(strings.TrimSpace(ctx.String("text")))

		switch {
		case isGreeting(text):
			ctx.Set("intent", "greeting")
			return illygen.Result{Next: "action", Confidence: 0.95}
		case isFarewell(text):
			ctx.Set("intent", "farewell")
			return illygen.Result{Next: "action", Confidence: 0.95}
		case isQuestion(text):
			ctx.Set("intent", "question")
			ctx.Set("query", text)
			return illygen.Result{Next: "action", Confidence: 0.80}
		default:
			ctx.Set("intent", "unknown")
			return illygen.Result{Next: "action", Confidence: 0.30}
		}
	})

	// ActionNode: uses intent + knowledge to build the response
	actionNode := illygen.NewNode("action", func(ctx illygen.Context) illygen.Result {
		intent := ctx.String("intent")
		store := illygen.Knowledge(ctx)

		switch intent {

		case "greeting":
			units := store.Domain("greetings")
			if len(units) > 0 {
				return illygen.Result{
					Value:      units[0].Fact("response"),
					Confidence: 0.95,
				}
			}

		case "farewell":
			units := store.Domain("farewells")
			if len(units) > 0 {
				return illygen.Result{
					Value:      units[0].Fact("response"),
					Confidence: 0.95,
				}
			}

		case "question":
			query := ctx.String("query")
			units := store.Domain("facts")
			for _, unit := range units {
				topic, _ := unit.Fact("topic").(string)
				if topic != "" && strings.Contains(query, topic) {
					return illygen.Result{
						Value:      unit.Fact("response"),
						Confidence: 0.85,
					}
				}
			}
			return illygen.Result{
				Value:      "I don't have knowledge about that yet — but I'm built to learn.",
				Confidence: 0.30,
			}
		}

		return illygen.Result{
			Value:      "I'm not sure how to respond to that. Try asking about Illygen, nodes, or flows.",
			Confidence: 0.20,
		}
	})

	// ── Flow ───────────────────────────────────────────────────────
	flow := illygen.NewFlow().
		Add(inputNode).
		Add(actionNode).
		Link("input", "action", 1.0)

	// ── Engine ─────────────────────────────────────────────────────
	engine := illygen.NewEngine(store)

	// ── REPL ───────────────────────────────────────────────────────
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║   Illygen v0.1 — Intent Detection Demo     ║")
	fmt.Println("║                                             ║")
	fmt.Println("║   Try: hi / what is illygen / bye          ║")
	fmt.Println("║   Type 'exit' to quit                       ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		if text == "exit" {
			fmt.Println("Illygen: Goodbye!")
			break
		}

		result, err := engine.Run(flow, illygen.Context{"text": text})
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Printf("Illygen [%.0f%%]: %v\n\n", result.Confidence*100, result.Value)
	}
}

// ── Helpers ────────────────────────────────────────────────────────

func isGreeting(s string) bool {
	for _, g := range []string{"hi", "hello", "hey", "yo", "howdy", "greetings"} {
		if s == g || strings.HasPrefix(s, g+" ") {
			return true
		}
	}
	return false
}

func isFarewell(s string) bool {
	for _, f := range []string{"bye", "goodbye", "see you", "later", "farewell", "ciao"} {
		if strings.Contains(s, f) {
			return true
		}
	}
	return false
}

func isQuestion(s string) bool {
	for _, q := range []string{"what", "who", "how", "why", "when", "where", "tell me", "explain"} {
		if strings.HasPrefix(s, q) {
			return true
		}
	}
	return false
}
