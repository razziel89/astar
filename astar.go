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
// This implementation modifies the original nodes during execution! In the end, the nodes are
// reverted to their original states, which allows you to use the same input graph again.
//
// It also takes a heuristic that estimates the cost for moving from a node to the end. In the
// easiest case, this can be built using ConstantHeuristic.
func FindPath(graph Graph, start *Node, end *Node, heuristic Heuristic) ([]*Node, error) {

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

	// The closed list is empty at the beginning.
	closed := Graph{}

	err := FindReversePath(&open, &closed, end, heuristic)
	if err != nil {
		return []*Node{}, fmt.Errorf("error during path finding: %s", err.Error())
	}
	// The only time the prev member of the end node is set is when a path has been found.
	if end.prev == nil {
		err := fmt.Errorf("no path found: no connection to end node found from start node")
		return []*Node{}, err
	}
	// Extract a path from end to start in the order from start to end.
	path, err := ExtractPath(end, start, true)
	if err != nil {
		return []*Node{}, fmt.Errorf("internal error during path extraction: %s", err.Error())
	}

	// Set the prev pointer back to nil. That way, the input graph can be used again. Also set the
	// tracked cost back to zero.
	for node := range graph {
		node.prev = nil
		node.trackedCost = defaultCost
	}

	return path, nil
}

// ExtractPath follows the connection from the end to the beginning and returns it. It begins at end
// and follows the prev member until it reaches start or until there is no prev member. In the
// latter case, an error is returned. Specify whether you want the original path, which is in
// reverse order, or the path from the original start to the end.
func ExtractPath(end, start *Node, orgOrder bool) ([]*Node, error) {
	invPath := []*Node{}
	for currNode := end; currNode != nil && currNode != start; currNode = currNode.prev {
		// Somehow, one node didn't have its prev member set correctly. Fail in that case.
		if currNode.prev == nil {
			return []*Node{}, fmt.Errorf("prev member of node %s not set", currNode.ToString())
		}
		invPath = append(invPath, currNode)
	}
	// Don't forget the starting node.
	invPath = append(invPath, start)
	if !orgOrder {
		return invPath, nil
	}
	// Reverse the path to restore the original order if that was desired.
	path := make([]*Node, 0, len(invPath))
	for idx := len(invPath) - 1; idx >= 0; idx-- {
		path = append(path, invPath[idx])
	}
	return path, nil
}

// FindReversePath finds a reverse path from the start node to the end node. Follow the prev member
// of the end node to traverse the path backwards. To use this function, in the beginning, the open
// list must contain the start node and the closed list must be empty.
func FindReversePath(open, closed *Graph, end *Node, heuristic Heuristic) error {
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
