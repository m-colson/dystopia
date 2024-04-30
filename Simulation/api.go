package main

import (
	"fmt"
	"net/http"

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

func queryParam(r *http.Request, name string) (string, psi.StatusCode) {
	value := r.URL.Query()[name]
	if len(value) == 0 {
		return "", &MissingQueryParam{name: name}
	}
	return value[0], nil
}

func parseNodeIDParam(r *http.Request, name string) (NodeID, psi.StatusCode) {
	value, queryErr := queryParam(r, name)
	if queryErr != nil {
		return SENTINAL_NODE, queryErr
	}
	id, err := ParseNodeID(value)
	if err != nil {
		return SENTINAL_NODE, &psi.BadRequestError{Inner: err}
	}
	return id, nil
}

type CarNotFound struct {
	ID CarID
}

func (e *CarNotFound) StatusCode() int {
	return http.StatusNotFound
}

func (e *CarNotFound) Error() string {
	return fmt.Sprintf("car with ID %d not found", e.ID)
}

func findCar(r *http.Request, carIDRaw string) (*Car, psi.StatusCode) {
	carID, err := ParseCarID(carIDRaw)
	if err != nil {
		return nil, &psi.BadRequestError{Inner: err}
	}

	carsMap := r.Context().Value(CarsKey{}).(*CarsMap)
	carsMap.Lock.Lock()
	defer carsMap.Lock.Unlock()

	car, ok := carsMap.Cars[carID]
	if !ok {
		return nil, &CarNotFound{ID: carID}
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
					Pos CarPos `json:"pos"`
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

				if car.Trip.To != SENTINAL_NODE {
					return &psi.BadRequestError{Inner: fmt.Errorf("car %d is already on a trip", car.ID)}
				}

				car.Trip = Trip{From: from, To: to}

				return &psi.NoData{}
			})
		})

		r.WithGroup("/cars", func(r psi.Router) {
			r.Get("/", func(r *http.Request) psi.StatusCode {
				carsMap := r.Context().Value(CarsKey{}).(*CarsMap)
				carsMap.Lock.Lock()
				defer carsMap.Lock.Unlock()

				cars := make([]*Car, 0, len(carsMap.Cars))
				for _, car := range carsMap.Cars {
					cars = append(cars, car)
				}

				type Response struct {
					psi.OkData
					Cars []*Car `json:"cars"`
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
