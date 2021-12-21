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

// Heap is a collection of nodes on a minimum heap. It implements Go's heap.Interface. It is used by
// the heapable graph to store the actual nodes. Don't use this data structre on its own.
type Heap []*Node

// Len provides the length of the heap. This is needed for Go's heap interface.
func (h *Heap) Len() int {
	return len(*h)
}

// Less determines whether one value is smaller than another one. This is needed for Go's heap
// interface.
func (h *Heap) Less(i, j int) bool {
	return (*h)[i].trackedCost < (*h)[j].trackedCost
}

// Swap swaps two values in the heap. This is needed for Go's heap interface.
func (h *Heap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

// Push adds a value to the heap. This is needed for Go's heap interface. Don't use it directly, use
// Add to add nodes. This one will panic if you provide an incorrect type thanks to Go's lack of
// generics.
func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(*Node))
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

// Has determines whether a heap contains a specific node. This is inefficient for a heap as we need
// to operate over much of the underlying slice.
func (h *Heap) Has(findNode *Node) bool {
	for _, node := range *h {
		if node == findNode {
			return true
		}
	}
	return false
}

// Add adds a node to the heap. If the node already exists, this a no-op. Note that repeated calls
// to this are inefficient for the heap since they query Has. Use Push if you are sure the node does
// not yet exist.
func (h *Heap) Add(node *Node) {
	if !h.Has(node) {
		h.Push(node)
	}
}

// Remove removes a node from the heap. If the node does not exist, this a no-op. Note that repeated
// calls to this are inefficient for the heap as they will have to iterate over the entire data
// structure. Even when removing a value, all data after the removed datum will be copied.
func (h *Heap) Remove(findNode *Node) {
	for idx, node := range *h {
		if node == findNode {
			*h = append((*h)[0:idx], (*h)[idx+1:]...)
			return
		}
	}
}

// PopCheapest retrieves one of the cheapest nodes and removes it. This will return nil if the heap
// is empty.
func (h *Heap) PopCheapest() *Node {
	if val, ok := h.Pop().(*Node); ok {
		return val
	}
	return nil
}

// ToString provides a string representation of the heap. The nodes are sorted according to their
// user-defined names. If you provide a heuristic != nil, the value that heuristic provides for each
// node is also provided at the end of a line. Provide nil to disable.
func (h *Heap) ToString(heuristic Heuristic) string {
	nodes := append([]*Node{}, (*h)...)

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
