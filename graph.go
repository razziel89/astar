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

package astar

import (
	"fmt"
	"sort"
)

// GraphOps is an interface needed for a graph to be usable with the path finding functions in this
// module. See the method documentation of the actual Graph type for what the individual methods do
// in detail.
type GraphOps interface {
	// Len specifies how many elements there are in the graph.
	Len() int
	// Has checks whether a node is in the graph.
	Has(node *Node) bool
	// Add adds a node to a graph without an estimate.
	Add(node *Node)
	// Push adds a node to a graph with an estimate.
	Push(node *Node, estimate int)
	// Remove removes a specific node from a graph.
	Remove(node *Node)
	// PopCheapest retrieves and removes the cheapest node from the graph. Cost is equal to the
	// node's cost added to its estimate.
	PopCheapest() *Node
	// Apply applies a function to all nodes in the graph. That function may error out.
	Apply(func(*Node) error) error
}

// Graph is a collection of nodes. Note that there are no guarantees for the nodes to be connected.
// Ensuring that is the user's task. Each nodes is assigned to its estimate. That means a node's
// estimate will never be able to change once added.
type Graph map[*Node]int

// This is the default value for the graph. Specifying it once here simplifies the code.
var graphVal = 0

// GraphVal is a convenience wrapper to return the default graph value.
func GraphVal() int {
	return graphVal
}

// Len determines the number of elements.
func (g *Graph) Len() int {
	return len(*g)
}

// Has determines whether a graph contains a specific node.
func (g *Graph) Has(node *Node) bool {
	_, ok := (*g)[node]
	return ok
}

// Add adds a node to the graph. If the node already exists, this a no-op.
func (g *Graph) Add(node *Node) {
	(*g)[node] = graphVal
}

// Push adds a node to the graph, including its estimate. If the node already exists, this a no-op.
func (g *Graph) Push(node *Node, estimate int) {
	(*g)[node] = estimate
}

// Remove removes a node from the graph. If the node does not exist, this a no-op.
func (g *Graph) Remove(node *Node) {
	delete(*g, node)
}

// PopCheapest retrieves one of the cheapest nodes and removes it. This will return nil if the graph
// is empty.
func (g *Graph) PopCheapest() *Node {
	found := false
	cost := 0
	var result *Node
	for node, estimatedCost := range *g {
		if !found || node.trackedCost+estimatedCost < cost {
			found = true
			result = node
			cost = node.trackedCost + estimatedCost
		}
	}
	g.Remove(result)
	return result
}

// Apply apples a function to all nodes in the graph.
func (g *Graph) Apply(fn func(*Node) error) error {
	for node := range *g {
		err := fn(node)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToString provides a string representation of the graph. The nodes are sorted according to their
// user-defined names. If you provide a heuristic != nil, the value that heuristic provides for each
// node is also provided at the end of a line. Providing nil will use the stored estimates.
func (g *Graph) ToString(heuristic Heuristic) string {
	nodes := make([]*Node, 0, len(*g))
	for node := range *g {
		nodes = append(nodes, node)
	}

	lessFn := func(idx1, idx2 int) bool {
		return nodes[idx1].ID < nodes[idx2].ID
	}
	sort.SliceStable(nodes, lessFn)
	str := ""
	for idx, node := range nodes {
		str += node.ToString()
		estimate := (*g)[node]
		if heuristic != nil {
			estimate = heuristic(node)
		}
		str += fmt.Sprintf(" -> %d", estimate)
		if idx < len(nodes)-1 {
			str += "\n"
		}
	}
	return str
}
