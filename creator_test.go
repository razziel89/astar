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

func TestDist2D(t *testing.T) {
	positions := [][2]int{
		[2]int{0, 0},
		[2]int{1, 1},
		[2]int{-10, 5},
		[2]int{5, -10},
		[2]int{2, 0},
		[2]int{0, 2},
	}
	expectedDistances := []int{
		0, 1, 11, 11, 2, 2,
		1, 0, 11, 11, 1, 1,
		11, 11, 0, 21, 13, 10,
		11, 11, 21, 0, 10, 13,
		2, 1, 13, 10, 0, 2,
		2, 1, 10, 13, 2, 0,
	}
	distIdx := 0
	for _, pos1 := range positions {
		for _, pos2 := range positions {
			dist := dist2D(pos1, pos2)
			assert.Equal(t, expectedDistances[distIdx], dist)
			distIdx++
		}
	}
}

func TestCreateRegular2DGridSuccess(t *testing.T) {
	for _, graphType := range []string{"default", "heaped"} {
		size := [2]int{10, 10}
		connections := [][2]int{
			[2]int{-1, 0},
			[2]int{0, -1},
			[2]int{1, 0},
			[2]int{0, 1},
		}

		graph, posMap, err := CreateRegular2DGrid(size, connections, graphType, 0)

		assert.NoError(t, err)
		assert.Equal(t, 100, graph.Len())
		assert.Equal(t, 100, len(posMap))
	}
}

func TestCreateRegular2DGridWrongTypeFailure(t *testing.T) {
	size := [2]int{10, 10}
	connections := [][2]int{}

	_, _, err := CreateRegular2DGrid(size, connections, "unknownType", 0)

	assert.Error(t, err)
}

func TestCreateRegular2DGridNodeCreationFailure(t *testing.T) {
	for _, graphType := range []string{"default", "heaped"} {
		size := [2]int{10, 10}
		connections := [][2]int{}

		_, _, err := CreateRegular2DGrid(size, connections, graphType, -1)

		assert.Error(t, err)
	}
}

func TestCreateConstantHeuristic2DSuccess(t *testing.T) {
	knownNode, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)
	unknownNode, err := NewNode("unknown", 0, 0, nil)
	assert.NoError(t, err)
	endPos := [2]int{10, 0}
	posMap := map[[2]int]*Node{
		[2]int{0, 0}: knownNode,
	}

	heuristic, err := CreateConstantHeuristic2D(posMap, endPos, 0)

	assert.NoError(t, err)
	assert.Equal(t, 10, heuristic(knownNode))
	assert.Equal(t, 0, heuristic(unknownNode))
}

func TestCreateConstantHeuristic2DFailure(t *testing.T) {
	node, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)
	endPos := [2]int{10, 0}
	posMap := map[[2]int]*Node{
		[2]int{0, 0}: node,
		[2]int{5, 5}: node, // No node must be added more than once to a heuristic.
	}

	_, err = CreateConstantHeuristic2D(posMap, endPos, 0)

	assert.Error(t, err)
}
