package main

import (
	"fmt"
	"log"
	"testing"
)

func TestPathfind(t *testing.T) {
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
	fmt.Println(out)

	log.Fatal("test")
}
