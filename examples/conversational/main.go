package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	illygen "github.com/leraniode/illygen"
)

// Simple conversational demo showing how to wire a KnowledgeStore and Flow.
// Run: `go run ./examples/conversational` and type messages.
func main() {
	// Create a knowledge store and add a couple of units.
	store := illygen.NewKnowledgeStore()
	_ = store.Add("greet-1", "greetings", map[string]any{
		"keywords": []string{"hello", "hi", "hey"},
		"response": "Hello! How can I help you today?",
	})
	_ = store.Add("greet-2", "greetings", map[string]any{
		"keywords": []string{"bye", "goodbye", "see ya"},
		"response": "Goodbye! Take care.",
	})
	_ = store.Add("smalltalk", "smalltalk", map[string]any{
		"keywords": []string{"how are you", "how're you", "how are things"},
		"response": "I'm a tiny reasoning engine — feeling deterministic today!",
	})

	// Define a single chat node that consults the KnowledgeStore.
	chat := illygen.NewNode("chat", func(ctx illygen.Context) illygen.Result {
		input := strings.ToLower(strings.TrimSpace(ctx.String("input")))
		ks := illygen.Knowledge(ctx)
		if ks != nil {
			// Domain returns units sorted by Weight (highest first)
			for _, u := range ks.Domain("greetings") {
				// facts stored as []string for keywords
				if kws, ok := u.Facts["keywords"].([]string); ok {
					for _, k := range kws {
						if strings.Contains(input, k) {
							return illygen.Result{Value: u.Facts["response"], Confidence: u.Weight}
						}
					}
				}
			}

			for _, u := range ks.Domain("smalltalk") {
				if kws, ok := u.Facts["keywords"].([]string); ok {
					for _, k := range kws {
						if strings.Contains(input, k) {
							return illygen.Result{Value: u.Facts["response"], Confidence: u.Weight}
						}
					}
				}
			}
		}

		// fallback
		return illygen.Result{Value: "Sorry, I don't understand — try saying 'hello' or 'bye'", Confidence: 0.5}
	})

	// Build flow and engine (attach knowledge store to engine so nodes can access it)
	flow := illygen.NewFlow().Add(chat)
	engine := illygen.NewEngine(store)

	// Run a small REPL
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Illygen conversational demo — type messages, 'quit' to exit")
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading input:", err)
			break
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		if strings.EqualFold(text, "quit") || strings.EqualFold(text, "exit") {
			fmt.Println("bye")
			break
		}

		res, err := engine.Run(flow, illygen.Context{"input": text})
		if err != nil {
			fmt.Println("engine error:", err)
			continue
		}
		fmt.Println(res.Value)
	}
}
