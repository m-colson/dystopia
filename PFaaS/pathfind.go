package main

import (
	"container/heap"
	"fmt"

	"github.com/m-colson/dystopia/shared/graph"
)

func dijkstraAlgorithm(
	prev map[graph.NodeID]graph.Edge,
	heapLut map[graph.NodeID]*graph.DataEntry,
	bestHeap *graph.DataHeap,
	g *graph.Graph,
	dests []graph.NodeID,
) {
	for bestHeap.Len() > 0 {
		best := heap.Pop(bestHeap).(*graph.DataEntry)
		for _, dest := range dests {
			if best.Id == dest {
				return
			}
		}

		for _, edge := range g.Nodes[best.Id].Next {
			alt := best.Cost + edge.Cost

			edgeEntry, alreadyExists := heapLut[edge.To]
			if alreadyExists && (edgeEntry.Done() || alt >= edgeEntry.Cost) {
				continue
			}

			if !alreadyExists {
				edgeEntry = &graph.DataEntry{Id: edge.To, Cost: 1<<63 - 1}
				heap.Push(bestHeap, edgeEntry)
				heapLut[edge.To] = edgeEntry
			}

			prev[edge.To] = graph.Edge{To: best.Id, Cost: edge.Cost}
			edgeEntry.Cost = alt
			heap.Fix(bestHeap, edgeEntry.Index)
		}
	}
}

func Dijkstra(g graph.Graph, source graph.NodeID, dests ...graph.NodeID) []graph.Edge {
	prev := make(map[graph.NodeID]graph.Edge)

	heapLut := make(map[graph.NodeID]*graph.DataEntry)
	bestHeap := &graph.DataHeap{{Id: source, Cost: 0}}
	heap.Init(bestHeap)

	dijkstraAlgorithm(prev, heapLut, bestHeap, &g, dests)

	closestDest := graph.SENTINAL_NODE
	for _, dest := range dests {
		if _, ok := prev[dest]; ok {
			closestDest = dest
			break
		}
	}

	if closestDest == graph.SENTINAL_NODE {
		return []graph.Edge{}
	}

	out := make([]graph.Edge, 0)

	temp := closestDest
	for temp != source {
		out = append(out, graph.Edge{To: temp, Cost: prev[temp].Cost})
		temp = prev[temp].To
	}

	return out
}

func DijkstraClosest(g graph.Graph, source graph.NodeID, dests ...graph.NodeID) (closest graph.NodeID, ok bool) {
	fmt.Println(source, dests)

	prev := make(map[graph.NodeID]graph.Edge)

	heapLut := make(map[graph.NodeID]*graph.DataEntry)
	bestHeap := &graph.DataHeap{{Id: source, Cost: 0}}
	heap.Init(bestHeap)

	dijkstraAlgorithm(prev, heapLut, bestHeap, &g, dests)

	fmt.Println(prev)

	closest = graph.SENTINAL_NODE
	for _, dest := range dests {
		if _, ok := prev[dest]; ok {
			closest = dest
			break
		}
	}

	return closest, closest != graph.SENTINAL_NODE
}
