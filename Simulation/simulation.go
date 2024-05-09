package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/m-colson/dystopia/shared/graph"
)

const SIMULATION_RATE = 0.1

// NOTE: cars must be locked
func SimulateOnce(cars *graph.CarsMap, g graph.Graph) (flaggedCars map[graph.CarID]struct{}) {
	flaggedCars = make(map[graph.CarID]struct{})

	for _, car := range cars.Cars {
		if car.Trip.To == graph.SENTINAL_NODE {
			continue
		}

		if car.Pos.To == graph.SENTINAL_NODE {
			if len(car.CurrPath) == 0 {
				// go to start
				dest := car.Trip.From
				if car.Pos.From == car.Trip.From {
					// go to end
					dest = car.Trip.To
				} else if car.Pos.From == car.Trip.To {
					// done
					flaggedCars[car.ID] = struct{}{}
					continue
				}

				path, err := pfaasPath(car.Pos.From, dest)
				if err != nil || len(path) == 0 {
					log.Printf("failed to find path from %d to %d because: %s", car.Pos.From, dest, err)
					continue
				}

				car.CurrPath = path
			}

			nextNode := car.CurrPath[len(car.CurrPath)-1]
			car.Pos.To = nextNode.To
			car.Pos.Cost = nextNode.Cost
			car.CurrPath = car.CurrPath[:len(car.CurrPath)-1]
		}

		if car.Pos.Ratio >= 1 {
			car.Pos = graph.CarPos{
				From:  car.Pos.To,
				To:    graph.SENTINAL_NODE,
				Ratio: 0.0,
				// Cost: set in next iteration
			}
		} else {
			car.Pos.Ratio += SIMULATION_RATE / float64(car.Pos.Cost)
		}
	}

	return
}

func pfaasPath(from graph.NodeID, to graph.NodeID) ([]graph.Edge, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(
		"%s/api/path?from=%d&to=%d",
		PFAAS_HOST,
		from,
		to,
	), nil)

	if err != nil {
		return nil, err
	}

	resp, err := (&http.Client{}).Do(req)

	if err != nil {
		return nil, err
	}

	type message struct {
		Path []graph.Edge `json:"path"`
	}

	out := message{}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out.Path, nil
}
