package main

import "github.com/m-colson/dystopia/shared/graph"

func SimulateOnce(cars *CarsMap, g graph.Graph) {
	cars.Lock.Lock()
	defer cars.Lock.Unlock()
	for _, car := range cars.Cars {
		if car.Trip.To == SENTINAL_NODE {
			continue
		}
	}
}
