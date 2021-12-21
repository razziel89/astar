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
	goheap "container/heap"
	"testing"
)

const (
	tenK     = 1000
	hundredK = 10000
)

func BenchmarkGraphAdd10KNodes(b *testing.B) {
	graph := Graph{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", tenK-j, 0, nil)
			graph.Add(node)
		}
	}
}

func BenchmarkGraphAdd100KNodes(b *testing.B) {
	graph := Graph{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < hundredK; j++ {
			node, _ := NewNode("", hundredK-j, 0, nil)
			graph.Add(node)
		}
	}
}

func BenchmarkHeapPush10KNodes(b *testing.B) {
	heap := Heap{}
	goheap.Init(&heap)
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", tenK-j, 0, nil)
			goheap.Push(&heap, node)
		}
	}
}

func BenchmarkHeapPush100KNodes(b *testing.B) {
	heap := Heap{}
	goheap.Init(&heap)
	for i := 0; i < b.N; i++ {
		for j := 0; j < hundredK; j++ {
			node, _ := NewNode("", hundredK-j, 0, nil)
			goheap.Push(&heap, node)
		}
	}
}

func BenchmarkGraphPopCheapest10KNodes(b *testing.B) {
	graph := Graph{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", tenK-j, 0, nil)
			graph.Add(node)
		}
		for j := 0; j < tenK; j++ {
			_ = graph.PopCheapest()
		}
	}
}

func BenchmarkGraphPopCheapest100KNodes(b *testing.B) {
	graph := Graph{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < hundredK; j++ {
			node, _ := NewNode("", hundredK-j, 0, nil)
			graph.Add(node)
		}
		for j := 0; j < hundredK; j++ {
			_ = graph.PopCheapest()
		}
	}
}

func BenchmarkHeapPop10KNodes(b *testing.B) {
	heap := Heap{}
	goheap.Init(&heap)
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", 0, 0, nil)
			goheap.Push(&heap, node)
		}
		for j := 0; j < tenK; j++ {
			_ = goheap.Pop(&heap)
		}
	}
}

func BenchmarkHeapPop100KNodes(b *testing.B) {
	heap := Heap{}
	goheap.Init(&heap)
	for i := 0; i < b.N; i++ {
		for j := 0; j < hundredK; j++ {
			node, _ := NewNode("", 0, 0, nil)
			goheap.Push(&heap, node)
		}
		for j := 0; j < hundredK; j++ {
			_ = goheap.Pop(&heap)
		}
	}
}

func BenchmarkHeapPopCheapest10KNodes(b *testing.B) {
	heap := Heap{}
	goheap.Init(&heap)
	for i := 0; i < b.N; i++ {
		for j := 0; j < tenK; j++ {
			node, _ := NewNode("", 0, 0, nil)
			goheap.Push(&heap, node)
		}
		for j := 0; j < tenK; j++ {
			_ = goheap.Pop(&heap)
		}
	}
}

func BenchmarkHeapPopCheapest100KNodes(b *testing.B) {
	heap := Heap{}
	goheap.Init(&heap)
	for i := 0; i < b.N; i++ {
		for j := 0; j < hundredK; j++ {
			node, _ := NewNode("", 0, 0, nil)
			goheap.Push(&heap, node)
		}
		for j := 0; j < hundredK; j++ {
			_ = goheap.Pop(&heap)
		}
	}
}
