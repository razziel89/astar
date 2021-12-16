package astar

import "sort"

// Graph is a collection of nodes. Note that there are no guarantees for the nodes to be connected.
// Ensuring that is the user's task.
type Graph map[*Node]struct{}

// This is the default value for the graph. Specifying it once here simplifies the code.
var graphVal = struct{}{}

// Has determines whether a graph contains a specific node.
func (g *Graph) Has(node *Node) bool {
	_, ok := (*g)[node]
	return ok
}

// Add adds a node to the graph. If the node already exists, this a no-op.
func (g *Graph) Add(node *Node) {
	(*g)[node] = graphVal
}

// Remove removes a node from the graph. If the node does not exist, this a no-op.
func (g *Graph) Remove(node *Node) {
	delete(*g, node)
}

// PopCheapest retrieves one of the cheapest nodes and removes it. This will return nil if the graph
// is empty.
func (g *Graph) PopCheapest(heuristic Heuristic) *Node {
	found := false
	cost := 0
	var result *Node
	for node := range *g {
		estimatedCost := heuristic(node)
		if !found || node.Cost+estimatedCost < cost {
			found = true
			result = node
		}
	}
	g.Remove(result)
	return result
}

// ToString provides a string representation of the graph. The nodes are sorted according to their
// names.
func (g *Graph) ToString() string {
	nodes := make([]*Node, 0, len(*g))
	for node := range *g {
		nodes = append(nodes, node)
	}

	lessFn := func(idx1, idx2 int) bool {
		return nodes[idx1].ID < nodes[idx2].ID
	}
	sort.SliceStable(nodes, lessFn)
	str := ""
	for _, node := range nodes {
		str += node.ToString() + "\n"
	}
	return str
}
