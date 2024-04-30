package main

import (
	"testing"

	"github.com/m-colson/dystopia/shared/graph"
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
	g := graph.New([]graph.Link{
		{From: 1, Cost: 5, To: 2},
		{From: 1, Cost: 1, To: 3},
		{From: 2, Cost: 1, To: 4},
		{From: 2, Cost: 1, To: 5},
		{From: 3, Cost: 1, To: 5},
		{From: 3, Cost: 2, To: 7},
		{From: 4, Cost: 1, To: 6},
		{From: 5, Cost: 2, To: 6},
		{From: 6, Cost: 1, To: 8},
		{From: 7, Cost: 6, To: 8},
	}...)

	out := Dijkstra(g, 1, 8)
	t.Log(out)

	if len(out) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(out))
	}

	exp := []graph.NodeID{8, 6, 5, 3}
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
	g := graph.New([]graph.Link{
		{From: 1, Cost: 2, To: 2},
		{From: 1, Cost: 2, To: 3},
		{From: 2, Cost: 1, To: 4},
		{From: 2, Cost: 3, To: 5},
		{From: 3, Cost: 1, To: 5},
		{From: 3, Cost: 2, To: 7},
		{From: 4, Cost: 1, To: 6},
		{From: 5, Cost: 3, To: 6},
		{From: 6, Cost: 1, To: 8},
		{From: 7, Cost: 6, To: 8},
	}...)

	out := Dijkstra(g, 1, 8)
	t.Log(out)

	if len(out) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(out))
	}

	exp := []graph.NodeID{8, 6, 4, 2}
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
	g := graph.New([]graph.Link{
		{From: 1, Cost: 2, To: 2},
		{From: 1, Cost: 2, To: 3},
		{From: 2, Cost: 1, To: 4},
		{From: 2, Cost: 3, To: 5},
		{From: 3, Cost: 1, To: 5},
		{From: 3, Cost: 1, To: 7},
		{From: 4, Cost: 1, To: 6},
		{From: 5, Cost: 3, To: 6},
		{From: 6, Cost: 1, To: 8},
		{From: 7, Cost: 1, To: 8},
	}...)

	out := Dijkstra(g, 1, 8)
	t.Log(out)

	if len(out) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(out))
	}

	exp := []graph.NodeID{8, 7, 3}
	for i, node := range out {
		if node != exp[i] {
			t.Fatalf("expected node %d to be %d, got %d", i, exp[i], node)
		}
	}
}
