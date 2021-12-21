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

// HeapElement is an element of a heap.
type HeapElement struct {
	Node     *Node
	Estimate int
}

// Heap is a collection of nodes on a minimum heap. It implements Go's heap.Interface. It is used by
// the heapable graph to store the actual nodes. Don't use this data structre on its own.
type Heap []HeapElement

// Len provides the length of the heap. This is needed for Go's heap interface.
func (h *Heap) Len() int {
	return len(*h)
}

// Less determines whether one value is smaller than another one. This is needed for Go's heap
// interface.
func (h *Heap) Less(i, j int) bool {
	return (*h)[i].Node.trackedCost < (*h)[j].Node.trackedCost
}

// Swap swaps two values in the heap. This is needed for Go's heap interface.
func (h *Heap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

// Push adds a value to the heap. This is needed for Go's heap interface. Don't use it directly, use
// Add to add nodes. This one will panic if you provide an incorrect type thanks to Go's lack of
// generics.
func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(HeapElement))
}

// Pop removes a value from the heap. This is needed for Go's heap interface.
func (h *Heap) Pop() interface{} {
	if len(*h) == 0 {
		return nil
	}
	length := len(*h)
	last := (*h)[length-1]
	*h = (*h)[0 : length-1]
	return last
}

// ToString provides a string representation of the heap. The nodes are sorted according to their
// user-defined names. If you provide a heuristic != nil, the value that heuristic provides for each
// node is also provided at the end of a line. Provide nil to disable.
func (h *Heap) ToString(heuristic Heuristic) string {
	nodes := make([]*Node, 0, len(*h))
	for _, nodeElem := range *h {
		nodes = append(nodes, nodeElem.Node)
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
