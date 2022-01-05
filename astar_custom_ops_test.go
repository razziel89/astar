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

// A custom graph ops implementation used for testing of errors due to custom graph ops.
type mockGraphOps struct{}

func (mo *mockGraphOps) Len() int {
	return 0
}
func (mo *mockGraphOps) Has(_ *Node) bool {
	return true // Simulate to contain all nodes.
}
func (mo *mockGraphOps) Add(_ *Node)         {}
func (mo *mockGraphOps) Push(_ *Node, _ int) {}
func (mo *mockGraphOps) Remove(*Node)        {}
func (mo *mockGraphOps) PopCheapest() *Node {
	return nil
}
func (mo *mockGraphOps) Apply(func(*Node) error) error {
	return nil
}
func (mo *mockGraphOps) UpdateIfBetter(*Node, *Node, int) {}

func TestFindPathFailureCustomGraphOps(t *testing.T) {

	mockStart, _ := NewNode("start", 0, 0, nil)
	mockEnd, _ := NewNode("end", 0, 0, nil)

	mockStart.AddPairwiseConnection(mockEnd)

	graph := mockGraphOps{}
	graph.Add(mockStart)
	graph.Add(mockEnd)

	_, err := FindPath(&graph, mockStart, mockEnd, mockHeuristic)
	assert.Error(t, err)
}
