package psi

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

var builders = make(map[reflect.Type]any)

func RegisterBuild[T Server](cb func(*http.Server) T) {
	if _, ok := builders[reflect.TypeFor[T]()]; ok {
		panic(fmt.Errorf("builder already registered for type %T", cb))
	}
	builders[reflect.TypeFor[T]()] = cb
}

func init() {
	RegisterBuild(func(s *http.Server) *PsiServer {
		return &PsiServer{Server: s}
	})
}

var _ TLSBuilder[*PsiBuilder[*PsiServer], *PsiServer] = (*PsiBuilder[*PsiServer])(nil)

type PsiBuilder[T Server] struct {
	errors    []error
	router    Router
	server    *http.Server
	tlsConfig *tls.Config
}

func New[T Server](inits ...any) *PsiBuilder[T] {
	out := new(PsiBuilder[T])
	initBuilderDefaults(out)
	for _, init := range inits {
		out.WithInit(init)
	}

	return out
}

func initBuilderDefaults[T Server](m *PsiBuilder[T]) {
	m.errors = make([]error, 0, 16)
	m.server = &http.Server{
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	m.tlsConfig = &tls.Config{
		Certificates:             []tls.Certificate{},
		CipherSuites:             nil,
		PreferServerCipherSuites: false,
		MinVersion:               tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
		},
	}
}

func (m *PsiBuilder[T]) AddError(err error) {
	m.errors = append(m.errors, err)
}

func (m *PsiBuilder[T]) Complete(addr string) (out T, err error) {
	if len(m.errors) > 0 {
		err = errors.Join(m.errors...)
		return
	}

	if m.tlsConfig != nil {
		m.server.TLSConfig = m.tlsConfig
		m.server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
	}

	m.server.Handler = m.router
	m.server.Addr = addr

	if cb, ok := builders[reflect.TypeFor[T]()]; ok {
		return cb.(func(*http.Server) T)(m.server), nil
	}

	panic(fmt.Errorf("no builder registered for type %T", out))
}

type RouterCreater = func(*Router) error
type RouterIniter = func(Router) error
type ServerIniter = func(*http.Server) error
type TLSIniter = func(*tls.Config) error

func (m *PsiBuilder[T]) WithInit(cb any) *PsiBuilder[T] {
	switch cb := cb.(type) {
	case RouterCreater:
		if m.router != nil {
			panic(errors.New("router already initialized"))
		}
		if err := cb(&m.router); err != nil {
			m.AddError(err)
		}
		return m
	case RouterIniter:
		if m.router == nil {
			panic(errors.New("router must be initialized before calling using it"))
		}
		if err := cb(m.router); err != nil {
			m.AddError(err)
		}
		return m
	case ServerIniter:
		if err := cb(m.server); err != nil {
			m.AddError(err)
		}
		return m
	case TLSIniter:
		if err := cb(m.tlsConfig); err != nil {
			m.AddError(err)
		}
		return m
	default:
		panic(fmt.Errorf("invalid init function type (%T)", cb))
	}
}

func (m *PsiBuilder[T]) Serve(addr string) ServeError {
	s, err := m.Complete(addr)
	if err != nil {
		return s.LiftError(err)
	}

	return s.Serve()
}

// func (m *PsiBuilder[T]) ServeTLS(addr string) ServeError {
// 	s, err := m.Complete(addr)
// 	if err != nil {
// 		return s.LiftError(err)
// 	}

// 	if m.tlsConfig == nil {
// 		return s.LiftError(err)
// 	}

// 	return s.ServeTLS()
// }
