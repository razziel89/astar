# General

This is an implementation of the `A*` path finding algoritm written in plain
Golang.
Please see [wikipedia](https://en.wikipedia.org/wiki/A*_search_algorithm) for a
description of the algorithm.

# State

This software is in alpha state.
Development started in Dec 2021.

# Overview

The algorithm has been developed during the amazing [Advent of
Code](https://adventofcode.com) evenent in 2021.
Test coverage is currently limited but will be expanded.

# How to use

The main function most users will want to use is `FindPath` in
[`astar.go`](./astar.go).
It requires a graph of nodes as input.
Assuming you want to create a regular grid of 100 nodes, you can create one like
this and find a path from the top left to the bottom right:

```go
package main

import (
	"fmt"
	"log"

	"github.com/razziel89/astar"
)

func main() {
	// 100 nodes means a 10x10 grid.
	gridSize := 10
	endX, endY := 9, 9
	// We disallow diagonal movements. Thus, each node has 4
	// neighbours, 2 in x direction and 2 in y direction. This value
	// just improves performances during allocation.
	neighbours := 4
	// This is the graph we will be using to find the path.
	graph := astar.Graph{}
	// We remember the node for each position. This makes creating
	// connections easier.
	posToNode := map[[2]int]*astar.Node{}
	// We also need to create a heuristic that estimates the distance
	// to the end node. We use a simple one based on the line of sight
	// distance.
	heuristic := astar.ConstantHeuristic{}

	// Create the actual nodes, connections will be created later.
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {

			// Create a name for the node, the algorithm does not use
			// this value.
			nodeName := fmt.Sprintf("x:%d,y:%d", x, y)
			// Note that we increase cost to larger values for x and
			// y, but twice as much for y for this particular example.
			nodeCost := x + 2*y
			// Nodes can take an arbitrary payload. You can attach the
			// position, for example. We don't use any payload for
			// this example. The default for an empty interface is
			// nil.
			var payload interface{}
			// Create the node.
			node, err := astar.NewNode(
				nodeName, nodeCost, neighbours, payload,
			)
			// Error handling. Costs must always be positive, for
			// example.
			if err != nil {
				log.Fatal(err.Error())
			}
			// Remember the node for this position.
			posToNode[[2]int{x, y}] = node
			// Add the node to the graph!
			graph.Add(node)
			// Add the line of sight distance in Manhattan metric to
			// the end node to the heuristic. Since we disallow
			// diagonal connections, this is realistic.
			err = heuristic.AddNode(node, (endX-x)+(endY-y))
			// Error handling. Estimates must always be positive, for
			// example. Furthermore, you cannot overwrite a node's
			// estimate after adding it.
			if err != nil {
				log.Fatal(err.Error())
			}

		}
	}

	// Create connections.
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {

			// Add all connections for this node.
			// Extract the node first.
			node := posToNode[[2]int{x, y}]
			// Then, iterate over all neighbours in x and y
			// directions. There are four different displacements
			// here.
			neighDisp := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
			for _, disp := range neighDisp {
				neighPos := [2]int{x + disp[0], y + disp[1]}
				// Only add connections if the neighbour actually
				// exists!
				if neigh, exists := posToNode[neighPos]; exists {
					// Add pairwise connections. Adding a connection
					// multiple times is no problem. All but the first
					// calls are no-ops.
					node.AddConnection(neigh)
					neigh.AddConnection(node)
				}
			}

		}
	}

	// Extract start and end nodes.
	start := posToNode[[2]int{0, 0}]
	end := posToNode[[2]int{endX, endY}]

	// Find the path! As you can see, creating the data structures is
	// the biggest headache. But this algorithm permits arbitrary
	// connections, even one-directional ones.
	path, err := astar.FindPath(
		graph, start, end, heuristic.Heuristic(0),
	)
	// Error handling.
	if err != nil {
		log.Fatal(err.Error())
	}

	// Output the path we found!
	for _, node := range path {
		fmt.Println(node.ToString())
	}
}
```

The output will be this:

```
{id: x:0,y:0, cost: 0, con: ['x:0,y:1', 'x:1,y:0']}
{id: x:1,y:0, cost: 1, con: ['x:1,y:1', 'x:2,y:0', 'x:0,y:0']}
{id: x:2,y:0, cost: 2, con: ['x:1,y:0', 'x:2,y:1', 'x:3,y:0']}
{id: x:3,y:0, cost: 3, con: ['x:2,y:0', 'x:3,y:1', 'x:4,y:0']}
{id: x:4,y:0, cost: 4, con: ['x:3,y:0', 'x:4,y:1', 'x:5,y:0']}
{id: x:5,y:0, cost: 5, con: ['x:4,y:0', 'x:5,y:1', 'x:6,y:0']}
{id: x:6,y:0, cost: 6, con: ['x:5,y:0', 'x:6,y:1', 'x:7,y:0']}
{id: x:7,y:0, cost: 7, con: ['x:6,y:0', 'x:7,y:1', 'x:8,y:0']}
{id: x:8,y:0, cost: 8, con: ['x:7,y:0', 'x:8,y:1', 'x:9,y:0']}
{id: x:9,y:0, cost: 9, con: ['x:8,y:0', 'x:9,y:1']}
{id: x:9,y:1, cost: 11, con: ['x:8,y:1', 'x:9,y:0', 'x:9,y:2']}
{id: x:9,y:2, cost: 13, con: ['x:8,y:2', 'x:9,y:1', 'x:9,y:3']}
{id: x:9,y:3, cost: 15, con: ['x:8,y:3', 'x:9,y:2', 'x:9,y:4']}
{id: x:9,y:4, cost: 17, con: ['x:8,y:4', 'x:9,y:3', 'x:9,y:5']}
{id: x:9,y:5, cost: 19, con: ['x:8,y:5', 'x:9,y:4', 'x:9,y:6']}
{id: x:9,y:6, cost: 21, con: ['x:8,y:6', 'x:9,y:5', 'x:9,y:7']}
{id: x:9,y:7, cost: 23, con: ['x:8,y:7', 'x:9,y:6', 'x:9,y:8']}
{id: x:9,y:8, cost: 25, con: ['x:9,y:9', 'x:8,y:8', 'x:9,y:7']}
{id: x:9,y:9, cost: 27, con: ['x:8,y:9', 'x:9,y:8']}
```

As you can see, you get a nice string representation of the path.

# Installation

Simply add `github.com/razziel89/astar` as a dependency to your project by
running
```bash
go get github.com/razziel89/astar@latest
```

# How to contribute

If you have found a bug and want to fix it, please simply go ahead and fork the
repository, fix the bug, and open a pull request to this repository!
Bug fixes are always welcome.

In all other cases, please open an issue on GitHub first to discuss the
contribution.
The feature you would like to introduce might already be in development.

# Licence

[GPLv3](./LICENCE)

If you want to use this piece of software under a different, more permissive
open-source licence, please contact me.
I am very open to discussing this point.
