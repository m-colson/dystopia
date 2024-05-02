package psi

import (
	"fmt"
	"log"
	"net/http"
)

var _ TLSServer = (*PsiServer)(nil)

type PsiServer struct {
	Server *http.Server
}

func (m *PsiServer) LiftError(err error) ServeError {
	return &nothingServeError{&ChainError{err}}
}

func (m *PsiServer) Serve() ServeError {
	log.Println("serving on", m.Server.Addr)
	return m.LiftError(m.Server.ListenAndServe())
}

func (m *PsiServer) ServeTLS() ServeError {
	if m.Server.TLSConfig == nil {
		return m.LiftError(fmt.Errorf("a tls config is required to serve tls"))
	}

	return m.LiftError(m.Server.ListenAndServeTLS("", ""))
}

type RedirectHTTPAddr string

func (addr RedirectHTTPAddr) StartBackground() {
	go func() {
		err := http.ListenAndServe(string(addr), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
		}))

		if err != nil {
			log.Fatal(err)
		}
	}()
}
