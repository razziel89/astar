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

	goheap "container/heap"
)

// HeapedGraph is a collection of nodes. Note that there are no guarantees for the nodes to be
// connected. Ensuring that is the user's task. Each nodes is assigned to its estimate. That means a
// node's estimate will never be able to change once added. Get a heaped graph via NewHeapedGraph.
// This uses a Heap as storage backend.
type HeapedGraph struct {
	Heap Heap
}

// NewHeapedGraph obtains a new heaped graph. Specify the estimated number of nodes as argument to
// boost performance.
func NewHeapedGraph(estimatedSize int) GraphOps {
	heap := make(Heap, 0, estimatedSize)
	self := &HeapedGraph{Heap: heap}
	goheap.Init(&self.Heap)
	return self
}

// Len determines the number of elements.
func (g *HeapedGraph) Len() int {
	return g.Heap.Len()
}

// Has determines whether a graph contains a specific node.
func (g *HeapedGraph) Has(node *Node) bool {
	return node.graph == g
}

// Add adds a node to the graph. If the node already exists, this a no-op. This panics if the node
// already has a different graph set.
func (g *HeapedGraph) Add(node *Node) {
	if !g.Has(node) {
		if node.graph != nil {
			panic(fmt.Errorf("different graph already set"))
		}
		elem := HeapElement{Node: node, Estimate: graphVal}
		goheap.Push(&g.Heap, elem)
		node.graph = g
	}
}

// Push adds a node to the graph, including its estimate. If the node already exists, this a no-op.
func (g *HeapedGraph) Push(node *Node, estimate int) {
	if !g.Has(node) {
		elem := HeapElement{Node: node, Estimate: estimate}
		goheap.Push(&g.Heap, elem)
		node.graph = g
	}
}

// Remove removes a node from the graph. If the node does not exist, this a no-op. This function is
// inefficient but not needed for the algorithm in general.
func (g *HeapedGraph) Remove(findNode *Node) {
	for idx, node := range g.Heap {
		if node.Node == findNode {
			g.Heap[idx].Node.graph = nil
			g.Heap[idx].Node = nil
			g.Heap = append(g.Heap[0:idx], g.Heap[idx+1:]...)
			return
		}
	}
}

// PopCheapest retrieves one of the cheapest nodes and removes it. This will return nil if the graph
// is empty.
func (g *HeapedGraph) PopCheapest() *Node {
	if len(g.Heap) > 0 {
		val := goheap.Pop(&g.Heap).(HeapElement)
		val.Node.graph = nil
		return val.Node
	}
	return nil
}

// Apply apples a function to all nodes in the graph.
func (g *HeapedGraph) Apply(fn func(*Node) error) error {
	for _, elem := range g.Heap {
		err := fn(elem.Node)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToString provides a string representation of the graph. The nodes are sorted according to their
// user-defined names. If you provide a heuristic != nil, the value that heuristic provides for each
// node is also provided at the end of a line. Providing nil will use the stored estimates.
func (g *HeapedGraph) ToString(heuristic Heuristic) string {
	nodes := make([]*Node, 0, len(g.Heap))
	estimates := make([]int, 0, len(g.Heap))
	for _, elem := range g.Heap {
		nodes = append(nodes, elem.Node)
		estimates = append(estimates, elem.Estimate)
	}

	lessFn := func(idx1, idx2 int) bool {
		return nodes[idx1].ID < nodes[idx2].ID
	}
	sort.SliceStable(nodes, lessFn)
	str := ""
	for idx, node := range nodes {
		str += node.ToString()
		estimate := estimates[idx]
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
