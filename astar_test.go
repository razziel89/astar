package astar

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// A type useful for testing.
type intNode struct {
	posX int
	posY int
	cost int
}

// Function nodeListToMap converts a list of nodes to a map to simplify access. If the nodes in the
// list are not unique, an error is returned.
func nodeListToMap(graph []*Node) (Graph, error) {
	result := make(Graph, len(graph))
	for _, node := range graph {
		result.Add(node)
	}
	if len(result) != len(graph) {
		err := fmt.Errorf("duplicate nodes in input")
		return Graph{}, err
	}
	return result, nil
}

func (n *intNode) name() string {
	return fmt.Sprintf("%d-%d", n.posX, n.posY)
}

func TestBasicConnectionStraightLine(t *testing.T) {
	nodes := []*Node{}
	start := intNode{0, 0, 1}
	end := intNode{2, 2, 1}
	heuristic := ConstantHeuristic{}
	// Three nodes, one straight line connecting them.
	for _, datum := range []intNode{
		start,
		intNode{1, 1, 1},
		end,
	} {
		newNode, err := NewNode(datum.name(), datum.cost, 0, nil)
		assert.NoError(t, err)
		nodes = append(nodes, newNode)
		estimate := (end.posX - datum.posX) + (end.posY - datum.posY)
		err = heuristic.AddNode(newNode, estimate)
		assert.NoError(t, err)
	}
	for _, cons := range [][]*Node{
		[]*Node{nodes[0], nodes[1]},
		[]*Node{nodes[2], nodes[1]},
	} {
		init := cons[0]
		for _, con := range cons[1:] {
			// Add pairwise connections.
			init.AddConnection(con)
			con.AddConnection(init)
		}
	}
	nodeMap, err := nodeListToMap(nodes)
	assert.NoError(t, err)
	path, err := FindPath(nodeMap, nodes[0], nodes[len(nodes)-1], heuristic.Heuristic(0))
	assert.NoError(t, err)
	expectedPath := []*Node{nodes[0], nodes[1], nodes[2]}
	for _, node := range path {
		t.Log(node.ToString())
	}
	for idx, node := range path {
		assert.Equal(
			t, expectedPath[idx], node,
			fmt.Sprintf(
				"nodes not equal, expected %s found %s",
				expectedPath[idx].ToString(), node.ToString(),
			),
		)
	}
}

func TestBasicConnectionStraightLineWithEndingBranches(t *testing.T) {
	nodes := []*Node{}
	start := intNode{0, 0, 1}
	end := intNode{2, 2, 1}
	heuristic := ConstantHeuristic{}
	// Nine nodes, one straight line connecting start and end, but some nodes branch off from start
	// and end.
	for _, datum := range []intNode{
		// First row.
		start,
		intNode{1, 0, 1},
		intNode{2, 0, 1},
		// Second row.
		intNode{0, 1, 1},
		intNode{1, 1, 1},
		intNode{2, 1, 1},
		// Third row.
		intNode{0, 2, 1},
		intNode{1, 2, 1},
		end,
	} {
		newNode, err := NewNode(datum.name(), datum.cost, 0, nil)
		assert.NoError(t, err)
		nodes = append(nodes, newNode)
		estimate := (end.posX - datum.posX) + (end.posY - datum.posY)
		err = heuristic.AddNode(newNode, estimate)
		assert.NoError(t, err)
	}
	for _, cons := range [][]*Node{
		// Connections including start node.
		[]*Node{nodes[0], nodes[1], nodes[4], nodes[3]},
		// Connections including end node.
		[]*Node{nodes[8], nodes[7], nodes[5], nodes[4]},
		// The center node is connected to both start and end node.
	} {
		init := cons[0]
		for _, con := range cons[1:] {
			// Add pairwise connections.
			init.AddConnection(con)
			con.AddConnection(init)
		}
	}
	nodeMap, err := nodeListToMap(nodes)
	assert.NoError(t, err)
	path, err := FindPath(nodeMap, nodes[0], nodes[len(nodes)-1], heuristic.Heuristic(0))
	assert.NoError(t, err)
	expectedPath := []*Node{nodes[0], nodes[4], nodes[8]}
	for _, node := range path {
		t.Log(node.ToString())
	}
	for idx, node := range path {
		assert.Equal(
			t, expectedPath[idx], node,
			fmt.Sprintf(
				"nodes not equal, expected %s found %s",
				expectedPath[idx].ToString(), node.ToString(),
			),
		)
	}
}

//nolint:funlen
func TestBasicConnectionSquareEqualCost(t *testing.T) {
	nodes := []*Node{}
	start := intNode{0, 0, 1}
	end := intNode{2, 2, 1}
	heuristic := ConstantHeuristic{}
	for _, datum := range []intNode{
		// First row.
		start,
		intNode{1, 0, 1},
		intNode{2, 0, 1},
		// Second row.
		intNode{0, 1, 1},
		intNode{1, 1, 1},
		intNode{2, 1, 1},
		// Third row.
		intNode{0, 2, 1},
		intNode{1, 2, 1},
		end,
	} {
		newNode, err := NewNode(datum.name(), datum.cost, 0, nil)
		assert.NoError(t, err)
		nodes = append(nodes, newNode)
		estimate := (end.posX - datum.posX) + (end.posY - datum.posY)
		err = heuristic.AddNode(newNode, estimate)
		assert.NoError(t, err)
	}
	// Add pairwise connections but leave some out. This way, we always expect the very same path.
	// The connections look this way, with # being points and | or -- being connections, S is the
	// starting point and E is the end point:
	//
	//     #  #--E
	//        |
	//     #  #--#
	//     |  |
	//     S--#--#
	//
	// The only possible connection then looks like this:
	//
	//     #  #--E
	//        |
	//     #  #  #
	//        |
	//     S--#  #
	//
	for _, cons := range [][]*Node{
		// First row.
		[]*Node{nodes[0], nodes[1], nodes[3]},
		[]*Node{nodes[1], nodes[0], nodes[2], nodes[4]},
		[]*Node{nodes[2], nodes[1]},
		// Second row.
		[]*Node{nodes[3], nodes[0]},
		[]*Node{nodes[4], nodes[5], nodes[1], nodes[7]},
		[]*Node{nodes[5], nodes[4]},
		// Third row.
		[]*Node{nodes[6]},
		[]*Node{nodes[7], nodes[8]},
		[]*Node{nodes[8], nodes[7]},
	} {
		init := cons[0]
		for _, con := range cons[1:] {
			// Add pairwise connections.
			init.AddConnection(con)
			con.AddConnection(init)
		}
	}
	nodeMap, err := nodeListToMap(nodes)
	assert.NoError(t, err)
	path, err := FindPath(nodeMap, nodes[0], nodes[len(nodes)-1], heuristic.Heuristic(0))
	assert.NoError(t, err)
	expectedPath := []*Node{nodes[0], nodes[1], nodes[4], nodes[7], nodes[8]}
	for _, node := range path {
		t.Log(node.ToString())
	}
	for idx, node := range path {
		assert.Equal(
			t, expectedPath[idx], node,
			fmt.Sprintf(
				"nodes not equal, expected %s found %s",
				expectedPath[idx].ToString(), node.ToString(),
			),
		)
	}
}

//nolint:funlen
func TestBasicConnectionSquareVaryingCost(t *testing.T) {
	nodes := []*Node{}
	start := intNode{0, 0, 1}
	end := intNode{2, 2, 1}
	heuristic := ConstantHeuristic{}
	for _, datum := range []intNode{
		// First row.
		start,
		intNode{posX: 1, posY: 0, cost: 10},
		intNode{posX: 2, posY: 0, cost: 10},
		// Second row.
		intNode{posX: 0, posY: 1, cost: 1},
		intNode{posX: 1, posY: 1, cost: 10},
		intNode{posX: 2, posY: 1, cost: 1},
		// Third row.
		intNode{posX: 0, posY: 2, cost: 10},
		intNode{posX: 1, posY: 2, cost: 10},
		end,
	} {
		newNode, err := NewNode(datum.name(), datum.cost, 0, nil)
		assert.NoError(t, err)
		nodes = append(nodes, newNode)
		estimate := (end.posX - datum.posX) + (end.posY - datum.posY)
		err = heuristic.AddNode(newNode, estimate)
		assert.NoError(t, err)
	}
	// Add all pairwise connections. Some nods are so costly that they will never be visited. Those
	// are marked with X. The graph looks like this, with -- and | being connections, # being normal
	// nodes, and S (E) being the start (end) node:
	//
	//     X--X--E
	//     |  |  |
	//     #--X--#
	//     |  |  |
	//     S--X--X
	//
	// The cheapest connection then looks like this. One expensive node always has to be traversed.
	//
	//     #  #  E
	//           |
	//     #--#--#
	//     |
	//     S  #  #
	//
	for _, cons := range [][]*Node{
		// First row.
		[]*Node{nodes[0], nodes[1], nodes[3]},
		[]*Node{nodes[1], nodes[2], nodes[4]},
		[]*Node{nodes[2], nodes[5]},
		// Second row.
		[]*Node{nodes[3], nodes[4], nodes[6]},
		[]*Node{nodes[4], nodes[5], nodes[3], nodes[1], nodes[7]},
		[]*Node{nodes[5], nodes[8]},
		// Third row.
		[]*Node{nodes[6], nodes[7]},
		[]*Node{nodes[7], nodes[8], nodes[4]},
		[]*Node{nodes[8], nodes[7]},
	} {
		init := cons[0]
		for _, con := range cons[1:] {
			// Add pairwise connections.
			init.AddConnection(con)
			con.AddConnection(init)
		}
	}
	nodeMap, err := nodeListToMap(nodes)
	assert.NoError(t, err)
	path, err := FindPath(nodeMap, nodes[0], nodes[len(nodes)-1], heuristic.Heuristic(0))
	assert.NoError(t, err)
	expectedPath := []*Node{nodes[0], nodes[3], nodes[4], nodes[5], nodes[8]}
	for _, node := range path {
		t.Log(node.ToString())
	}
	for idx, node := range path {
		assert.Equal(
			t, expectedPath[idx], node,
			fmt.Sprintf(
				"nodes not equal, expected %s found %s",
				expectedPath[idx].ToString(), node.ToString(),
			),
		)
	}
}
