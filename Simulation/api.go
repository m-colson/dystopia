package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/m-colson/dystopia/shared/graph"
	"github.com/m-colson/psi"
)

type MissingQueryParam struct {
	name string
}

func (e *MissingQueryParam) StatusCode() int {
	return http.StatusBadRequest
}

func (e *MissingQueryParam) Error() string {
	return fmt.Sprintf("missing URL parameter '%s'", e.name)
}

type IllegalQueryParam struct {
	name  string
	value string
	inner error
}

func (e *IllegalQueryParam) StatusCode() int {
	return http.StatusBadRequest
}

func (e *IllegalQueryParam) Error() string {
	return fmt.Sprintf("illegal value '%s' for query parameter '%s': %s", e.value, e.name, e.inner)
}

func queryParam(r *http.Request, name string) (string, psi.StatusCode) {
	value := r.URL.Query()[name]
	if len(value) == 0 {
		return "", &MissingQueryParam{name: name}
	}
	return value[0], nil
}

func parseNodeIDParam(r *http.Request, name string) (graph.NodeID, psi.StatusCode) {
	value, queryErr := queryParam(r, name)
	if queryErr != nil {
		return graph.SENTINAL_NODE, queryErr
	}
	id, err := graph.ParseID(value)
	if err != nil {
		return graph.SENTINAL_NODE, &IllegalQueryParam{name: name, value: value, inner: err}
	}
	return id, nil
}

type CarNotFound struct {
	ID graph.CarID
}

func (e *CarNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e *CarNotFound) Error() string {
	return fmt.Sprintf("car with ID %d not found", e.ID)
}

func findCar(r *http.Request, carIDRaw string) (*graph.Car, psi.StatusCode) {
	carID, err := graph.ParseCarID(carIDRaw)
	if err != nil {
		return nil, &IllegalQueryParam{name: "(.findCar.)", value: carIDRaw, inner: err}
	}

	carsMap := r.Context().Value(graph.CarsKey{}).(*graph.CarsMap)
	carsMap.Lock.Lock()
	defer carsMap.Lock.Unlock()

	car, ok := carsMap.Cars[carID]
	if !ok {
		carsMap.Cars[carID] = &graph.Car{
			ID:       carID,
			Trip:     graph.Trip{From: graph.SENTINAL_NODE, To: graph.SENTINAL_NODE},
			Pos:      graph.CarPos{From: 0, To: graph.SENTINAL_NODE, Ratio: 0.0},
			CurrPath: nil,
		}

		return carsMap.Cars[carID], nil
	}

	return car, nil
}

func AddApiRoutes(r psi.Router) error {
	r.WithGroup("/api", func(r psi.Router) {
		r.WithGroup("/car", func(r psi.Router) {
			r.Get("/{car_id}", func(r *http.Request) psi.StatusCode {
				car, queryErr := findCar(r, chi.URLParam(r, "car_id"))
				if queryErr != nil {
					return queryErr
				}

				type Response struct {
					psi.OkData
					Pos graph.CarPos `json:"pos"`
				}

				return &Response{Pos: car.Pos}
			})
			r.Post("/{car_id}/trip", func(r *http.Request) psi.StatusCode {
				car, queryErr := findCar(r, chi.URLParam(r, "car_id"))
				if queryErr != nil {
					return queryErr
				}

				from, queryErr := parseNodeIDParam(r, "from")
				if queryErr != nil {
					return queryErr
				}

				to, queryErr := parseNodeIDParam(r, "to")
				if queryErr != nil {
					return queryErr
				}

				if car.Trip.To != graph.SENTINAL_NODE {
					return &psi.BadRequestError{Inner: fmt.Errorf("car %d is already on a trip", car.ID)}
				}

				car.Trip = graph.Trip{From: from, To: to}

				return &psi.NoData{}
			})
		})

		r.WithGroup("/cars", func(r psi.Router) {
			r.Get("/", func(r *http.Request) psi.StatusCode {
				cars := make([]*graph.Car, 0, 16)

				idsStr := r.URL.Query().Get("ids")
				if idsStr == "" {
					carsMap := r.Context().Value(graph.CarsKey{}).(*graph.CarsMap)
					carsMap.Lock.Lock()
					defer carsMap.Lock.Unlock()

					for _, car := range carsMap.Cars {
						cars = append(cars, car)
					}
				} else {
					ids := strings.Split(idsStr, ",")
					carsMap := r.Context().Value(graph.CarsKey{}).(*graph.CarsMap)
					carsMap.Lock.Lock()
					defer carsMap.Lock.Unlock()

					for _, idStr := range ids {
						id, err := graph.ParseCarID(idStr)
						if err != nil {
							return &IllegalQueryParam{name: "ids", value: idsStr, inner: err}
						}
						if car, ok := carsMap.Cars[id]; ok {
							cars = append(cars, car)
						} else {
							cars = append(cars, &graph.Car{
								ID:       id,
								Trip:     graph.Trip{From: graph.SENTINAL_NODE, To: graph.SENTINAL_NODE},
								Pos:      graph.CarPos{From: 0, To: graph.SENTINAL_NODE, Ratio: 0.0},
								CurrPath: nil,
							})
						}
					}
				}

				type Response struct {
					psi.OkData
					Cars []*graph.Car `json:"cars"`
				}

				return &Response{Cars: cars}
			})
		})

		r.WithGroup("/graph", func(r psi.Router) {
			r.Get("/", func(r *http.Request) psi.StatusCode {
				g := ParseMap("./map.txt")

				type Response struct {
					psi.OkData
					Graph graph.Graph `json:"graph"`
				}

				return &Response{Graph: g}
			})
		})
	})
	return nil
}
