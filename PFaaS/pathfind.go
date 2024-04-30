package main

import (
	"container/heap"
)

func dijkstraAlgorithm(prev map[NodeID]NodeID, heapLut map[NodeID]*DataEntry, bestHeap *DataHeap, g *Graph, dests []NodeID) {
	for bestHeap.Len() > 0 {
		best := heap.Pop(bestHeap).(*DataEntry)
		for _, dest := range dests {
			if best.Id == dest {
				break
			}
		}

		for _, edge := range g.Nodes[best.Id].Next {
			alt := best.Cost + edge.Cost

			edgeEntry, alreadyExists := heapLut[edge.To]
			if alreadyExists && (edgeEntry.Done() || alt >= edgeEntry.Cost) {
				continue
			}

			if !alreadyExists {
				edgeEntry = &DataEntry{Id: edge.To, Cost: 1<<63 - 1}
				heap.Push(bestHeap, edgeEntry)
				heapLut[edge.To] = edgeEntry
			}

			prev[edge.To] = best.Id
			edgeEntry.Cost = alt
			heap.Fix(bestHeap, edgeEntry.index)
		}
	}
}

func Dijkstra(g Graph, source NodeID, dests ...NodeID) []NodeID {
	prev := make(map[NodeID]NodeID)

	heapLut := make(map[NodeID]*DataEntry)
	bestHeap := &DataHeap{{Id: source, Cost: 0}}
	heap.Init(bestHeap)

	dijkstraAlgorithm(prev, heapLut, bestHeap, &g, dests)

	closestDest := SENTINAL_NODE
	for _, dest := range dests {
		if _, ok := prev[dest]; ok {
			closestDest = dest
			break
		}
	}

	if closestDest == SENTINAL_NODE {
		return []NodeID{}
	}

	out := make([]NodeID, 0)

	temp := closestDest
	for temp != source {
		out = append(out, temp)
		temp = prev[temp]
	}

	return out
}

func DijkstraClosest(g Graph, source NodeID, dests ...NodeID) (closest NodeID, ok bool) {
	prev := make(map[NodeID]NodeID)

	heapLut := make(map[NodeID]*DataEntry)
	bestHeap := &DataHeap{{Id: source, Cost: 0}}
	heap.Init(bestHeap)

	dijkstraAlgorithm(prev, heapLut, bestHeap, &g, dests)

	closest = SENTINAL_NODE
	for _, dest := range dests {
		if _, ok := prev[dest]; ok {
			closest = dest
			break
		}
	}

	return closest, closest != SENTINAL_NODE
}
