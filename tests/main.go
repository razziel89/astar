// Package main describes one large tests.

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
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

var quiet = os.Getenv("QUIET") == "1"

func logStr(str string) {
	if !quiet {
		log.Print(str)
	}
}

func main() {
	// Generate test data.
	rand.Seed(seed)
	endX, endY := gridSize-1, gridSize-1
	nodes := make([]*astar.Node, 0, gridSize*gridSize)
	heuristic := astar.ConstantHeuristic{}

	logStr("inititalising")

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

	logStr("created nodes")

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

	logStr("connected nodes")

	// Convert to graph.
	var graph astar.GraphOps
	switch os.Getenv("GRAPH_TYPE") {
	case "HEAPED":
		graph = astar.NewHeapedGraph(gridSize * gridSize)
	case "MAPPED":
		graph = astar.NewGraph()
	default:
		log.Fatal("Env var GRAPH_TYPE not set to a supported value.")
	}

	for _, node := range nodes {
		graph.Add(node)
	}

	logStr("converted to graph")

	startTime := time.Now()
	// Run the test.
	path, err := astar.FindPath(graph, start, end, heuristic.Heuristic(0))
	if err != nil {
		log.Fatal(err.Error())
	}
	duration := time.Since(startTime)

	logStr("path is")

	cost := 0
	for _, node := range path {
		logStr(node.ToString())
		cost += node.Cost
	}

	logStr(fmt.Sprintf("total cost is %d", cost))

	if len(path) != expectedLength {
		log.Fatalf(
			"path does not have the expected length (want: %d, has: %d)",
			expectedLength, len(path),
		)
	}

	if cost != expectedCost {
		log.Fatalf("path does not have the expected cost (want: %d, has: %d)", expectedCost, cost)
	}

	logStr("obtained path")

	fmt.Println(duration)
}
