// Package astar implements the A* path finding algorithm. The main function is FindPath. Please
// find a description here: https://en.wikipedia.org/wiki/A*_search_algorithm
package astar

import "fmt"

// Function nodeListToMap converts a list of nodes to a map to simplify access. If the nodes in the
// list are not unique, an error is returned.
func nodeListToMap(graph []*Node) (Graph, error) {
	result := make(Graph, len(graph))
	for _, node := range graph {
		result[node] = graphVal
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
func FindPath(inputGraph []*Node, startNode *Node, endNode *Node) ([]*Node, error) {
	// Swap nodes to return path in expected order.
	start, end := endNode, startNode

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

	return []*Node{}, nil
}
