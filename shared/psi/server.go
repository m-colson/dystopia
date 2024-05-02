package psi

type Server interface {
	LiftError(error) ServeError
	Serve() ServeError
}

type TLSServer interface {
	Server
	ServeTLS() ServeError
}

type Builder[Self any, T Server] interface {
	Complete(addr string) (T, error)
	WithInit(any) Self
}

type TLSBuilder[Self any, T TLSServer] interface {
	Builder[Self, T]
}

type ImmBuilder[T Server] interface {
	Builder[ImmBuilder[T], T]
	Serve(addr string) ServeError
}

type ImmTLSBuilder[T TLSServer] interface {
	ImmBuilder[T]
	ServeTLS(addr string) ServeError
}
