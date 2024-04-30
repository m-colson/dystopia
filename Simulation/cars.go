package main

import (
	"strconv"
	"sync"
)

type CarsKey struct{}
type CarsMap struct {
	Cars map[CarID]*Car
	Lock sync.Mutex
}

type Car struct {
	ID   CarID
	Trip Trip
	Pos  CarPos
}

type CarID uint64

const SENTINAL_CAR = CarID(1<<63 - 1)

func ParseCarID(s string) (out CarID, err error) {
	idNum, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return SENTINAL_CAR, err
	}
	return CarID(idNum), nil
}

type Trip struct {
	From NodeID
	To   NodeID
}

type NodeID uint64

const SENTINAL_NODE = NodeID(1<<63 - 1)

func ParseNodeID(s string) (out NodeID, err error) {
	idNum, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return SENTINAL_NODE, err
	}
	return NodeID(idNum), nil
}

type CarPos struct {
	From  NodeID  `json:"from"`
	To    NodeID  `json:"to"`
	Ratio float64 `json:"ratio"`
}
