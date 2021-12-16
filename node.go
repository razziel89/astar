package astar

import (
	"fmt"
	"strings"
)

// Node is a node for a connected graph along which to travel. Use NewNode to create one. It *will*
// be modified while the algorithm is being executed.
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
func NewNode(id string, cost int, numExpectedNeighbours int, payload interface{}) (*Node, error) {
	if cost < 0 {
		return nil, fmt.Errorf("cannot apply negative cost")
	}
	if numExpectedNeighbours < 0 {
		numExpectedNeighbours = 0
	}
	startCost := 0
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
		"{id: %s, cost: %d, t: %d, con: ['%s']}",
		n.ID, n.Cost, n.trackedCost, conString,
	)
}
