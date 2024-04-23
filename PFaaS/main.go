package main

import (
	"net/http"

	"github.com/m-colson/psi"
	backend "github.com/m-colson/psi/backend-chi"
)

func AddRoutes(r psi.Router) error {
	r.Use(psi.LogRecoverer)

	r.WithGroup("/api", func(r psi.Router) {
		r.Get("/", func(r *http.Request) psi.StatusCode {
			type Metadata struct {
				psi.OkData
				Versions []string `json:"versions"`
			}

			return &Metadata{Versions: []string{"0.0.1"}}
		})
		r.Get("/solve", func(r *http.Request) psi.StatusCode {
			return &psi.NotAcceptableError{}
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
