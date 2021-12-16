package astar

import "fmt"

// Heuristic is a function that estimates the remaining cost to reach the end for a node. It must
// always return an integer cost value, even for nodes it does not know.
type Heuristic = func(*Node) int

// ConstantHeuristic can be used to construct a simple heuristic function with constant costs for
// reaching the end node. Use AddNode to add a node with estimated cost and use Heuristic to
// retrieve the heuristic function.
type ConstantHeuristic struct {
	data *map[*Node]int
}

// AddNode adds a new node with estimated constant cost data. If the node is already there with a
// different estimate, this will error out.
func (s *ConstantHeuristic) AddNode(node *Node, estimate int) error {
	if s.data == nil {
		// Initialise the data if needed.
		s.data = &map[*Node]int{}
	}
	if val, found := (*s.data)[node]; found && val != estimate {
		return fmt.Errorf(
			"estimate for node %s deviates: old %d, new %d",
			node.ToString(), val, estimate,
		)
	}
	(*s.data)[node] = estimate
	return nil
}

// Heuristic obtains a heuristic function for the provided data. If a node cannot be found in the
// available data, return the default value. This is suitable for use with FindPath. The heuristic
// is cleared so that adding new nodes won't influence the function's data.
func (s *ConstantHeuristic) Heuristic(defaultValue int) Heuristic {
	data := *s.data
	heuristic := func(node *Node) int {
		if datum, found := data[node]; found {
			return datum
		}
		return defaultValue
	}
	// Clear heuristic by pointing to new, empty data.
	s.data = &map[*Node]int{}
	return heuristic
}
