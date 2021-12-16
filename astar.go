// Package astar implements the A* path finding algorithm. The main function is FindPath. Please
// find a description here: https://en.wikipedia.org/wiki/A*_search_algorithm
package astar

import (
	"fmt"
)

// FindPath finds the path between the start and end node. It takes a graph in the form of a set of
// nodes, a start node, and an end node. It returns errors in case there are problems with the input
// or during execution. The path is returned in the correct order. This is achieved by using the
// normal algorithm and reversing the path at the end.
//
// This implementation modifies the original nodes!
//
// It also takes a heuristic that estimates the cost for moving from a node to the end. In the
// easiest case, this can be built using SimpleHeuristic.
func FindPath(
	inputGraph map[*Node]struct{}, start *Node, end *Node, heuristic Heuristic,
) ([]*Node, error) {

	// Convert to internally-used data structure. Both are the same, really, so we can convert
	// with almost no cost.
	graph := Graph(inputGraph)

	// Sanity checks
	if !graph.Has(start) {
		return []*Node{}, fmt.Errorf("input sanitation: start node not in graph")
	}
	if !graph.Has(end) {
		return []*Node{}, fmt.Errorf("input sanitation: end node not in graph")
	}

	// Variable open is our open list containing all nodes that should still be checked. At the
	// beginning, this is only the start node.
	open := Graph{start: graphVal}

	err := findPath(&open, end, heuristic)
	if err != nil {
		return []*Node{}, fmt.Errorf("error during path finding: %s", err.Error())
	}
	// The only time the prev member of the end node is set is when a path has been found.
	if end.prev == nil {
		if len(open) == 0 {
			return []*Node{}, fmt.Errorf("no path found: all nodes exhaused")
		}
		return []*Node{}, fmt.Errorf(
			"no path found: no connection to end node found from start node",
		)
	}
	// Extract a path from end to start.
	invPath, err := extractPath(end, start)
	if err != nil {
		return []*Node{}, fmt.Errorf("error during path extraction: %s", err.Error())
	}
	// Reverse the path to restore the original order.
	path := make([]*Node, 0, len(invPath))
	for idx := len(invPath) - 1; idx >= 0; idx-- {
		path = append(path, invPath[idx])
	}

	return path, nil
}

// Function extractPath follows the connection from the end to the beginning and returns it. It
// begins at end and follows the prev member until it reaches start or until there is no prev
// member. In the latter case, an error is returned.
func extractPath(end, start *Node) ([]*Node, error) {
	result := []*Node{}
	for currNode := end; currNode != nil && currNode != start; currNode = currNode.prev {
		// Somehow, one node didn't have its prev member set correctly. Fail in that case.
		if currNode.prev == nil {
			return []*Node{}, fmt.Errorf("prev member of node %s not set", currNode.ToString())
		}
		result = append(result, currNode)
	}
	// Don't forget the starting node.
	result = append(result, start)
	return result, nil
}

// Internal function that controls the pathfinding.
func findPath(open *Graph, end *Node, heuristic Heuristic) error {
	// Variable closed is our closed list. At the beginning, it is empty.
	closed := Graph{}
	for len(*open) != 0 && !closed.Has(end) {
		// Find the next cheapest node from the open list. This removes it as well as return it.
		nextCheckNode := open.PopCheapest(heuristic)
		// Add this node to the closed list.
		closed.Add(nextCheckNode)
		// Process each of the neighbours.
		for neigh := range nextCheckNode.connections {
			// If a neighbour is already on the closed list, skip it. Don't modify it at all.
			if closed.Has(neigh) {
				continue
			}
			if open.Has(neigh) {
				// Update the node in case we found a better path to it.
				newCost := nextCheckNode.trackedCost + neigh.Cost
				if newCost < neigh.trackedCost {
					neigh.prev = nextCheckNode
					neigh.trackedCost = newCost
				}
			} else {
				if neigh.prev != nil {
					return fmt.Errorf("node %s already has a predecessor", neigh.ToString())
				}
				// Add the new, as yet unknown node to the open list.
				open.Add(neigh)
				neigh.prev = nextCheckNode
				neigh.trackedCost = nextCheckNode.trackedCost + neigh.Cost
			}
		}
	}
	return nil
}
