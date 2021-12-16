package astar

// Graph is a collection of nodes. Note that there are no guarantees for the nodes to be connected.
// Ensuring that is the user's task.
type Graph map[*Node]struct{}

// This is the default value for the graph. It simplifies the code.
var graphVal = struct{}{}

// Has determines whether a graph contains a specific node.
func (g *Graph) Has(node *Node) bool {
	_, ok := (*g)[node]
	return ok
}
