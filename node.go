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
	"strings"
)

const (
	defaultCost = 0
)

// Node is a node for a connected graph along which to travel. Use NewNode to create one. It *will*
// be modified while the algorithm is being executed but FindPath reverts it to its original state
// at the end.
type Node struct {
	// Public members follow.
	// ID identifies the node. It is just a nice representation for the user and not used by the
	// algorithm.
	ID string
	// Cost specifies the cost of accessing this node from one connected to it.
	Cost int
	// Payload is some arbitrary user-defined payload that can be used with the heuristic, for
	// example. Type checks are the user's obligation.
	Payload interface{}
	// Private members follow.
	// Member connections determines which nodes this one is connected to.
	connections Graph
	// Member trackedCost tracks the accumulated minimal cost for reaching this node. If the prev
	// member is still nil, the algorithm assumes that no costs have been tracked yet.
	trackedCost int
	// Member prev tracks the previous node on the minimal cost connection.
	prev *Node
	// Member heap tracks which graph this node is in.
	graph GraphOps
}

// NewNode creates a new node. Provide an id string that describes this node for the user. Also
// provide a non-negative cost value. If the cost is negative, an error is returned. For performance
// reasons, specify the number of expected neighbours. A non-positive value means you are not sure
// or don't want to optimise this part, which is fine, too.
func NewNode(id string, cost int, numExpectedNeighbours int, payload interface{}) (*Node, error) {
	if cost < 0 {
		return nil, fmt.Errorf("cannot apply negative cost")
	}
	if numExpectedNeighbours < 0 {
		numExpectedNeighbours = 0
	}
	startCost := defaultCost
	newNode := Node{
		ID:          id,
		Cost:        cost,
		Payload:     payload,
		connections: make(Graph, numExpectedNeighbours),
		trackedCost: startCost,
		prev:        nil,
	}
	return &newNode, nil
}

// AddConnection adds a connection to a node. If the connection already exists, this is a no-op.
func (n *Node) AddConnection(neighbour *Node) {
	n.connections[neighbour] = graphVal
}

// AddPairwiseConnection adds a connection to a node and from that node back to the receiver.. If
// the connection already exists, this is a no-op.
func (n *Node) AddPairwiseConnection(neighbour *Node) {
	n.AddConnection(neighbour)
	neighbour.AddConnection(n)
}

// RemoveConnection removes a connection to a node. If the speicified node does not connect to this
// node, this is a no-op.
func (n *Node) RemoveConnection(neighbour *Node) {
	delete(n.connections, neighbour)
}

// ToString provides a nice string representation for this node. Not all members are used.
func (n *Node) ToString() string {
	conStrings := make([]string, 0, len(n.connections))
	for con := range n.connections {
		conStrings = append(conStrings, con.ID)
	}
	conString := strings.Join(conStrings, "', '")
	return fmt.Sprintf(
		"{id: %s, cost: %d, con: ['%s']}",
		n.ID, n.Cost, conString,
	)
}
