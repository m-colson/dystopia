package main

import (
	"container/heap"
	"fmt"
)

func Dijkstra(g Graph, source NodeID, dest NodeID) []NodeID {
	prev := make(map[NodeID]NodeID)

	heapLut := make(map[NodeID]*DataEntry)
	bestHeap := &DataHeap{{Id: source, Cost: 0}}
	heap.Init(bestHeap)

	for bestHeap.Len() > 0 {
		best := heap.Pop(bestHeap).(*DataEntry)
		fmt.Println(best)
		if best.Id == dest {
			break
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

	if _, ok := prev[dest]; !ok {
		return nil
	}

	out := make([]NodeID, 0)

	temp := dest
	for temp != source {
		out = append(out, temp)
		temp = prev[temp]
	}

	return out
}
