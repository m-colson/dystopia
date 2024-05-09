package main

import (
	"net/http"

	"github.com/m-colson/psi"
)

func AddFrontendRoutes(r psi.Router) error {
	r.GetRaw("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	// r.MountRaw("/proxy/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	url := r.URL.Path[len("/proxy/"):]

	// 	proxyReq, err := http.NewRequest(r.Method, url, r.Body)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	proxyReq.Header = r.Header

	// 	client := http.Client{}
	// 	resp, err := client.Do(proxyReq)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// 	defer resp.Body.Close()

	// 	for k, v := range resp.Header {
	// 		for _, v := range v {
	// 			w.Header().Add(k, v)
	// 		}
	// 	}
	// 	w.WriteHeader(resp.StatusCode)
	// 	io.Copy(w, resp.Body)
	// }))

	return nil
}
