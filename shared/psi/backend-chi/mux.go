package psichi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/m-colson/psi"
)

var _ psi.Router = (*Mux)(nil)

type Mux struct {
	*chi.Mux
}

func NewMux() *Mux {
	r := chi.NewRouter()

	r.NotFound(psi.Responder(func(r *http.Request) psi.StatusCode {
		return &psi.NotFoundError{}
	}))
	r.MethodNotAllowed(psi.Responder(func(r *http.Request) psi.StatusCode {
		return &psi.MethodNotAllowedError{}
	}))

	return &Mux{r}
}

func (mx *Mux) MethodFunc(method, pattern string, h psi.HandlerFunc) {
	mx.Mux.MethodFunc(method, pattern, psi.Responder(h))
}
func (mx *Mux) Connect(pattern string, h psi.HandlerFunc) {
	mx.ConnectRaw(pattern, psi.Responder(h))
}
func (mx *Mux) Delete(pattern string, h psi.HandlerFunc) {
	mx.DeleteRaw(pattern, psi.Responder(h))
}
func (mx *Mux) Get(pattern string, h psi.HandlerFunc) {
	mx.GetRaw(pattern, psi.Responder(h))
}
func (mx *Mux) Head(pattern string, h psi.HandlerFunc) {
	mx.HeadRaw(pattern, psi.Responder(h))
}
func (mx *Mux) Options(pattern string, h psi.HandlerFunc) {
	mx.OptionsRaw(pattern, psi.Responder(h))
}
func (mx *Mux) Patch(pattern string, h psi.HandlerFunc) {
	mx.PatchRaw(pattern, psi.Responder(h))
}
func (mx *Mux) Post(pattern string, h psi.HandlerFunc) {
	mx.PostRaw(pattern, psi.Responder(h))
}
func (mx *Mux) Put(pattern string, h psi.HandlerFunc) {
	mx.PutRaw(pattern, psi.Responder(h))
}
func (mx *Mux) Trace(pattern string, h psi.HandlerFunc) {
	mx.TraceRaw(pattern, psi.Responder(h))
}

func (mx *Mux) MethodFuncRaw(method, pattern string, h http.HandlerFunc) {
	mx.Mux.MethodFunc(method, pattern, h)
}
func (mx *Mux) ConnectRaw(pattern string, h http.HandlerFunc) {
	mx.Mux.Connect(pattern, h)
}
func (mx *Mux) DeleteRaw(pattern string, h http.HandlerFunc) {
	mx.Mux.Delete(pattern, h)
}
func (mx *Mux) GetRaw(pattern string, h http.HandlerFunc) {
	mx.Mux.Get(pattern, h)
}
func (mx *Mux) HeadRaw(pattern string, h http.HandlerFunc) {
	mx.Mux.Head(pattern, h)
}
func (mx *Mux) OptionsRaw(pattern string, h http.HandlerFunc) {
	mx.Mux.Options(pattern, h)
}
func (mx *Mux) PatchRaw(pattern string, h http.HandlerFunc) {
	mx.Mux.Patch(pattern, h)
}
func (mx *Mux) PostRaw(pattern string, h http.HandlerFunc) {
	mx.Mux.Post(pattern, h)
}
func (mx *Mux) PutRaw(pattern string, h http.HandlerFunc) {
	mx.Mux.Put(pattern, h)
}
func (mx *Mux) TraceRaw(pattern string, h http.HandlerFunc) {
	mx.Mux.Trace(pattern, h)
}

func (mx *Mux) MountRaw(pattern string, h http.Handler) {
	mx.Mux.Mount(pattern, h)
}

func (mx *Mux) WithGroup(pattern string, cb func(psi.Router)) {
	group := NewMux()
	cb(group)
	mx.Mux.Mount(pattern, group)
}
