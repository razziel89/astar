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

func TestHeapPop(t *testing.T) {
	heap := Heap{}
	goheap.Init(&heap)
	var expected *Node
	for idx := 0; idx < 3; idx++ {
		node, err := NewNode(fmt.Sprint(idx), idx, 0, nil)
		if expected == nil {
			expected = node
		}
		assert.NoError(t, err)
		goheap.Push(&heap, node)
	}
	popped := goheap.Pop(&heap).(*Node)
	assert.Equal(t, expected, popped)
}
