package main

import (
	"net/http"

	"github.com/m-colson/psi"
)

func AddFrontendRoutes(r psi.Router) error {
	r.GetRaw("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	return nil
}
