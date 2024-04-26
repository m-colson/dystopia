package main

import (
	"testing"
)

func TestPathfind1(t *testing.T) {
	t.Parallel()

	//       1
	//      / \
	//      2  3
	//     /| / \
	//    4 5   7
	//    \ |   |
	//     6 - 8
	graph := NewGraph([]Link{
		{1, 5, 2},
		{1, 1, 3},
		{2, 1, 4},
		{2, 1, 5},
		{3, 1, 5},
		{3, 2, 7},
		{4, 1, 6},
		{5, 2, 6},
		{6, 1, 8},
		{7, 6, 8},
	}...)

	out := Dijkstra(graph, 1, 8)
	t.Log(out)

	if len(out) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(out))
	}

	exp := []NodeID{8, 6, 5, 3}
	for i, node := range out {
		if node != exp[i] {
			t.Fatalf("expected node %d to be %d, got %d", i, exp[i], node)
		}
	}
}

func TestPathfind2(t *testing.T) {
	t.Parallel()

	//       1
	//      / \
	//      2  3
	//     /| / \
	//    4 5   7
	//    \ |   |
	//     6 - 8
	graph := NewGraph([]Link{
		{1, 2, 2},
		{1, 2, 3},
		{2, 1, 4},
		{2, 3, 5},
		{3, 1, 5},
		{3, 2, 7},
		{4, 1, 6},
		{5, 3, 6},
		{6, 1, 8},
		{7, 6, 8},
	}...)

	out := Dijkstra(graph, 1, 8)
	t.Log(out)

	if len(out) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(out))
	}

	exp := []NodeID{8, 6, 4, 2}
	for i, node := range out {
		if node != exp[i] {
			t.Fatalf("expected node %d to be %d, got %d", i, exp[i], node)
		}
	}
}

func TestPathfind3(t *testing.T) {
	t.Parallel()

	//       1
	//      / \
	//      2  3
	//     /| / \
	//    4 5   7
	//    \ |   |
	//     6 - 8
	graph := NewGraph([]Link{
		{1, 2, 2},
		{1, 2, 3},
		{2, 1, 4},
		{2, 3, 5},
		{3, 1, 5},
		{3, 1, 7},
		{4, 1, 6},
		{5, 3, 6},
		{6, 1, 8},
		{7, 1, 8},
	}...)

	out := Dijkstra(graph, 1, 8)
	t.Log(out)

	if len(out) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(out))
	}

	exp := []NodeID{8, 7, 3}
	for i, node := range out {
		if node != exp[i] {
			t.Fatalf("expected node %d to be %d, got %d", i, exp[i], node)
		}
	}
}
