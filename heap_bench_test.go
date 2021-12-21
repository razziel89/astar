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

import "testing"

const (
	tenK = 1000
)

func BenchmarkGraphAdd10KNodes(b *testing.B) {
	graph := Graph{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", 0, 0, nil)
			graph.Add(node)
		}
	}
}

func BenchmarkHeapPush10KNodes(b *testing.B) {
	heap := Heap{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", 0, 0, nil)
			heap.Push(node)
		}
	}
}

func BenchmarkHeapAdd10KNodes(b *testing.B) {
	heap := Heap{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", 0, 0, nil)
			heap.Add(node)
		}
	}
}

func BenchmarkGraphPopCheapest10KNodes(b *testing.B) {
	graph := Graph{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", 0, 0, nil)
			graph.Add(node)
		}
		for j := 0; j < tenK; j++ {
			_ = graph.PopCheapest(nil)
		}
	}
}

func BenchmarkHeapPop10KNodes(b *testing.B) {
	heap := Heap{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", 0, 0, nil)
			heap.Push(node)
		}
		for j := 0; j < tenK; j++ {
			_ = heap.Pop()
		}
	}
}

func BenchmarkHeapPopCheapest10KNodes(b *testing.B) {
	heap := Heap{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", 0, 0, nil)
			heap.Push(node)
		}
		for j := 0; j < tenK; j++ {
			_ = heap.PopCheapest()
		}
	}
}
