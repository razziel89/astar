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
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockPath []*Node
var mockGraph Graph
var mockStart *Node
var mockEnd *Node

// Set up test cases for find path. The two functions used by it will return the provided error
// values.
func setUpFindPath(errFindReverse, errExtract error, connect bool) func() {

	node1, _ := NewNode("node1", 0, 0, nil)
	node2, _ := NewNode("node2", 0, 0, nil)

	mockPath = []*Node{node1, node2}

	mockStart = node1
	mockEnd = node2

	mockGraph = Graph{}
	mockGraph.Add(mockStart)
	mockGraph.Add(mockEnd)

	// Ensure there is an actual connection if desired.
	if connect {
		mockEnd.prev = mockStart
	}

	extractPath = func(_, _ *Node, _ bool) ([]*Node, error) {
		return mockPath, errExtract
	}

	findReversePath = func(_, _ *Graph, _ *Node, _ Heuristic) error {
		return errFindReverse
	}

	return func() {
		// Revert changes.
		extractPath = ExtractPath
		findReversePath = FindReversePath
		mockPath = []*Node{}
		mockGraph = Graph{}
		mockStart = nil
		mockEnd = nil
	}
}

func TestFindPathSuccess(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, true)
	defer tearDown()

	path, err := FindPath(mockGraph, mockStart, mockEnd, mockHeuristic)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(path))
	assert.Equal(t, path[0], mockStart)
	assert.Equal(t, path[1], mockEnd)
}
