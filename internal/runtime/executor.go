// Package runtime contains the internal execution engine for Illygen.
// Users interact with Engine in the illygen package — not this directly.
package runtime

import (
	"fmt"
)

// maxVisits is the maximum number of times a single node can be visited
// in one flow execution before the engine declares a cycle and returns an error.
// This guards against infinite loops caused by circular Next routing.
const maxVisits = 50

// Step records what happened at a single node during execution.
type Step struct {
	NodeID     string
	Value      any
	Confidence float64
	Next       string
}

// ExecutionTrace is the complete record of a flow execution.
// It is used by the engine to return results and, in future versions,
// by the learning logic to adjust weights.
// ExecutionTrace is the complete record of a flow execution.
// Steps holds every node visited in order.
// Final holds the last step — its Value and Confidence are returned to the caller.
// In future versions, the learning logic will read the trace to adjust weights.
type ExecutionTrace struct {
	Steps []Step
	Final Step
	Done  bool
}

// NodeExecutor is a function that runs a node — the engine calls this
// to decouple the executor from the illygen package types.
type NodeExecutor func(nodeID string) (value any, confidence float64, next string, err error)

// Execute runs the flow from the entry node, walking the graph
// until a node returns an empty Next or no outgoing edges exist.
//
// This is the core algorithm:
//
//	Start at entry node
//	Loop:
//	  execute node → get result
//	  choose next node (from result.Next or highest-weight edge)
//	  move to next node
//	Stop when no next node
func Execute(entry string, executor NodeExecutor) (*ExecutionTrace, error) {
	trace := &ExecutionTrace{}
	current := entry

	visited := make(map[string]int)

	for current != "" {
		visited[current]++
		if visited[current] > maxVisits {
			return nil, fmt.Errorf(
				"illygen/runtime: node %q visited %d times — possible cycle detected",
				current, visited[current],
			)
		}

		value, confidence, next, err := executor(current)
		if err != nil {
			return nil, fmt.Errorf("illygen/runtime: node %q failed: %w", current, err)
		}

		step := Step{
			NodeID:     current,
			Value:      value,
			Confidence: confidence,
			Next:       next,
		}
		trace.Steps = append(trace.Steps, step)
		trace.Final = step

		current = next
	}

	trace.Done = true
	return trace, nil
}
