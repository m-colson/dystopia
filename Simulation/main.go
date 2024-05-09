package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lampctl/go-sse"

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

func AddEventsRoute(server *sse.Handler) func(psi.Router) error {
	return func(r psi.Router) error {
		r.MountRaw("/events", server)

		return nil
	}
}

func FixWriteTimeout(cfg *http.Server) error {
	cfg.WriteTimeout = 0

	return nil
}

const SCHEDULER_HOST = "http://scheduler:5000"
const PFAAS_HOST = "http://pfaas:9080"

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

	g := ParseMap("./map.txt")

	eventServer := sse.NewHandler(nil)
	defer eventServer.Close()

	purgatory := make(chan graph.CarID, 64)

	simClock := time.Tick(100 * time.Millisecond)
	go func() {
		for range simClock {
			func() {
				cars.Lock.Lock()
				defer cars.Lock.Unlock()

				prevCars := len(cars.Cars)

				flaggedCars := SimulateOnce(&cars, g)
				for carId := range flaggedCars {
					purgatory <- carId
					delete(cars.Cars, carId)
				}

				type message struct {
					Cars map[graph.CarID]*graph.Car `json:"cars"`
				}

				carData := strings.Builder{}
				if err := json.NewEncoder(&carData).Encode(message{cars.Cars}); err != nil {
					log.Printf("car serialization failed because: %s\n", err)
					return
				}

				if prevCars != 0 {
					eventServer.Send(&sse.Event{
						Type: "tick",
						Data: carData.String(),
					})
				}
			}()
		}
	}()

	go func() {
		for id := range purgatory {
			req, err := http.NewRequest("PUT", fmt.Sprintf(
				"%s/ride/at/dest?id=%d",
				SCHEDULER_HOST,
				id), nil)
			if err != nil {
				log.Printf("request creation failed because: %s\n", err)
				purgatory <- id
				continue
			}

			resp, err := (&http.Client{}).Do(req)
			if err != nil {
				log.Printf("car %d failed to escape purgatory because: %s\n", id, err)
				purgatory <- id
				continue
			}

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				continue
			}

			msg := strings.Builder{}
			io.Copy(&msg, resp.Body)

			log.Printf(
				"car %d failed to escape purgatory because code %d: %s\n",
				id, resp.StatusCode, msg.String())
			purgatory <- id
		}
	}()

	psi.New[*psi.PsiServer](
		backend.Register,
		psi.AddTLSCert("cert.pem", "key.pem"),
		psi.Use(psi.LogRecoverer, RequestInsertCars(&cars), RequestInsertGraph(&g)),
		AddApiRoutes,
		AddEventsRoute(eventServer),
		FixWriteTimeout,
		AddFrontendRoutes,
	).ServeTLS(":9081").OrFatal()
}
