package astar

import (
	"fmt"
	"strings"
)

// Node is a node for a connected graph along which to travel. Use NewNode to create one.
type Node struct {
	// Public members follow.
	// ID identifies the node. It is just a nice representation for the user and not used by the
	// algorithm.
	ID string
	// Cost specifies the cost of accessing this node from one connected to it.
	Cost int
	// Private members follow.
	// Member connections determines which nodes this one is connected to.
	connections map[*Node]struct{}
	// Member trackedCost tracks the accumulated minimal cost for reaching this node. If the cost is
	// <0, it means no cost has yet been assigned.
	trackedCost int
	// Member prev tracks the previous node on the minimal cost connection.
	prev *Node
}

// NewNode creates a new node. Provide an id string that describes this node for the user. Also
// provide a non-negative cost value. If the cost is negative, an error is returned. For performance
// reasons, specify the number of expected neighbours. A non-positive value means you are not sure
// or don't want to optimise this part.
func NewNode(id string, cost int, numExpectedNeighbours int) (*Node, error) {
	if cost < 0 {
		return nil, fmt.Errorf("cannot apply negative cost")
	}
	if numExpectedNeighbours < 0 {
		numExpectedNeighbours = 0
	}
	newNode := Node{
		ID:          id,
		Cost:        cost,
		connections: make(map[*Node]struct{}, numExpectedNeighbours),
		trackedCost: -1,
		prev:        nil,
	}
	return &newNode, nil
}

// AddConnection adds a connection to a node. If the connection already exists, this is a no-op.
func (n *Node) AddConnection(neighbour *Node) {
	n.connections[neighbour] = struct{}{}
}

// RemoveConnection removes a connection to a node. If the speicified node does not connect to this
// node, this is a no-op.
func (n *Node) RemoveConnection(neighbour *Node) {
	delete(n.connections, neighbour)
}

// ToString provides a nice string representation for this node.
func (n *Node) ToString() string {
	conStrings := make([]string, 0, len(n.connections))
	for con := range n.connections {
		conStrings = append(conStrings, con.ID)
	}
	conString := strings.Join(conStrings, "', '")
	return fmt.Sprintf("{id: %s, cost: %d, con: ['%s']}", n.ID, n.Cost, conString)
}
