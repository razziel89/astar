// Package astar implements the A* path finding algorithm. The main function is FindPath. Please
// find a description here: https://en.wikipedia.org/wiki/A*_search_algorithm
package astar

import "fmt"

// Function nodeListToMap converts a list of nodes to a map to simplify access. If the nodes in the
// list are not unique, an error is returned.
func nodeListToMap(graph []*Node) (Graph, error) {
	result := make(Graph, len(graph))
	for _, node := range graph {
		result.Add(node)
	}
	if len(result) != len(graph) {
		err := fmt.Errorf("duplicate nodes in input")
		return Graph{}, err
	}
	return result, nil
}

// FindPath finds the path between the start and end node. It takes a graph in the form of a list of
// nodes, a start node, and an end node. It returns errors in case there are problems with the input
// or during execution. The path is returned in the correct order. This is achieved by using the
// normal algorithm and swapping start and end node at the beginning.
// This implementation modifies the original nodes!
func FindPath(inputGraph []*Node, startNode *Node, endNode *Node) ([]*Node, error) {
	// Swap nodes to return path in expected order.
	start, end := endNode, startNode

	// Variable graph is our open list containing all nodes that should be checked.
	graph, err := nodeListToMap(inputGraph)

	// Sanity checks
	if err != nil {
		return []*Node{}, fmt.Errorf("input sanitation: %s", err.Error())
	}
	if !graph.Has(start) {
		return []*Node{}, fmt.Errorf("input sanitation: start node not in graph")
	}
	if !graph.Has(end) {
		return []*Node{}, fmt.Errorf("input sanitation: end node not in graph")
	}

	closed := Graph{}

	err = findPath(&graph, &closed)
	if err != nil {
		return []*Node{}, fmt.Errorf("error during path finding: %s", err.Error())
	}
	if len(graph) == 0 {
		return []*Node{}, fmt.Errorf("no path found: all nodes exhaused")
	}
	// The only time the prev member of the end node is set is when a path has been found.
	if end.prev == nil {
		return []*Node{}, fmt.Errorf(
			"no path found: no connection to end node found from start node",
		)
	}
	// Extract a path from end to start. The result will contain the path in the order desired by
	// the user.
	path, err := extractPath(end, start)
	if err != nil {
		return []*Node{}, fmt.Errorf("error during path extraction: %s", err.Error())
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
	return result, nil
}

// Internal function that controls the pathfinding.
func findPath(open, closed *Graph) error {
	return nil
}
