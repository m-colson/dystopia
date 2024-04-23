module github.com/m-colson/dystopia/pfaas

go 1.22.0

replace github.com/m-colson/psi => ./psi

replace github.com/m-colson/psi/backend-chi => ./psi/backend-chi

require (
	github.com/m-colson/psi v0.0.0-00010101000000-000000000000
	github.com/m-colson/psi/backend-chi v0.0.0-00010101000000-000000000000
)

require github.com/go-chi/chi/v5 v5.0.12 // indirect
