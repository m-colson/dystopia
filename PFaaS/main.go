package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func RequestInsertGraph() func(http.Handler) http.Handler {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/graph", WORLD_STATE_HOST), nil)
	if err != nil {
		panic(err)
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	type message struct {
		Graph graph.Graph `json:"graph"`
	}

	g := message{}

	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		panic(err)
	}

	fmt.Println("loaded", g.Graph)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), GraphKey{}, g.Graph))
			next.ServeHTTP(w, r)
		})
	}
}

type OutputLogger struct {
	http.ResponseWriter
	Output *bytes.Buffer
}

func (ol *OutputLogger) Write(buf []byte) (written int, err error) {
	written, err = ol.ResponseWriter.Write(buf)
	if err == nil {
		ol.Output.Write(buf)
	}

	return
}

func LogOutput(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := &OutputLogger{w, new(bytes.Buffer)}
		next.ServeHTTP(logger, r)

		fmt.Println(">>> " + logger.Output.String())
	})
}

func AddRoutes(extraMiddlewares ...func(http.Handler) http.Handler) func(r psi.Router) error {
	return func(r psi.Router) error {
		r.Use(psi.LogRecoverer)
		r.Use(extraMiddlewares...)
		r.Use(LogOutput)

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
					Path []graph.Edge `json:"path"`
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
					Id graph.CarID `json:"id"`
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

				fmt.Println("given", carOptions)

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

				for _, car := range cars {
					if car.Pos.From == node {
						return &ClosestResponse{Id: car.ID}
					}
				}

				panic("found a node that was not a car's position?")
			})
		})
		return nil
	}
}

const WORLD_STATE_HOST = "http://frontend:9081"

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

	type message struct {
		Cars []graph.Car `json:"cars"`
	}

	out := message{}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out.Cars, nil
}

func main() {
	psi.New[*psi.PsiServer](
		backend.Register,
		AddRoutes(RequestInsertGraph()),
	).Serve("pfaas:9080").OrFatal()
}
