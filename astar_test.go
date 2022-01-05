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
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockPath []*Node
var mockGraph Graph
var mockStart *Node
var mockEnd *Node
var errMock = fmt.Errorf("some error")

// Set up test cases for find path. The two functions used by it will return the provided error
// values.
func setUpFindPath(errFindReverse, errExtract, errCleanUp error, connect bool) func() {

	node1, _ := NewNode("start", 0, 0, nil)
	node2, _ := NewNode("end", 0, 0, nil)

	node1.AddPairwiseConnection(node2)

	mockPath = []*Node{node1, node2}

	mockStart = node1
	mockEnd = node2

	mockGraph = Graph{}
	mockGraph.Add(mockStart)
	mockGraph.Add(mockEnd)

	// Ensure there is an actual connection if desired.
	if connect {
		mockEnd.prev = mockStart
	}

	extractPath = func(_, _ *Node, _ bool) ([]*Node, error) {
		return mockPath, errExtract
	}

	findReversePath = func(_, _ GraphOps, _ *Node, _ Heuristic) error {
		return errFindReverse
	}

	resetFnGetter = func(_ GraphOps) func(_ *Node) error {
		return func(_ *Node) error {
			return errCleanUp
		}
	}

	return func() {
		// Revert changes.
		extractPath = ExtractPath
		findReversePath = FindReversePath

		mockPath = []*Node{}
		mockGraph = Graph{}
		mockStart = nil
		mockEnd = nil
	}
}

func TestFindPathSuccess(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, nil, true)
	defer tearDown()

	path, err := FindPath(&mockGraph, mockStart, mockEnd, mockHeuristic)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(path))
	assert.Equal(t, path[0], mockStart)
	assert.Equal(t, path[1], mockEnd)
}

func TestFindPathFailurePathExtraction(t *testing.T) {
	tearDown := setUpFindPath(errMock, nil, nil, true)
	defer tearDown()

	_, err := FindPath(&mockGraph, mockStart, mockEnd, mockHeuristic)
	assert.Error(t, err)
}

func TestFindPathFailurePathFinding(t *testing.T) {
	tearDown := setUpFindPath(nil, errMock, nil, true)
	defer tearDown()

	_, err := FindPath(&mockGraph, mockStart, mockEnd, mockHeuristic)
	assert.Error(t, err)
}

func TestFindPathFailureNoEnd(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, nil, true)
	defer tearDown()

	_, err := FindPath(&mockGraph, mockStart, nil, mockHeuristic)
	assert.Error(t, err)
}

func TestFindPathFailureNoStart(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, nil, true)
	defer tearDown()

	_, err := FindPath(&mockGraph, nil, mockEnd, mockHeuristic)
	assert.Error(t, err)
}

func TestFindPathFailureAlreadyConnection(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, nil, true)
	defer tearDown()

	findReversePath = FindReversePath

	_, err := FindPath(&mockGraph, mockStart, mockEnd, mockHeuristic)
	assert.Error(t, err)
}

func TestFindPathFailureNoConnectionToEnd(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, nil, false)
	defer tearDown()

	mockEnd.RemoveConnection(mockStart)
	mockStart.RemoveConnection(mockEnd)

	findReversePath = FindReversePath

	_, err := FindPath(&mockGraph, mockStart, mockEnd, mockHeuristic)
	assert.Error(t, err)
}

func TestExtractPathSuccessNoReverse(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, nil, true)
	defer tearDown()

	path, err := ExtractPath(mockEnd, mockStart, false)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(path))
	assert.Equal(t, mockEnd, path[0])
	assert.Equal(t, mockStart, path[1])
}

func TestExtractPathSuccessReverse(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, nil, true)
	defer tearDown()

	path, err := ExtractPath(mockEnd, mockStart, true)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(path))
	assert.Equal(t, mockStart, path[0])
	assert.Equal(t, mockEnd, path[1])
}

func TestExtractPathFailureNoConnection(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, nil, false)
	defer tearDown()

	_, err := ExtractPath(mockEnd, mockStart, true)
	assert.Error(t, err)
}

func TestFindPathBetterConnection(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, nil, true)
	defer tearDown()

	mockMid, _ := NewNode("mid", 0, 0, nil)

	mockMid.AddConnection(mockStart)
	mockStart.AddConnection(mockMid)

	mockMid.AddConnection(mockEnd)
	mockEnd.AddConnection(mockMid)

	orgCost := 1000
	mockMid.trackedCost = orgCost
	mockEnd.prev = nil

	err := FindReversePath(
		&Graph{mockStart: graphVal, mockMid: graphVal},
		&Graph{}, mockEnd, mockHeuristic,
	)

	assert.NoError(t, err)
	assert.NotEqual(t, orgCost, mockMid.trackedCost)
}

func TestFindPathFailureNodeReset(t *testing.T) {
	tearDown := setUpFindPath(nil, nil, errMock, true)
	defer tearDown()

	_, err := FindPath(&mockGraph, mockStart, mockEnd, mockHeuristic)
	assert.Error(t, err)
}

func TestPanicHandlerNoPanic(t *testing.T) {
	callMe := func() (err error) {
		defer getPanicHandler(&err)()
		return nil
	}

	err := callMe()
	assert.NoError(t, err)
}

func TestPanicHandlerPanicInAstarNoOtherError(t *testing.T) {
	callMe := func() (err error) {
		defer getPanicHandler(&err)()
		panic(Error{"astar panic"})
	}

	err := callMe()
	assert.Error(t, err)
	assert.NotEmpty(t, err.Error())
}

func TestPanicHandlerPanicInAstarOtherError(t *testing.T) {
	callMe := func() (err error) {
		defer getPanicHandler(&err)()
		err = fmt.Errorf("some other error")
		panic(Error{"astar panic"})
	}

	err := callMe()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "some other error")
	assert.NotContains(t, err.Error(), "astar panic")
}

func TestPanicHandlerPanicOutsideAstar(t *testing.T) {
	callMe := func() (err error) {
		defer getPanicHandler(&err)()
		panic("panic outside astar signified via a different error type")
	}

	defer func() {
		errStr, wasError := recover().(string)
		assert.True(t, wasError)
		assert.NotEmpty(t, errStr)
	}()

	_ = callMe()
}
