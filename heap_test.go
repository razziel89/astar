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
		heap.Push(node)
	}
	assert.Equal(t, 3, heap.Len())
}

func TestHeapAddRemove(t *testing.T) {
	heap := Heap{}
	var expected *Node
	for idx := 0; idx < 3; idx++ {
		node, err := NewNode("", 0, 0, nil)
		if expected == nil {
			expected = node
		}
		assert.NoError(t, err)
		heap.Add(node)
	}
	assert.Equal(t, 3, heap.Len())
	heap.Remove(expected)
	assert.Equal(t, 2, heap.Len())
}

func TestHeapPopHas(t *testing.T) {
	heap := Heap{}
	goheap.Init(&heap)
	var expected *Node
	for _, cost := range []int{5, 2, 4, 6, 0, 4, 6, 2} {
		node, err := NewNode(fmt.Sprint(cost), cost, 0, nil)
		node.trackedCost = cost
		if expected == nil {
			expected = node
		}
		assert.NoError(t, err)
		goheap.Push(&heap, node)
	}
	assert.True(t, heap.Has(expected))
	popped := goheap.Pop(&heap).(*Node)
	assert.Equal(t, 0, popped.Cost)
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
		goheap.Push(&heap, node)
	}
	popped := goheap.Pop(&heap).(*Node)
	assert.NotNil(t, popped)
	popped = goheap.Pop(&heap).(*Node)
	assert.NotNil(t, popped)
	popped = heap.PopCheapest()
	assert.NotNil(t, popped)
	popped = heap.PopCheapest()
	assert.Nil(t, popped)
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
		goheap.Push(&heap, node)
	}
	expectedString := "{id: node2, cost: 2, con: ['']} -> 3\n" +
		"{id: node4, cost: 4, con: ['']} -> 5\n" +
		"{id: node5, cost: 5, con: ['']} -> 6"
	assert.Equal(t, expectedString, heap.ToString(mockHeuristic))
}
