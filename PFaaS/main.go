package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/m-colson/psi"
	backend "github.com/m-colson/psi/backend-chi"

	"github.com/m-colson/dystopia/shared/graph"
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

func parseQueryParam(r *http.Request, name string) (graph.NodeID, psi.StatusCode) {
	value, queryErr := queryParam(r, name)
	if queryErr != nil {
		return 0, queryErr
	}
	id, err := graph.ParseID(value)
	if err != nil {
		return 0, &IllegalQueryParam{
			name:  name,
			value: value,
			inner: err,
		}
	}
	return id, nil

}

type GraphKey struct{}

func RequestInsertGraph(next http.Handler) http.Handler {
	graph := graph.New([]graph.Link{
		{From: 1, Cost: 5, To: 2},
		{From: 1, Cost: 1, To: 3},
		{From: 2, Cost: 1, To: 4},
		{From: 2, Cost: 1, To: 5},
		{From: 3, Cost: 1, To: 5},
		{From: 3, Cost: 2, To: 7},
		{From: 4, Cost: 1, To: 6},
		{From: 5, Cost: 2, To: 6},
		{From: 6, Cost: 1, To: 8},
		{From: 7, Cost: 6, To: 8},
	}...)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), GraphKey{}, graph))
		next.ServeHTTP(w, r)
	})
}

func AddRoutes(r psi.Router) error {
	r.Use(psi.LogRecoverer, RequestInsertGraph)

	r.WithGroup("/api", func(r psi.Router) {
		r.Get("/", func(r *http.Request) psi.StatusCode {
			type Metadata struct {
				psi.OkData
				Versions []string `json:"versions"`
			}

			return &Metadata{Versions: []string{"0.0.1"}}
		})
		r.Get("/path", func(r *http.Request) psi.StatusCode {
			type PathResponse struct {
				psi.OkData
				Path []graph.NodeID `json:"path"`
				// TotalTime float64  `json:"totalTime"`
			}

			from, queryErr := parseQueryParam(r, "from")
			if queryErr != nil {
				return queryErr
			}

			to, queryErr := parseQueryParam(r, "to")
			if queryErr != nil {
				return queryErr
			}

			graph := r.Context().Value(GraphKey{}).(graph.Graph)
			path := Dijkstra(graph, from, to)

			return &PathResponse{Path: path}
		})
		r.Get("/closest", func(r *http.Request) psi.StatusCode {
			type ClosestResponse struct {
				psi.OkData
				Id graph.NodeID `json:"id"`
			}

			to, queryErr := parseQueryParam(r, "to")
			if queryErr != nil {
				return queryErr
			}

			carOptionsRaw, queryErr := queryParam(r, "options")
			if queryErr != nil {
				return queryErr
			}
			carOptionsStrs := strings.Split(carOptionsRaw, ",")

			carOptions := make([]graph.CarID, 0, len(carOptionsStrs))
			for _, optionStr := range carOptionsStrs {
				option, err := graph.ParseCarID(optionStr)
				if err != nil {
					return &IllegalQueryParam{name: "options", value: carOptionsRaw, inner: err}
				}
				carOptions = append(carOptions, option)
			}

			cars, err := findCarIDS(carOptions...)
			if err != nil {
				panic(err)
			}

			options := make([]graph.NodeID, 0, len(cars))
			for _, car := range cars {
				options = append(options, car.Pos.From)
			}

			graph := r.Context().Value(GraphKey{}).(graph.Graph)
			node, ok := DijkstraClosest(graph, to, options...)
			if !ok {
				return &psi.NotAcceptableError{}
			}

			return &ClosestResponse{Id: node}
		})
	})
	return nil
}

const WORLD_STATE_HOST = "http://localhost:9081"

func findCarIDS(ids ...graph.CarID) ([]graph.Car, error) {
	client := http.Client{}

	idsStr := strings.Builder{}
	for i, id := range ids {
		if i > 0 {
			idsStr.WriteString(",")
		}
		fmt.Fprintf(&idsStr, "%d", id)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/cars?ids=%s", WORLD_STATE_HOST, idsStr.String()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	out := make([]graph.Car, 0, len(ids))

	if respBytes, err := io.ReadAll(resp.Body); err == nil {
		fmt.Printf("%s\n", respBytes)
	} else {
		fmt.Println(err)
	}

	// if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
	// 	return nil, err
	// }

	return out, nil
}

func main() {
	psi.New[*psi.PsiServer](
		backend.Register,
		AddRoutes,
	).Serve("localhost:9080").OrFatal()
}
