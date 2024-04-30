package main

import "strconv"

type NodeID uint64

const SENTINAL_NODE = NodeID(1<<63 - 1)

func ParseNodeID(s string) (out NodeID, err error) {
	idNum, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return SENTINAL_NODE, err
	}
	return NodeID(idNum), nil
}

type Graph struct {
	Nodes map[NodeID]Node
}

type Link struct {
	From NodeID
	Cost Cost
	To   NodeID
}

func NewGraph(links ...Link) Graph {
	nodes := make(map[NodeID]Node)

	for _, link := range links {
		n, ok := nodes[link.From]
		if !ok {
			n = Node{make([]Edge, 0)}
		}
		n.Next = append(n.Next, Edge{To: link.To, Cost: link.Cost})
		nodes[link.From] = n

		if _, ok := nodes[link.To]; !ok {
			nodes[link.To] = Node{make([]Edge, 0)}
		}
	}

	return Graph{Nodes: nodes}
}

func NewGraphRaw(start ...Node) Graph {
	nodes := make(map[NodeID]Node)

	idCount := 1
	for _, node := range start {
		nodes[NodeID(idCount)] = node
		idCount++
	}

	return Graph{Nodes: nodes}
}

type Node struct {
	Next []Edge
}

func NewNode(edges ...Edge) Node {
	return Node{Next: edges}
}

type Edge struct {
	To   NodeID
	Cost Cost
}

type Cost uint64

type DataEntry struct {
	Id    NodeID
	Cost  Cost
	index int
}

func (h *DataEntry) Done() bool {
	return h.index == -1
}

type DataHeap []*DataEntry

func (h DataHeap) Len() int           { return len(h) }
func (h DataHeap) Less(i, j int) bool { return h[i].Cost < h[j].Cost }
func (h DataHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *DataHeap) Push(x any) {
	n := len(*h)
	entry := x.(*DataEntry)
	entry.index = n
	*h = append(*h, entry)
}

func (h *DataHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = &DataEntry{}
	item.index = -1
	*h = old[0 : n-1]
	return item
}
