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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGraph(t *testing.T) {
	graph := NewHeapedGraph(100).(*HeapedGraph)
	assert.Equal(t, 100, cap(graph.Heap))
}

func TestHeapedGraphAddRemoveSuccess(t *testing.T) {
	node, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)
	graph := NewHeapedGraph(0).(*HeapedGraph)
	assert.Equal(t, 0, len(graph.Heap))
	graph.Add(node)
	assert.Equal(t, 1, len(graph.Heap))
	graph.Remove(node)
	assert.Equal(t, 0, len(graph.Heap))
	// Another removal is no problem.
	graph.Remove(node)
	assert.Equal(t, 0, len(graph.Heap))
	// Adding a node twice is no problem.
	graph.Add(node)
	assert.Equal(t, 1, len(graph.Heap))
	graph.Add(node)
	assert.Equal(t, 1, len(graph.Heap))
}

func TestHeapedGraphPushPopCheapestNoExist(t *testing.T) {
	node, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)
	graph := NewHeapedGraph(0).(*HeapedGraph)
	assert.Equal(t, 0, len(graph.Heap))
	graph.Push(node, 1)
	assert.Equal(t, 1, len(graph.Heap))
	popped := graph.PopCheapest()
	assert.Equal(t, 0, len(graph.Heap))
	assert.Equal(t, node, popped)
	popped = graph.PopCheapest()
	assert.Nil(t, popped)
}

func TestHeapedGraphPopCheapest(t *testing.T) {
	graph := NewHeapedGraph(0)
	var expectedCheapest *Node
	for idx, cost := range []int{1, 2, 0, 3} {
		node, err := NewNode(fmt.Sprintf("node%d", idx), cost, 0, nil)
		node.trackedCost = node.Cost
		assert.NoError(t, err)
		graph.Add(node)
		if cost == 0 {
			expectedCheapest = node
		}
	}
	cheapest := graph.PopCheapest()
	assert.Equal(t, expectedCheapest, cheapest)
}

func TestHeapedGraphHas(t *testing.T) {
	node, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)
	graph := NewHeapedGraph(0)
	assert.False(t, graph.Has(node))
	graph.Add(node)
	assert.True(t, graph.Has(node))
}

func TestHeapedGraphApplySuccess(t *testing.T) {
	node, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)

	graph := NewHeapedGraph(0)
	graph.Add(node)

	mockApplyFn := func(node *Node) error {
		node.prev = node
		return nil
	}

	assert.Nil(t, node.prev)
	err = graph.Apply(mockApplyFn)
	assert.NoError(t, err)
	assert.Equal(t, node, node.prev)
}

func TestHeapedGraphApplyFail(t *testing.T) {
	node1, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)
	node2, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)

	graph := NewHeapedGraph(0).(*HeapedGraph)
	graph.Add(node1)
	graph.Add(node2)

	mockApplyFn := func(node *Node) error {
		node.prev = node
		return errMock
	}

	assert.Nil(t, node1.prev)
	assert.Nil(t, node2.prev)
	err = graph.Apply(mockApplyFn)
	assert.Error(t, err)
	// Ensure the apply fn has not been applied to all nodes but only until an error happened.
	nilCount := 0
	for _, elem := range graph.Heap {
		if elem.Node.prev == nil {
			nilCount++
		}
	}
	assert.Equal(t, 1, nilCount)
}

func TestHeapedGraphLen(t *testing.T) {
	node, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)
	graph := NewHeapedGraph(0)
	assert.Equal(t, 0, graph.Len())
	graph.Add(node)
	assert.Equal(t, 1, graph.Len())
}

func TestHeapedGraphGraphMemberSetPanic(t *testing.T) {
	node, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)
	otherGraph := NewHeapedGraph(0)
	graph := NewHeapedGraph(0)

	assert.Zero(t, graph.Len())
	node.graph = otherGraph
	assert.NotNil(t, node.graph)
	defer func() {
		err, wasError := recover().(error)
		assert.True(t, wasError)
		assert.Error(t, err)
	}()
	graph.Add(node)
}

func TestHeapedGraphToString(t *testing.T) {
	graph := NewHeapedGraph(0).(*HeapedGraph)
	for idx, cost := range []int{1, 2, 0, 3} {
		node, err := NewNode(fmt.Sprintf("node%d", idx), cost, 0, nil)
		assert.NoError(t, err)
		graph.Add(node)
	}
	str := graph.ToString(mockHeuristic)
	expectedStr := "{id: node0, cost: 1, con: ['']} -> 2\n" +
		"{id: node1, cost: 2, con: ['']} -> 3\n" +
		"{id: node2, cost: 0, con: ['']} -> 1\n" +
		"{id: node3, cost: 3, con: ['']} -> 4"
	assert.Equal(t, expectedStr, str)
}
