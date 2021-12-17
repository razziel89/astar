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

// Graph is a collection of nodes. Note that there are no guarantees for the nodes to be connected.
// Ensuring that is the user's task.
type Graph map[*Node]struct{}

// This is the default value for the graph. Specifying it once here simplifies the code.
var graphVal = struct{}{}

// GraphVal is a convenience wrapper to return the default graph value.
func GraphVal() struct{} {
	return graphVal
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
		if !found || node.trackedCost+estimatedCost < cost {
			found = true
			result = node
			cost = node.trackedCost + estimatedCost
		}
	}
	g.Remove(result)
	return result
}

// ToString provides a string representation of the graph. The nodes are sorted according to their
// user-defined names. If you provide a heuristic != nil, the value that heuristic provides for each
// node is also provided at the end of a line. Provide nil to disable.
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
		if heuristic != nil {
			str += fmt.Sprintf(" -> %d", heuristic(node))
		}
		if idx != len(nodes)-1 {
			str += "\n"
		}
	}
	return str
}
