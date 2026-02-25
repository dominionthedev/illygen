package illygen

import (
	"fmt"
	"sync"
	"time"
)

// KnowledgeUnit is the atomic piece of knowledge in Illygen.
// It is lightweight — similar in concept to a tensor in AI, but built around
// structured facts rather than numerical matrices.
//
// Training produces units with high initial Weight.
// Exploring refines them over time.
type KnowledgeUnit struct {
	ID      string
	Domain  string
	Facts   map[string]any
	Weight  float64
	Updated time.Time
}

// Fact returns a single fact value by key. Returns nil if not found.
func (u *KnowledgeUnit) Fact(key string) any {
	return u.Facts[key]
}

// KnowledgeStore holds all KnowledgeUnits for an Illygen engine.
// Nodes query it by domain to retrieve relevant knowledge during execution.
type KnowledgeStore struct {
	mu    sync.RWMutex
	units map[string]*KnowledgeUnit
}

// NewKnowledgeStore creates an empty KnowledgeStore.
//
// Example:
//
//	store := illygen.NewKnowledgeStore()
//	store.Add("k1", "greetings", map[string]any{"response": "Hi! I'm Illygen."})
func NewKnowledgeStore() *KnowledgeStore {
	return &KnowledgeStore{
		units: make(map[string]*KnowledgeUnit),
	}
}

// Add inserts a new KnowledgeUnit into the store.
// Returns an error if a unit with the same ID already exists.
func (s *KnowledgeStore) Add(id, domain string, facts map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.units[id]; exists {
		return fmt.Errorf("illygen: knowledge unit %q already exists", id)
	}
	s.units[id] = &KnowledgeUnit{
		ID:      id,
		Domain:  domain,
		Facts:   facts,
		Weight:  1.0,
		Updated: time.Now(),
	}
	return nil
}

// Get retrieves a KnowledgeUnit by ID.
func (s *KnowledgeStore) Get(id string) (*KnowledgeUnit, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.units[id]
	return u, ok
}

// Domain returns all KnowledgeUnits in a given domain, sorted by weight descending.
// This is how nodes query knowledge — by domain, not by ID.
func (s *KnowledgeStore) Domain(domain string) []*KnowledgeUnit {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*KnowledgeUnit
	for _, u := range s.units {
		if u.Domain == domain {
			result = append(result, u)
		}
	}
	sortUnitsByWeight(result)
	return result
}

// Size returns the total number of units in the store.
func (s *KnowledgeStore) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.units)
}

func sortUnitsByWeight(units []*KnowledgeUnit) {
	for i := 1; i < len(units); i++ {
		for j := i; j > 0 && units[j].Weight > units[j-1].Weight; j-- {
			units[j], units[j-1] = units[j-1], units[j]
		}
	}
}
