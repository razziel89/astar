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
	goheap "container/heap"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapLenPush(t *testing.T) {
	heap := Heap{}
	for idx := 0; idx < 3; idx++ {
		node, err := NewNode("", 0, 0, nil)
		assert.NoError(t, err)
		heap.Push(HeapElement{Node: node, Estimate: 0})
	}
	assert.Equal(t, 3, heap.Len())
}

func TestHeapPop(t *testing.T) {
	heap := Heap{}
	goheap.Init(&heap)
	var expected *Node
	for _, cost := range []int{5, 2, 4, 6, 0, 4, 6, 2} {
		node, err := NewNode(fmt.Sprint(cost), cost, 0, nil)
		node.trackedCost = cost
		if cost == 0 {
			expected = node
		}
		assert.NoError(t, err)
		goheap.Push(&heap, HeapElement{Node: node, Estimate: 0})
	}
	popped := goheap.Pop(&heap).(HeapElement)
	assert.Equal(t, 0, popped.Node.Cost)
	assert.Equal(t, expected, popped.Node)
}

func TestHeapPushPopAndPopCheapest(t *testing.T) {
	heap := Heap{}
	goheap.Init(&heap)
	var expected *Node
	for _, cost := range []int{0, 1, 2} {
		node, err := NewNode(fmt.Sprint(cost), cost, 0, nil)
		node.trackedCost = cost
		if expected == nil {
			expected = node
		}
		assert.NoError(t, err)
		goheap.Push(&heap, HeapElement{Node: node, Estimate: 0})
	}
	popped := goheap.Pop(&heap).(HeapElement)
	assert.NotNil(t, popped)
	popped = goheap.Pop(&heap).(HeapElement)
	assert.NotNil(t, popped)
	popped = goheap.Pop(&heap).(HeapElement)
	assert.NotNil(t, popped)
	assert.Zero(t, heap.Len())
}

func TestHeapToString(t *testing.T) {
	heap := Heap{}
	goheap.Init(&heap)
	var expected *Node
	for _, cost := range []int{5, 2, 4} {
		node, err := NewNode(fmt.Sprintf("node%d", cost), cost, 0, nil)
		node.trackedCost = cost
		if expected == nil {
			expected = node
		}
		assert.NoError(t, err)
		goheap.Push(&heap, HeapElement{Node: node, Estimate: 0})
	}
	expectedString := "{id: node2, cost: 2, con: ['']} -> 3\n" +
		"{id: node4, cost: 4, con: ['']} -> 5\n" +
		"{id: node5, cost: 5, con: ['']} -> 6"
	assert.Equal(t, expectedString, heap.ToString(mockHeuristic))
}
