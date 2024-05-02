module dystopia

go 1.22.0

replace github.com/m-colson/psi => ../shared/psi

replace github.com/m-colson/psi/backend-chi => ../shared/psi/backend-chi

replace github.com/m-colson/dystopia/shared/graph => ../shared/graph

require (
	github.com/m-colson/psi v0.0.0-00010101000000-000000000000
	github.com/m-colson/psi/backend-chi v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-chi/chi/v5 v5.0.12
	github.com/lampctl/go-sse v1.1.4
	github.com/m-colson/dystopia/shared/graph v0.0.0-00010101000000-000000000000
)
