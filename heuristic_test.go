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

func TestConstantHeuristicAddNode(t *testing.T) {
	heuristic := ConstantHeuristic{}
	node, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)

	err = heuristic.AddNode(node, 1)
	assert.NoError(t, err)

	// Adding a node multiple times with the same estimate works.
	err = heuristic.AddNode(node, 1)
	assert.NoError(t, err)
}

func TestConstantHeuristicAddNodeFailure(t *testing.T) {
	heuristic := ConstantHeuristic{}
	node, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)

	err = heuristic.AddNode(node, 1)
	assert.NoError(t, err)

	// Adding a node multiple times with different estimates fails.
	err = heuristic.AddNode(node, -1)
	assert.Error(t, err)
}

func TestConstantHeuristicHeuristicFunction(t *testing.T) {
	heuristic := ConstantHeuristic{}
	node, err := NewNode("node", 0, 0, nil)
	assert.NoError(t, err)
	err = heuristic.AddNode(node, 1)
	assert.NoError(t, err)

	fn := heuristic.Heuristic(0)
	// Return the stored value for known nodes.
	assert.Equal(t, 1, fn(node))
	// Return the default value for unknown nodes. Here, nil serves as a placeholder for a pointer
	// to an unknown node.
	assert.Equal(t, 0, fn(nil))
}
