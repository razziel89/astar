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
	"math"
)

func dist2D(pos1, pos2 [2]int) int {
	xDist := pos1[0] - pos2[0]
	yDist := pos1[1] - pos2[1]
	return int(math.Floor(math.Sqrt(float64(xDist*xDist) + float64(yDist*yDist))))
}

// CreateRegular2DGrid creates a regular 2D grid of connected nodes in a graph suitable for path
// finding. All nodes have a default cost of zero assigned.
//
// Provide the number of nodes in each direction via the `size` argument. Provide a list of
// displacements that connect a node to its neighbours via the `connections` argument. For example,
// to connect horizontally and vertically only, pass the following as input:
//  [][2]int{
//      [2]int{-1, 0},
//      [2]int{0, -1},
//      [2]int{1, 0},
//      [2]int{0, 1},
//  }
//
// Also provide the name of the type of graph you want to obtain as input. This can be "default" or
// "heaped". The heaped graph has a better performance but is a more complex data structure.
//
// CreateRegular2DGrid returns three values:
// 1. A graph object suitable for paht finding via FindPath.
// 2. A map from node positions to node pointers, this makes modifying the nodes easy, e.g. to set
//	  the costs.
// 3. An error value in case there were problems.
func CreateRegular2DGrid(
	size [2]int, connections [][2]int, dest [2]int, graphType string,
) (GraphOps, map[[2]int]*Node, error) {

	// These values just improve performances during allocation.
	neighbours := len(connections)
	gridSize := size[0] * size[1]

	var graph GraphOps
	switch graphType {
	case "default":
		graph = NewGraph(gridSize)
	case "heaped":
		graph = NewHeapedGraph(gridSize)
	default:
		return nil, nil, fmt.Errorf("unknown graph type, need 'default' or 'heaped'")
	}

	// We remember the node for each position. This makes creating connections easier later on.
	posToNode := map[[2]int]*Node{}

	// Create the actual nodes, connections will be created later.
	for x := 0; x < size[0]; x++ {
		for y := 0; y < size[1]; y++ {
			// Create a name for the node.
			nodeName := fmt.Sprintf("x:%d,y:%d", x, y)
			// Create the node with empty payload. The payload is for the user to use. Don't provide
			// a cost for now. The user needs to set that later.
			node, err := NewNode(nodeName, 0, neighbours, nil)
			if err != nil {
				return nil, nil, err
			}
			// Remember the node's position and add it to the graph.
			posToNode[[2]int{x, y}] = node
			graph.Add(node)
		}
	}

	// Create connections.
	for x := 0; x < size[0]; x++ {
		for y := 0; y < size[1]; y++ {
			// Add all connections for this node.
			node := posToNode[[2]int{x, y}]
			for _, disp := range connections {
				neighPos := [2]int{x + disp[0], y + disp[1]}
				// Only add connections if the neighbour actually exists!
				if neigh, exists := posToNode[neighPos]; exists {
					node.AddPairwiseConnection(neigh)
				}
			}
		}
	}

	return graph, posToNode, nil
}

// CreateConstantHeuristic2D creates a constant heuristic for a regular 2D grid. It takes a map from
// node positions to node pointers and the desired end position and returns a simple, constant
// heuristic that estimates the remaining cost as the line-of-sight distance to the desired
// destination. Heuristics must return a value for all nodes, even ones they don't remember. For
// such nodes, it returns `defaultVal`.
func CreateConstantHeuristic2D(
	posMap map[[2]int]*Node, dest [2]int, defaultVal int,
) (Heuristic, error) {
	// Create heuristic that remembers distance estimates.
	heuristic := ConstantHeuristic{}

	for pos, node := range posMap {
		// Add the line of sight distance to the end node to the heuristic.
		err := heuristic.AddNode(node, dist2D(pos, dest))
		if err != nil {
			return nil, err
		}
	}

	return heuristic.Heuristic(0), nil
}
