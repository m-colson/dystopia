package psi

import "net/http"

type Router interface {
	http.Handler

	Use(...func(http.Handler) http.Handler)

	MethodFunc(method, pattern string, h HandlerFunc)
	Connect(pattern string, h HandlerFunc)
	Delete(pattern string, h HandlerFunc)
	Get(pattern string, h HandlerFunc)
	Head(pattern string, h HandlerFunc)
	Options(pattern string, h HandlerFunc)
	Patch(pattern string, h HandlerFunc)
	Post(pattern string, h HandlerFunc)
	Put(pattern string, h HandlerFunc)
	Trace(pattern string, h HandlerFunc)

	MethodFuncRaw(method, pattern string, h http.HandlerFunc)
	ConnectRaw(pattern string, h http.HandlerFunc)
	DeleteRaw(pattern string, h http.HandlerFunc)
	GetRaw(pattern string, h http.HandlerFunc)
	HeadRaw(pattern string, h http.HandlerFunc)
	OptionsRaw(pattern string, h http.HandlerFunc)
	PatchRaw(pattern string, h http.HandlerFunc)
	PostRaw(pattern string, h http.HandlerFunc)
	PutRaw(pattern string, h http.HandlerFunc)
	TraceRaw(pattern string, h http.HandlerFunc)

	MountRaw(pattern string, h http.Handler)
	WithGroup(pattern string, cb func(Router))
}
