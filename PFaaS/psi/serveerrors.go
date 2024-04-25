package psi

import "log"

type ChainableError interface {
	error
	OrFatal()
}

type ServeError interface {
	ChainableError
	Cleanup() ServeError
}

type ChainError struct {
	error
}

func (e ChainError) OrFatal() {
	if e.error != nil {
		log.Fatalf("fatal error: %v", e.error)
	}
}

type nothingServeError struct {
	ChainableError
}

func (n *nothingServeError) Cleanup() ServeError {
	return n
}
