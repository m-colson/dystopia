package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/m-colson/psi"
	backend "github.com/m-colson/psi/backend-chi"
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

func parseQueryParam(r *http.Request, name string) (NodeID, psi.StatusCode) {
	value, queryErr := queryParam(r, name)
	if queryErr != nil {
		return 0, queryErr
	}
	id, err := ParseNodeID(value)
	if err != nil {
		return 0, &psi.BadRequestError{Inner: err}
	}
	return id, nil

}

type GraphKey struct{}

func RequestInsertGraph(next http.Handler) http.Handler {
	graph := NewGraph([]Link{
		{1, 5, 2},
		{1, 1, 3},
		{2, 1, 4},
		{2, 1, 5},
		{3, 1, 5},
		{3, 2, 7},
		{4, 1, 6},
		{5, 2, 6},
		{6, 1, 8},
		{7, 6, 8},
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
				Path []NodeID `json:"path"`
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

			graph := r.Context().Value(GraphKey{}).(Graph)
			path := Dijkstra(graph, from, to)

			return &PathResponse{Path: path}
		})
		r.Get("/closest", func(r *http.Request) psi.StatusCode {
			type ClosestResponse struct {
				psi.OkData
				Id NodeID `json:"id"`
			}

			to, queryErr := parseQueryParam(r, "to")
			if queryErr != nil {
				return queryErr
			}

			optionsRaw, queryErr := queryParam(r, "options")
			if queryErr != nil {
				return queryErr
			}
			optionsStrs := strings.Split(optionsRaw, ",")

			options := make([]NodeID, 0, len(optionsStrs))
			for _, optionStr := range optionsStrs {
				option, err := ParseNodeID(optionStr)
				if err != nil {
					return &psi.BadRequestError{Inner: err}
				}
				options = append(options, option)
			}

			graph := r.Context().Value(GraphKey{}).(Graph)
			node, ok := DijkstraClosest(graph, to, options...)
			if !ok {
				return &psi.NotAcceptableError{}
			}

			return &ClosestResponse{Id: node}
		})
	})
	return nil
}

func main() {
	psi.New[*psi.PsiServer](
		backend.Register,
		AddRoutes,
	).Serve("localhost:9080").OrFatal()
}
