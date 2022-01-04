/* An implementation of the A* algorithm in plain Golang.
Copyright (C) 2021  Torsten Sachse

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package astar implements the A* path finding algorithm. The main function is FindPath. Please
// find a description here: https://en.wikipedia.org/wiki/A*_search_algorithm
package astar

import (
	"fmt"
)

// These simplify tests by replacing them with mock implementations.
var (
	extractPath     = ExtractPath
	findReversePath = FindReversePath
	resetFn         = nodeResetFn
)

func nodeResetFn(node *Node) error {
	node.prev = nil
	node.trackedCost = defaultCost
	return nil
}

// FindPath finds the path between the start and end node. It is a convenience wrapper around
// FindReversePath. FindPath takes a graph in the form of a set of nodes, a start node, and an end
// node. It returns errors in case there are problems with the input or during execution. The path
// is returned in the correct order. This is achieved by using the normal algorithm and reversing
// the path at the end.
//
// This function requires you to provide one of the graph types from this package as the `graph`
// argument. If you want to use your own implementation of a GraphOps, you will need to use
// FindReversePath instead. Use NewGraph or NewHeapedGraph to obtain a suitable data structure.
//
// This implementation modifies the original nodes during execution! In the end, the nodes are
// reverted to null states (private members set to the appropriate zero values), which allows you to
// use the same input graph again with FindPath.
//
// FindPath also takes a heuristic that estimates the cost for moving from a node to the end. In the
// easiest case, this can be built using ConstantHeuristic. This heuristic is evaluated exactly once
// for a node when adding that node to an internal graph.
func FindPath(graph GraphOps, start *Node, end *Node, heuristic Heuristic) ([]*Node, error) {

	// Sanity checks
	if !graph.Has(start) {
		return []*Node{}, fmt.Errorf("input sanitation: start node not in graph")
	}
	if !graph.Has(end) {
		return []*Node{}, fmt.Errorf("input sanitation: end node not in graph")
	}

	// Open and closed lists will be of the same type as the input graph. To support that, we assert
	// the type here and initialise appropriately.
	var open, closed GraphOps
	switch graph.(type) {
	case *Graph:
		open = NewGraph()
		closed = NewGraph()
	case *HeapedGraph:
		open = NewHeapedGraph(1)
		closed = NewHeapedGraph(1)
	default:
		err := fmt.Errorf(
			"unknown input GraphOps type, if you provided your own, use FindReversePath directly",
		)
		return []*Node{}, err
	}

	// Variable open is our open list containing all nodes that should still be checked. At the
	// beginning, this is only the start node.
	// The closed list is empty at the beginning.
	open.Push(start, graphVal)

	err := findReversePath(open, closed, end, heuristic)
	if err != nil {
		return []*Node{}, fmt.Errorf("error during path finding: %s", err.Error())
	}
	// The only time the prev member of the end node is set is when a path has been found.
	if end.prev == nil {
		err := fmt.Errorf("no path found: no connection to end node found from start node")
		return []*Node{}, err
	}
	// Extract a path from end to start in the order from start to end.
	path, err := extractPath(end, start, true)
	if err != nil {
		return []*Node{}, fmt.Errorf("internal error during path extraction: %s", err.Error())
	}

	// Set the prev pointer back to nil. That way, the input graph can be used again. Also set the
	// tracked cost back to zero.
	err = graph.Apply(resetFn)
	if err != nil {
		return []*Node{}, fmt.Errorf("internal error during node reset: %s", err.Error())
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
func FindReversePath(open, closed GraphOps, end *Node, heuristic Heuristic) error {
	for open.Len() != 0 && !closed.Has(end) {
		// Find the next cheapest node from the open list. This removes it as well as return it.
		nextCheckNode := open.PopCheapest()
		// Add this node to the closed list.
		closed.Push(nextCheckNode, heuristic(nextCheckNode))
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
				open.Push(neigh, heuristic(neigh))
				neigh.prev = nextCheckNode
				neigh.trackedCost = nextCheckNode.trackedCost + neigh.Cost
			}
		}
	}
	return nil
}
