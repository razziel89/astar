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

func TestNewNode(t *testing.T) {
	node, err := NewNode("node", 1, 0, "some payload")
	assert.NoError(t, err)
	// Test public members are set correctly.
	assert.Equal(t, node.ID, "node")
	assert.Equal(t, node.Cost, 1)
	assert.Equal(t, node.Payload, "some payload")
	// Test private members.
	assert.Zero(t, len(node.connections))
	assert.Zero(t, node.trackedCost)
	assert.Nil(t, node.prev)
}

func TestNewNodeNegativeCost(t *testing.T) {
	_, err := NewNode("node", -1, 0, nil)
	assert.Error(t, err)
}

func TestNewNodeNegativeNumNeighbours(t *testing.T) {
	_, err := NewNode("node", 1, -1, nil)
	assert.NoError(t, err)
}

func TestNodeConnection(t *testing.T) {
	node1, err := NewNode("node1", 1, 0, nil)
	assert.NoError(t, err)
	node2, err := NewNode("node2", 1, 0, nil)
	assert.NoError(t, err)

	node1.AddConnection(node2)
	// Adding a connection only adds it in one way.
	assert.NotZero(t, len(node1.connections))
	assert.Zero(t, len(node2.connections))

	// Adding a connection multiple times is no problem.
	node1.AddConnection(node2)
	assert.NotZero(t, len(node1.connections))

	// Removiung a connection multiple times is no problem.
	node1.RemoveConnection(node2)
	assert.Zero(t, len(node1.connections))
	node1.RemoveConnection(node2)
	assert.Zero(t, len(node1.connections))
}

func TestToString(t *testing.T) {
	node1, err := NewNode("node1", 1, 0, nil)
	assert.NoError(t, err)
	node2, err := NewNode("node2", 1, 0, nil)
	assert.NoError(t, err)

	// Present connections are shown.
	node1.AddConnection(node2)
	assert.Equal(t, "{id: node1, cost: 1, con: ['node2']}", node1.ToString())
	node1.RemoveConnection(node2)
	assert.Equal(t, "{id: node1, cost: 1, con: ['']}", node1.ToString())
}
