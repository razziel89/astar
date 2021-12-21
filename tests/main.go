// Package main describes one large tests.

package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/razziel89/astar"
)

const (
	seed           = 42
	maxRand        = 10
	gridSize       = 200
	numNeigh       = 4
	expectedLength = 417
	expectedCost   = 827
)

func main() {
	// Generate test data.
	rand.Seed(seed)
	endX, endY := gridSize-1, gridSize-1
	nodes := make([]*astar.Node, 0, gridSize*gridSize)
	heuristic := astar.ConstantHeuristic{}

	log.Println("inititalising")

	// Create nodes and compute constant heuristic.
	for xIdx := 0; xIdx < gridSize; xIdx++ {
		for yIdx := 0; yIdx < gridSize; yIdx++ {
			node, err := astar.NewNode(
				fmt.Sprintf("{%d,%d}", xIdx, yIdx), rand.Intn(maxRand), numNeigh, nil,
			)
			if err != nil {
				log.Fatal(err.Error())
			}
			nodes = append(nodes, node)
			heuristic.AddNode(node, (endX-xIdx)+(endY-yIdx))
		}
	}

	log.Println("created nodes")

	start := nodes[0]
	end := nodes[len(nodes)-1]

	// Add connections.
	for idx, node := range nodes {
		if idx >= gridSize {
			up := nodes[idx-gridSize]
			node.AddPairwiseConnection(up)
		}
		if idx+gridSize < len(nodes) {
			down := nodes[idx+gridSize]
			node.AddPairwiseConnection(down)
		}
		if idx%gridSize > 0 {
			left := nodes[idx-1]
			node.AddPairwiseConnection(left)
		}
		if idx%gridSize < gridSize-1 {
			right := nodes[idx+1]
			node.AddPairwiseConnection(right)
		}
	}

	log.Println("connected nodes")

	// Convert to graph.
	graph := astar.Graph{}
	for _, node := range nodes {
		graph.Add(node)
	}

	log.Println("converted to graph")

	startTime := time.Now()
	// Run the test.
	path, err := astar.FindPath(&graph, start, end, heuristic.Heuristic(0))
	if err != nil {
		log.Fatal(err.Error())
	}
	duration := time.Since(startTime)

	log.Println("path is")

	cost := 0
	for _, node := range path {
		log.Println(node.ToString())
		cost += node.Cost
	}

	if len(path) != expectedLength {
		log.Fatal(
			"path does not have the expected length (want: %d, has: %d",
			expectedLength, len(path),
		)
	}

	if cost != expectedCost {
		log.Fatalf("path does not have the expected cost (want: %d, has: %d)", expectedCost, cost)
	}

	log.Println("obtained path")

	fmt.Println(duration)
}
