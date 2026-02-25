// Package graph provides the internal weighted graph that powers Illygen flows.
// This package is private — users interact with Flow, not Graph directly.
package graph

import (
	"fmt"
	"sync"
)

// Edge is a directed weighted connection between two nodes.
type Edge struct {
	From   string
	To     string
	Weight float64
}

// Graph is a directed weighted graph of node connections.
type Graph struct {
	mu    sync.RWMutex
	edges map[string][]*Edge // keyed by From node ID
}

// New creates an empty Graph.
func New() *Graph {
	return &Graph{edges: make(map[string][]*Edge)}
}

// Add creates a directed edge from → to with the given weight.
// Returns an error if the edge already exists.
func (g *Graph) Add(from, to string, weight float64) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, e := range g.edges[from] {
		if e.To == to {
			return fmt.Errorf("graph: edge %q → %q already exists", from, to)
		}
	}
	g.edges[from] = append(g.edges[from], &Edge{From: from, To: to, Weight: weight})
	return nil
}

// From returns all edges outgoing from a node, sorted by weight descending.
func (g *Graph) From(id string) []*Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()

	edges := make([]*Edge, len(g.edges[id]))
	copy(edges, g.edges[id])
	sortEdges(edges)
	return edges
}

// Has reports whether a node ID has any outgoing edges.
func (g *Graph) Has(id string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.edges[id]
	return ok
}

func sortEdges(edges []*Edge) {
	for i := 1; i < len(edges); i++ {
		for j := i; j > 0 && edges[j].Weight > edges[j-1].Weight; j-- {
			edges[j], edges[j-1] = edges[j-1], edges[j]
		}
	}
}
