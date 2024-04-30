package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/m-colson/dystopia/shared/graph"
	"github.com/m-colson/psi"
	backend "github.com/m-colson/psi/backend-chi"
)

func RequestInsertCars(cars *graph.CarsMap) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), graph.CarsKey{}, cars))
			next.ServeHTTP(w, r)
		})
	}
}

type GraphKey struct{}

func RequestInsertGraph(g *graph.Graph) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), GraphKey{}, g))
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	generate := flag.Bool("g", false, "generate new map")
	flag.Parse()
	if *generate {
		file, err := os.Create("./map.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		GenerateMap(file, 50, 100, 3)
		return
	}

	cars := graph.CarsMap{
		Cars: make(map[graph.CarID]*graph.Car),
		Lock: sync.Mutex{},
	}

	graph := ParseMap("./map.txt")

	simClock := time.Tick(1 * time.Second)
	go func() {
		for range simClock {
			SimulateOnce(&cars, graph)
		}
	}()

	psi.New[*psi.PsiServer](
		backend.Register,
		psi.Use(psi.LogRecoverer, RequestInsertCars(&cars), RequestInsertGraph(&graph)),
		AddApiRoutes,
		AddFrontendRoutes,
	).Serve("localhost:9081").OrFatal()
}
